package utils

import (
	"reflect"
	"testing"

	"github.com/dop251/goja"
)

func TestBindXORMProxy_EngineNotSet(t *testing.T) {
	rt := goja.New()
	obj, _ := BindXORMProxy(rt)
	if err := rt.Set("db", obj); err != nil {
		t.Fatalf("set db: %v", err)
	}
	// Calling any method without setting engine should raise a JS TypeError
	if _, err := rt.RunString("db.DriverName()"); err == nil {
		t.Fatalf("expected error when engine not set")
	}
}

func TestBindXORMProxy_ArgMismatch(t *testing.T) {
	rt := goja.New()
	obj, set := BindXORMProxy(rt)
	if err := rt.Set("db", obj); err != nil {
		t.Fatalf("set db: %v", err)
	}
	eng, err := NewXORM("mysql", "")
	if err != nil {
		t.Fatalf("new engine: %v", err)
	}
	set(eng)
	// SetMaxOpenConns expects 1 argument, calling with none should error
	if _, err := rt.RunString("db.SetMaxOpenConns()"); err == nil {
		t.Fatalf("expected type error for arg mismatch")
	}
}

func TestTryConvert_FallbackUnsupported_NoPanic(t *testing.T) {
	// try to convert a string to a struct type – unsupported – should return (zero,false) without panic
	typ := reflect.TypeOf(struct{ A int }{})
	if v, ok := tryConvert("abc", typ); ok {
		t.Fatalf("expected conversion to fail, got ok=true with %v", v)
	}
}
