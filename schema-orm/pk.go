package schema_orm

import (
	"bytes"
	"encoding/gob"
	"reflect"
)

// PK mirrors xorm.io/xorm/schemas.PK with JSON/YAML tags
// In JSON/YAML it's naturally represented as an array
// so we don't need custom marshal methods for those formats.
type PK []interface{}

func NewPK(pks ...interface{}) *PK {
	p := PK(pks)
	return &p
}

func (p *PK) IsZero() bool {
	for _, k := range *p {
		if isZero(k) {
			return true
		}
	}
	return false
}

// ToString convert to a gob-encoded string
func (p *PK) ToString() (string, error) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(*p)
	return buf.String(), err
}

// FromString decode from gob string
func (p *PK) FromString(content string) error {
	dec := gob.NewDecoder(bytes.NewBufferString(content))
	return dec.Decode(p)
}

// isZero determines whether v is zero-like (simplified utils.IsZero)
func isZero(v interface{}) bool {
	switch x := v.(type) {
	case nil:
		return true
	case string:
		return x == ""
	case bool:
		return !x
	case int:
		return x == 0
	case int8:
		return x == 0
	case int16:
		return x == 0
	case int32:
		return x == 0
	case int64:
		return x == 0
	case uint, uint8, uint16, uint32, uint64:
		return reflect.ValueOf(v).Uint() == 0
	case float32:
		return x == 0
	case float64:
		return x == 0
	}
	return false
}
