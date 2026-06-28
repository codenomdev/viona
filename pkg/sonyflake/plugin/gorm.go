package plugin

import (
	"reflect"

	"github.com/codenomdev/viona/pkg/sonyflake"
	"gorm.io/gorm"
)

// Plugin struct
type SonyflakePlugin struct{}

// Register to GORM
func (SonyflakePlugin) Name() string {
	return "sonyflake_plugin"
}

// Initialize
func (SonyflakePlugin) Initialize(db *gorm.DB) error {
	db.Callback().Create().Before("gorm:before_create").Register("sonyflake:before_create", func(tx *gorm.DB) {
		stmt := tx.Statement
		if stmt.Schema == nil {
			return
		}

		field := stmt.Schema.LookUpField("ID")
		if field == nil {
			return
		}

		if field.FieldType.Kind() != reflect.Int64 {
			return
		}

		rv := stmt.ReflectValue

		switch rv.Kind() {

		// BULK INSERT (slice)
		case reflect.Slice, reflect.Array:
			for i := 0; i < rv.Len(); i++ {
				elem := rv.Index(i)

				if elem.Kind() == reflect.Ptr {
					elem = elem.Elem()
				}

				if elem.Kind() != reflect.Struct {
					continue
				}

				val, isZero := field.ValueOf(tx.Statement.Context, elem)
				if !isZero && !reflect.ValueOf(val).IsZero() {
					continue
				}

				newID := sonyflake.Generate(tx.Statement.Context)
				_ = field.Set(tx.Statement.Context, elem, newID)
			}

		// SINGLE INSERT
		case reflect.Struct:
			val, isZero := field.ValueOf(tx.Statement.Context, rv)
			if !isZero && !reflect.ValueOf(val).IsZero() {
				return
			}

			newID := sonyflake.Generate(tx.Statement.Context)
			_ = field.Set(tx.Statement.Context, rv, newID)

		// POINTER
		case reflect.Ptr:
			elem := rv.Elem()
			if elem.Kind() != reflect.Struct {
				return
			}

			val, isZero := field.ValueOf(tx.Statement.Context, elem)
			if !isZero && !reflect.ValueOf(val).IsZero() {
				return
			}

			newID := sonyflake.Generate(tx.Statement.Context)
			_ = field.Set(tx.Statement.Context, elem, newID)
		}
	})

	return nil
}
