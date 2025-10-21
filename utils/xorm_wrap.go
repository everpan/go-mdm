package utils

import (
	"database/sql"

	"xorm.io/xorm"
)

type XormWrap struct {
	engine *xorm.Engine
}

// Exec raw sql
func (wrap *XormWrap) Exec(sqlOrArgs ...interface{}) (sql.Result, error) {
	return wrap.engine.Exec(sqlOrArgs...)
}

// Query a raw SQL and return records as []map[string][]byte
func (wrap *XormWrap) Query(sqlOrArgs ...interface{}) (resultsSlice []map[string][]byte, err error) {
	return wrap.engine.Query(sqlOrArgs...)
}

// QueryString runs a raw sql and return records as []map[string]string
func (wrap *XormWrap) QueryString(sqlOrArgs ...interface{}) ([]map[string]string, error) {
	return wrap.engine.QueryString(sqlOrArgs...)
}

// QueryInterface runs a raw sql and return records as []map[string]interface{}
func (wrap *XormWrap) QueryInterface(sqlOrArgs ...interface{}) ([]map[string]interface{}, error) {
	return wrap.engine.QueryInterface(sqlOrArgs...)
}

// Insert one or more records
func (wrap *XormWrap) Insert(beans ...interface{}) (int64, error) {
	return wrap.engine.Insert(beans...)
}

// InsertOne insert only one record
func (wrap *XormWrap) InsertOne(bean interface{}) (int64, error) {
	return wrap.engine.Insert(bean)
}

// Update records, bean's non-empty fields are updated contents
func (wrap *XormWrap) Update(bean interface{}, condiBeans ...interface{}) (int64, error) {
	return wrap.engine.Update(bean, condiBeans...)
}

// Delete records, bean's non-empty fields are conditions
// At least one condition must be set.
func (wrap *XormWrap) Delete(beans ...interface{}) (int64, error) {
	return wrap.engine.Delete(beans...)
}

// Truncate records, bean's non-empty fields are conditions
// In contrast to Delete, this method allows deletes without conditions.
func (wrap *XormWrap) Truncate(beans ...interface{}) (int64, error) {
	return wrap.engine.Truncate(beans...)
}
