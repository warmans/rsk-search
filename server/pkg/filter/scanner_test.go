package filter

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestScan(t *testing.T) {

	tests := map[string]struct {
		expectTokens []token
		expectError  bool
	}{
		"()": {
			expectTokens: []token{
				{tag: tagLParen, lexeme: "("},
				{tag: tagRParen, lexeme: ")"},
				{tag: tagEOF},
			},
		},
		// bool comparisons
		"foo = true": {
			expectTokens: []token{
				{tag: tagField, lexeme: "foo"},
				{tag: tagEq, lexeme: "="},
				{tag: tagBool, lexeme: "true"},
				{tag: tagEOF},
			},
		},
		"foo != false": {
			expectTokens: []token{
				{tag: tagField, lexeme: "foo"},
				{tag: tagNeq, lexeme: "!="},
				{tag: tagBool, lexeme: "false"},
				{tag: tagEOF},
			},
		},
		// int comparisons
		"foo = 1": {
			expectTokens: []token{
				{tag: tagField, lexeme: "foo"},
				{tag: tagEq, lexeme: "="},
				{tag: tagInt, lexeme: "1"},
				{tag: tagEOF},
			},
		},
		"foo > 1": {
			expectTokens: []token{
				{tag: tagField, lexeme: "foo"},
				{tag: tagGt, lexeme: ">"},
				{tag: tagInt, lexeme: "1"},
				{tag: tagEOF},
			},
		},
		"foo >= 1": {
			expectTokens: []token{
				{tag: tagField, lexeme: "foo"},
				{tag: tagGe, lexeme: ">="},
				{tag: tagInt, lexeme: "1"},
				{tag: tagEOF},
			},
		},
		"foo < 1": {
			expectTokens: []token{
				{tag: tagField, lexeme: "foo"},
				{tag: tagLt, lexeme: "<"},
				{tag: tagInt, lexeme: "1"},
				{tag: tagEOF},
			},
		},
		"foo <= 1": {
			expectTokens: []token{
				{tag: tagField, lexeme: "foo"},
				{tag: tagLe, lexeme: "<="},
				{tag: tagInt, lexeme: "1"},
				{tag: tagEOF},
			},
		},
		// floats
		"foo = 1.2": {
			expectTokens: []token{
				{tag: tagField, lexeme: "foo"},
				{tag: tagEq, lexeme: "="},
				{tag: tagFloat, lexeme: "1.2"},
				{tag: tagEOF},
			},
		},
		"foo = null": {
			expectTokens: []token{
				{tag: tagField, lexeme: "foo"},
				{tag: tagEq, lexeme: "="},
				{tag: tagNull, lexeme: "null"},
				{tag: tagEOF},
			},
		},
		// and/or
		"foo = 1 and bar = 2": {
			expectTokens: []token{
				{tag: tagField, lexeme: "foo"},
				{tag: tagEq, lexeme: "="},
				{tag: tagInt, lexeme: "1"},
				{tag: tagAnd, lexeme: "and"},
				{tag: tagField, lexeme: "bar"},
				{tag: tagEq, lexeme: "="},
				{tag: tagInt, lexeme: "2"},
				{tag: tagEOF},
			},
		},
		"foo = 1 or bar = 2": {
			expectTokens: []token{
				{tag: tagField, lexeme: "foo"},
				{tag: tagEq, lexeme: "="},
				{tag: tagInt, lexeme: "1"},
				{tag: tagOr, lexeme: "or"},
				{tag: tagField, lexeme: "bar"},
				{tag: tagEq, lexeme: "="},
				{tag: tagInt, lexeme: "2"},
				{tag: tagEOF},
			},
		},
		// groups
		"(foo = 1) or bar = 2": {
			expectTokens: []token{
				{tag: tagLParen, lexeme: "("},
				{tag: tagField, lexeme: "foo"},
				{tag: tagEq, lexeme: "="},
				{tag: tagInt, lexeme: "1"},
				{tag: tagRParen, lexeme: ")"},
				{tag: tagOr, lexeme: "or"},
				{tag: tagField, lexeme: "bar"},
				{tag: tagEq, lexeme: "="},
				{tag: tagInt, lexeme: "2"},
				{tag: tagEOF},
			},
		},

		// string
		`foo = "bar"`: {
			expectTokens: []token{
				{tag: tagField, lexeme: "foo"},
				{tag: tagEq, lexeme: "="},
				{tag: tagString, lexeme: "bar"},
				{tag: tagEOF},
			},
		},
		`foo ~ "bar"`: {
			expectTokens: []token{
				{tag: tagField, lexeme: "foo"},
				{tag: tagFuzzy, lexeme: "~"},
				{tag: tagString, lexeme: "bar"},
				{tag: tagEOF},
			},
		},
		`foo ~= "bar"`: {
			expectTokens: []token{
				{tag: tagField, lexeme: "foo"},
				{tag: tagLike, lexeme: "~="},
				{tag: tagString, lexeme: "bar"},
				{tag: tagEOF},
			},
		},
		`foo = "bar`: {
			expectError: true, // unclosed quote
		},

	}

	for str, test := range tests {
		t.Run(str, func(t *testing.T) {
			res, err := Scan(str)
			if test.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.EqualValues(t, test.expectTokens, res)
		})
	}
}
