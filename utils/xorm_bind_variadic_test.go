package utils

import (
	"testing"

	"github.com/dop251/goja"
)

type vdm struct{}

func (vdm) Vsum(prefix string, nums ...int) int {
	s := 0
	for _, n := range nums {
		s += n
	}
	return s
}

func TestBindAllMethods_Variadic_WithArgs(t *testing.T) {
	rt := goja.New()
	obj := BindAllMethods(rt, &vdm{})
	_ = rt.Set("v", obj)
	v, err := rt.RunString("v.Vsum('x', 1,2,3,4)")
	if err != nil {
		t.Fatalf("variadic call failed: %v", err)
	}
	if got := v.Export(); got != int64(10) && got != 10 {
		t.Fatalf("unexpected sum: %v", got)
	}
}

func TestBindAllMethods_Variadic_WithArray(t *testing.T) {
	rt := goja.New()
	obj := BindAllMethods(rt, &vdm{})
	_ = rt.Set("v", obj)
	v, err := rt.RunString("v.Vsum('x', [1,2,3])")
	if err != nil {
		t.Fatalf("variadic array call failed: %v", err)
	}
	if got := v.Export(); got != int64(6) && got != 6 {
		t.Fatalf("unexpected sum: %v", got)
	}
}

func TestBindAllMethods_Variadic_Mismatch(t *testing.T) {
	rt := goja.New()
	obj := BindAllMethods(rt, &vdm{})
	_ = rt.Set("v", obj)
	// Missing required non-variadic argument should error
	if _, err := rt.RunString("v.Vsum()"); err == nil {
		t.Fatalf("expected type error for missing fixed arg")
	}
}
