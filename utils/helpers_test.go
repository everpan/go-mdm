package utils

import (
	"errors"
	"reflect"
	"testing"
)

type h struct{}

type vv struct{}

func (h) A()          {}
func (h) B(x int) int { return x + 1 }
func (h) C(x int, y string) (int, error) {
	if y == "bad" {
		return 0, errors.New("bad")
	}
	return x, nil
}
func (h) D() (int, string) { return 5, "ok" }
func (vv) V(a ...int)      {}

func TestFuncNameList(t *testing.T) {
	names := FuncNameList(&h{})
	if len(names) == 0 {
		t.Fatalf("expected some method names")
	}
	// Should contain A, B, C, D
	m := map[string]bool{}
	for _, n := range names {
		m[n] = true
	}
	for _, want := range []string{"A", "B", "C", "D"} {
		if !m[want] {
			t.Fatalf("missing method %s in %v", want, names)
		}
	}
}

func TestCallMethod_SuccessAndErrors(t *testing.T) {
	// Success no return
	if res, err := CallMethod(&h{}, "A"); err != nil || len(res) != 0 {
		t.Fatalf("A err=%v res=%v", err, res)
	}
	// Success single return
	if res, err := CallMethod(&h{}, "B", 1); err != nil || !reflect.DeepEqual(res, []any{2}) {
		t.Fatalf("B err=%v res=%v", err, res)
	}
	// Success multi-return without error
	if res, err := CallMethod(&h{}, "D"); err != nil || !reflect.DeepEqual(res, []any{5, "ok"}) {
		t.Fatalf("D err=%v res=%v", err, res)
	}
	// Error return from method
	if _, err := CallMethod(&h{}, "C", 1, "bad"); err == nil {
		t.Fatalf("expected error from C")
	}
	// Method not found
	if _, err := CallMethod(&h{}, "Z"); err == nil {
		t.Fatalf("expected not found")
	}
	// Bad arg count
	if _, err := CallMethod(&h{}, "B"); err == nil {
		t.Fatalf("expected bad arg count")
	}
	// Variadic unsupported check using a small type
	_, err := CallMethod(&vv{}, "V", []int{})
	if err == nil {
		t.Fatalf("expected variadic unsupported")
	}
	// Nil target
	_, err = CallMethod(nil, "A")
	if err == nil {
		t.Fatalf("expected nil target error")
	}
}
