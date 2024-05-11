package searchterms

import (
	"github.com/warmans/rsk-search/pkg/filter"
	"reflect"
	"testing"
)

func TestMustParse(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want []Term
	}{
		{
			name: "parse word",
			args: args{s: "foo"},
			want: []Term{{Field: "content", Value: "foo", Op: filter.CompOpFuzzyLike}},
		},
		{
			name: "parse words",
			args: args{s: "foo bar baz"},
			want: []Term{{Field: "content", Value: "foo bar baz", Op: filter.CompOpFuzzyLike}},
		},
		{
			name: "parse quoted string",
			args: args{s: `"foo bar"`},
			want: []Term{{Field: "content", Value: "foo bar", Op: filter.CompOpEq}},
		},
		{
			name: "parse quoted strings",
			args: args{s: `"foo bar" "baz"`},
			want: []Term{
				{Field: "content", Value: "foo bar", Op: filter.CompOpEq},
				{Field: "content", Value: "baz", Op: filter.CompOpEq},
			},
		},
		{
			name: "parse publication",
			args: args{s: `~xfm`},
			want: []Term{
				{Field: "publication", Value: "xfm", Op: filter.CompOpEq},
			},
		},
		{
			name: "parse mention",
			args: args{s: `@steve`},
			want: []Term{
				{Field: "actor", Value: "steve", Op: filter.CompOpEq},
			},
		},
		{
			name: "parse all",
			args: args{s: `@steve ~xfm "man alive" karl`},
			want: []Term{
				{Field: "actor", Value: "steve", Op: filter.CompOpEq},
				{Field: "publication", Value: "xfm", Op: filter.CompOpEq},
				{Field: "content", Value: "man alive", Op: filter.CompOpEq},
				{Field: "content", Value: "karl", Op: filter.CompOpFuzzyLike},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MustParse(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MustParse() = %v, want %v", got, tt.want)
			}
		})
	}
}
