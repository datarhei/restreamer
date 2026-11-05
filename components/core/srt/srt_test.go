package srt

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseStreamId(t *testing.T) {
	streamids := map[string]streamInfo{
		"bla":                                 {resource: "bla", mode: "request"},
		"bla,mode:publish":                    {resource: "bla", mode: "publish"},
		"123456789":                           {resource: "123456789", mode: "request"},
		"bla,token:foobar":                    {resource: "bla", token: "foobar", mode: "request"},
		"bla,token:foo,bar":                   {resource: "bla", token: "foo,bar", mode: "request"},
		"123456789,mode:publish,token:foobar": {resource: "123456789", token: "foobar", mode: "publish"},
		"mode:publish":                        {resource: "mode:publish", mode: "request"},
	}

	for streamid, wantsi := range streamids {
		si, err := parseStreamId(streamid)

		require.NoError(t, err)
		require.Equal(t, wantsi, si)
	}
}

func TestParseOldStreamId(t *testing.T) {
	streamids := map[string]streamInfo{
		"#!:":                                   {},
		"#!:key=value":                          {},
		"#!:m=publish":                          {mode: "publish"},
		"#!:r=123456789":                        {resource: "123456789"},
		"#!:token=foobar":                       {token: "foobar"},
		"#!:token=foo,bar":                      {token: "foo"},
		"#!:m=publish,r=123456789,token=foobar": {mode: "publish", resource: "123456789", token: "foobar"},
	}

	for streamid, wantsi := range streamids {
		si, _ := parseOldStreamId(streamid)

		require.Equal(t, wantsi, si)
	}
}
