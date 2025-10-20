//go:build integration

package schema_orm

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	_ "github.com/lib/pq"
	"xorm.io/xorm"
)

func TestImportTablesFromJSON_File(t *testing.T) {
	// This test requires a running local PostgreSQL instance. Skip unless explicitly enabled.
	if os.Getenv("WIZ_PG_DSN") == "" {
		// Allow enabling via WIZ_PG_ENABLE=1 without full DSN; build a local DSN then.
		if os.Getenv("WIZ_PG_ENABLE") == "" {
			t.Skip("skip: no PostgreSQL available (set WIZ_PG_DSN or WIZ_PG_ENABLE=1 to run)")
		}
	}
	// Load the bundled sample schema.json if present
	p := filepath.Join("", "schema.json")
	b, err := os.ReadFile(p)
	if err != nil {
		// If file missing, create a minimal JSON content and test with it
		b = []byte(`[{"name":"demo","columns":[{"name":"id","SQLType":{"name":"INT"},"isPrimaryKey":true}]}]`)
	}
	tables, err := ImportTablesFromJSON(string(b))
	if err != nil {
		t.Fatalf("ImportTablesFromJSON failed: %v", err)
	}
	postgresDSN := BuildPostgresDSN("localhost", "ever", "", "postgres")
	eng, err := xorm.NewEngine("postgres", postgresDSN)
	if nil != err {
		t.Fatalf("NewEngine failed: %v", err)
	}
	ctx := context.Background()
	for i, tb := range tables {
		if tb == nil || tb.Name == "" {
			t.Fatalf("table %d has empty name", i)
		}
		xtb := ToXormTable(tb)
		sql, b2, err := eng.Dialect().CreateTableSQL(ctx, eng.DB(), xtb, tb.Name)
		if err != nil {
			t.Fatalf("CreateTableSQL failed: %v", err)
		}
		t.Logf("sql: %s, b2: %v\n", sql, b2)
		if ret, err := eng.Exec(sql); err != nil {
			b1, _ := json.Marshal(tb)
			b2, _ := json.Marshal(xtb)
			t.Logf("\n tb: %s\n, \nxtb: %s", string(b1), string(b2))
			t.Fatalf("Exec failed: %v", err)
		} else {
			t.Logf("create table return %v", ret)
		}
	}
}

func TestBuildPostgresDSN(t *testing.T) {
	d := BuildPostgresDSN("localhost", "ever", "", "postgres")
	if !strings.Contains(d, "postgres://ever@localhost:5432/postgres") {
		t.Fatalf("unexpected dsn: %s", d)
	}
	d2 := BuildPostgresDSN("localhost:5433", "user", "p@ss", "db1")
	if !strings.Contains(d2, "postgres://user:p@ss@localhost:5433/db1") {
		t.Fatalf("unexpected dsn2: %s", d2)
	}
}

// TestExportPostgresToJSON_EnvOrLocal tries to export postgres schema using xorm.
// It skips if no local DB available. Defaults: host localhost, user ever, db postgres, password from PG_PASSWORD.
func TestExportPostgresToJSON_EnvOrLocal(t *testing.T) {
	if dsn := os.Getenv("WIZ_PG_DSN"); dsn != "" {
		jsonStr, err := ExportPostgresSchemaToJSONWithDSN(dsn)
		if err != nil {
			t.Skipf("skip: cannot connect using WIZ_PG_DSN: %v", err)
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
		return
	}

	pass := os.Getenv("PG_PASSWORD")
	jsonStr, err := ExportPostgresToJSON("localhost", "ever", pass, "postgres")
	if err != nil {
		t.Skipf("skip: cannot connect to local PostgreSQL: %v (set WIZ_PG_DSN to enable)", err)
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
