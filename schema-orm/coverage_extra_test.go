package schema_orm

import (
	"reflect"
	"testing"
	"time"

	xs "xorm.io/xorm/schemas"
)

type sample2 struct {
	UID uint
}

type blobHolder struct{ B []byte }

type boolHolder struct{ V bool }

func TestPK_isZero_VariousTypes(t *testing.T) {
	if !isZero(nil) {
		t.Fatalf("nil should be zero")
	}
	if !isZero("") {
		t.Fatalf("empty string zero")
	}
	if isZero("x") {
		t.Fatalf("non-empty string not zero")
	}
	if isZero(true) {
		t.Fatalf("true not zero")
	}
	if !isZero(false) {
		t.Fatalf("false is zero")
	}
	if isZero(1) {
		t.Fatalf("int non-zero")
	}
	if !isZero(0) {
		t.Fatalf("int zero")
	}
	if isZero(uint64(5)) {
		t.Fatalf("uint64 non-zero")
	}
	if !isZero(uint16(0)) {
		t.Fatalf("uint16 zero")
	}
	if isZero(float32(1.5)) {
		t.Fatalf("float32 non-zero")
	}
	if !isZero(float64(0)) {
		t.Fatalf("float64 zero")
	}
}

func TestTypes_SQLType2Type_AllCases(t *testing.T) {
	cases := []struct {
		in   SQLType
		kind reflect.Kind
	}{
		{SQLType{Name: "VARCHAR"}, reflect.String},
		{SQLType{Name: "INT"}, reflect.Int},
		{SQLType{Name: "BOOLEAN"}, reflect.Bool},
		{SQLType{Name: "FLOAT"}, reflect.Float64},
		{SQLType{Name: "BLOB"}, reflect.Slice},
		{SQLType{Name: "UNKNOWN"}, reflect.String},
	}
	for _, c := range cases {
		if got := SQLType2Type(c.in).Kind(); got != c.kind {
			t.Fatalf("%s -> kind %v, want %v", c.in.Name, got, c.kind)
		}
	}
}

func TestTypes_Type2SQLType_MoreKinds(t *testing.T) {
	if Type2SQLType(reflect.TypeOf(true)).Name != "BOOLEAN" {
		t.Fatalf("bool mapping")
	}
	if Type2SQLType(reflect.TypeOf(float64(0))).Name != "FLOAT" {
		t.Fatalf("float mapping")
	}
	// default mapping for struct goes to TEXT in simplified impl
	type S struct{ A int }
	if Type2SQLType(reflect.TypeOf(S{})).Name != "TEXT" {
		t.Fatalf("default struct mapping")
	}
}

func TestTypes_IsBlob_IsTime_FalsePaths(t *testing.T) {
	st := SQLType{Name: "TEXT"}
	if st.IsBlob() {
		t.Fatalf("TEXT is not blob")
	}
	if st.IsTime() {
		t.Fatalf("TEXT is not time")
	}
}

func TestConvert_NilInputs(t *testing.T) {
	if ToXormColumn(nil) != nil {
		t.Fatalf("nil column")
	}
	if FromXormColumn(nil) != nil {
		t.Fatalf("nil column from")
	}
	if ToXormIndex(nil) != nil {
		t.Fatalf("nil index")
	}
	if FromXormIndex(nil) != nil {
		t.Fatalf("nil index from")
	}
	if ToXormTable(nil) != nil {
		t.Fatalf("nil table to")
	}
	if FromXormTable(nil) != nil {
		t.Fatalf("nil table from")
	}
}

func TestTable_GetColumn_Negative_and_UintPKPath(t *testing.T) {
	// negative
	typ := reflect.TypeOf(sample{})
	tb := NewTable("t", typ)
	if tb.GetColumn("missing") != nil {
		t.Fatalf("expected nil get column")
	}
	if tb.GetColumnIdx("missing", 0) != nil {
		t.Fatalf("expected nil get column idx")
	}
	// uint PK path in IDOfV
	typ2 := reflect.TypeOf(sample2{})
	tb2 := NewTable("t2", typ2)
	col := &Column{Name: "UID", FieldIndex: []int{0}, IsPrimaryKey: true, SQLType: SQLType{Name: "INT"}}
	tb2.AddColumn(col)
	obj := sample2{UID: 42}
	pk, err := tb2.IDOfV(reflect.ValueOf(obj))
	if err != nil {
		t.Fatal(err)
	}
	if len(pk) != 1 || pk[0].(int64) != 42 {
		t.Fatalf("uint pk conversion failed: %v", pk)
	}
}

func TestIndex_XName_SchemaQualified(t *testing.T) {
	i := NewIndex("nm", UniqueType)
	name := i.XName("public.\"tbl\"")
	if name[:4] != "UQE_" {
		t.Fatalf("expected unique prefix for schema-qualified table")
	}
}

func TestConvert_Table_MapsAndIndexesFilled(t *testing.T) {
	tbl := NewTable("conv", reflect.TypeOf(sample{}))
	c := &Column{Name: "ID"}
	tbl.AddColumn(c)
	tbl.Created["ID"] = true
	i := NewIndex("idx", IndexType)
	i.AddColumn("ID")
	tbl.AddIndex(i)
	x := ToXormTable(tbl)
	if x.Indexes["idx"] == nil {
		t.Fatalf("indexes not copied")
	}
	if !x.Created["ID"] {
		t.Fatalf("created not copied")
	}
	// back
	b := FromXormTable(x)
	if len(b.Columns()) != 1 || len(b.Indexes) != 1 || !b.Created["ID"] {
		t.Fatalf("round conversion failed")
	}
}

func TestColumn_ConvertID_InvalidNumber(t *testing.T) {
	c := &Column{SQLType: SQLType{Name: "INT"}}
	if _, err := c.ConvertID("not-a-number"); err == nil {
		t.Fatalf("expected parse error")
	}
}

func TestMarshal_Table_EmptyCollections(t *testing.T) {
	// Ensure Unmarshal paths that directly assign maps/slices also work on empties
	// Create an empty table DTO via json/yaml roundtrip through our methods already covered elsewhere.
	// Here, just make sure the internal zero-value maps are usable after NewTable.
	tbl := NewEmptyTable()
	if tbl.Created == nil || tbl.Indexes == nil {
		t.Fatalf("zero table should have maps initialized")
	}
}

func TestConvertSQLType_Preservation(t *testing.T) {
	s := SQLType{Name: "DECIMAL", DefaultLength: 10, DefaultLength2: 2}
	x := ToXormSQLType(s)
	if x.DefaultLength != 10 || x.DefaultLength2 != 2 {
		t.Fatalf("to xorm preserve lengths")
	}
	back := FromXormSQLType(xs.SQLType{Name: s.Name, DefaultLength: 5, DefaultLength2: 1})
	if back.DefaultLength != 5 || back.DefaultLength2 != 1 {
		t.Fatalf("from xorm preserve lengths")
	}
}

func TestColumn_ValueOfV_InterfaceBranch(t *testing.T) {
	// Craft a value path to hit the interface branch inside ValueOfV
	type inner struct{ N int }
	type wrap struct{ I interface{} }
	w := wrap{I: &inner{N: 7}}
	v := reflect.ValueOf(&w)                // pointer so first branch handles pointer
	col := &Column{FieldIndex: []int{0, 0}} // first selects field I, second triggers interface branch
	val, err := col.ValueOfV(&v)
	if err != nil {
		t.Fatal(err)
	}
	_ = val
}

func TestColumn_ValueOfV_PointerAllocBranch(t *testing.T) {
	// Make v.Kind() == Ptr at loop time and ensure pointer deref path is executed safely
	type node struct{ X int }
	type holderPtr struct{ P *node }
	h := holderPtr{}
	v := reflect.ValueOf(&h) // *holderPtr
	col := &Column{FieldIndex: []int{0}}
	val, err := col.ValueOfV(&v)
	if err != nil {
		t.Fatal(err)
	}
	_ = val
	_ = time.UTC // touch time package to avoid unused in some builds
}
