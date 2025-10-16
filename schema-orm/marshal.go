package schema_orm

import (
	"encoding/json"

	"gopkg.in/yaml.v3"
)

// Aliases to avoid recursion

type sqlTypeAlias SQLType

type columnAlias Column

type indexAlias Index

type tableAlias Table

type pkAlias PK

// MarshalJSON SQLType JSON/YAML
func (s *SQLType) MarshalJSON() ([]byte, error) { return json.Marshal(sqlTypeAlias(*s)) }
func (s *SQLType) UnmarshalJSON(b []byte) error {
	var a sqlTypeAlias
	if err := json.Unmarshal(b, &a); err != nil {
		return err
	}
	*s = SQLType(a)
	return nil
}
func (s *SQLType) MarshalYAML() (interface{}, error) { return sqlTypeAlias(*s), nil }
func (s *SQLType) UnmarshalYAML(value *yaml.Node) error {
	var a sqlTypeAlias
	if err := value.Decode(&a); err != nil {
		return err
	}
	*s = SQLType(a)
	return nil
}

// Column JSON/YAML
func (col *Column) MarshalJSON() ([]byte, error) { return json.Marshal(columnAlias(*col)) }
func (col *Column) UnmarshalJSON(b []byte) error {
	var a columnAlias
	if err := json.Unmarshal(b, &a); err != nil {
		return err
	}
	*col = Column(a)
	return nil
}
func (col *Column) MarshalYAML() (interface{}, error) { return columnAlias(*col), nil }
func (col *Column) UnmarshalYAML(value *yaml.Node) error {
	var a columnAlias
	if err := value.Decode(&a); err != nil {
		return err
	}
	*col = Column(a)
	return nil
}

// Index JSON/YAML
func (index *Index) MarshalJSON() ([]byte, error) { return json.Marshal(indexAlias(*index)) }
func (index *Index) UnmarshalJSON(b []byte) error {
	var a indexAlias
	if err := json.Unmarshal(b, &a); err != nil {
		return err
	}
	*index = Index(a)
	return nil
}
func (index *Index) MarshalYAML() (interface{}, error) { return indexAlias(*index), nil }
func (index *Index) UnmarshalYAML(value *yaml.Node) error {
	var a indexAlias
	if err := value.Decode(&a); err != nil {
		return err
	}
	*index = Index(a)
	return nil
}

// Table JSON/YAML uses a DTO to include unexported fields

type tableDTO struct {
	Name          string            `json:"name" yaml:"name"`
	ColumnsSeq    []string          `json:"columnsSeq" yaml:"columnsSeq"`
	Columns       []*Column         `json:"columns" yaml:"columns"`
	Indexes       map[string]*Index `json:"indexes" yaml:"indexes"`
	PrimaryKeys   []string          `json:"primaryKeys" yaml:"primaryKeys"`
	AutoIncrement string            `json:"autoIncrement" yaml:"autoIncrement"`
	Created       map[string]bool   `json:"created" yaml:"created"`
	Updated       string            `json:"updated" yaml:"updated"`
	Deleted       string            `json:"deleted" yaml:"deleted"`
	Version       string            `json:"version" yaml:"version"`
	StoreEngine   string            `json:"storeEngine" yaml:"storeEngine"`
	Charset       string            `json:"charset" yaml:"charset"`
	Comment       string            `json:"comment" yaml:"comment"`
	Collation     string            `json:"collation" yaml:"collation"`
}

func (table *Table) MarshalJSON() ([]byte, error) {
	d := tableDTO{
		Name:          table.Name,
		ColumnsSeq:    append([]string(nil), table.columnsSeq...),
		Columns:       append([]*Column(nil), table.columns...),
		Indexes:       table.Indexes,
		PrimaryKeys:   append([]string(nil), table.PrimaryKeys...),
		AutoIncrement: table.AutoIncrement,
		Created:       table.Created,
		Updated:       table.Updated,
		Deleted:       table.Deleted,
		Version:       table.Version,
		StoreEngine:   table.StoreEngine,
		Charset:       table.Charset,
		Comment:       table.Comment,
		Collation:     table.Collation,
	}
	return json.Marshal(d)
}

func (table *Table) UnmarshalJSON(b []byte) error {
	var d tableDTO
	if err := json.Unmarshal(b, &d); err != nil {
		return err
	}
	nt := NewTable(d.Name, nil)
	nt.columnsSeq = append(nt.columnsSeq, d.ColumnsSeq...)
	for _, c := range d.Columns {
		nt.AddColumn(c)
	}
	nt.Indexes = d.Indexes
	nt.PrimaryKeys = append(nt.PrimaryKeys, d.PrimaryKeys...)
	nt.AutoIncrement = d.AutoIncrement
	nt.Created = d.Created
	nt.Updated = d.Updated
	nt.Deleted = d.Deleted
	nt.Version = d.Version
	nt.StoreEngine = d.StoreEngine
	nt.Charset = d.Charset
	nt.Comment = d.Comment
	nt.Collation = d.Collation
	*table = *nt
	return nil
}

func (table *Table) MarshalYAML() (interface{}, error) {
	d := tableDTO{
		Name:          table.Name,
		ColumnsSeq:    append([]string(nil), table.columnsSeq...),
		Columns:       append([]*Column(nil), table.columns...),
		Indexes:       table.Indexes,
		PrimaryKeys:   append([]string(nil), table.PrimaryKeys...),
		AutoIncrement: table.AutoIncrement,
		Created:       table.Created,
		Updated:       table.Updated,
		Deleted:       table.Deleted,
		Version:       table.Version,
		StoreEngine:   table.StoreEngine,
		Charset:       table.Charset,
		Comment:       table.Comment,
		Collation:     table.Collation,
	}
	return d, nil
}

func (table *Table) UnmarshalYAML(value *yaml.Node) error {
	var d tableDTO
	if err := value.Decode(&d); err != nil {
		return err
	}
	nt := NewTable(d.Name, nil)
	nt.columnsSeq = append(nt.columnsSeq, d.ColumnsSeq...)
	for _, c := range d.Columns {
		nt.AddColumn(c)
	}
	nt.Indexes = d.Indexes
	nt.PrimaryKeys = append(nt.PrimaryKeys, d.PrimaryKeys...)
	nt.AutoIncrement = d.AutoIncrement
	nt.Created = d.Created
	nt.Updated = d.Updated
	nt.Deleted = d.Deleted
	nt.Version = d.Version
	nt.StoreEngine = d.StoreEngine
	nt.Charset = d.Charset
	nt.Comment = d.Comment
	nt.Collation = d.Collation
	*table = *nt
	return nil
}

// PK JSON/YAML
func (p *PK) MarshalJSON() ([]byte, error) { return json.Marshal(pkAlias(*p)) }
func (p *PK) UnmarshalJSON(b []byte) error {
	var a pkAlias
	if err := json.Unmarshal(b, &a); err != nil {
		return err
	}
	*p = PK(a)
	return nil
}
func (p *PK) MarshalYAML() (interface{}, error) { return pkAlias(*p), nil }
func (p *PK) UnmarshalYAML(value *yaml.Node) error {
	var a pkAlias
	if err := value.Decode(&a); err != nil {
		return err
	}
	*p = PK(a)
	return nil
}
