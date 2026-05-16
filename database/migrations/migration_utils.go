package migrations

import (
	"database/sql"
	"reflect"
	"unsafe"
)

// DBFromTx extracts *sql.DB from *sql.Tx using reflection
func DBFromTx(tx *sql.Tx) *sql.DB {
	val := reflect.ValueOf(tx).Elem()
	field := val.FieldByName("db")
	// unsafe pointer trick to read unexported field
	field = reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem()
	return field.Interface().(*sql.DB)
}
