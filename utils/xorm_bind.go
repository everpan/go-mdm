package utils

import (
	"fmt"
	"reflect"

	"github.com/dop251/goja"
	"xorm.io/xorm"
)

// NewXORM creates a new xorm Engine using the given driver and DSN.
// It is a light wrapper that keeps the utils package self-contained for tests.
func NewXORM(driver, dsn string) (*xorm.Engine, error) {
	return xorm.NewEngine(driver, dsn)
}

// BindAllMethods exposes all exported methods of the target (usually *xorm.Engine)
// as functions on a new JavaScript object. Methods are invoked via reflection.
//
// Behavior:
//   - Only exported methods are exposed.
//   - Arguments are converted using goja's ExportTo to the method parameter types.
//   - If the last return value is an error, and it's non-nil, a JS Error is thrown.
//   - If there is a single non-error return value, it is returned directly.
//   - If there are multiple return values (excluding an error), an Array is returned.
func BindAllMethods(rt *goja.Runtime, target any) *goja.Object {
	obj := rt.NewObject()
	if target == nil {
		return obj
	}
	val := reflect.ValueOf(target)
	typ := val.Type()

	for i := 0; i < typ.NumMethod(); i++ {
		m := typ.Method(i)
		// Skip unexported
		if m.PkgPath != "" {
			continue
		}

		name := m.Name // capture
		meth := m      // capture
		_ = obj.Set(name, func(call goja.FunctionCall) goja.Value {
			// Prepare arguments
			mt := meth.Type
			inTotal := mt.NumIn() // includes receiver
			isVar := mt.IsVariadic()
			// Exclude receiver for argument expectations
			fixedParams := inTotal - 1
			if isVar {
				fixedParams = inTotal - 2 // last is the variadic slice type
			}
			if (!isVar && len(call.Arguments) != fixedParams) || (isVar && len(call.Arguments) < fixedParams) {
				err := fmt.Errorf("method %s requires %d arguments%v, got %d", name, fixedParams, ternary(isVar, "+variadic", ""), len(call.Arguments))
				panic(rt.NewTypeError(err.Error()))
			}
			args := make([]reflect.Value, 0, inTotal)
			args = append(args, val)
			// Fixed params
			for pi := 1; pi <= fixedParams; pi++ {
				pt := mt.In(pi)
				// Create a zero-value pointer to use ExportTo
				argPtr := reflect.New(pt)
				gojaVal := call.Arguments[pi-1]
				if err := rt.ExportTo(gojaVal, argPtr.Interface()); err != nil {
					// Fallback: try using the generic Export and Convert if types are compatible
					exported := gojaVal.Export()
					a, convOK := tryConvert(exported, pt)
					if !convOK {
						panic(rt.NewTypeError("convert argument %d for %s failed: %v", pi-1, name, err))
					}
					args = append(args, a)
					continue
				}
				args = append(args, argPtr.Elem())
			}

			// Variadic
			var results []reflect.Value
			if isVar {
				sliceT := mt.In(inTotal - 1)
				elemT := sliceT.Elem()
				nVar := len(call.Arguments) - fixedParams
				slice := reflect.MakeSlice(sliceT, 0, nVar)
				// If a single JS array is passed for varargs, expand it
				if nVar == 1 {
					if obj, ok := call.Arguments[fixedParams].Export().([]interface{}); ok {
						for _, it := range obj {
							v, ok := tryConvert(it, elemT)
							if !ok {
								panic(rt.NewTypeError("convert variadic element for %s failed", name))
							}
							slice = reflect.Append(slice, v)
						}
						args = append(args, slice)
						results = meth.Func.CallSlice(args)
						goto AFTER_CALL
					}
				}
				for i := fixedParams; i < len(call.Arguments); i++ {
					elPtr := reflect.New(elemT)
					gv := call.Arguments[i]
					if err := rt.ExportTo(gv, elPtr.Interface()); err != nil {
						exported := gv.Export()
						v, ok := tryConvert(exported, elemT)
						if !ok {
							panic(rt.NewTypeError("convert variadic element for %s failed: %v", name, err))
						}
						slice = reflect.Append(slice, v)
					} else {
						slice = reflect.Append(slice, elPtr.Elem())
					}
				}
				args = append(args, slice)
				results = meth.Func.CallSlice(args)
			} else {
				// Call method
				results = meth.Func.Call(args)
			}
		AFTER_CALL:
			n := len(results)
			if n == 0 {
				return goja.Undefined()
			}

			// Handle trailing error
			if _, hasErr := trailingErrorType(meth.Type); hasErr {
				errVal := results[n-1]
				if !errVal.IsZero() && !errVal.IsNil() {
					err := errVal.Interface().(error)
					panic(rt.NewGoError(err))
				}
				results = results[:n-1]
			}

			// Convert results to JS
			sz := len(results)
			if sz == 0 {
				return goja.Undefined()
			}
			wrap := func(iv any) goja.Value {
				switch v := iv.(type) {
				case *xorm.Engine:
					return BindAllMethods(rt, v)
				case *xorm.Session:
					return BindAllMethods(rt, v)
				default:
					return rt.ToValue(iv)
				}
			}
			if sz == 1 {
				return wrap(results[0].Interface())
			}
			vals := make([]goja.Value, sz)
			for i := 0; i < sz; i++ {
				vals[i] = wrap(results[i].Interface())
			}
			return rt.ToValue(vals)
		})
	}

	return obj
}

// BindXORM binds an existing *xorm.Engine into a JS object exposing all its methods.
func BindXORM(rt *goja.Runtime, eng *xorm.Engine) *goja.Object {
	return BindAllMethods(rt, eng)
}

// RegisterXORM registers a global function `xorm(driver, dsn)` in the given runtime.
// It constructs a new *xorm.Engine and returns a JS object with all its methods bound.
func RegisterXORM(rt *goja.Runtime) error {
	creator := func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 2 {
			panic(rt.NewTypeError("xorm(driver, dsn) requires 2 arguments"))
		}
		driver, _ := call.Arguments[0].Export().(string)
		dsn, _ := call.Arguments[1].Export().(string)
		if driver == "" {
			panic(rt.NewTypeError("driver cannot be empty"))
		}
		eng, err := NewXORM(driver, dsn)
		if err != nil {
			panic(rt.NewGoError(err))
		}
		return BindXORM(rt, eng)
	}
	return rt.Set("xorm", creator)
}

// trailingErrorType checks if the last return type is error.
func ternary(cond bool, a, b string) string {
	if cond {
		return a
	}
	return b
}

func trailingErrorType(mt reflect.Type) (reflect.Type, bool) {
	n := mt.NumOut()
	if n == 0 {
		return nil, false
	}
	last := mt.Out(n - 1)
	if last == reflect.TypeOf((*error)(nil)).Elem() {
		return last, true
	}
	return nil, false
}

// tryConvert attempts a best-effort conversion from an exported interface{} to the target type.
// Returns the converted reflection.Value and whether the conversion succeeded.
func tryConvert(v any, t reflect.Type) (reflect.Value, bool) {
	if v == nil {
		// Return zero for a non-pointer, nil for pointers/interfaces
		if t.Kind() == reflect.Pointer || t.Kind() == reflect.Interface || t.Kind() == reflect.Slice || t.Kind() == reflect.Map || t.Kind() == reflect.Func {
			return reflect.Zero(t), true
		}
		return reflect.Zero(t), true
	}
	vv := reflect.ValueOf(v)
	if vv.Type().AssignableTo(t) {
		return vv, true
	}
	if vv.Type().ConvertibleTo(t) {
		return vv.Convert(t), true
	}
	// Simple numeric widening
	switch t.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if f, ok := v.(float64); ok {
			return reflect.ValueOf(int64(f)).Convert(t), true
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if f, ok := v.(float64); ok {
			return reflect.ValueOf(uint64(f)).Convert(t), true
		}
	case reflect.Float32, reflect.Float64:
		if i, ok := v.(int64); ok {
			return reflect.ValueOf(float64(i)).Convert(t), true
		}
		if i, ok := v.(int); ok {
			return reflect.ValueOf(float64(i)).Convert(t), true
		}
	default:
		// no supported conversion path
	}
	return reflect.Value{}, false
}

// Sentinel errors for testability

// BindXORMProxy creates a JS object that exposes all exported methods of *xorm.Engine,
// but dispatches each call to whatever current engine instance is installed via the returned setter.
// This allows swapping the underlying engine without re-binding the JS object.
func BindXORMProxy(rt *goja.Runtime) (*goja.Object, func(*xorm.Engine)) {
	var current *xorm.Engine
	setter := func(e *xorm.Engine) { current = e }

	obj := rt.NewObject()
	// We iterate methods from the type to avoid capturing a specific instance.
	typ := reflect.TypeOf((*xorm.Engine)(nil)) // *xorm.Engine
	for i := 0; i < typ.NumMethod(); i++ {
		m := typ.Method(i)
		if m.PkgPath != "" { // unexported
			continue
		}
		name := m.Name
		meth := m
		_ = obj.Set(name, func(call goja.FunctionCall) goja.Value {
			mt := meth.Type
			inTotal := mt.NumIn() // includes receiver
			isVar := mt.IsVariadic()
			fixedParams := inTotal - 1
			if isVar {
				fixedParams = inTotal - 2
			}
			if (!isVar && len(call.Arguments) != fixedParams) || (isVar && len(call.Arguments) < fixedParams) {
				err := fmt.Errorf("method %s requires %d arguments%v, got %d", name, fixedParams, ternary(isVar, "+variadic", ""), len(call.Arguments))
				panic(rt.NewTypeError(err.Error()))
			}
			if current == nil {
				panic(rt.NewTypeError("xorm engine not set for proxy"))
			}
			rev := reflect.ValueOf(current)
			args := make([]reflect.Value, 0, inTotal)
			args = append(args, rev)
			for pi := 1; pi <= fixedParams; pi++ {
				pt := mt.In(pi)
				argPtr := reflect.New(pt)
				gv := call.Arguments[pi-1]
				if err := rt.ExportTo(gv, argPtr.Interface()); err != nil {
					exported := gv.Export()
					v, ok := tryConvert(exported, pt)
					if !ok {
						panic(rt.NewTypeError("convert argument %d for %s failed: %v", pi-1, name, err))
					}
					args = append(args, v)
					continue
				}
				args = append(args, argPtr.Elem())
			}

			var results []reflect.Value
			if isVar {
				sliceT := mt.In(inTotal - 1)
				elemT := sliceT.Elem()
				nVar := len(call.Arguments) - fixedParams
				slice := reflect.MakeSlice(sliceT, 0, nVar)
				if nVar == 1 {
					if arr, ok := call.Arguments[fixedParams].Export().([]interface{}); ok {
						for _, it := range arr {
							v, ok := tryConvert(it, elemT)
							if !ok {
								panic(rt.NewTypeError("convert variadic element for %s failed", name))
							}
							slice = reflect.Append(slice, v)
						}
						args = append(args, slice)
						results = meth.Func.CallSlice(args)
					} else {
						for i := fixedParams; i < len(call.Arguments); i++ {
							elPtr := reflect.New(elemT)
							gv := call.Arguments[i]
							if err := rt.ExportTo(gv, elPtr.Interface()); err != nil {
								exported := gv.Export()
								v, ok := tryConvert(exported, elemT)
								if !ok {
									panic(rt.NewTypeError("convert variadic element for %s failed: %v", name, err))
								}
								slice = reflect.Append(slice, v)
							} else {
								slice = reflect.Append(slice, elPtr.Elem())
							}
						}
						args = append(args, slice)
						results = meth.Func.CallSlice(args)
					}
				} else {
					for i := fixedParams; i < len(call.Arguments); i++ {
						elPtr := reflect.New(elemT)
						gv := call.Arguments[i]
						if err := rt.ExportTo(gv, elPtr.Interface()); err != nil {
							exported := gv.Export()
							v, ok := tryConvert(exported, elemT)
							if !ok {
								panic(rt.NewTypeError("convert variadic element for %s failed: %v", name, err))
							}
							slice = reflect.Append(slice, v)
						} else {
							slice = reflect.Append(slice, elPtr.Elem())
						}
					}
					args = append(args, slice)
					results = meth.Func.CallSlice(args)
				}
			} else {
				results = meth.Func.Call(args)
			}
			n := len(results)
			if n == 0 {
				return goja.Undefined()
			}
			if _, hasErr := trailingErrorType(meth.Type); hasErr {
				errVal := results[n-1]
				if !errVal.IsZero() && !errVal.IsNil() {
					panic(rt.NewGoError(errVal.Interface().(error)))
				}
				results = results[:n-1]
			}
			sz := len(results)
			if sz == 0 {
				return goja.Undefined()
			}
			wrap := func(iv any) goja.Value {
				switch v := iv.(type) {
				case *xorm.Engine:
					return BindAllMethods(rt, v)
				case *xorm.Session:
					return BindAllMethods(rt, v)
				default:
					return rt.ToValue(iv)
				}
			}
			if sz == 1 {
				return wrap(results[0].Interface())
			}
			vals := make([]goja.Value, sz)
			for i := 0; i < sz; i++ {
				vals[i] = wrap(results[i].Interface())
			}
			return rt.ToValue(vals)
		})
	}
	return obj, setter
}
