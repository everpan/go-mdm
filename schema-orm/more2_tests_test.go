package schema_orm

import (
	"encoding/json"
	"reflect"
	"testing"

	yaml "gopkg.in/yaml.v3"
)

func TestIndexXNamePrefixed(t *testing.T) {
	i1 := &Index{Name: "IDX_foo", Type: IndexType}
	if i1.XName("bar") != "IDX_foo" {
		t.Fatalf("prefixed idx")
	}
	i2 := &Index{Name: "UQE_foo", Type: UniqueType}
	if i2.XName("bar") != "UQE_foo" {
		t.Fatalf("prefixed uqe")
	}
}

func TestPKIsZeroUint(t *testing.T) {
	p := NewPK(uint(0))
	if !p.IsZero() {
		t.Fatalf("uint zero")
	}
}

func TestSQLTypeYAMLJSON(t *testing.T) {
	s := SQLType{Name: "VARCHAR", DefaultLength: 255, DefaultLength2: 0}
	b, err := json.Marshal(s)
	if err != nil {
		t.Fatal(err)
	}
	var s2 SQLType
	if err := json.Unmarshal(b, &s2); err != nil {
		t.Fatal(err)
	}
	if s2.Name != s.Name || s2.DefaultLength != 255 {
		t.Fatalf("json sqltype")
	}
	y, err := yaml.Marshal(s)
	if err != nil {
		t.Fatal(err)
	}
	var s3 SQLType
	if err := yaml.Unmarshal(y, &s3); err != nil {
		t.Fatal(err)
	}
	if s3.Name != s.Name {
		t.Fatalf("yaml sqltype")
	}
}

func TestColumnYAMLJSON(t *testing.T) {
	c := NewColumn("name", "Name", SQLType{Name: "VARCHAR"}, 50, 0, true)
	b, err := json.Marshal(c)
	if err != nil {
		t.Fatal(err)
	}
	var c2 Column
	if err := json.Unmarshal(b, &c2); err != nil {
		t.Fatal(err)
	}
	if c2.Name != c.Name || !c2.Nullable {
		t.Fatalf("json column")
	}
	y, err := yaml.Marshal(c)
	if err != nil {
		t.Fatal(err)
	}
	var c3 Column
	if err := yaml.Unmarshal(y, &c3); err != nil {
		t.Fatal(err)
	}
	if c3.SQLType.Name != "VARCHAR" {
		t.Fatalf("yaml column")
	}
}

func TestNewColumnDefaults(t *testing.T) {
	c := NewColumn("id", "ID", SQLType{Name: "INT"}, 10, 0, false)
	if c.MapType != TWOSIDES || c.Default != "" || c.Indexes == nil || !c.DefaultIsEmpty {
		t.Fatalf("defaults")
	}
}

func TestNewEmptyTable(t *testing.T) {
	tbl := NewEmptyTable()
	if tbl.Name != "" || tbl.Type != nil || len(tbl.Columns) != 0 || len(tbl.ColumnsSeq) != 0 {
		t.Fatalf("empty table")
	}
}

func TestType2SQLTypeByteSlice(t *testing.T) {
	if Type2SQLType(reflect.TypeOf([]byte{})).Name != "BLOB" {
		t.Fatalf("[]byte mapping")
	}
}
