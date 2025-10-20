package schema_orm

import "encoding/json"

// ImportTablesFromJSON parses a JSON string into a slice of Table pointers.
// The JSON is expected to represent an array of Table objects (as produced by the exporter).
func ImportTablesFromJSON(s string) ([]*Table, error) {
	var tables []*Table
	if err := json.Unmarshal([]byte(s), &tables); err != nil {
		return nil, err
	}
	return tables, nil
}
