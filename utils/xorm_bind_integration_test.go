package utils

import (
	"testing"

	"github.com/dop251/goja"
)

func TestRegisterXORM_MySQL_Smoke(t *testing.T) {
	rt := goja.New()
	if err := RegisterXORM(rt); err != nil {
		t.Fatalf("register: %v", err)
	}
	// Create engine; mysql driver is imported in go.mod, no actual connection attempted until Ping
	_, err := rt.RunString("var db = xorm('mysql', ''); db.SetMaxOpenConns(5); db.SetMaxIdleConns(2);")
	if err != nil {
		t.Fatalf("set conns failed: %v", err)
	}
	v, err := rt.RunString("db.DriverName()")
	if err != nil {
		t.Fatalf("DriverName: %v", err)
	}
	if got := v.Export(); got != "mysql" {
		t.Fatalf("DriverName got %v", got)
	}
	// Ping should fail without a running DB, exercising error propagation
	_, err = rt.RunString("db.Ping()")
	if err == nil {
		t.Fatalf("expected ping error")
	}
}
