package schema_orm

import (
	"reflect"
	"testing"
)

// Additional tests to push coverage over 95%
func TestPK_isZero_MoreIntegerTypes(t *testing.T) {
	if !isZero(int8(0)) { t.Fatalf("int8 zero") }
	if isZero(int16(2)) { t.Fatalf("int16 non-zero") }
	if !isZero(int32(0)) { t.Fatalf("int32 zero") }
	if isZero(uint(3)) { t.Fatalf("uint non-zero") }
	if !isZero(uint32(0)) { t.Fatalf("uint32 zero") }
}

func TestIndex_XName_NonUnique(t *testing.T) {
	i := NewIndex("nm", IndexType)
	if name := i.XName("tbl"); name[:4] != "IDX_" { t.Fatalf("expect IDX_ prefix: %s", name) }
}

type sample3 struct{ SID string }

func TestTable_IDOfV_StringPK(t *testing.T) {
	typ := reflect.TypeOf(sample3{})
	tb := NewTable("s", typ)
	tb.AddColumn(&Column{Name: "SID", FieldIndex: []int{0}, IsPrimaryKey: true, SQLType: SQLType{Name: "VARCHAR"}})
	obj := sample3{SID: "abc"}
	pk, err := tb.IDOfV(reflect.ValueOf(obj))
	if err != nil { t.Fatal(err) }
	if len(pk) != 1 || pk[0].(string) != "abc" { t.Fatalf("string pk failed: %v", pk) }
}

func TestTable_AddColumn_AppendPath(t *testing.T) {
	tb := NewTable("dup", reflect.TypeOf(sample{}))
	c1 := &Column{Name: "Dup"}
	c2 := &Column{Name: "Dup"}
	tb.AddColumn(c1)
	tb.AddColumn(c2)
	if tb.GetColumnIdx("Dup", 1) == nil { t.Fatalf("expected second column with same name") }
}
