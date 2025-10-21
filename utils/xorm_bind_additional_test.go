package utils

import (
	"testing"

	"github.com/dop251/goja"
)

// Covers: RegisterXORM arity error, BindXORM wrapping of *xorm.Session return, and simple bool/map/slice handling
func TestXorm_Bind_SessionWrapped_And_ArityError(t *testing.T) {
	rt := goja.New()
	if err := RegisterXORM(rt); err != nil {
		t.Fatalf("register: %v", err)
	}
	// Missing second argument should be a type error
	if _, err := rt.RunString("xorm('mysql')"); err == nil {
		t.Fatalf("expected type error for missing dsn")
	}

	// Create an engine object and ensure session object is wrapped with methods
	if _, err := rt.RunString("var db = xorm('mysql', '');"); err != nil {
		t.Fatalf("create engine: %v", err)
	}
	// NewSession returns a *xorm.Session which should be wrapped into a JS object with methods
	v, err := rt.RunString("var s = db.NewSession(); typeof s === 'object' && typeof s.Close === 'function'")
	if err != nil {
		t.Fatalf("NewSession wrapping check failed: %v", err)
	}
	if got := v.Export(); got != true {
		t.Fatalf("expected session object with Close method, got %v (%T)", got, got)
	}
}

// Additional coverage exercising map/slice/bool pathways.
type cs struct{}

func (cs) TakesMap(m map[string]string) int {
	if m == nil {
		return -1
	}
	return len(m)
}
func (cs) TakesSlice(a []int) int {
	if a == nil {
		return -1
	}
	s := 0
	for _, v := range a {
		s += v
	}
	return s
}
func (cs) BoolNeg(b bool) bool { return !b }

func TestBindAllMethods_CollectionsAndBool(t *testing.T) {
	rt := goja.New()
	obj := BindAllMethods(rt, &cs{})
	_ = rt.Set("c", obj)

	v, err := rt.RunString("c.TakesMap(null)")
	if err != nil {
		t.Fatalf("TakesMap null: %v", err)
	}
	if got := v.Export(); got != -1 && got != int64(-1) {
		t.Fatalf("TakesMap(null) got %v", got)
	}

	v, err = rt.RunString("c.TakesSlice(null)")
	if err != nil {
		t.Fatalf("TakesSlice null: %v", err)
	}
	if got := v.Export(); got != -1 && got != int64(-1) {
		t.Fatalf("TakesSlice(null) got %v", got)
	}

	v, err = rt.RunString("c.TakesSlice([1,2,3])")
	if err != nil {
		t.Fatalf("TakesSlice array: %v", err)
	}
	if got := v.Export(); got != 6 && got != int64(6) {
		t.Fatalf("TakesSlice([1,2,3]) got %v", got)
	}

	v, err = rt.RunString("c.BoolNeg(false)")
	if err != nil {
		t.Fatalf("BoolNeg: %v", err)
	}
	if got := v.Export(); got != true {
		t.Fatalf("BoolNeg(false) got %v", got)
	}
}
