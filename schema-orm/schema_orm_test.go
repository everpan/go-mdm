package schema_orm

import (
	"encoding/json"
	"reflect"
	"testing"

	yaml "gopkg.in/yaml.v3"
	xs "xorm.io/xorm/schemas"
)

type sample struct {
	ID   int
	Name string
}

func TestSQLTypeHelpers(t *testing.T) {
	st := SQLType{Name: "INT"}
	if !st.IsNumeric() {
		t.Fatalf("expected numeric")
	}
	if st.IsText() {
		t.Fatalf("not text")
	}
	st = SQLType{Name: "VARCHAR"}
	if !st.IsText() {
		t.Fatalf("expected text")
	}
	st = SQLType{Name: "JSON"}
	if !st.IsJson() {
		t.Fatalf("expected json")
	}
}

func TestColumnConvertID(t *testing.T) {
	c := NewColumn("id", "ID", SQLType{Name: "INT"}, 0, 0, true)
	v, err := c.ConvertID("123")
	if err != nil || v.(int64) != 123 {
		t.Fatalf("convert int failed: %v %v", v, err)
	}
	c.SQLType = SQLType{Name: "VARCHAR"}
	v, err = c.ConvertID("abc")
	if err != nil || v.(string) != "abc" {
		t.Fatalf("convert text failed: %v %v", v, err)
	}
}

func TestColumnValueOf(t *testing.T) {
	c := &Column{FieldIndex: []int{0}, FieldName: "ID"}
	obj := &sample{ID: 7}
	val, err := c.ValueOf(obj)
	if err != nil {
		t.Fatal(err)
	}
	if val.Int() != 7 {
		t.Fatalf("expected 7 got %v", val.Int())
	}
}

func TestIndexMethods(t *testing.T) {
	i := NewIndex("name", UniqueType)
	i.AddColumn("a", "b")
	if len(i.Cols) != 2 {
		t.Fatalf("cols")
	}
	name := i.XName("public.\"foo\"")
	if name[:4] != "UQE_" {
		t.Fatalf("xname prefix")
	}
	j := NewIndex("name", UniqueType)
	j.AddColumn("b", "a")
	if !i.Equal(j) {
		t.Fatalf("equal")
	}
}

func TestTableMethods(t *testing.T) {
	typ := reflect.TypeOf(sample{})
	tbl := NewTable("sample", typ)
	colID := &Column{Name: "ID", FieldIndex: []int{0}, IsPrimaryKey: true, SQLType: SQLType{Name: "INT"}}
	tbl.AddColumn(colID)
	colName := &Column{Name: "Name", FieldIndex: []int{1}}
	tbl.AddColumn(colName)
	if tbl.GetColumn("ID") == nil || tbl.GetColumnIdx("Name", 0) == nil {
		t.Fatalf("get column")
	}
	if len(tbl.Columns) != 2 || len(tbl.ColumnsSeq) != 2 {
		t.Fatalf("columns length")
	}
	if len(tbl.PKColumns()) != 1 {
		t.Fatalf("pk len")
	}
	obj := sample{ID: 9}
	pk, err := tbl.IDOfV(reflect.ValueOf(obj))
	if err != nil {
		t.Fatal(err)
	}
	if len(pk) != 1 || pk[0].(int64) != 9 {
		t.Fatalf("pk value: %v", pk)
	}
}

func TestPKMethods(t *testing.T) {
	p := NewPK(0)
	if !p.IsZero() {
		t.Fatalf("expected zero")
	}
	p = NewPK(1, "a")
	if p.IsZero() {
		t.Fatalf("not zero")
	}
	s, err := p.ToString()
	if err != nil {
		t.Fatal(err)
	}
	var p2 PK
	if err := p2.FromString(s); err != nil {
		t.Fatal(err)
	}
	if len(p2) != 2 {
		t.Fatalf("decode")
	}
}

func TestJSONRoundTrip(t *testing.T) {
	c := NewColumn("id", "ID", SQLType{Name: "INT"}, 0, 0, true)
	b, err := json.Marshal(c)
	if err != nil {
		t.Fatal(err)
	}
	var c2 Column
	if err := json.Unmarshal(b, &c2); err != nil {
		t.Fatal(err)
	}
	if c2.Name != c.Name || c2.SQLType.Name != c.SQLType.Name {
		t.Fatalf("json round-trip")
	}
	// Table
	typ := reflect.TypeOf(sample{})
	tb := NewTable("sample", typ)
	tb.AddColumn(&Column{Name: "ID"})
	bb, err := json.Marshal(tb)
	if err != nil {
		t.Fatal(err)
	}
	var tb2 Table
	if err := json.Unmarshal(bb, &tb2); err != nil {
		t.Fatal(err)
	}
	if tb2.Name != tb.Name || len(tb2.Columns) != 1 {
		t.Fatalf("table json")
	}
}

func TestYAMLRoundTrip(t *testing.T) {
	i := NewIndex("idx", IndexType)
	i.AddColumn("a")
	b, err := yaml.Marshal(i)
	if err != nil {
		t.Fatal(err)
	}
	var i2 Index
	if err := yaml.Unmarshal(b, &i2); err != nil {
		t.Fatal(err)
	}
	if i2.Name != i.Name || len(i2.Cols) != 1 {
		t.Fatalf("yaml round-trip")
	}
}

func TestConversions(t *testing.T) {
	// SQLType
	ast := SQLType{Name: "INT", DefaultLength: 11, DefaultLength2: 0}
	xst := ToXormSQLType(ast)
	bst := FromXormSQLType(xst)
	if bst != ast {
		t.Fatalf("sqltype conv")
	}
	// Column
	ac := &Column{Name: "ID", SQLType: ast, Nullable: true, Indexes: map[string]int{"a": 1}}
	xc := ToXormColumn(ac)
	bc := FromXormColumn(xc)
	if bc.Name != ac.Name || bc.SQLType.Name != ac.SQLType.Name || !bc.Nullable {
		t.Fatalf("column conv")
	}
	// Index
	ai := NewIndex("i", UniqueType)
	ai.AddColumn("A")
	xi := ToXormIndex(ai)
	bi := FromXormIndex(xi)
	if !ai.Equal(bi) {
		t.Fatalf("index conv")
	}
	// PK
	apk := PK{"a", 1}
	xpk := ToXormPK(apk)
	bpk := FromXormPK(xpk)
	if len(bpk) != 2 {
		t.Fatalf("pk conv")
	}
	// Table
	typ := reflect.TypeOf(sample{})
	at := NewTable("t", typ)
	at.AddColumn(&Column{Name: "ID"})
	at.AddIndex(ai)
	xt := ToXormTable(at)
	bt := FromXormTable(xt)
	if bt.Name != at.Name || len(bt.Columns) != 1 || len(bt.Indexes) != 1 {
		t.Fatalf("table conv")
	}
	// Ensure xorm.Table methods still work with our conversions
	if len(xt.Columns()) != 1 {
		t.Fatalf("xorm table columns")
	}
	// sanity use in engine-related packages avoided; only struct mapping
	_ = xs.Table{Name: "x"}
}
