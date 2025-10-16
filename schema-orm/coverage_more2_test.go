package schema_orm

import (
	"reflect"
	"testing"
)

func TestType2SQLType_SliceNonByte(t *testing.T) {
	if Type2SQLType(reflect.TypeOf([]int{})).Name != "TEXT" { t.Fatalf("[]int should map to TEXT") }
}

func TestPK_isZero_RemainingCases(t *testing.T) {
	if !isZero(int64(0)) { t.Fatalf("int64 zero") }
	if isZero(uint8(1)) { t.Fatalf("uint8 non-zero") }
}
