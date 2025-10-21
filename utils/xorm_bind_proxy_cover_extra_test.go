package utils

import (
	"testing"

	"github.com/dop251/goja"
)

func TestBindXORMProxy_ErrorPropagationAndWraps(t *testing.T) {
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

	// Error propagation: Ping should fail without a live DB
	if _, err := rt.RunString("db.Ping()"); err == nil {
		t.Fatalf("expected ping error to propagate")
	}

	// Wrapped return types: NewSession returns *xorm.Session which should be wrapped with methods.
	v, err := rt.RunString("var s = db.NewSession(); s !== undefined")
	if err != nil {
		t.Fatalf("NewSession eval: %v", err)
	}
	if v.Export() == false {
		t.Fatalf("expected session object to be returned")
	}
	// Call a simple method on session to ensure binding works; Close() has no return
	if _, err := rt.RunString("s.Close(); undefined"); err != nil {
		t.Fatalf("session.Close failed: %v", err)
	}
}
