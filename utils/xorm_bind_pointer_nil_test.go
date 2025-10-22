package utils

import (
	"testing"

	"github.com/dop251/goja"
)

type ptrdm struct{}

func (ptrdm) WantsPtr(p *int) bool { return p == nil }

func TestBindAllMethods_PointerParam_NullConvertsToNil(t *testing.T) {
	rt := goja.New()
	obj := BindAllMethods(rt, &ptrdm{})
	_ = rt.Set("p", obj)
	v, err := rt.RunString("p.WantsPtr(null)")
	if err != nil {
		t.Fatalf("call failed: %v", err)
	}
	if got, ok := v.Export().(bool); !ok || !got {
		t.Fatalf("expected true, got %v (%T)", v.Export(), v.Export())
	}
}
