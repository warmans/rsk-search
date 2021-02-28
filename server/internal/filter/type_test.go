package filter

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestType_Equal(t *testing.T) {
	require.True(t, BoolType.Equal(BoolType))
	require.True(t, StringType.Equal(StringType))
	require.True(t, BoolType.Equal(BoolType))
	require.True(t, IntType.Equal(IntType))
	require.True(t, FloatType.Equal(FloatType))

	require.False(t, BoolType.Equal(FloatType))
}

func TestType_String(t *testing.T) {
	require.EqualValues(t, "float", FloatType.String())
	require.EqualValues(t, "int", IntType.String())
	require.EqualValues(t, "bool", BoolType.String())
}
