package utils

import (
	"errors"
	"strings"
	"testing"

	"github.com/dop251/goja"
)

type dummy struct{}

func (dummy) NoArg() int           { return 42 }
func (dummy) Echo(s string) string { return "hi:" + s }
func (dummy) Sum(a int, b int) (int, error) {
	if a < 0 || b < 0 {
		return 0, errors.New("neg")
	}
	return a + b, nil
}
func (dummy) Multi() (int, string) { return 7, "ok" }

func TestBindAllMethods_NoArgAndEcho(t *testing.T) {
	rt := goja.New()
	obj := BindAllMethods(rt, &dummy{})
	if err := rt.Set("d", obj); err != nil {
		t.Fatalf("set d: %v", err)
	}

	v, err := rt.RunString("d.NoArg()")
	if err != nil {
		t.Fatalf("NoArg run: %v", err)
	}
	if got := v.Export(); got != int64(42) && got != 42 {
		t.Fatalf("NoArg got = %v", got)
	}

	v, err = rt.RunString("d.Echo('js')")
	if err != nil {
		t.Fatalf("Echo run: %v", err)
	}
	if got := v.Export(); got != "hi:js" {
		t.Fatalf("Echo got = %v", got)
	}
}

func TestBindAllMethods_MultiAndError(t *testing.T) {
	rt := goja.New()
	obj := BindAllMethods(rt, &dummy{})
	_ = rt.Set("d", obj)

	v, err := rt.RunString("d.Multi()")
	if err != nil {
		t.Fatalf("Multi run: %v", err)
	}
	// Accept both []interface{} and []goja.Value
	if arr, ok := v.Export().([]interface{}); ok {
		if len(arr) != 2 || (arr[0] != int64(7) && arr[0] != 7) || arr[1] != "ok" {
			t.Fatalf("Multi export invalid: %#v", v.Export())
		}
	} else if arr2, ok := v.Export().([]goja.Value); ok {
		if len(arr2) != 2 {
			t.Fatalf("Multi export invalid: %#v", v.Export())
		}
		if g := arr2[0].Export(); g != int64(7) && g != 7 {
			t.Fatalf("bad first: %v", g)
		}
		if arr2[1].Export() != "ok" {
			t.Fatalf("bad second: %v", arr2[1].Export())
		}
	} else {
		t.Fatalf("unexpected export type: %T", v.Export())
	}

	_, err = rt.RunString("d.Sum(-1, 5)")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "neg") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBindAllMethods_ArgMismatch(t *testing.T) {
	rt := goja.New()
	obj := BindAllMethods(rt, &dummy{})
	_ = rt.Set("d", obj)
	_, err := rt.RunString("d.Echo()")
	if err == nil {
		t.Fatalf("expected type error for missing arg")
	}
}

func TestRegisterXORM_ErrorOnUnknownDriver(t *testing.T) {
	rt := goja.New()
	if err := RegisterXORM(rt); err != nil {
		t.Fatalf("register: %v", err)
	}
	_, err := rt.RunString("xorm('not-a-driver', '')")
	if err == nil {
		t.Fatalf("expected error for unknown driver")
	}
}
