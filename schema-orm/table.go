package schema_orm

import (
	"reflect"
	"strconv"
	"strings"
)

// Table mirrors xorm.io/xorm/schemas.Table with JSON/YAML tags
// Note: Type is excluded from serialization because reflect.Type isn't portable.
type Table struct {
	Name          string               `json:"name" yaml:"name"`
	Type          reflect.Type         `json:"-" yaml:"-"`
	columnsSeq    []string             `json:"columnsSeq" yaml:"columnsSeq"`
	columnsMap    map[string][]*Column `json:"columnsMap" yaml:"columnsMap"`
	columns       []*Column            `json:"columns" yaml:"columns"`
	Indexes       map[string]*Index    `json:"indexes" yaml:"indexes"`
	PrimaryKeys   []string             `json:"primaryKeys" yaml:"primaryKeys"`
	AutoIncrement string               `json:"autoIncrement" yaml:"autoIncrement"`
	Created       map[string]bool      `json:"created" yaml:"created"`
	Updated       string               `json:"updated" yaml:"updated"`
	Deleted       string               `json:"deleted" yaml:"deleted"`
	Version       string               `json:"version" yaml:"version"`
	StoreEngine   string               `json:"storeEngine" yaml:"storeEngine"`
	Charset       string               `json:"charset" yaml:"charset"`
	Comment       string               `json:"comment" yaml:"comment"`
	Collation     string               `json:"collation" yaml:"collation"`
}

func NewEmptyTable() *Table { return NewTable("", nil) }

func NewTable(name string, t reflect.Type) *Table {
	return &Table{
		Name:        name,
		Type:        t,
		columnsSeq:  make([]string, 0),
		columns:     make([]*Column, 0),
		columnsMap:  make(map[string][]*Column),
		Indexes:     make(map[string]*Index),
		Created:     make(map[string]bool),
		PrimaryKeys: make([]string, 0),
	}
}

func (table *Table) Columns() []*Column   { return table.columns }
func (table *Table) ColumnsSeq() []string { return table.columnsSeq }

func (table *Table) columnsByName(name string) []*Column {
	return table.columnsMap[strings.ToLower(name)]
}

func (table *Table) GetColumn(name string) *Column {
	cols := table.columnsByName(name)
	if cols != nil {
		return cols[0]
	}
	return nil
}

func (table *Table) GetColumnIdx(name string, idx int) *Column {
	cols := table.columnsByName(name)
	if cols != nil && idx < len(cols) {
		return cols[idx]
	}
	return nil
}

func (table *Table) PKColumns() []*Column {
	columns := make([]*Column, len(table.PrimaryKeys))
	for i, name := range table.PrimaryKeys {
		columns[i] = table.GetColumn(name)
	}
	return columns
}

func (table *Table) ColumnType(name string) reflect.Type {
	t, _ := table.Type.FieldByName(name)
	return t.Type
}

func (table *Table) AutoIncrColumn() *Column { return table.GetColumn(table.AutoIncrement) }
func (table *Table) VersionColumn() *Column  { return table.GetColumn(table.Version) }
func (table *Table) UpdatedColumn() *Column  { return table.GetColumn(table.Updated) }
func (table *Table) DeletedColumn() *Column  { return table.GetColumn(table.Deleted) }

func (table *Table) AddColumn(col *Column) {
	table.columnsSeq = append(table.columnsSeq, col.Name)
	table.columns = append(table.columns, col)
	colName := strings.ToLower(col.Name)
	if c, ok := table.columnsMap[colName]; ok {
		table.columnsMap[colName] = append(c, col)
	} else {
		table.columnsMap[colName] = []*Column{col}
	}

	if col.IsPrimaryKey {
		table.PrimaryKeys = append(table.PrimaryKeys, col.Name)
	}
	if col.IsAutoIncrement {
		table.AutoIncrement = col.Name
	}
	if col.IsCreated {
		table.Created[col.Name] = true
	}
	if col.IsUpdated {
		table.Updated = col.Name
	}
	if col.IsDeleted {
		table.Deleted = col.Name
	}
	if col.IsVersion {
		table.Version = col.Name
	}
}

func (table *Table) AddIndex(index *Index) { table.Indexes[index.Name] = index }

func (table *Table) IDOfV(rv reflect.Value) (PK, error) {
	v := reflect.Indirect(rv)
	pk := make([]interface{}, len(table.PrimaryKeys))
	for i, col := range table.PKColumns() {
		var err error
		pkField := v.FieldByIndex(col.FieldIndex)
		switch pkField.Kind() {
		case reflect.String:
			pk[i], err = col.ConvertID(pkField.String())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			pk[i], err = col.ConvertID(strconv.FormatInt(pkField.Int(), 10))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			pk[i], err = col.ConvertID(strconv.FormatUint(pkField.Uint(), 10))
		default:
			panic("unhandled default case")
		}
		if err != nil {
			return nil, err
		}
	}
	return PK(pk), nil
}
