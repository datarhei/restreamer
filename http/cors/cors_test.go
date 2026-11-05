package cors

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAllowAll(t *testing.T) {
	origins := []string{"*"}

	err := Validate(origins)
	require.Equal(t, nil, err)
}

func TestAllowURL(t *testing.T) {
	origins := []string{"http://example.com", "https://foobar.com"}

	err := Validate(origins)
	require.Equal(t, nil, err)
}

func TestInvalid(t *testing.T) {
	origins := []string{"rtmp://example.com"}

	err := Validate(origins)
	require.NotEqual(t, nil, err)
}
