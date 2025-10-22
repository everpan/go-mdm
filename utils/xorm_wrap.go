package utils

import (
	"database/sql"

	"xorm.io/xorm"
)

type XormWrap struct {
	engine *xorm.Engine
}

func (wrap *XormWrap) SetEngine(engine *xorm.Engine) {
	wrap.engine = engine
}

// NewXORM creates a new xorm Engine using the given driver and DSN.
// It is a light wrapper that keeps the utils package self-contained for tests.
func NewXORM(driver, dsn string) (*xorm.Engine, error) {
	return xorm.NewEngine(driver, dsn)
}

// NewXORMWrap creates a new XormWrap Engine using the given driver and DSN.
func NewXORMWrap(driver, dsn string) (*XormWrap, error) {
	eng, err := NewXORM(driver, dsn)
	if err != nil {
		return nil, err
	}
	return &XormWrap{
		engine: eng,
	}, nil
}

// Exec raw sql
func (wrap *XormWrap) Exec(sqlOrArgs ...interface{}) (sql.Result, error) {
	return wrap.engine.Exec(sqlOrArgs...)
}

// Query a raw SQL and return records as []map[string][]byte
func (wrap *XormWrap) Query(sqlOrArgs ...interface{}) (resultsSlice []map[string][]byte, err error) {
	return wrap.engine.Query(sqlOrArgs...)
}

// QueryString runs a raw SQL and return records as []map[string]string
func (wrap *XormWrap) QueryString(sqlOrArgs ...interface{}) ([]map[string]string, error) {
	return wrap.engine.QueryString(sqlOrArgs...)
}

// QueryInterface runs a raw SQL and return records as []map[string]interface{}
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

// Get retrieve one record from the table, bean's non-empty fields
// are conditions
func (wrap *XormWrap) Get(beans ...interface{}) (bool, error) {
	return wrap.engine.Get(beans...)
}

// Exist returns true if the record exist otherwise return false
func (wrap *XormWrap) Exist(bean ...interface{}) (bool, error) {
	return wrap.engine.Exist(bean...)
}

// Find retrieve records from the table, condiBeans's non-empty fields
// are conditions. Beans could be []Struct, []*Struct, map[int64]Struct
// map[int64]*Struct
func (wrap *XormWrap) Find(beans interface{}, condiBeans ...interface{}) error {
	return wrap.engine.Find(beans, condiBeans...)
}

// FindAndCount find the results and also return the counts
func (wrap *XormWrap) FindAndCount(rowsSlicePtr interface{}, condiBean ...interface{}) (int64, error) {
	return wrap.engine.FindAndCount(rowsSlicePtr, condiBean...)
}

// Rows return sql.Rows compatible Rows obj, as a forward Iterator object for iterating record by record, bean's non-empty fields
// are conditions.
func (wrap *XormWrap) Rows(bean interface{}) (*xorm.Rows, error) {
	return wrap.engine.Rows(bean)
}

// Count counts the records. Bean's non-empty fields are conditions.
func (wrap *XormWrap) Count(bean ...interface{}) (int64, error) {
	return wrap.engine.Count(bean...)
}

// Sum the records by some column. Bean's non-empty fields are conditions.
func (wrap *XormWrap) Sum(bean interface{}, colName string) (float64, error) {
	return wrap.engine.Sum(bean, colName)
}

// SumInt sum the records by some column. Bean's non-empty fields are conditions.
func (wrap *XormWrap) SumInt(bean interface{}, colName string) (int64, error) {
	return wrap.engine.SumInt(bean, colName)
}

// Sums sum the records by some columns. Bean's non-empty fields are conditions.
func (wrap *XormWrap) Sums(bean interface{}, colNames ...string) ([]float64, error) {
	return wrap.engine.Sums(bean, colNames...)
}

// SumsInt like Sums but return slice of int64 instead of float64.
func (wrap *XormWrap) SumsInt(bean interface{}, colNames ...string) ([]int64, error) {
	return wrap.engine.SumsInt(bean, colNames...)
}
