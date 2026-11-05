package ffmpeg

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidatorNew(t *testing.T) {
	var err error

	_, err = NewValidator([]string{}, []string{})
	require.NoError(t, err)

	_, err = NewValidator([]string{"^https?://", "^rtmp://", ""}, []string{})
	require.NoError(t, err)

	_, err = NewValidator([]string{}, []string{"^https?://", "^rtmp://", ""})
	require.NoError(t, err)
}

func TestValidatorError(t *testing.T) {
	var err error

	_, err = NewValidator([]string{}, []string{})
	require.NoError(t, err)

	_, err = NewValidator([]string{"^(foobar"}, []string{})
	require.Error(t, err, "should not accept invalid expression")

	_, err = NewValidator([]string{}, []string{"^(foobar"})
	require.Error(t, err, "should not accept invalid expression")
}

func TestValidatorNoExpressions(t *testing.T) {
	valider, _ := NewValidator([]string{}, []string{})

	require.Equal(t, true, valider.IsValid("foobar"), "Expression should be allowed")
}

func TestValidatorAllowBlock(t *testing.T) {
	valider, _ := NewValidator([]string{"^rtmps?://"}, []string{"^https?://"})

	require.Equal(t, true, valider.IsValid("rtmp://"), "Expression should be allowed")
	if valider.IsValid("rtmp://") == false {
		t.Errorf("Allowed expression should be allowed")
	}

	require.Equal(t, false, valider.IsValid("http://"), "Unallowed expression should be blocked")
}
