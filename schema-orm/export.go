package schema_orm

import (
	"encoding/json"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
)

// BuildMySQLDSN builds a standard MySQL DSN for xorm/mysql driver.
// Example: root:pass@tcp(localhost:3306)/wiz_hr2?charset=utf8mb4&parseTime=True&loc=Local
func BuildMySQLDSN(host, user, password, db string) string {
	addr := host
	// allow passing host with port, default to 3306 if no port specified
	if addr == "" {
		addr = "localhost:3306"
	} else if !containsPort(addr) {
		addr = fmt.Sprintf("%s:3306", addr)
	}
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, addr, db)
}

func containsPort(h string) bool {
	for i := len(h) - 1; i >= 0; i-- {
		if h[i] == ':' {
			return true
		}
		if h[i] == ']' { // IPv6 literal end; port would be after this, so keep simple
			break
		}
	}
	return false
}

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

	// b, err := json.MarshalIndent(out, "", "  ")
	b, err := json.Marshal(out)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
