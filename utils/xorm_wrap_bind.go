package utils

import (
	"github.com/dop251/goja"
)

// BindXORMWrap binds a *XormWrap so that all its exported methods are
// accessible from the provided goja Runtime as properties on a new object.
// It reuses the generic BindAllMethods helper to provide:
//   - argument conversion
//   - variadic handling
//   - trailing error propagation
//   - wrapping of returned Go values when applicable
func BindXORMWrap(rt *goja.Runtime, wrap *XormWrap) *goja.Object {
	return BindAllMethods(rt, wrap)
}

// RegisterXORMWrap registers a global constructor function `xormWrap(driver, dsn)`
// in the given JS runtime. It creates a new *XormWrap using NewXORMWrap and
// returns a JS object exposing all methods of XormWrap.
func RegisterXORMWrap(rt *goja.Runtime) error {
	creator := func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 2 {
			panic(rt.NewTypeError("xormWrap(driver, dsn) requires 2 arguments"))
		}
		driver, _ := call.Arguments[0].Export().(string)
		dsn, _ := call.Arguments[1].Export().(string)
		if driver == "" {
			panic(rt.NewTypeError("driver cannot be empty"))
		}
		w, err := NewXORMWrap(driver, dsn)
		if err != nil {
			panic(rt.NewGoError(err))
		}
		return BindXORMWrap(rt, w)
	}
	return rt.Set("xormWrap", creator)
}
