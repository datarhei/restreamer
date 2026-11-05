package value

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestS3Value(t *testing.T) {
	filesystems := []S3Storage{}

	v := NewS3StorageListValue(&filesystems, nil, " ")
	require.Equal(t, "(empty)", v.String())

	v.Set("https://access_key_id1:secret_access_id1@region1.subdomain.endpoint1.com/bucket1?name=aaa1&mount=/abc1&username=xxx1&password=yyy1 http://access_key_id2:secret_access_id2@region2.endpoint2.com/bucket2?name=aaa2&mount=/abc2&username=xxx2&password=yyy2")
	require.Equal(t, []S3Storage{
		{
			Name:       "aaa1",
			Mountpoint: "/abc1",
			Auth: S3StorageAuth{
				Enable:   true,
				Username: "xxx1",
				Password: "yyy1",
			},
			Endpoint:        "subdomain.endpoint1.com",
			AccessKeyID:     "access_key_id1",
			SecretAccessKey: "secret_access_id1",
			Bucket:          "bucket1",
			Region:          "region1",
			UseSSL:          true,
		},
		{
			Name:       "aaa2",
			Mountpoint: "/abc2",
			Auth: S3StorageAuth{
				Enable:   true,
				Username: "xxx2",
				Password: "yyy2",
			},
			Endpoint:        "endpoint2.com",
			AccessKeyID:     "access_key_id2",
			SecretAccessKey: "secret_access_id2",
			Bucket:          "bucket2",
			Region:          "region2",
			UseSSL:          false,
		},
	}, filesystems)
	require.Equal(t, "https://access_key_id1:---@region1.subdomain.endpoint1.com/bucket1?mount=%2Fabc1&name=aaa1&password=---&username=xxx1 http://access_key_id2:---@region2.endpoint2.com/bucket2?mount=%2Fabc2&name=aaa2&password=---&username=xxx2", v.String())
	require.NoError(t, v.Validate())

	v.Set("https://access_key_id1:secret_access_id1@region1.endpoint1.com/bucket1?name=djk*;..&mount=/abc1&username=xxx1&password=yyy1")
	require.Error(t, v.Validate())
}
