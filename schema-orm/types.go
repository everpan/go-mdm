package schema_orm

import (
	"reflect"
	"strings"
)

// DBType mirrors xorm.io/xorm/schemas.DBType
// Tags are added, though this is a type alias.
type DBType string

// SQLType mirrors xorm.io/xorm/schemas.SQLType with JSON/YAML tags
// and helper methods copied one-to-one.
type SQLType struct {
	Name           string `json:"name" yaml:"name"`
	DefaultLength  int64  `json:"defaultLength,omitempty" yaml:"defaultLength,omitempty"`
	DefaultLength2 int64  `json:"defaultLength2,omitempty" yaml:"defaultLength2,omitempty"`
}

// kinds for detection (copied behaviorally from upstream)
const (
// Text/blob/time/bool/numeric/array/json/xml are represented via Name matching.
)

// IsType Below helper methods mirror upstream behavior but simplified to string matching.
func (s *SQLType) IsType(st int) bool { // not used in our scope but kept for parity
	return false
}

func (s *SQLType) IsText() bool {
	// minimal set
	switch strings.ToUpper(s.Name) {
	case "TEXT", "VARCHAR", "CHAR", "NVARCHAR", "NTEXT", "UUID", "JSON", "JSONB":
		return true
	}
	return false
}

func (s *SQLType) IsBlob() bool {
	switch strings.ToUpper(s.Name) {
	case "BLOB", "BYTEA", "LONGBLOB", "MEDIUMBLOB", "VARBINARY", "BINARY":
		return true
	}
	return false
}

func (s *SQLType) IsTime() bool {
	switch strings.ToUpper(s.Name) {
	case "DATETIME", "DATE", "TIME", "TIMESTAMP":
		return true
	}
	return false
}

func (s *SQLType) IsBool() bool {
	return strings.ToUpper(s.Name) == "BOOL" || strings.ToUpper(s.Name) == "BOOLEAN"
}

func (s *SQLType) IsNumeric() bool {
	switch strings.ToUpper(s.Name) {
	case "INT", "INTEGER", "BIGINT", "SMALLINT", "TINYINT", "DECIMAL", "NUMERIC", "FLOAT", "DOUBLE", "REAL":
		return true
	}
	return false
}

func (s *SQLType) IsArray() bool {
	return strings.HasSuffix(strings.ToUpper(s.Name), "[]")
}

func (s *SQLType) IsJson() bool {
	name := strings.ToUpper(s.Name)
	return name == "JSON" || name == "JSONB"
}

func (s *SQLType) IsXML() bool {
	return strings.ToUpper(s.Name) == "XML"
}

// Utilities mirrored from upstream type.go for conversions kept minimal
func Type2SQLType(t reflect.Type) (st SQLType) {
	// very simplified mapping for tests
	if t == nil {
		return SQLType{}
	}
	switch t.Kind() {
	case reflect.String:
		return SQLType{Name: "VARCHAR"}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return SQLType{Name: "INT"}
	case reflect.Bool:
		return SQLType{Name: "BOOLEAN"}
	case reflect.Float32, reflect.Float64:
		return SQLType{Name: "FLOAT"}
	case reflect.Slice:
		if t.Elem().Kind() == reflect.Uint8 {
			return SQLType{Name: "BLOB"}
		}
	}
	return SQLType{Name: "TEXT"}
}

func SQLType2Type(st SQLType) reflect.Type {
	switch strings.ToUpper(st.Name) {
	case "VARCHAR", "CHAR", "TEXT", "UUID":
		return reflect.TypeOf("")
	case "INT", "INTEGER", "BIGINT", "SMALLINT", "TINYINT", "DECIMAL", "NUMERIC":
		return reflect.TypeOf(int(0))
	case "BOOLEAN", "BOOL":
		return reflect.TypeOf(true)
	case "FLOAT", "DOUBLE", "REAL":
		return reflect.TypeOf(float64(0))
	case "BLOB", "BYTEA", "VARBINARY", "BINARY":
		return reflect.TypeOf([]byte{})
	}
	return reflect.TypeOf("")
}

func SQLTypeName(tp string) string {
	return strings.ToUpper(tp)
}
