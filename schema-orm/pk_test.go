package schema_orm

import (
	"encoding/gob"
	"testing"
)

func TestPK_New_IsZero_ToFromString(t *testing.T) {
	p := NewPK(0, "", false)
	if !p.IsZero() {
		t.Fatalf("expected zero pk")
	}
	// Non-zero
	p2 := NewPK(1, "x", true)
	if p2.IsZero() {
		t.Fatalf("not zero pk")
	}
	// ToString/FromString roundtrip
	gob.Register(map[string]any{})
	s, err := p2.ToString()
	if err != nil {
		t.Fatal(err)
	}
	var back PK
	if err := back.FromString(s); err != nil {
		t.Fatal(err)
	}
	if len(back) != 3 {
		t.Fatalf("back len")
	}
}

func TestIsZeroStrict(t *testing.T) {
	if !isZeroStrict(0.0) {
		t.Fatalf("0 should be zero strict")
	}
	if isZeroStrict(1e-6) {
		t.Fatalf("1e-6 should not be zero strict")
	}
}
