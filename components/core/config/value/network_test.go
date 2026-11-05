package value

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddressValue(t *testing.T) {
	var x string

	val := NewAddress(&x, ":8080")

	require.Equal(t, ":8080", val.String())
	require.Equal(t, nil, val.Validate())
	require.Equal(t, false, val.IsEmpty())

	x = "foobaz:9090"

	require.Equal(t, "foobaz:9090", val.String())
	require.Equal(t, nil, val.Validate())
	require.Equal(t, false, val.IsEmpty())

	val.Set("fooboz:7070")

	require.Equal(t, "fooboz:7070", x)
}

func TestCIDRListValue(t *testing.T) {
	var x []string

	val := NewCIDRList(&x, []string{}, " ")

	require.Equal(t, "(empty)", val.String())
	require.Equal(t, nil, val.Validate())
	require.Equal(t, true, val.IsEmpty())

	x = []string{"127.0.0.1/32", "127.0.0.2/32"}

	require.Equal(t, "127.0.0.1/32 127.0.0.2/32", val.String())
	require.Equal(t, nil, val.Validate())
	require.Equal(t, false, val.IsEmpty())

	val.Set("129.0.0.1/32 129.0.0.2/32")

	require.Equal(t, []string{"129.0.0.1/32", "129.0.0.2/32"}, x)
}

func TestCORSOriginaValue(t *testing.T) {
	var x []string

	val := NewCORSOrigins(&x, []string{}, " ")

	require.Equal(t, "(empty)", val.String())
	require.Equal(t, nil, val.Validate())
	require.Equal(t, true, val.IsEmpty())

	x = []string{"*"}

	require.Equal(t, "*", val.String())
	require.Equal(t, nil, val.Validate())
	require.Equal(t, false, val.IsEmpty())

	val.Set("http://localhost")

	require.Equal(t, []string{"http://localhost"}, x)
}

func TestPortValue(t *testing.T) {
	var x int

	val := NewPort(&x, 11)

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

func TestURLValue(t *testing.T) {
	var x string

	val := NewURL(&x, "http://localhost/foobar")

	require.Equal(t, "http://localhost/foobar", val.String())
	require.Equal(t, nil, val.Validate())
	require.Equal(t, false, val.IsEmpty())

	x = "http://localhost:8080/foobar"

	require.Equal(t, "http://localhost:8080/foobar", val.String())
	require.Equal(t, nil, val.Validate())
	require.Equal(t, false, val.IsEmpty())

	val.Set("http://localhost:8080/fooboz/foobaz")

	require.Equal(t, "http://localhost:8080/fooboz/foobaz", x)
}

func TestEmailValue(t *testing.T) {
	var x string

	val := NewEmail(&x, "foobar@example.com")

	require.Equal(t, "foobar@example.com", val.String())
	require.Equal(t, nil, val.Validate())
	require.Equal(t, false, val.IsEmpty())

	x = "foobar+baz@example.com"

	require.Equal(t, "foobar+baz@example.com", val.String())
	require.Equal(t, nil, val.Validate())
	require.Equal(t, false, val.IsEmpty())

	val.Set("foobar@sub.example.com")

	require.Equal(t, "foobar@sub.example.com", x)
}
