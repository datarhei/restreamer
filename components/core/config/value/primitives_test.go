package value

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStringValue(t *testing.T) {
	var x string

	val := NewString(&x, "foobar")

	require.Equal(t, "foobar", val.String())
	require.Equal(t, nil, val.Validate())
	require.Equal(t, false, val.IsEmpty())

	x = "foobaz"

	require.Equal(t, "foobaz", val.String())
	require.Equal(t, nil, val.Validate())
	require.Equal(t, false, val.IsEmpty())

	val.Set("fooboz")

	require.Equal(t, "fooboz", x)
}

func TestStringListValue(t *testing.T) {
	var x []string

	val := NewStringList(&x, []string{"foobar"}, " ")

	require.Equal(t, "foobar", val.String())
	require.Equal(t, nil, val.Validate())
	require.Equal(t, false, val.IsEmpty())

	x = []string{"foobar", "foobaz"}

	require.Equal(t, "foobar foobaz", val.String())
	require.Equal(t, nil, val.Validate())
	require.Equal(t, false, val.IsEmpty())

	val.Set("fooboz foobar")

	require.Equal(t, []string{"fooboz", "foobar"}, x)
}

func TestStringMapStringValue(t *testing.T) {
	var x map[string]string

	val := NewStringMapString(&x, map[string]string{"a": "foobar"})

	require.Equal(t, "a:foobar", val.String())
	require.Equal(t, nil, val.Validate())
	require.Equal(t, false, val.IsEmpty())

	x = map[string]string{"a": "foobar", "b": "foobaz"}

	require.Equal(t, "a:foobar b:foobaz", val.String())
	require.Equal(t, nil, val.Validate())
	require.Equal(t, false, val.IsEmpty())

	val.Set("x:fooboz y:foobar")

	require.Equal(t, map[string]string{"x": "fooboz", "y": "foobar"}, x)
}

func TestBoolValue(t *testing.T) {
	var x bool

	val := NewBool(&x, false)

	require.Equal(t, "false", val.String())
	require.Equal(t, nil, val.Validate())
	require.Equal(t, true, val.IsEmpty())

	x = true

	require.Equal(t, "true", val.String())
	require.Equal(t, nil, val.Validate())
	require.Equal(t, false, val.IsEmpty())

	val.Set("false")

	require.Equal(t, false, x)
}

func TestIntValue(t *testing.T) {
	var x int

	val := NewInt(&x, 11)

	require.Equal(t, "11", val.String())
	require.Equal(t, nil, val.Validate())
	require.Equal(t, false, val.IsEmpty())

	x = 42

	require.Equal(t, "42", val.String())
	require.Equal(t, nil, val.Validate())
	require.Equal(t, false, val.IsEmpty())

	val.Set("77")

	require.Equal(t, int(77), x)
}

func TestInt64Value(t *testing.T) {
	var x int64

	val := NewInt64(&x, 11)

	require.Equal(t, "11", val.String())
	require.Equal(t, nil, val.Validate())
	require.Equal(t, false, val.IsEmpty())

	x = 42

	require.Equal(t, "42", val.String())
	require.Equal(t, nil, val.Validate())
	require.Equal(t, false, val.IsEmpty())

	val.Set("77")

	require.Equal(t, int64(77), x)
}

func TestUint64Value(t *testing.T) {
	var x uint64

	val := NewUint64(&x, 11)

	require.Equal(t, "11", val.String())
	require.Equal(t, nil, val.Validate())
	require.Equal(t, false, val.IsEmpty())

	x = 42

	require.Equal(t, "42", val.String())
	require.Equal(t, nil, val.Validate())
	require.Equal(t, false, val.IsEmpty())

	val.Set("77")

	require.Equal(t, uint64(77), x)
}
