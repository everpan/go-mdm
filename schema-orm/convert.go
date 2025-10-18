package schema_orm

import (
	xs "xorm.io/xorm/schemas"
)

// ToXormSQLType SQLType conversions
func ToXormSQLType(s SQLType) xs.SQLType {
	return xs.SQLType{Name: s.Name, DefaultLength: s.DefaultLength, DefaultLength2: s.DefaultLength2}
}
func FromXormSQLType(s xs.SQLType) SQLType {
	return SQLType{Name: s.Name, DefaultLength: s.DefaultLength, DefaultLength2: s.DefaultLength2}
}

// ToXormColumn Column conversions
func ToXormColumn(c *Column) *xs.Column {
	if c == nil {
		return nil
	}
	xc := &xs.Column{
		Name:            c.Name,
		TableName:       c.TableName,
		FieldName:       c.FieldName,
		FieldIndex:      c.FieldIndex,
		SQLType:         ToXormSQLType(c.SQLType),
		IsJSON:          c.IsJSON,
		IsJSONB:         c.IsJSONB,
		Length:          c.Length,
		Length2:         c.Length2,
		Nullable:        c.Nullable,
		Default:         c.Default,
		Indexes:         c.Indexes,
		IsPrimaryKey:    c.IsPrimaryKey,
		IsAutoIncrement: c.IsAutoIncrement,
		MapType:         c.MapType,
		IsCreated:       c.IsCreated,
		IsUpdated:       c.IsUpdated,
		IsDeleted:       c.IsDeleted,
		IsCascade:       c.IsCascade,
		IsVersion:       c.IsVersion,
		DefaultIsEmpty:  c.DefaultIsEmpty,
		EnumOptions:     c.EnumOptions,
		SetOptions:      c.SetOptions,
		DisableTimeZone: c.DisableTimeZone,
		TimeZone:        c.TimeZone,
		Comment:         c.Comment,
		Collation:       c.Collation,
	}
	return xc
}

func FromXormColumn(c *xs.Column) *Column {
	if c == nil {
		return nil
	}
	return &Column{
		Name:            c.Name,
		TableName:       c.TableName,
		FieldName:       c.FieldName,
		FieldIndex:      c.FieldIndex,
		SQLType:         FromXormSQLType(c.SQLType),
		IsJSON:          c.IsJSON,
		IsJSONB:         c.IsJSONB,
		Length:          c.Length,
		Length2:         c.Length2,
		Nullable:        c.Nullable,
		Default:         c.Default,
		Indexes:         c.Indexes,
		IsPrimaryKey:    c.IsPrimaryKey,
		IsAutoIncrement: c.IsAutoIncrement,
		MapType:         c.MapType,
		IsCreated:       c.IsCreated,
		IsUpdated:       c.IsUpdated,
		IsDeleted:       c.IsDeleted,
		IsCascade:       c.IsCascade,
		IsVersion:       c.IsVersion,
		DefaultIsEmpty:  c.DefaultIsEmpty,
		EnumOptions:     c.EnumOptions,
		SetOptions:      c.SetOptions,
		DisableTimeZone: c.DisableTimeZone,
		TimeZone:        c.TimeZone,
		Comment:         c.Comment,
		Collation:       c.Collation,
	}
}

// ToXormIndex Index conversions
func ToXormIndex(i *Index) *xs.Index {
	if i == nil {
		return nil
	}
	return &xs.Index{IsRegular: i.IsRegular, Name: i.Name, Type: i.Type, Cols: append([]string(nil), i.Cols...)}
}
func FromXormIndex(i *xs.Index) *Index {
	if i == nil {
		return nil
	}
	return &Index{IsRegular: i.IsRegular, Name: i.Name, Type: i.Type, Cols: append([]string(nil), i.Cols...)}
}

// ToXormPK PK conversions
func ToXormPK(p PK) xs.PK {
	out := make(xs.PK, len(p))
	copy(out, p)
	return out
}
func FromXormPK(p xs.PK) PK {
	out := make(PK, len(p))
	copy(out, p)
	return out
}

// ToXormTable Table conversions (Type cannot be fully mapped via JSON, keep reflection.Type as-is)
func ToXormTable(t *Table) *xs.Table {
	if t == nil {
		return nil
	}
	x := xs.NewTable(t.Name, t.Type)
	x.AutoIncrement = t.AutoIncrement
	x.Updated = t.Updated
	x.Deleted = t.Deleted
	x.Version = t.Version
	x.StoreEngine = t.StoreEngine
	x.Charset = t.Charset
	x.Comment = t.Comment
	x.Collation = t.Collation
	// columns
	for _, c := range t.Columns {
		x.AddColumn(ToXormColumn(c))
	}
	// the index map must be filled
	for k, v := range t.Indexes {
		x.Indexes[k] = ToXormIndex(v)
	}
	// created a map
	for k, v := range t.Created {
		x.Created[k] = v
	}
	return x
}

func FromXormTable(t *xs.Table) *Table {
	if t == nil {
		return nil
	}
	nt := NewTable(t.Name, t.Type)
	nt.AutoIncrement = t.AutoIncrement
	nt.PrimaryKeys = append(nt.PrimaryKeys, t.PrimaryKeys...)
	nt.Updated = t.Updated
	nt.Deleted = t.Deleted
	nt.Version = t.Version
	nt.StoreEngine = t.StoreEngine
	nt.Charset = t.Charset
	nt.Comment = t.Comment
	nt.Collation = t.Collation
	for k, v := range t.Created {
		nt.Created[k] = v
	}
	for _, c := range t.Columns() {
		nt.AddColumn(FromXormColumn(c))
	}
	for k, v := range t.Indexes {
		nt.Indexes[k] = FromXormIndex(v)
	}
	return nt
}
