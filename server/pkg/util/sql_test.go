package util

import (
	"reflect"
	"testing"
)

func TestCreatePlaceholdersForStrings(t *testing.T) {
	type args struct {
		ss []string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 []interface{}
	}{
		{
			name:  "single placeholder",
			args:  args{ss: []string{"foo"}},
			want:  "$1",
			want1: []interface{}{"foo"},
		},
		{
			name:  "multiple placeholders",
			args:  args{ss: []string{"foo", "bar", "baz"}},
			want:  "$1, $2, $3",
			want1: []interface{}{"foo", "bar", "baz"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := CreatePlaceholdersForStrings(tt.args.ss)
			if got != tt.want {
				t.Errorf("CreatePlaceholdersForStrings() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("CreatePlaceholdersForStrings() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
