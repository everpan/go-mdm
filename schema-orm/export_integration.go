//go:build integration

package schema_orm

import (
	"encoding/json"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"xorm.io/xorm"
)

// ExportMySQLToJSON connects to MySQL via xorm Engine, introspects the database schema,
// converts to schema-orm structures, and returns the JSON string.
//
// host example: "localhost" or "localhost:3306"
// user example: "root"
// password may be empty if not needed
// db example: "wiz_hr2"
func ExportMySQLToJSON(host, user, password, db string) (string, error) {
	dsn := BuildMySQLDSN(host, user, password, db)
	return ExportMySQLSchemaToJSONWithDSN(dsn)
}

// ExportMySQLSchemaToJSONWithDSN does the same as ExportMySQLToJSON but accepts a DSN directly.
func ExportMySQLSchemaToJSONWithDSN(dsn string) (string, error) {
	engine, err := xorm.NewEngine("mysql", dsn)
	if err != nil {
		return "", err
	}
	defer engine.Close()

	if err := engine.Ping(); err != nil {
		return "", err
	}

	metas, err := engine.DBMetas()
	if err != nil {
		return "", err
	}

	// Convert all xorm tables to our tables
	out := make([]*Table, 0, len(metas))
	for _, xt := range metas {
		out = append(out, FromXormTable(xt))
	}

	b, err := json.Marshal(out)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// ExportPostgresToJSON connects to PostgreSQL via xorm Engine, introspects the database schema,
// converts to schema-orm structures, and returns the JSON string.
// host example: "localhost" or "localhost:5432"
// user example: "ever"
// db example: "postgres"
func ExportPostgresToJSON(host, user, password, db string) (string, error) {
	dsn := BuildPostgresDSN(host, user, password, db)
	return ExportPostgresSchemaToJSONWithDSN(dsn)
}

// ExportPostgresSchemaToJSONWithDSN does the same as ExportPostgresToJSON but accepts a DSN directly.
func ExportPostgresSchemaToJSONWithDSN(dsn string) (string, error) {
	engine, err := xorm.NewEngine("postgres", dsn)
	if err != nil {
		return "", err
	}
	defer engine.Close()

	if err := engine.Ping(); err != nil {
		return "", err
	}

	metas, err := engine.DBMetas()
	if err != nil {
		return "", err
	}

	out := make([]*Table, 0, len(metas))
	for _, xt := range metas {
		out = append(out, FromXormTable(xt))
	}

	b, err := json.Marshal(out)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
