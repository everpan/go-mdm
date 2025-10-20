package schema_orm

import (
	"reflect"
	"testing"
)

type bean struct {
	ID   int
	Name string
}

func TestBuildDSN_MySQL_And_Postgres(t *testing.T) {
	m := BuildMySQLDSN("", "u", "p", "db")
	if m == "" || m[:2] == "//" {
		t.Fatalf("bad mysql dsn: %s", m)
	}
	m2 := BuildMySQLDSN("127.0.0.1:3307", "u", "p", "db")
	if m2 == m {
		t.Fatalf("expected different dsn with custom port")
	}

	p1 := BuildPostgresDSN("localhost", "u", "", "db")
	if p1 == "" || p1[:11] != "postgres://" {
		t.Fatalf("pg dsn1: %s", p1)
	}
	p2 := BuildPostgresDSN("127.0.0.1:5434", "u", "pass", "db")
	if p2 == "" || p2[:11] != "postgres://" {
		t.Fatalf("pg dsn2: %s", p2)
	}
}

func TestSQLType_Helpers_More(t *testing.T) {
	if !(&SQLType{Name: "timestamp"}).IsTime() {
		t.Fatalf("IsTime")
	}
	if !(&SQLType{Name: "bytea"}).IsBlob() {
		t.Fatalf("IsBlob")
	}
	if !(&SQLType{Name: "bool"}).IsBool() {
		t.Fatalf("IsBool")
	}
}

func TestTable_GetColumnIdx_Positive(t *testing.T) {
	tb := NewTable("t", reflect.TypeOf(bean{}))
	tb.AddColumn(&Column{Name: "ID", FieldIndex: []int{0}, SQLType: SQLType{Name: "INT"}})
	tb.AddColumn(&Column{Name: "Name", FieldIndex: []int{1}, SQLType: SQLType{Name: "VARCHAR"}})
	if c := tb.GetColumnIdx("Name", 0); c == nil || c.Name != "Name" {
		t.Fatalf("GetColumnIdx")
	}
}

func TestColumn_ValueOf_StraightPath(t *testing.T) {
	b := bean{ID: 1, Name: "abc"}
	col := &Column{FieldIndex: []int{1}}
	v, err := col.ValueOf(&b)
	if err != nil {
		t.Fatal(err)
	}
	if v.Kind() != reflect.String || v.String() != "abc" {
		t.Fatalf("unexpected value: %v", v)
	}
}
