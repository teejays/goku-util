package types

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/teejays/goku-util/naam"
)

type SampleField int

const (
	SampleField_INVALID SampleField = 0
	SampleField_1       SampleField = 1
	SampleField_2       SampleField = 2
	SampleField_3       SampleField = 3
)

func (f SampleField) String() string {
	switch f {
	case SampleField_INVALID:
		return "INVALID"
	case SampleField_1:
		return "1"
	case SampleField_2:
		return "2"
	case SampleField_3:
		return "3"
	default:
		panic(fmt.Sprintf("'%d' is not a valid type '%s'", f, "SampleField"))
	}
}

func (f SampleField) Name() naam.Name {
	switch f {
	case SampleField_INVALID:
		return naam.New("INVALID")
	case SampleField_1:
		return naam.New("1")
	case SampleField_2:
		return naam.New("2")
	case SampleField_3:
		return naam.New("3")
	default:
		panic(fmt.Sprintf("'%d' is not a valid type '%s'", f, "SampleField"))
	}
}

func (f SampleField) ToDatabaseColumn() string {
	switch f {
	case SampleField_INVALID:
		return "invalid"
	case SampleField_1:
		return "one"
	case SampleField_2:
		return "two"
	case SampleField_3:
		return "three"
	default:
		panic(fmt.Sprintf("'%d' is not a valid type '%s'", f, "SampleField"))
	}
}

// // Value implements them the `drive.Valuer` interface for this enum. It allows us to save these enum values to the DB as a string.
// func (f SampleField) Value() (driver.Value, error) {
// 	switch f {
// 	case SampleField_INVALID:
// 		return nil, nil
// 	case SampleField_1:
// 		return "1", nil
// 	case SampleField_2:
// 		return "2", nil
// 	case SampleField_3:
// 		return "3", nil

// 	default:
// 		return nil, fmt.Errorf("Cannot save enum SampleField to DB: '%d' is not a valid value for enum SampleField", f)
// 	}
// }

func TestPruneFields(t *testing.T) {
	tests := []struct {
		name          string
		columns       []Field
		includeFields []Field
		excludeFields []Field
		want          []Field
	}{
		{
			name:          "No fields provided",
			columns:       []Field{},
			includeFields: []Field{},
			excludeFields: []Field{},
			want:          []Field{},
		},
		{
			name:          "All cols with no include or exclude fields - all should be used",
			columns:       []Field{SampleField_1, SampleField_2, SampleField_3},
			includeFields: []Field{},
			excludeFields: []Field{},
			want:          []Field{SampleField_1, SampleField_2, SampleField_3},
		},
		{
			name:          "All cols with one include field and no exclude fields - one should be used",
			columns:       []Field{SampleField_1, SampleField_2, SampleField_3},
			includeFields: []Field{SampleField_2},
			excludeFields: []Field{},
			want:          []Field{SampleField_2},
		},
		{
			name:          "All cols with no include field and one exclude fields - all but one should be used",
			columns:       []Field{SampleField_1, SampleField_2, SampleField_3},
			includeFields: []Field{},
			excludeFields: []Field{SampleField_2},
			want:          []Field{SampleField_1, SampleField_3},
		},
		{
			name:          "All cols with two include field and one exclude fields - two but one should be used",
			columns:       []Field{SampleField_1, SampleField_2, SampleField_3},
			includeFields: []Field{SampleField_1, SampleField_2},
			excludeFields: []Field{SampleField_2},
			want:          []Field{SampleField_1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PruneFields(tt.columns, tt.includeFields, tt.excludeFields)
			assert.ElementsMatch(t, tt.want, got)
		})
	}
}
