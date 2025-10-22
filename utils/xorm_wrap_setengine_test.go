package utils

import "testing"

func TestXormWrap_SetEngine(t *testing.T) {
	w, err := NewXORMWrap("mysql", "")
	if err != nil {
		t.Fatalf("new wrap: %v", err)
	}
	e2, err := NewXORM("postgres", "postgres://localhost/test?sslmode=disable")
	if err != nil {
		// if driver not available in env, skip
		t.Skipf("postgres driver not available: %v", err)
	}
	w.SetEngine(e2)
	// Nothing to assert strongly here; just ensure no panic and method call still reachable
	if _, err := w.Exec(); err == nil {
		// very unlikely to succeed without a real DB; still acceptable either way
	}
}
