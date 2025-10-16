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
	TableName       string         `json:"tableName" yaml:"tableName"`
	FieldName       string         `json:"fieldName" yaml:"fieldName"`
	FieldIndex      []int          `json:"fieldIndex" yaml:"fieldIndex"`
	SQLType         SQLType        `json:"sqlType" yaml:"sqlType"`
	IsJSON          bool           `json:"isJSON" yaml:"isJSON"`
	IsJSONB         bool           `json:"isJSONB" yaml:"isJSONB"`
	Length          int64          `json:"length" yaml:"length"`
	Length2         int64          `json:"length2" yaml:"length2"`
	Nullable        bool           `json:"nullable" yaml:"nullable"`
	Default         string         `json:"default" yaml:"default"`
	Indexes         map[string]int `json:"indexes" yaml:"indexes"`
	IsPrimaryKey    bool           `json:"isPrimaryKey" yaml:"isPrimaryKey"`
	IsAutoIncrement bool           `json:"isAutoIncrement" yaml:"isAutoIncrement"`
	MapType         int            `json:"mapType" yaml:"mapType"`
	IsCreated       bool           `json:"isCreated" yaml:"isCreated"`
	IsUpdated       bool           `json:"isUpdated" yaml:"isUpdated"`
	IsDeleted       bool           `json:"isDeleted" yaml:"isDeleted"`
	IsCascade       bool           `json:"isCascade" yaml:"isCascade"`
	IsVersion       bool           `json:"isVersion" yaml:"isVersion"`
	DefaultIsEmpty  bool           `json:"defaultIsEmpty" yaml:"defaultIsEmpty"`
	EnumOptions     map[string]int `json:"enumOptions" yaml:"enumOptions"`
	SetOptions      map[string]int `json:"setOptions" yaml:"setOptions"`
	DisableTimeZone bool           `json:"disableTimeZone" yaml:"disableTimeZone"`
	TimeZone        *time.Location `json:"timeZone" yaml:"timeZone"`
	Comment         string         `json:"comment" yaml:"comment"`
	Collation       string         `json:"collation" yaml:"collation"`
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
