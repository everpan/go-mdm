package utils

import "testing"

func TestFuncNameList_NilTarget(t *testing.T) {
	if got := FuncNameList(nil); got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestCallMethod_BadArgType(t *testing.T) {
	// Method B expects int; pass string to trigger not assignable/convertible branch
	if _, err := CallMethod(&h{}, "B", "not-int"); err == nil {
		t.Fatalf("expected bad argument type error")
	}
}
