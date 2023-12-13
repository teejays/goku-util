package types

import (
	"database/sql/driver"
	"time"

	"github.com/teejays/goku-util/naam"
	"github.com/teejays/goku-util/scalars"
)

type BasicTypeMeta[T BasicType, F Field] interface {
	GetBasicTypeMetaBase() BasicTypeMetaBase[T, F]
	SetMetaFieldValues(T, time.Time) T

	ConvertTimestampColumnsToUTC(T) T
	SetDefaultFieldValues(T) T
}

type BasicTypeMetaBase[T BasicType, F Field] struct {
	Name   naam.Name
	Fields []F
}

type BasicType interface {
	GetID() scalars.ID
	GetUpdatedAt() scalars.Time
	SetUpdatedAt(scalars.Time)
}

type FilterType interface{}

type EntityMetaBase[T BasicType, F Field] struct {
	DbTableName   naam.Name
	BasicTypeMeta BasicTypeMeta[T, F]
}

type EntityMeta[T BasicType, F Field] interface{}

type Enum interface {
	String() string
	Name() naam.Name

	Value() (driver.Value, error)
	Scan(src interface{}) error

	MarshalJSON() ([]byte, error)
	UnmarshalJSON([]byte) error

	ImplementsGraphQLType(name string) bool
	UnmarshalGraphQL(input interface{}) error
}

type Field interface {
	String() string
	Name() naam.Name
}

// Helper function

// PruneFields syncs the list of fields with the list of allowed and excluded fields
func PruneFields[T Field](columns []T, includeFields []T, excludeFields []T) []T {
	var includeFilteredColumns []T

	// If include fields is provided, add those fields
	if len(includeFields) > 0 {
		for _, col := range columns {
			if IsFieldInFields(col, includeFields) {
				includeFilteredColumns = append(includeFilteredColumns, col)
			}
		}
	}

	// If no include fields provided, assume everything is halal
	if len(includeFields) < 1 {
		includeFilteredColumns = columns
	}

	var excludedFilteredColumns []T

	// If exclude fields is provided, remove those fields
	for _, col := range includeFilteredColumns {
		if !IsFieldInFields(col, excludeFields) {
			excludedFilteredColumns = append(excludedFilteredColumns, col)
		}
	}

	return excludedFilteredColumns

}

func IsFieldInFields[T Field](column T, fields []T) bool {
	for _, fld := range fields {
		if fld.Name().Equal(column.Name()) {
			return true
		}
	}
	return false
}

func RemoveFieldFromFields[T Field](column T, fields []T) []T {
	for i, fld := range fields {
		if fld.Name().Equal(column.Name()) {
			fields = append(fields[:i], fields[i+1:]...)
		}
	}
	return fields
}
