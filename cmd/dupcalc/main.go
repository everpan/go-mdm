package main

import (
	"bufio"
	"flag"
	"fmt"
	"hash/fnv"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// occurrence represents a window occurrence within a file (by normalized line index)
type occurrence struct {
	file string
	idx  int // start index in normalized lines
}

type fileData struct {
	lines    []string // normalized (no blanks/comments)
	dupMarks []bool   // marks for duplicated lines
}

func main() {
	var (
		root    string
		win     int
		exclude string
	)
	flag.StringVar(&root, "root", ".", "project root to scan")
	flag.IntVar(&win, "w", 5, "window size (number of lines) for duplication detection")
	flag.StringVar(&exclude, "exclude", "vendor,node_modules,.git", "comma-separated directory names to exclude")
	flag.Parse()

	excluded := map[string]struct{}{}
	for _, p := range strings.Split(exclude, ",") {
		p = strings.TrimSpace(p)
		if p != "" {
			excluded[p] = struct{}{}
		}
	}

	files := map[string]*fileData{}

	// Collect .go files under root (excluding common directories and module cache)
	walkErr := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		name := d.Name()
		if d.IsDir() {
			if _, ok := excluded[name]; ok {
				return filepath.SkipDir
			}
			// skip Go module cache dirs if present in repo layout
			if name == "pkg" && strings.Contains(path, string(os.PathSeparator)+"go"+string(os.PathSeparator)+"pkg"+string(os.PathSeparator)+"mod") {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(name, ".go") {
			return nil
		}
		// skip files under $GOMODCACHE that might be vendored accidentally
		if strings.Contains(path, string(os.PathSeparator)+"go"+string(os.PathSeparator)+"pkg"+string(os.PathSeparator)+"mod"+string(os.PathSeparator)) {
			return nil
		}
		// Skip generated cover files if any
		if strings.HasSuffix(name, "_gen.go") || strings.Contains(strings.ToLower(name), "generated") {
			return nil
		}
		fd, ferr := os.Open(path)
		if ferr != nil {
			return ferr
		}
		defer fd.Close()
		s := bufio.NewScanner(fd)
		var lines []string
		inBlockComment := false
		for s.Scan() {
			ln := s.Text()
			l := strings.TrimSpace(ln)
			if l == "" {
				continue
			}
			// handle block comments /* */ (naive, line-based)
			if inBlockComment {
				if i := strings.Index(l, "*/"); i >= 0 {
					l = strings.TrimSpace(l[i+2:])
					inBlockComment = false
				} else {
					continue
				}
			}
			for {
				start := strings.Index(l, "/*")
				if start >= 0 {
					end := strings.Index(l[start+2:], "*/")
					if end >= 0 {
						l = strings.TrimSpace(l[:start] + l[start+2+end+2:])
						continue // may be multiple blocks in a line
					}
					// starts a block comment and continues to next lines
					inBlockComment = true
					l = strings.TrimSpace(l[:start])
					break
				}
				break
			}
			// strip line comments //...
			if idx := strings.Index(l, "//"); idx >= 0 {
				l = strings.TrimSpace(l[:idx])
			}
			if l == "" {
				continue
			}
			lines = append(lines, l)
		}
		if err := s.Err(); err != nil {
			return err
		}
		files[path] = &fileData{lines: lines, dupMarks: make([]bool, len(lines))}
		return nil
	})
	if walkErr != nil {
		fmt.Fprintf(os.Stderr, "walk error: %v\n", walkErr)
		os.Exit(1)
	}

	// Build shingles map
	shingles := map[uint64][]occurrence{}
	for path, fd := range files {
		if len(fd.lines) < win {
			continue
		}
		for i := 0; i+win <= len(fd.lines); i++ {
			// Join window lines to a single string and hash
			w := strings.Join(fd.lines[i:i+win], "\n")
			h := hashString(w)
			shingles[h] = append(shingles[h], occurrence{file: path, idx: i})
		}
	}

	// Mark duplicated windows (occurrences with count >= 2)
	duplicateWindows := 0
	for _, occs := range shingles {
		if len(occs) < 2 {
			continue
		}
		duplicateWindows++
		for _, oc := range occs {
			fd := files[oc.file]
			for j := 0; j < win && oc.idx+j < len(fd.dupMarks); j++ {
				fd.dupMarks[oc.idx+j] = true
			}
		}
	}

	// Compute totals
	totalLines := 0
	duplLines := 0
	paths := make([]string, 0, len(files))
	for p := range files {
		paths = append(paths, p)
	}
	sort.Strings(paths)
	for _, p := range paths {
		fd := files[p]
		totalLines += len(fd.lines)
		for _, m := range fd.dupMarks {
			if m {
				duplLines++
			}
		}
	}

	dupPct := 0.0
	if totalLines > 0 {
		dupPct = float64(duplLines) / float64(totalLines) * 100
	}

	fmt.Printf("Files: %d\n", len(files))
	fmt.Printf("Total normalized lines: %d\n", totalLines)
	fmt.Printf("Duplicate windows (size=%d): %d\n", win, duplicateWindows)
	fmt.Printf("Duplicated lines (estimated): %d\n", duplLines)
	fmt.Printf("Duplication rate (%% of lines duplicated): %.2f%%\n", dupPct)
}

func hashString(s string) uint64 {
	h := fnv.New64a()
	_, _ = h.Write([]byte(s))
	return h.Sum64()
}
