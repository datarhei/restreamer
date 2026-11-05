package sdp

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/datarhei/joy4/av"
	"strconv"
	"strings"
)

type Session struct {
	Uri string
}

type Media struct {
	AVType             string
	Type               av.CodecType
	TimeScale          int
	Control            string
	Rtpmap             int
	Config             []byte
	SpropParameterSets [][]byte
	PayloadType        int
	SizeLength         int
	IndexLength        int
}

func Parse(content string) (sess Session, medias []Media) {
	var media *Media

	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		typeval := strings.SplitN(line, "=", 2)
		if len(typeval) == 2 {
			fields := strings.SplitN(typeval[1], " ", 2)

			switch typeval[0] {
			case "m":
				if len(fields) > 0 {
					switch fields[0] {
					case "audio", "video":
						medias = append(medias, Media{AVType: fields[0]})
						media = &medias[len(medias)-1]
						mfields := strings.Split(fields[1], " ")
						if len(mfields) >= 3 {
							media.PayloadType, _ = strconv.Atoi(mfields[2])
						}
					}
				}

			case "u":
				sess.Uri = typeval[1]

			case "a":
				if media != nil {
					for _, field := range fields {
						keyval := strings.SplitN(field, ":", 2)
						if len(keyval) >= 2 {
							key := keyval[0]
							val := keyval[1]
							switch key {
							case "control":
								media.Control = val
							case "rtpmap":
								media.Rtpmap, _ = strconv.Atoi(val)
							}
						}
						keyval = strings.Split(field, "/")
						if len(keyval) >= 2 {
							key := keyval[0]
							switch strings.ToUpper(key) {
							case "MPEG4-GENERIC":
								media.Type = av.AAC
							case "H264":
								media.Type = av.H264
							}
							if i, err := strconv.Atoi(keyval[1]); err == nil {
								media.TimeScale = i
							}
							if false {
								fmt.Println("sdp:", keyval[1], media.TimeScale)
							}
						}
						keyval = strings.Split(field, ";")
						if len(keyval) > 1 {
							for _, field := range keyval {
								keyval := strings.SplitN(field, "=", 2)
								if len(keyval) == 2 {
									key := strings.TrimSpace(keyval[0])
									val := keyval[1]
									switch key {
									case "config":
										media.Config, _ = hex.DecodeString(val)
									case "sizelength":
										media.SizeLength, _ = strconv.Atoi(val)
									case "indexlength":
										media.IndexLength, _ = strconv.Atoi(val)
									case "sprop-parameter-sets":
										fields := strings.Split(val, ",")
										for _, field := range fields {
											val, _ := base64.StdEncoding.DecodeString(field)
											media.SpropParameterSets = append(media.SpropParameterSets, val)
										}
									}
								}
							}
						}
					}
				}

			}
		}
	}
	return
}
