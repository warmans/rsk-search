package filter

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPrint(t *testing.T) {
	tests := []struct {
		filter       Filter
		expectString string
		expectError  bool
	}{
		{
			filter:       Eq("foo", String("bar")),
			expectString: `foo = "bar"`,
			expectError:  false,
		}, {
			filter:       And(Eq("foo", String("bar")), Gt("bar", Int(1))),
			expectString: `foo = "bar" and bar > 1`,
			expectError:  false,
		}, {
			filter:       Or(And(Eq("foo", String("bar")), Gt("bar", Int(1))), Neq("baz", Int(2))),
			expectString: `foo = "bar" and bar > 1 or baz != 2`,
			expectError:  false,
		}, {
			filter:       And(Eq("foo", String("bar")), Or(Gt("bar", Int(1)), Neq("baz", Int(2)))),
			expectString: `foo = "bar" and (bar > 1 or baz != 2)`,
			expectError:  false,
		},
	}

	for k, test := range tests {
		t.Run(fmt.Sprintf("test %d", k), func(t *testing.T) {
			s, err := Print(test.filter)
			if test.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.EqualValues(t, test.expectString, s)
		})
	}
}

func TestParsePrinted(t *testing.T) {
	in := And(Eq("foo", String("bar")), Or(Gt("bar", Int(1)), Neq("baz", Int(2))))
	out := MustParse(MustPrint(in))
	require.EqualValues(t, in, out)
}
