package schema_orm

import (
	"reflect"
	"testing"

	xs "xorm.io/xorm/schemas"
)

type kb struct {
	ID   int64
	Code string
}

func TestToXormColumn_BIT_To_BOOL_and_Default(t *testing.T) {
	c := &Column{SQLType: SQLType{Name: "BIT"}, Default: "", Name: "A"}
	xc := ToXormColumn(c)
	if xc.SQLType.Name != "BOOL" {
		t.Fatalf("expected BOOL, got %s", xc.SQLType.Name)
	}
	if xc.Default != "true" {
		t.Fatalf("expected default true for BOOL")
	}
}

func TestTable_FlagColumnsAndAccessors(t *testing.T) {
	tbl := NewTable("t", reflect.TypeOf(kb{}))
	pk := &Column{Name: "ID", FieldIndex: []int{0}, IsPrimaryKey: true, SQLType: SQLType{Name: "INT"}}
	auto := &Column{Name: "ID", IsAutoIncrement: true}
	ver := &Column{Name: "VER", IsVersion: true}
	upd := &Column{Name: "UPD", IsUpdated: true}
	del := &Column{Name: "DEL", IsDeleted: true}
	cre := &Column{Name: "CRT", IsCreated: true}
	tbl.AddColumn(pk)
	tbl.AddColumn(auto)
	tbl.AddColumn(ver)
	tbl.AddColumn(upd)
	tbl.AddColumn(del)
	tbl.AddColumn(cre)
	if tbl.AutoIncrColumn() == nil || tbl.VersionColumn() == nil || tbl.UpdatedColumn() == nil || tbl.DeletedColumn() == nil {
		t.Fatalf("accessors should return columns")
	}
	if !tbl.Created["CRT"] {
		t.Fatalf("created map not set")
	}
}

func TestTable_ColumnType_And_PKColumns(t *testing.T) {
	tbl := NewTable("t", reflect.TypeOf(kb{}))
	pk := &Column{Name: "ID", FieldIndex: []int{0}, IsPrimaryKey: true, SQLType: SQLType{Name: "INT"}}
	name := &Column{Name: "Code", FieldIndex: []int{1}, SQLType: SQLType{Name: "VARCHAR"}}
	tbl.AddColumn(pk)
	tbl.AddColumn(name)
	if got := tbl.ColumnType("ID"); got.Kind() != reflect.Int64 {
		t.Fatalf("ColumnType ID kind = %v", got.Kind())
	}
	if len(tbl.PKColumns()) != 1 || tbl.PKColumns()[0].Name != "ID" {
		t.Fatalf("PKColumns failed")
	}
}

func TestIndex_Equal(t *testing.T) {
	ia := NewIndex("ix", IndexType)
	ib := NewIndex("ix", IndexType)
	ia.AddColumn("A", "B")
	ib.AddColumn("B", "A")
	if !ia.Equal(ib) {
		t.Fatalf("indexes should be equal disregarding order")
	}
	ic := NewIndex("ix", UniqueType)
	ic.AddColumn("A")
	if ia.Equal(ic) {
		t.Fatalf("should not be equal (different type/cols)")
	}
}

func TestSQLType_Helpers(t *testing.T) {
	if !(&SQLType{Name: "xml"}).IsXML() {
		t.Fatalf("IsXML")
	}
	if !(&SQLType{Name: "int[]"}).IsArray() {
		t.Fatalf("IsArray")
	}
	if SQLTypeName("varChar") != "VARCHAR" {
		t.Fatalf("SQLTypeName")
	}
}

func TestType2SQLType_PanicBranch(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic for unhandled default case kind")
		}
	}()
	// Channel kind will hit default panic branch
	_ = Type2SQLType(reflect.TypeOf(make(chan int)))
}

func TestToFromXormConversions(t *testing.T) {
	col := &Column{Name: "C1", FieldName: "C1", SQLType: SQLType{Name: "VARCHAR"}, Length: 10}
	xcol := ToXormColumn(col)
	backCol := FromXormColumn(xcol)
	if backCol.Name != col.Name || backCol.SQLType.Name != col.SQLType.Name {
		t.Fatalf("column conversion mismatch")
	}
	idx := NewIndex("i1", IndexType)
	idx.AddColumn("C1")
	xidx := ToXormIndex(idx)
	backIdx := FromXormIndex(xidx)
	if !idx.Equal(backIdx) {
		t.Fatalf("index conversion mismatch")
	}

	tbl := NewTable("t2", reflect.TypeOf(kb{}))
	tbl.AddColumn(&Column{Name: "ID"})
	tbl.AddIndex(idx)
	tbl.Created["ID"] = true
	x := ToXormTable(tbl)
	back := FromXormTable(x)
	if len(back.Columns) != 1 || len(back.Indexes) != 1 || !back.Created["ID"] {
		t.Fatalf("table conversion mismatch")
	}
	// Ensure PrimaryKeys preserved in FromXormTable
	x.PrimaryKeys = []string{"ID"}
	back2 := FromXormTable(x)
	if len(back2.PrimaryKeys) != 1 {
		t.Fatalf("primary keys not preserved")
	}

	// Also validate ToXormPK/FromXormPK do copy
	pk := PK{"a", 2}
	xpk := ToXormPK(pk)
	backpk := FromXormPK(xpk)
	if len(backpk) != 2 || backpk[0] != "a" || backpk[1] != 2 {
		t.Fatalf("pk roundtrip")
	}

	_ = xs.SQLType{} // touch import
}
