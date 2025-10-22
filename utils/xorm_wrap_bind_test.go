package utils

import (
	"strings"
	"testing"

	"github.com/dop251/goja"
)

func TestRegisterXORMWrap_CreateAndExecErrorsPropagate(t *testing.T) {
	rt := goja.New()
	if err := RegisterXORMWrap(rt); err != nil {
		t.Fatalf("register: %v", err)
	}
	// Create a wrap with a known driver but likely invalid/empty DSN.
	// Calling Exec should fail at runtime in Go, and the error should
	// be propagated to JS as a thrown exception.
	_, err := rt.RunString(`
		var w = xormWrap('mysql', '');
		w.Exec();
	`)
	if err == nil {
		t.Fatalf("expected Exec to throw due to invalid connection")
	}
}

func TestBindXORMWrap_SetEngineArgTypeError(t *testing.T) {
	rt := goja.New()
	w, err := NewXORMWrap("mysql", "")
	if err != nil {
		t.Fatalf("new wrap: %v", err)
	}
	obj := BindXORMWrap(rt, w)
	_ = rt.Set("w", obj)
	// Pass the wrong type to SetEngine(*xorm.Engine)
	_, err = rt.RunString("w.SetEngine(1)")
	if err == nil {
		t.Fatalf("expected type error when passing wrong argument to SetEngine")
	}
}

func TestRegisterXORMWrap_EmptyDriverTypeError(t *testing.T) {
	rt := goja.New()
	_ = RegisterXORMWrap(rt)
	_, err := rt.RunString("xormWrap('', '')")
	if err == nil {
		t.Fatalf("expected type error for empty driver")
	}
	if got := err.Error(); !strings.Contains(got, "driver") {
		t.Fatalf("unexpected error: %v", got)
	}
}

func TestBindXORMWrap_NilTargetSafe(t *testing.T) {
	rt := goja.New()
	obj := BindXORMWrap(rt, nil)
	_ = rt.Set("nw", obj)
	v, err := rt.RunString("nw.NonExisting")
	if err != nil {
		t.Fatalf("reading property should not error: %v", err)
	}
	if v.Export() != nil {
		t.Fatalf("expected undefined or nil, got %v", v.Export())
	}
}
