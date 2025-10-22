package main

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestHashString_Deterministic(t *testing.T) {
	a := hashString("abc")
	b := hashString("abc")
	c := hashString("abcd")
	if a != b {
		t.Fatalf("expected same hash for identical strings")
	}
	if a == c {
		t.Fatalf("expected different hash for different strings")
	}
}

// NOTE: main() can only be invoked once due to global flag registration.
// The scenario-based execution is covered in TestMain_CommentsExcludeAndSmallWindow.
//func TestMain_RunOnTempProject(t *testing.T) { /* deprecated */ }

func Disabled_TestMain_RunOnTempProject(t *testing.T) {
	dir := t.TempDir()
	// create two small go files with duplicated content windows
	code1 := `package p

// comment
func A() {
	// dup-start
	println(1)
	println(2)
	// dup-end
}
`
	code2 := `package p

func B() {
	// dup-start
	println(1)
	println(2)
	// dup-end
}
`
	if err := os.WriteFile(filepath.Join(dir, "a.go"), []byte(code1), 0o644); err != nil {
		t.Fatalf("write a.go: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "b.go"), []byte(code2), 0o644); err != nil {
		t.Fatalf("write b.go: %v", err)
	}

	// Prepare flags for main
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"dupcalc", "-root", dir, "-w", "2", "-exclude", ""}

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	defer func() { os.Stdout = oldStdout }()

	// Run main
	main()

	// Read output
	_ = w.Close()
	outBytes, _ := io.ReadAll(r)
	out := string(outBytes)
	if !strings.Contains(out, "Files:") {
		t.Fatalf("unexpected output: %s", out)
	}
	if !strings.Contains(out, "Duplicate windows") {
		t.Fatalf("expected duplicate windows line in output: %s", out)
	}
}
