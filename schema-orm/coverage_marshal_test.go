package schema_orm

import (
	"encoding/json"
	"reflect"
	"testing"

	"gopkg.in/yaml.v3"
)

type rt struct {
	A int
}

type multi struct {
	A int64
	B string
	S rt
}

func TestMarshal_JSON_YAML_Roundtrip_AllTypes(t *testing.T) {
	// SQLType
	st := &SQLType{Name: "DECIMAL", DefaultLength: 10, DefaultLength2: 2}
	b, err := json.Marshal(st)
	if err != nil {
		t.Fatal(err)
	}
	var st2 SQLType
	if err := json.Unmarshal(b, &st2); err != nil {
		t.Fatal(err)
	}
	if st2.Name != st.Name || st2.DefaultLength != 10 || st2.DefaultLength2 != 2 {
		t.Fatalf("sqltype mismatch")
	}
	// YAML
	by, err := yaml.Marshal(st)
	if err != nil {
		t.Fatal(err)
	}
	var st3 SQLType
	if err := yaml.Unmarshal(by, &st3); err != nil {
		t.Fatal(err)
	}

	// Column
	col := &Column{Name: "C", FieldName: "C", SQLType: SQLType{Name: "VARCHAR"}, Length: 5, Nullable: true}
	cb, _ := json.Marshal(col)
	var col2 Column
	_ = json.Unmarshal(cb, &col2)
	if col2.Name != col.Name || col2.SQLType.Name != "VARCHAR" {
		t.Fatalf("column mismatch")
	}
	// YAML
	cby, _ := yaml.Marshal(col)
	var col3 Column
	_ = yaml.Unmarshal(cby, &col3)

	// Index
	idx := NewIndex("nm", UniqueType)
	idx.AddColumn("C")
	ib, _ := json.Marshal(idx)
	var idx2 Index
	_ = json.Unmarshal(ib, &idx2)
	if !idx.Equal(&idx2) {
		t.Fatalf("index roundtrip")
	}
	// YAML
	iby, _ := yaml.Marshal(idx)
	var idx3 Index
	_ = yaml.Unmarshal(iby, &idx3)

	// Table with everything
	tbl := NewTable("t", reflect.TypeOf(multi{}))
	tbl.AddColumn(&Column{Name: "A", FieldIndex: []int{0}, SQLType: SQLType{Name: "INT"}, IsPrimaryKey: true})
	tbl.AddColumn(&Column{Name: "B", FieldIndex: []int{1}, SQLType: SQLType{Name: "VARCHAR"}, IsPrimaryKey: true})
	tbl.AddIndex(idx)
	tbl.Created["A"] = true
	tbl.Updated = "B"
	tbl.Deleted = "B"
	tbl.Version = "A"
	tbj, err := json.Marshal(tbl)
	if err != nil {
		t.Fatal(err)
	}
	var tbl2 Table
	if err := json.Unmarshal(tbj, &tbl2); err != nil {
		t.Fatal(err)
	}
	if tbl2.Name != "t" || len(tbl2.Columns) != 2 || len(tbl2.Indexes) != 1 {
		t.Fatalf("table roundtrip json")
	}
	// YAML
	tby, err := yaml.Marshal(tbl)
	if err != nil {
		t.Fatal(err)
	}
	var tbl3 Table
	if err := yaml.Unmarshal(tby, &tbl3); err != nil {
		t.Fatal(err)
	}
	if tbl3.Name != "t" || len(tbl3.Columns) != 2 || len(tbl3.Indexes) != 1 {
		t.Fatalf("table roundtrip yaml")
	}

	// PK
	pk := PK{"x", 3}
	pb, _ := json.Marshal(&pk)
	var pk2 PK
	_ = json.Unmarshal(pb, &pk2)
	if len(pk2) != 2 {
		t.Fatalf("pk len")
	}
	pby, _ := yaml.Marshal(&pk)
	var pk3 PK
	_ = yaml.Unmarshal(pby, &pk3)
}

func TestImportTablesFromJSON_ErrorAndSuccess(t *testing.T) {
	if _, err := ImportTablesFromJSON("{"); err == nil {
		t.Fatalf("expected json error")
	}
	// success
	in := `[{"name":"t","columns":[{"name":"id","sqlType":{"name":"INT"}}]}]`
	ts, err := ImportTablesFromJSON(in)
	if err != nil || len(ts) != 1 {
		t.Fatalf("import tables")
	}
}

func TestTable_IDOfV_MultiPK_and_PanicPath(t *testing.T) {
	// multi-pk happy path
	tbl := NewTable("t", reflect.TypeOf(multi{}))
	tbl.AddColumn(&Column{Name: "A", FieldIndex: []int{0}, SQLType: SQLType{Name: "INT"}, IsPrimaryKey: true})
	tbl.AddColumn(&Column{Name: "B", FieldIndex: []int{1}, SQLType: SQLType{Name: "VARCHAR"}, IsPrimaryKey: true})
	obj := multi{A: 9, B: "k"}
	pk, err := tbl.IDOfV(reflect.ValueOf(obj))
	if err != nil || len(pk) != 2 {
		t.Fatalf("multi pk idofv: %v %v", pk, err)
	}
	// panic path: unsupported kind (struct)
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic for unsupported kind")
		}
	}()
	// Add a primary key with struct kind to hit default panic
	tbl2 := NewTable("t2", reflect.TypeOf(multi{}))
	tbl2.AddColumn(&Column{Name: "S", FieldIndex: []int{2}, SQLType: SQLType{Name: "JSON"}, IsPrimaryKey: true})
	_, _ = tbl2.IDOfV(reflect.ValueOf(obj))
}
