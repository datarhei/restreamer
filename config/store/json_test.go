package store

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/datarhei/core/v16/config"

	"github.com/stretchr/testify/require"
)

func TestMigrationV1ToV3(t *testing.T) {
	jsondatav1, err := os.ReadFile("./fixtures/config_v1.json")
	require.NoError(t, err)

	jsondatav3, err := os.ReadFile("./fixtures/config_v1_v3.json")
	require.NoError(t, err)

	datav3 := config.New(nil)
	json.Unmarshal(jsondatav3, datav3)

	data, err := migrate(jsondatav1)
	require.NoError(t, err)

	datav3.Data.CreatedAt = time.Time{}
	data.CreatedAt = time.Time{}

	require.Equal(t, datav3.Data, *data)
}

func TestMigrationV2ToV3(t *testing.T) {
	jsondatav2, err := os.ReadFile("./fixtures/config_v2.json")
	require.NoError(t, err)

	jsondatav3, err := os.ReadFile("./fixtures/config_v2_v3.json")
	require.NoError(t, err)

	datav3 := config.New(nil)
	json.Unmarshal(jsondatav3, datav3)

	data, err := migrate(jsondatav2)
	require.NoError(t, err)

	datav3.Data.CreatedAt = time.Time{}
	data.CreatedAt = time.Time{}

	require.Equal(t, datav3.Data, *data)
}
