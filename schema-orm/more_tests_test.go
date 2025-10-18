package schema_orm

import (
	"encoding/json"
	"reflect"
	"testing"

	"gopkg.in/yaml.v3"
)

type nested struct{ V int }

type holder struct{ P *nested }

type ifaceHolder struct{ I interface{} }

func TestSQLTypeMore(t *testing.T) {
	cases := []SQLType{
		{Name: "BLOB"}, {Name: "DATETIME"}, {Name: "BOOLEAN"}, {Name: "INT[]"}, {Name: "XML"},
	}
	if !cases[0].IsBlob() {
		t.Fatalf("blob")
	}
	if !cases[1].IsTime() {
		t.Fatalf("time")
	}
	if !cases[2].IsBool() {
		t.Fatalf("bool")
	}
	if !cases[3].IsArray() {
		t.Fatalf("array")
	}
	if !cases[4].IsXML() {
		t.Fatalf("xml")
	}
	// Type mappings
	if Type2SQLType(reflect.TypeOf("")).Name != "VARCHAR" {
		t.Fatalf("type2sql string")
	}
	if Type2SQLType(reflect.TypeOf(1)).Name != "INT" {
		t.Fatalf("type2sql int")
	}
	if SQLType2Type(SQLType{Name: "BOOLEAN"}).Kind() != reflect.Bool {
		t.Fatalf("sql2type bool")
	}
	if SQLTypeName("varchar") != "VARCHAR" {
		t.Fatalf("sqlTypeName")
	}
}

func TestColumnValueOfWithPointerAndInterface(t *testing.T) {
	c := &Column{FieldIndex: []int{0}}
	h := &holder{}
	val, err := c.ValueOf(h)
	if err != nil {
		t.Fatal(err)
	}
	if val.Kind() != reflect.Struct && val.Kind() != reflect.Ptr {
		t.Fatalf("unexpected kind: %v", val.Kind())
	}
	// interface case
	ci := &Column{FieldIndex: []int{0}}
	ih := &ifaceHolder{I: &nested{V: 3}}
	v := reflect.ValueOf(ih)
	// emulate interface indirection
	ci.FieldIndex = []int{0}
	vi, err := ci.ValueOfV(&v)
	if err != nil {
		t.Fatal(err)
	}
	_ = vi // ensure a path executed
}

func TestConvertIDError(t *testing.T) {
	c := &Column{SQLType: SQLType{Name: "BLOB"}}
	if _, err := c.ConvertID("1"); err == nil {
		t.Fatalf("expected error")
	}
}

func TestTableUtilityMethods(t *testing.T) {
	typ := reflect.TypeOf(sample{})
	tbl := NewTable("tb", typ)
	c1 := &Column{Name: "ID", FieldIndex: []int{0}, IsPrimaryKey: true, IsAutoIncrement: true}
	tbl.AddColumn(c1)
	c2 := &Column{Name: "Name", FieldIndex: []int{1}, IsUpdated: true}
	tbl.AddColumn(c2)
	c3 := &Column{Name: "Deleted", IsDeleted: true}
	tbl.AddColumn(c3)
	c4 := &Column{Name: "Version", IsVersion: true}
	tbl.AddColumn(c4)
	if tbl.AutoIncrColumn() != c1 {
		t.Fatalf("autoincr")
	}
	if tbl.UpdatedColumn() != c2 {
		t.Fatalf("updated")
	}
	if tbl.DeletedColumn() != c3 {
		t.Fatalf("deleted")
	}
	if tbl.VersionColumn() != c4 {
		t.Fatalf("version")
	}
	if tbl.ColumnType("ID") != typ.Field(0).Type {
		t.Fatalf("col type")
	}
}

func TestTableJSONYAMLFull(t *testing.T) {
	typ := reflect.TypeOf(sample{})
	tbl := NewTable("full", typ)
	c := &Column{Name: "ID"}
	tbl.AddColumn(c)
	idx := NewIndex("ix", IndexType)
	idx.AddColumn("ID")
	tbl.AddIndex(idx)
	tbl.Created["ID"] = true
	tbl.AutoIncrement = "ID"
	tbl.Updated = "ID"
	tbl.Deleted = "ID"
	tbl.Version = "ID"
	tbl.StoreEngine = "InnoDB"
	tbl.Charset = "utf8mb4"
	tbl.Comment = "comment"
	tbl.Collation = "utf8mb4_0900_ai_ci"
	// JSON
	b, err := json.Marshal(tbl)
	if err != nil {
		t.Fatal(err)
	}
	var out Table
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatal(err)
	}
	if out.Name != tbl.Name || out.AutoIncrement != "ID" || len(out.Columns) != 1 {
		t.Fatalf("json full")
	}
	// YAML
	y, err := yaml.Marshal(tbl)
	if err != nil {
		t.Fatal(err)
	}
	var outY Table
	if err := yaml.Unmarshal(y, &outY); err != nil {
		t.Fatal(err)
	}
	if outY.Comment != tbl.Comment || len(outY.Indexes) != 1 {
		t.Fatalf("yaml full")
	}
}

func TestIndexNotEqualBranches(t *testing.T) {
	a := NewIndex("a", IndexType)
	a.AddColumn("x")
	b := NewIndex("a", UniqueType)
	b.AddColumn("x")
	if a.Equal(b) {
		t.Fatalf("type should differ")
	}
	c := NewIndex("a", IndexType)
	c.AddColumn("x", "y")
	if a.Equal(c) {
		t.Fatalf("len should differ")
	}
	d := NewIndex("a", IndexType)
	d.AddColumn("y")
	if a.Equal(d) {
		t.Fatalf("cols should differ")
	}
}

func TestPKJSONYAML(t *testing.T) {
	p := PK{1, "a"}
	b, err := json.Marshal(p)
	if err != nil {
		t.Fatal(err)
	}
	var p2 PK
	if err := json.Unmarshal(b, &p2); err != nil {
		t.Fatal(err)
	}
	if len(p2) != 2 {
		t.Fatalf("json pk")
	}
	y, err := yaml.Marshal(p)
	if err != nil {
		t.Fatal(err)
	}
	var p3 PK
	if err := yaml.Unmarshal(y, &p3); err != nil {
		t.Fatal(err)
	}
	if len(p3) != 2 {
		t.Fatalf("yaml pk")
	}
}
