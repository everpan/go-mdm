package utils

import (
	"reflect"
	"testing"
)

func TestTryConvert_NilSliceAndMap(t *testing.T) {
	// Ensure nil can be converted to slice and map types without panic and reports success
	sliceT := reflect.TypeOf([]int{})
	if v, ok := tryConvert(nil, sliceT); !ok || !v.IsValid() || v.Kind() != reflect.Slice {
		t.Fatalf("nil to slice failed: ok=%v v=%v kind=%v", ok, v, v.Kind())
	}

	mapT := reflect.TypeOf(map[string]string{})
	if v, ok := tryConvert(nil, mapT); !ok || !v.IsValid() || v.Kind() != reflect.Map {
		t.Fatalf("nil to map failed: ok=%v v=%v kind=%v", ok, v, v.Kind())
	}
}
