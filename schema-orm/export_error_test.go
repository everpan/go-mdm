//go:build integration

package schema_orm

import "testing"

func TestExportMySQLSchemaToJSONWithDSN_Error(t *testing.T) {
	// Intentionally invalid DSN to trigger connection error and cover error path
	if _, err := ExportMySQLSchemaToJSONWithDSN("invalid:pass@tcp(127.0.0.1:0)/db"); err == nil {
		t.Fatalf("expected error for invalid mysql dsn")
	}
}

func TestExportPostgresSchemaToJSONWithDSN_Error(t *testing.T) {
	if _, err := ExportPostgresSchemaToJSONWithDSN("postgres://bad:bad@127.0.0.1:0/db?sslmode=disable"); err == nil {
		t.Fatalf("expected error for invalid postgres dsn")
	}
}
