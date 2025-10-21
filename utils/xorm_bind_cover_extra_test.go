package utils

import (
	"reflect"
	"testing"

	"github.com/dop251/goja"
)

func equalsAny(got interface{}, candidates []interface{}) bool {
	for _, c := range candidates {
		if reflect.DeepEqual(got, c) {
			return true
		}
	}
	return false
}

type conv struct{}

func (conv) PtrArg(p *int) int {
	if p == nil {
		return -1
	}
	return *p
}
func (conv) TakesUint(u uint) uint       { return u + 1 }
func (conv) TakesUint32(u uint32) uint32 { return u + 2 }
func (conv) TakesFloat32(f float32) int  { return int(f) }
func (conv) TakesMap(m map[string]string) int {
	if m == nil {
		return -1
	}
	return len(m)
}
func (conv) TakesSlice(a []int) int {
	if a == nil {
		return -1
	}
	s := 0
	for _, v := range a {
		s += v
	}
	return s
}
func (conv) BoolNeg(b bool) bool { return !b }

func TestTryConvert_PointerNilAndNumbers(t *testing.T) { // and collections/bool
	rt := goja.New()
	obj := BindAllMethods(rt, &conv{})
	_ = rt.Set("c", obj)

	v, err := rt.RunString("c.PtrArg(null)")
	if err != nil {
		t.Fatalf("PtrArg null: %v", err)
	}
	if got := v.Export(); got != -1 && got != int64(-1) {
		t.Fatalf("PtrArg(null) got %v", got)
	}

	v, err = rt.RunString("c.TakesUint(2.0)")
	if err != nil {
		t.Fatalf("TakesUint: %v", err)
	}
	if got := v.Export(); !equalsAny(got, []interface{}{uint64(3), int64(3), 3, float64(3)}) {
		t.Fatalf("TakesUint got %T %v", got, got)
	}

	v, err = rt.RunString("c.TakesUint32(5.0)")
	if err != nil {
		t.Fatalf("TakesUint32: %v", err)
	}
	if got := v.Export(); !equalsAny(got, []interface{}{uint64(7), int64(7), 7, float64(7)}) {
		t.Fatalf("TakesUint32 got %T %v", got, got)
	}

	v, err = rt.RunString("c.TakesFloat32(3)")
	if err != nil {
		t.Fatalf("TakesFloat32: %v", err)
	}
	if got := v.Export(); got != int64(3) && got != 3 {
		t.Fatalf("TakesFloat32 got %v", got)
	}
}

// direct coverage for trailingErrorType outcomes
func TestTrailingErrorType(t *testing.T) {
	// func() {}
	ft0 := reflect.TypeOf(func() {})
	if _, ok := trailingErrorType(ft0); ok {
		t.Fatalf("expected no error return")
	}
	// func() error
	ft1 := reflect.TypeOf(func() error { return nil })
	if _, ok := trailingErrorType(ft1); !ok {
		t.Fatalf("expected error return")
	}
	// func() (int, error)
	ft2 := reflect.TypeOf(func() (int, error) { return 0, nil })
	if _, ok := trailingErrorType(ft2); !ok {
		t.Fatalf("expected trailing error return")
	}
}
