package utils

import (
	"testing"

	"github.com/dop251/goja"
)

func TestRegisterXORM_Postgres_Smoke(t *testing.T) {
	rt := goja.New()
	if err := RegisterXORM(rt); err != nil {
		t.Fatalf("register: %v", err)
	}
	_, err := rt.RunString("var db = xorm('postgres', 'postgres://localhost/test?sslmode=disable'); db.SetMaxOpenConns(1); db.DriverName()")
	if err != nil {
		t.Fatalf("postgres smoke: %v", err)
	}
}
