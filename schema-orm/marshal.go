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

// SQLType JSON/YAML
func (s SQLType) MarshalJSON() ([]byte, error) { return json.Marshal(sqlTypeAlias(s)) }
func (s *SQLType) UnmarshalJSON(b []byte) error { var a sqlTypeAlias; if err := json.Unmarshal(b, &a); err != nil { return err }; *s = SQLType(a); return nil }
func (s SQLType) MarshalYAML() (interface{}, error) { return sqlTypeAlias(s), nil }
func (s *SQLType) UnmarshalYAML(value *yaml.Node) error { var a sqlTypeAlias; if err := value.Decode(&a); err != nil { return err }; *s = SQLType(a); return nil }

// Column JSON/YAML
func (c Column) MarshalJSON() ([]byte, error) { return json.Marshal(columnAlias(c)) }
func (c *Column) UnmarshalJSON(b []byte) error { var a columnAlias; if err := json.Unmarshal(b, &a); err != nil { return err }; *c = Column(a); return nil }
func (c Column) MarshalYAML() (interface{}, error) { return columnAlias(c), nil }
func (c *Column) UnmarshalYAML(value *yaml.Node) error { var a columnAlias; if err := value.Decode(&a); err != nil { return err }; *c = Column(a); return nil }

// Index JSON/YAML
func (i Index) MarshalJSON() ([]byte, error) { return json.Marshal(indexAlias(i)) }
func (i *Index) UnmarshalJSON(b []byte) error { var a indexAlias; if err := json.Unmarshal(b, &a); err != nil { return err }; *i = Index(a); return nil }
func (i Index) MarshalYAML() (interface{}, error) { return indexAlias(i), nil }
func (i *Index) UnmarshalYAML(value *yaml.Node) error { var a indexAlias; if err := value.Decode(&a); err != nil { return err }; *i = Index(a); return nil }

// Table JSON/YAML uses a DTO to include unexported fields

type tableDTO struct {
	Name          string               `json:"name" yaml:"name"`
	ColumnsSeq    []string             `json:"columnsSeq" yaml:"columnsSeq"`
	Columns       []*Column            `json:"columns" yaml:"columns"`
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

func (t Table) MarshalJSON() ([]byte, error) {
	d := tableDTO{
		Name: t.Name,
		ColumnsSeq: append([]string(nil), t.columnsSeq...),
		Columns: append([]*Column(nil), t.columns...),
		Indexes: t.Indexes,
		PrimaryKeys: append([]string(nil), t.PrimaryKeys...),
		AutoIncrement: t.AutoIncrement,
		Created: t.Created,
		Updated: t.Updated,
		Deleted: t.Deleted,
		Version: t.Version,
		StoreEngine: t.StoreEngine,
		Charset: t.Charset,
		Comment: t.Comment,
		Collation: t.Collation,
	}
	return json.Marshal(d)
}

func (t *Table) UnmarshalJSON(b []byte) error {
	var d tableDTO
	if err := json.Unmarshal(b, &d); err != nil { return err }
	nt := NewTable(d.Name, nil)
	nt.columnsSeq = append(nt.columnsSeq, d.ColumnsSeq...)
	for _, c := range d.Columns { nt.AddColumn(c) }
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
	*t = *nt
	return nil
}

func (t Table) MarshalYAML() (interface{}, error) {
	d := tableDTO{
		Name: t.Name,
		ColumnsSeq: append([]string(nil), t.columnsSeq...),
		Columns: append([]*Column(nil), t.columns...),
		Indexes: t.Indexes,
		PrimaryKeys: append([]string(nil), t.PrimaryKeys...),
		AutoIncrement: t.AutoIncrement,
		Created: t.Created,
		Updated: t.Updated,
		Deleted: t.Deleted,
		Version: t.Version,
		StoreEngine: t.StoreEngine,
		Charset: t.Charset,
		Comment: t.Comment,
		Collation: t.Collation,
	}
	return d, nil
}

func (t *Table) UnmarshalYAML(value *yaml.Node) error {
	var d tableDTO
	if err := value.Decode(&d); err != nil { return err }
	nt := NewTable(d.Name, nil)
	nt.columnsSeq = append(nt.columnsSeq, d.ColumnsSeq...)
	for _, c := range d.Columns { nt.AddColumn(c) }
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
	*t = *nt
	return nil
}

// PK JSON/YAML
func (p PK) MarshalJSON() ([]byte, error) { return json.Marshal(pkAlias(p)) }
func (p *PK) UnmarshalJSON(b []byte) error { var a pkAlias; if err := json.Unmarshal(b, &a); err != nil { return err }; *p = PK(a); return nil }
func (p PK) MarshalYAML() (interface{}, error) { return pkAlias(p), nil }
func (p *PK) UnmarshalYAML(value *yaml.Node) error { var a pkAlias; if err := value.Decode(&a); err != nil { return err }; *p = PK(a); return nil }
