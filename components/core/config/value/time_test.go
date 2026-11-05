package value

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTimeValue(t *testing.T) {
	var x time.Time

	tm := time.Unix(1257894000, 0).UTC()

	val := NewTime(&x, tm)

	require.Equal(t, "2009-11-10T23:00:00Z", val.String())
	require.Equal(t, nil, val.Validate())
	require.Equal(t, false, val.IsEmpty())

	x = time.Unix(1257894001, 0).UTC()

	require.Equal(t, "2009-11-10T23:00:01Z", val.String())
	require.Equal(t, nil, val.Validate())
	require.Equal(t, false, val.IsEmpty())

	val.Set("2009-11-11T23:00:00Z")

	require.Equal(t, time.Time(time.Date(2009, time.November, 11, 23, 0, 0, 0, time.UTC)), x)
}
