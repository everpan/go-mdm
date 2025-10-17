package schema_orm

import (
	"errors"
	"reflect"
	"strconv"
	"time"
)

// mapping types
const (
	TWOSIDES = iota + 1
	ONLYTODB
	ONLYFROMDB
)

// Column mirrors xorm.io/xorm/schemas.Column with JSON/YAML tags
// and the same field names /types.
type Column struct {
	Name            string         `json:"name" yaml:"name"`
	TableName       string         `json:"tableName,omitempty" yaml:"tableName,omitempty"`
	FieldName       string         `json:"fieldName,omitempty" yaml:"fieldName,omitempty"`
	FieldIndex      []int          `json:"fieldIndex,omitempty" yaml:"fieldIndex,omitempty"`
	SQLType         SQLType        `json:"sqlType" yaml:"sqlType"`
	IsJSON          bool           `json:"isJSON,omitempty" yaml:"isJSON,omitempty"`
	IsJSONB         bool           `json:"isJSONB,omitempty" yaml:"isJSONB,omitempty"`
	Length          int64          `json:"length,omitempty" yaml:"length,omitempty"`
	Length2         int64          `json:"length2,omitempty" yaml:"length2,omitempty"`
	Nullable        bool           `json:"nullable,omitempty" yaml:"nullable,omitempty"`
	Default         string         `json:"default,omitempty" yaml:"default,omitempty"`
	Indexes         map[string]int `json:"indexes,omitempty" yaml:"indexes,omitempty"`
	IsPrimaryKey    bool           `json:"isPrimaryKey,omitempty" yaml:"isPrimaryKey,omitempty"`
	IsAutoIncrement bool           `json:"isAutoIncrement,omitempty" yaml:"isAutoIncrement,omitempty"`
	MapType         int            `json:"mapType,omitempty" yaml:"mapType,omitempty"`
	IsCreated       bool           `json:"isCreated,omitempty" yaml:"isCreated,omitempty"`
	IsUpdated       bool           `json:"isUpdated,omitempty" yaml:"isUpdated,omitempty"`
	IsDeleted       bool           `json:"isDeleted,omitempty" yaml:"isDeleted,omitempty"`
	IsCascade       bool           `json:"isCascade,omitempty" yaml:"isCascade,omitempty"`
	IsVersion       bool           `json:"isVersion,omitempty" yaml:"isVersion,omitempty"`
	DefaultIsEmpty  bool           `json:"defaultIsEmpty,omitempty" yaml:"defaultIsEmpty,omitempty"`
	EnumOptions     map[string]int `json:"enumOptions,omitempty" yaml:"enumOptions,omitempty"`
	SetOptions      map[string]int `json:"setOptions,omitempty" yaml:"setOptions,omitempty"`
	DisableTimeZone bool           `json:"disableTimeZone,omitempty" yaml:"disableTimeZone,omitempty"`
	TimeZone        *time.Location `json:"timeZone,omitempty" yaml:"timeZone,omitempty"`
	Comment         string         `json:"comment,omitempty" yaml:"comment,omitempty"`
	Collation       string         `json:"collation,omitempty" yaml:"collation,omitempty"`
}

func NewColumn(name, fieldName string, sqlType SQLType, len1, len2 int64, nullable bool) *Column {
	return &Column{
		Name:            name,
		IsJSON:          sqlType.IsJson(),
		TableName:       "",
		FieldName:       fieldName,
		SQLType:         sqlType,
		Length:          len1,
		Length2:         len2,
		Nullable:        nullable,
		Default:         "",
		Indexes:         make(map[string]int),
		IsPrimaryKey:    false,
		IsAutoIncrement: false,
		MapType:         TWOSIDES,
		IsCreated:       false,
		IsUpdated:       false,
		IsDeleted:       false,
		IsCascade:       false,
		IsVersion:       false,
		DefaultIsEmpty:  true,
		EnumOptions:     make(map[string]int),
		Comment:         "",
	}
}

func (col *Column) ValueOf(bean interface{}) (*reflect.Value, error) {
	dataStruct := reflect.Indirect(reflect.ValueOf(bean))
	return col.ValueOfV(&dataStruct)
}

func (col *Column) ValueOfV(dataStruct *reflect.Value) (*reflect.Value, error) {
	v := *dataStruct
	for _, i := range col.FieldIndex {
		if v.Kind() == reflect.Ptr {
			if v.IsNil() {
				v.Set(reflect.New(v.Type().Elem()))
			}
			v = v.Elem()
		} else if v.Kind() == reflect.Interface {
			v = reflect.Indirect(v.Elem())
		}
		v = v.FieldByIndex([]int{i})
	}
	return &v, nil
}

func (col *Column) ConvertID(sid string) (interface{}, error) {
	if col.SQLType.IsNumeric() {
		n, err := strconv.ParseInt(sid, 10, 64)
		if err != nil {
			return nil, err
		}
		return n, nil
	} else if col.SQLType.IsText() {
		return sid, nil
	}
	return nil, errors.New("not supported")
}
