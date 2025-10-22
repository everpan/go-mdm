package utils

import (
	"reflect"
	"testing"
)

func TestTryConvert_FloatToUint32(t *testing.T) {
	typ := reflect.TypeOf(uint32(0))
	if v, ok := tryConvert(5.0, typ); !ok || !v.IsValid() || v.Uint() != uint64(5) {
		t.Fatalf("float->uint32 failed: ok=%v v=%v", ok, v)
	}
}

func TestTryConvert_IntToFloat64(t *testing.T) {
	typ := reflect.TypeOf(float64(0))
	if v, ok := tryConvert(int(7), typ); !ok || !v.IsValid() || v.Float() != float64(7) {
		t.Fatalf("int->float64 failed: ok=%v v=%v", ok, v)
	}
}
