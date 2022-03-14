package filter

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParseCompFilter(t *testing.T) {

	tests := map[string]struct {
		expectFilter Filter
		expectError  bool
	}{
		"": {
			expectFilter: nil,
		},
		// types
		`foo = true`: {
			expectFilter: Eq("foo", Bool(true)),
		},
		`foo = "bar"`: {
			expectFilter: Eq("foo", String("bar")),
		},
		`foo = 1`: {
			expectFilter: Eq("foo", Int(1)),
		},
		`foo = 1.5`: {
			expectFilter: Eq("foo", Float(1.5)),
		},
		`foo = null`: {
			expectFilter: Eq("foo", Null()),
		},
		// comparison operators
		`foo > 1`: {
			expectFilter: Gt("foo", Int(1)),
		},
		`foo >= 1`: {
			expectFilter: Ge("foo", Int(1)),
		},
		`foo < 1`: {
			expectFilter: Lt("foo", Int(1)),
		},
		`foo <= 1`: {
			expectFilter: Le("foo", Int(1)),
		},
		`foo != "bar"`: {
			expectFilter: Neq("foo", String("bar")),
		},
		`foo ~= "bar"`: {
			expectFilter: Like("foo", String("bar")),
		},
		`foo ~ "bar"`: {
			expectFilter: FuzzyLike("foo", String("bar")),
		},
	}
	for condition, test := range tests {
		t.Run(condition, func(t *testing.T) {
			result, err := Parse(condition)
			if test.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.EqualValues(t, test.expectFilter, result)
		})
	}
}

func TestParseBoolFilter(t *testing.T) {

	tests := map[string]struct {
		expectFilter Filter
		expectError  bool
	}{
		"": {
			expectFilter: nil,
		},
		`(foo = true) and (bar = false)`: {
			expectFilter: And(Eq("foo", Bool(true)), Eq("bar", Bool(false))),
		},
		`(foo = true) or (bar = false)`: {
			expectFilter: Or(Eq("foo", Bool(true)), Eq("bar", Bool(false))),
		},
		`((foo = true) or (bar = false)) and baz = 1.0`: {
			expectFilter: And(Or(Eq("foo", Bool(true)), Eq("bar", Bool(false))), Eq("baz", Float(1.0))),
		},
	}
	for condition, test := range tests {
		t.Run(condition, func(t *testing.T) {
			result, err := Parse(condition)
			if test.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.EqualValues(t, test.expectFilter, result)
		})
	}
}
