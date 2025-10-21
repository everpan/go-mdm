package utils

import (
	"errors"
	"fmt"
	"reflect"
)

// FuncNameList returns a sorted list of exported method names for the given target.
func FuncNameList(target any) []string {
	if target == nil {
		return nil
	}
	t := reflect.TypeOf(target)
	n := t.NumMethod()
	names := make([]string, 0, n)
	for i := 0; i < n; i++ {
		m := t.Method(i)
		if m.PkgPath != "" { // unexported
			continue
		}
		names = append(names, m.Name)
	}
	// stable order by natural method set order; tests can sort if needed
	return names
}

// CallMethod calls an exported method by name using reflection.
// If the method returns an error as the last value, and it is non-nil, it will be returned.
// All non-error return values are returned as a slice of interface{}.
func CallMethod(target any, name string, args ...any) ([]any, error) {
	if target == nil {
		return nil, errors.New("target is nil")
	}
	v := reflect.ValueOf(target)
	m := v.MethodByName(name)
	if !m.IsValid() {
		return nil, fmt.Errorf("method %s not found", name)
	}
	mt := m.Type()
	if mt.IsVariadic() {
		// For simplicity, this helper does not support variadic methods.
		return nil, fmt.Errorf("variadic methods not supported: %s", name)
	}
	if len(args) != mt.NumIn() {
		return nil, fmt.Errorf("method %s requires %d args, got %d", name, mt.NumIn(), len(args))
	}
	callArgs := make([]reflect.Value, 0, len(args))
	for i := 0; i < len(args); i++ {
		av := reflect.ValueOf(args[i])
		pt := mt.In(i)
		if !av.IsValid() {
			callArgs = append(callArgs, reflect.Zero(pt))
			continue
		}
		if av.Type().AssignableTo(pt) {
			callArgs = append(callArgs, av)
			continue
		}
		if av.Type().ConvertibleTo(pt) {
			callArgs = append(callArgs, av.Convert(pt))
			continue
		}
		return nil, fmt.Errorf("arg %d not assignable to %s", i, pt.String())
	}
	outs := m.Call(callArgs)
	if len(outs) == 0 {
		return nil, nil
	}
	// If last return is error
	if outs[len(outs)-1].Type() == reflect.TypeOf((*error)(nil)).Elem() {
		errV := outs[len(outs)-1]
		outs = outs[:len(outs)-1]
		if !errV.IsZero() && !errV.IsNil() {
			return nil, errV.Interface().(error)
		}
	}
	res := make([]any, len(outs))
	for i := range outs {
		res[i] = outs[i].Interface()
	}
	return res, nil
}
