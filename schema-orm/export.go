package schema_orm

import (
	"fmt"
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

// BuildPostgresDSN builds a standard PostgreSQL DSN for xorm/postgres driver.
// Example: postgres://user:pass@localhost:5432/postgres?sslmode=disable
func BuildPostgresDSN(host, user, password, db string) string {
	addr := host
	if addr == "" {
		addr = "localhost:5432"
	} else if !containsPort(addr) {
		addr = fmt.Sprintf("%s:5432", addr)
	}
	// Prefer URL form for lib/pq
	if password == "" {
		return fmt.Sprintf("postgres://%s@%s/%s?sslmode=disable", user, addr, db)
	}
	return fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", user, password, addr, db)
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
