package schema_orm

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
)

func TestBuildMySQLDSN(t *testing.T) {
	dsn := BuildMySQLDSN("localhost", "root", "", "wiz_hr2")
	if !strings.Contains(dsn, "@tcp(localhost:3306)/wiz_hr2") {
		t.Fatalf("unexpected dsn: %s", dsn)
	}

	dsn2 := BuildMySQLDSN("localhost:3307", "u", "p", "db1")
	if !strings.Contains(dsn2, "u:p@tcp(localhost:3307)/db1") {
		t.Fatalf("unexpected dsn2: %s", dsn2)
	}

	dsn3 := BuildMySQLDSN("", "u", "p", "db1")
	if !strings.Contains(dsn3, "@tcp(localhost:3306)/db1") {
		t.Fatalf("unexpected dsn3: %s", dsn3)
	}
}

// TestExportMySQLToJSON_EnvOrLocal tries to export wiz_hr2 schema.
// It will be skipped unless a local MySQL is available or WIZ_HR2_DSN is provided.
func TestExportMySQLToJSON_EnvOrLocal(t *testing.T) {
	// Allow overriding full DSN
	if dsn := os.Getenv("WIZ_HR2_DSN"); dsn != "" {
		jsonStr, err := ExportMySQLSchemaToJSONWithDSN(dsn)
		if err != nil {
			t.Skipf("skip: cannot connect using WIZ_HR2_DSN: %v", err)
		}
		// Validate JSON shape
		var tables []*Table
		if err := json.Unmarshal([]byte(jsonStr), &tables); err != nil {
			t.Fatalf("invalid json output: %v", err)
		}
		// Not asserting non-empty to be lenient; DB might be empty in CI
		// But ensure structure items (if any) have names.
		for i, tb := range tables {
			if tb.Name == "" {
				t.Fatalf("table %d has empty name", i)
			}
		}
		return
	}

	// Try localhost defaults: user root, password from env MYSQL_PASSWORD, db wiz_hr2
	pass := os.Getenv("MYSQL_PASSWORD")
	jsonStr, err := ExportMySQLToJSON("localhost", "root", pass, "wiz_hr2")
	if err := os.WriteFile("schema.json", []byte(jsonStr), 0644); err != nil {
		t.Fatal(err)
	}
	// t.Log(jsonStr)

	if err != nil {
		t.Skipf("skip: cannot connect to local MySQL wiz_hr2: %v (set WIZ_HR2_DSN to enable)", err)
	}

	var tables []*Table
	if err := json.Unmarshal([]byte(jsonStr), &tables); err != nil {
		t.Fatalf("invalid json output: %v", err)
	}
	for i, tb := range tables {
		if tb.Name == "" {
			t.Fatalf("table %d has empty name", i)
		}
	}
}
