package utils

import (
	"testing"

	"github.com/dop251/goja"
)

type numdm struct{}

func (numdm) Uint32Add1(u uint32) uint32     { return u + 1 }
func (numdm) FloatAddHalf(f float64) float64 { return f + 0.5 }

func TestBindAllMethods_NumericConversions(t *testing.T) {
	rt := goja.New()
	obj := BindAllMethods(rt, &numdm{})
	_ = rt.Set("n", obj)
	v, err := rt.RunString("n.Uint32Add1(5.0)")
	if err != nil {
		t.Fatalf("Uint32Add1 failed: %v", err)
	}
	if got := v.Export(); got != int64(6) && got != uint32(6) && got != 6 {
		t.Fatalf("unexpected: %v (%T)", v.Export(), v.Export())
	}
	v, err = rt.RunString("n.FloatAddHalf(3)")
	if err != nil {
		t.Fatalf("FloatAddHalf failed: %v", err)
	}
	if f, ok := v.Export().(float64); !ok || f != 3.5 {
		t.Fatalf("unexpected float: %v (%T)", v.Export(), v.Export())
	}
}

func TestRegisterXORMWrap_SuccessReturnsObject(t *testing.T) {
	rt := goja.New()
	_ = RegisterXORMWrap(rt)
	v, err := rt.RunString("(function(){ var w = xormWrap('mysql',''); return typeof w.Exec === 'function'; })()")
	if err != nil {
		t.Fatalf("register xormWrap failed: %v", err)
	}
	if got, ok := v.Export().(bool); !ok || !got {
		t.Fatalf("expected true, got %v", v.Export())
	}
}

func TestRegisterXORMWrap_SetEngineFromGoIntoJS(t *testing.T) {
	rt := goja.New()
	_ = RegisterXORMWrap(rt)
	eng, err := NewXORM("mysql", "")
	if err != nil {
		t.Fatalf("new engine: %v", err)
	}
	_ = rt.Set("eng", eng)
	_, err = rt.RunString("var w = xormWrap('mysql',''); w.SetEngine(eng);")
	if err != nil {
		t.Fatalf("SetEngine via JS failed: %v", err)
	}
}
