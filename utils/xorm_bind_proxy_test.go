package utils

import (
	"testing"

	"github.com/dop251/goja"
)

// This test verifies that we can switch the underlying *xorm.Engine instance
// without re-binding the JavaScript object, and that subsequent method calls
// reflect the new engine.
func TestBindXORMProxy_SwapInstanceAffectsResults(t *testing.T) {
	rt := goja.New()

	obj, setEngine := BindXORMProxy(rt)
	if err := rt.Set("db", obj); err != nil {
		t.Fatalf("set db: %v", err)
	}

	// Start with MySQL engine
	emy, err := NewXORM("mysql", "")
	if err != nil {
		t.Fatalf("mysql new engine: %v", err)
	}
	setEngine(emy)
	v, err := rt.RunString("db.DriverName()")
	if err != nil {
		t.Fatalf("DriverName mysql: %v", err)
	}
	if got := v.Export(); got != "mysql" {
		t.Fatalf("expected mysql, got %v", got)
	}

	// Swap to Postgres engine without re-binding db
	epq, err := NewXORM("postgres", "postgres://localhost/test?sslmode=disable")
	if err != nil {
		t.Fatalf("postgres new engine: %v", err)
	}
	setEngine(epq)
	v, err = rt.RunString("db.DriverName()")
	if err != nil {
		t.Fatalf("DriverName postgres: %v", err)
	}
	if got := v.Export(); got != "postgres" {
		t.Fatalf("expected postgres, got %v", got)
	}
}
