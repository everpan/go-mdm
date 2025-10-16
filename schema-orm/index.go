package schema_orm

import (
	"fmt"
	"strings"
)

// Index types
const (
	IndexType = iota + 1
	UniqueType
)

// Index mirrors xorm.io/xorm/schemas.Index
// with JSON and YAML tags
type Index struct {
	IsRegular bool     `json:"isRegular" yaml:"isRegular"`
	Name      string   `json:"name" yaml:"name"`
	Type      int      `json:"type" yaml:"type"`
	Cols      []string `json:"cols" yaml:"cols"`
}

func NewIndex(name string, indexType int) *Index {
	return &Index{IsRegular: true, Name: name, Type: indexType, Cols: make([]string, 0)}
}

// XName returns the special index name for the table
func (index *Index) XName(tableName string) string {
	if !strings.HasPrefix(index.Name, "UQE_") &&
		!strings.HasPrefix(index.Name, "IDX_") {
		tableParts := strings.Split(strings.ReplaceAll(tableName, "\"", ""), ".")
		tableName = tableParts[len(tableParts)-1]
		if index.Type == UniqueType {
			return fmt.Sprintf("UQE_%v_%v", tableName, index.Name)
		}
		return fmt.Sprintf("IDX_%v_%v", tableName, index.Name)
	}
	return index.Name
}

// AddColumn add columns which will be composite index
func (index *Index) AddColumn(cols ...string) {
	index.Cols = append(index.Cols, cols...)
}

// Equal return true if the two Index is equal
func (index *Index) Equal(dst *Index) bool {
	if index.Type != dst.Type {
		return false
	}
	if len(index.Cols) != len(dst.Cols) {
		return false
	}

	for i := 0; i < len(index.Cols); i++ {
		var found bool
		for j := 0; j < len(dst.Cols); j++ {
			if index.Cols[i] == dst.Cols[j] {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}
