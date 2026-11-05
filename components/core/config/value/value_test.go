package value

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type testdata struct {
	value1 int
	value2 int
}

func TestCopyStruct(t *testing.T) {
	data1 := testdata{}

	NewInt(&data1.value1, 1)
	NewInt(&data1.value2, 2)

	require.Equal(t, int(1), data1.value1)
	require.Equal(t, int(2), data1.value2)

	data2 := testdata{}

	val21 := NewInt(&data2.value1, 3)
	val22 := NewInt(&data2.value2, 4)

	require.Equal(t, int(3), data2.value1)
	require.Equal(t, int(4), data2.value2)

	data2 = data1

	require.Equal(t, int(1), data2.value1)
	require.Equal(t, int(2), data2.value2)

	require.Equal(t, "1", val21.String())
	require.Equal(t, "2", val22.String())
}
