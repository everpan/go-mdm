package utils

import (
	"errors"
	"testing"

	"github.com/dop251/goja"
)

type extra struct{}

type small struct{ A int }

func (extra) Void() {}
func (extra) OnlyErr(ok bool) error {
	if !ok {
		return errors.New("bad")
	}
	return nil
}
func (extra) TakesInt64(i int64) int64    { return i + 1 }
func (extra) WantsStruct(s small) int     { return s.A }
func (extra) TwoValsNoErr() (string, int) { return "x", 3 }

func TestBindAllMethods_VoidAndOnlyErr(t *testing.T) {
	rt := goja.New()
	obj := BindAllMethods(rt, &extra{})
	_ = rt.Set("e", obj)

	v, err := rt.RunString("e.Void()")
	if err != nil {
		t.Fatalf("Void run: %v", err)
	}
	if v != goja.Undefined() {
		t.Fatalf("Void should be undefined, got %v", v)
	}

	v, err = rt.RunString("e.OnlyErr(true)")
	if err != nil {
		t.Fatalf("OnlyErr ok: %v", err)
	}
	if v != goja.Undefined() {
		t.Fatalf("OnlyErr ok should return undefined, got %v", v)
	}

	_, err = rt.RunString("e.OnlyErr(false)")
	if err == nil {
		t.Fatalf("OnlyErr expected error")
	}
}

func TestBindAllMethods_ConvertFallbackAndMismatch(t *testing.T) {
	rt := goja.New()
	obj := BindAllMethods(rt, &extra{})
	_ = rt.Set("e", obj)

	// Numbers in JS are float64; ensure conversion to int64 works
	v, err := rt.RunString("e.TakesInt64(40.0)")
	if err != nil {
		t.Fatalf("TakesInt64: %v", err)
	}
	if got := v.Export(); got != int64(41) && got != 41 {
		t.Fatalf("TakesInt64 got %v", got)
	}

	// Force an arg conversion failure by passing wrong type for struct
	_, err = rt.RunString("e.WantsStruct(1)")
	if err == nil {
		t.Fatalf("expected conversion error")
	}
}

func TestBindAllMethods_MultiNoErr_ReturnsArray(t *testing.T) {
	rt := goja.New()
	obj := BindAllMethods(rt, &extra{})
	_ = rt.Set("e", obj)
	v, err := rt.RunString("e.TwoValsNoErr()")
	if err != nil {
		t.Fatalf("TwoValsNoErr: %v", err)
	}
	if s, ok := v.Export().([]interface{}); ok {
		if len(s) != 2 || s[0] != "x" || (s[1] != 3 && s[1] != int64(3)) {
			t.Fatalf("unexpected return: %#v", v.Export())
		}
	} else if s2, ok := v.Export().([]goja.Value); ok {
		if len(s2) != 2 {
			t.Fatalf("unexpected len: %#v", v.Export())
		}
		if s2[0].Export() != "x" {
			t.Fatalf("unexpected first: %v", s2[0].Export())
		}
		if g := s2[1].Export(); g != 3 && g != int64(3) {
			t.Fatalf("unexpected second: %v", g)
		}
	} else {
		t.Fatalf("unexpected export type: %T", v.Export())
	}
}

func TestRegisterXORM_EmptyDriverTypeError(t *testing.T) {
	rt := goja.New()
	_ = RegisterXORM(rt)
	_, err := rt.RunString("xorm('', '')")
	if err == nil {
		t.Fatalf("expected type error for empty driver")
	}
}

func TestBindAllMethods_NilTarget(t *testing.T) {
	rt := goja.New()
	obj := BindAllMethods(rt, nil)
	_ = rt.Set("n", obj)
	v, err := rt.RunString("n.NonExistent")
	if err != nil {
		t.Fatalf("reading property should not error: %v", err)
	}
	if v.Export() != nil {
		t.Fatalf("expected undefined or nil, got %v", v.Export())
	}
}
