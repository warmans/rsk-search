package psql

import (
	"github.com/warmans/rsk-search/pkg/filter"
	"reflect"
	"testing"
)

func TestFilterToQuery1(t *testing.T) {

	filterMapping := map[string]string{
		"foo": "foo",
		"bar": "b.bar",
		"baz": "baz",
	}

	tests := []struct {
		f          filter.Filter
		wantSQL    string
		wantParams []interface{}
		wantErr    bool
	}{
		// comp filters
		{
			f:          filter.Eq("foo", filter.String("bar")),
			wantSQL:    "foo = $1",
			wantParams: []interface{}{"bar"},
			wantErr:    false,
		},
		{
			f:          filter.Eq("bar", filter.String("bar")),
			wantSQL:    "b.bar = $1",
			wantParams: []interface{}{"bar"},
			wantErr:    false,
		},
		{
			f:          filter.Eq("foo", filter.Null()),
			wantSQL:    "foo IS NULL",
			wantParams: []interface{}{},
			wantErr:    false,
		},
		{
			f:          filter.Neq("foo", filter.Null()),
			wantSQL:    "foo IS NOT NULL",
			wantParams: []interface{}{},
			wantErr:    false,
		},
		{
			f:          filter.Neq("foo", filter.String("baz")),
			wantSQL:    "foo != $1",
			wantParams: []interface{}{"baz"},
			wantErr:    false,
		},
		{
			f:          filter.Gt("foo", filter.Int(10)),
			wantSQL:    "foo > $1",
			wantParams: []interface{}{int64(10)},
			wantErr:    false,
		},
		{
			f:          filter.Ge("foo", filter.Int(10)),
			wantSQL:    "foo >= $1",
			wantParams: []interface{}{int64(10)},
			wantErr:    false,
		},
		{
			f:          filter.Lt("foo", filter.Int(10)),
			wantSQL:    "foo < $1",
			wantParams: []interface{}{int64(10)},
			wantErr:    false,
		},
		{
			f:          filter.Le("foo", filter.Int(10)),
			wantSQL:    "foo <= $1",
			wantParams: []interface{}{int64(10)},
			wantErr:    false,
		},
		{
			f:          filter.Le("foo", filter.Int(10)),
			wantSQL:    "foo <= $1",
			wantParams: []interface{}{int64(10)},
			wantErr:    false,
		},
		{
			f:          filter.Like("foo", filter.String("bar")),
			wantSQL:    "foo LIKE $1",
			wantParams: []interface{}{"%bar%"},
			wantErr:    false,
		},
		// bool filters
		{
			f:          filter.And(filter.Eq("foo", filter.String("dog")), filter.Eq("baz", filter.String("cat"))),
			wantSQL:    "(foo = $1) and (baz = $2)",
			wantParams: []interface{}{"dog", "cat"},
			wantErr:    false,
		},
		{
			f: filter.And(
				filter.And(filter.Eq("foo", filter.String("dog")), filter.Eq("baz", filter.String("cat"))),
				filter.Eq("foo", filter.String("bar")),
			),
			wantSQL:    "((foo = $1) and (baz = $2)) and (foo = $3)",
			wantParams: []interface{}{"dog", "cat", "bar"},
			wantErr:    false,
		},
		{
			f: filter.Or(
				filter.And(filter.Eq("foo", filter.String("dog")), filter.Eq("baz", filter.String("cat"))),
				filter.And(filter.Eq("foo", filter.String("dog")), filter.Eq("baz", filter.String("cat"))),
			),
			wantSQL:    "((foo = $1) and (baz = $2)) or ((foo = $3) and (baz = $4))",
			wantParams: []interface{}{"dog", "cat", "dog", "cat"},
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(filter.MustPrint(tt.f), func(t *testing.T) {
			got, got1, err := FilterToQuery(tt.f, filterMapping)
			if (err != nil) != tt.wantErr {
				t.Errorf("FilterToQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.wantSQL {
				t.Errorf("FilterToQuery() got = %v, want %v", got, tt.wantSQL)
			}
			if !reflect.DeepEqual(got1, tt.wantParams) {
				t.Errorf("FilterToQuery() got1 = %v, want %v", got1, tt.wantParams)
			}
		})
	}
}
