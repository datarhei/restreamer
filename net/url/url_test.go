package url

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLookup(t *testing.T) {
	ip, err := Lookup("/localhost:8080/foobar")

	require.NoError(t, err)
	require.Equal(t, "", ip)

	ip, err = Lookup("http://")

	require.NoError(t, err)
	require.Equal(t, "", ip)

	ip, err = Lookup("https://www.google.com")

	require.NoError(t, err)
	require.NotEmpty(t, ip)
}

func TestLocalhost(t *testing.T) {
	ip, err := Lookup("http://localhost:8080/foobar")

	require.NoError(t, err)
	require.Subset(t, []string{"127.0.0.1", "::1"}, []string{ip})
}

func TestValidate(t *testing.T) {
	err := Validate("http://localhost/foobar")
	require.NoError(t, err)

	err = Validate("foobar")
	require.NoError(t, err)

	err = Validate("http://localhost/foobar_%25v")
	require.NoError(t, err)

	err = Validate("http://localhost/foobar_%v")
	require.NoError(t, err)
}

func TestScheme(t *testing.T) {
	r := HasScheme("http://localhost/foobar")
	require.True(t, r)

	r = HasScheme("iueriherfd://localhost/foobar")
	require.True(t, r)

	r = HasScheme("//localhost/foobar")
	require.False(t, r)
}

func TestParse(t *testing.T) {
	u, err := Parse("http://localhost/foobar")
	require.NoError(t, err)
	require.Equal(t, &URL{
		Scheme:      "http",
		Opaque:      "",
		User:        nil,
		Host:        "localhost",
		RawPath:     "/foobar",
		RawQuery:    "",
		RawFragment: "",
	}, u)

	u, err = Parse("iueriherfd://localhost/foobar")
	require.NoError(t, err)
	require.Equal(t, &URL{
		Scheme:      "iueriherfd",
		Opaque:      "",
		User:        nil,
		Host:        "localhost",
		RawPath:     "/foobar",
		RawQuery:    "",
		RawFragment: "",
	}, u)

	u, err = Parse("//localhost/foobar")
	require.NoError(t, err)
	require.Equal(t, &URL{
		Scheme:      "",
		Opaque:      "",
		User:        nil,
		Host:        "localhost",
		RawPath:     "/foobar",
		RawQuery:    "",
		RawFragment: "",
	}, u)

	u, err = Parse("http://localhost/foobar_%v?foo=bar#foobar")
	require.NoError(t, err)
	require.Equal(t, &URL{
		Scheme:      "http",
		Opaque:      "",
		User:        nil,
		Host:        "localhost",
		RawPath:     "/foobar_%v",
		RawQuery:    "foo=bar",
		RawFragment: "foobar",
	}, u)

	u, err = Parse("http:localhost/foobar_%v?foo=bar#foobar")
	require.NoError(t, err)
	require.Equal(t, &URL{
		Scheme:      "http",
		Opaque:      "localhost/foobar_%v",
		User:        nil,
		Host:        "",
		RawPath:     "",
		RawQuery:    "foo=bar",
		RawFragment: "foobar",
	}, u)

	u, err = Parse("http:/localhost/foobar_%v?foo=bar#foobar")
	require.NoError(t, err)
	require.Equal(t, &URL{
		Scheme:      "http",
		Opaque:      "",
		User:        nil,
		Host:        "",
		RawPath:     "/localhost/foobar_%v",
		RawQuery:    "foo=bar",
		RawFragment: "foobar",
	}, u)

	u, err = Parse("http:///localhost/foobar_%v?foo=bar#foobar")
	require.NoError(t, err)
	require.Equal(t, &URL{
		Scheme:      "http",
		Opaque:      "",
		User:        nil,
		Host:        "",
		RawPath:     "/localhost/foobar_%v",
		RawQuery:    "foo=bar",
		RawFragment: "foobar",
	}, u)

	u, err = Parse("foo:bar://localhost/foobar_%v?foo=bar#foobar")
	require.NoError(t, err)
	require.Equal(t, &URL{
		Scheme:      "foo:bar",
		Opaque:      "",
		User:        nil,
		Host:        "localhost",
		RawPath:     "/foobar_%v",
		RawQuery:    "foo=bar",
		RawFragment: "foobar",
	}, u)

	u, err = Parse("http://localhost:8080/foobar")
	require.NoError(t, err)
	require.Equal(t, &URL{
		Scheme:      "http",
		Opaque:      "",
		User:        nil,
		Host:        "localhost:8080",
		RawPath:     "/foobar",
		RawQuery:    "",
		RawFragment: "",
	}, u)
	require.Equal(t, "localhost", u.Hostname())
	require.Equal(t, "8080", u.Port())

	u, err = Parse("https://www.google.com")
	require.NoError(t, err)
	require.Equal(t, &URL{
		Scheme:      "https",
		Opaque:      "",
		User:        nil,
		Host:        "www.google.com",
		RawPath:     "/",
		RawQuery:    "",
		RawFragment: "",
	}, u)
	require.Equal(t, "www.google.com", u.Hostname())
	require.Equal(t, "", u.Port())
}
