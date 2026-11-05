package rtsp

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/datarhei/joy4/av"
	"github.com/datarhei/joy4/av/avutil"
	"github.com/datarhei/joy4/codec"
	"github.com/datarhei/joy4/codec/aacparser"
	"github.com/datarhei/joy4/codec/h264parser"
	"github.com/datarhei/joy4/format/rtsp/sdp"
	"github.com/datarhei/joy4/utils/bits/pio"
	"io"
	"net"
	"net/textproto"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var ErrCodecDataChange = fmt.Errorf("rtsp: codec data change, please call HandleCodecDataChange()")

var DebugRtp = false
var DebugRtsp = false
var SkipErrRtpBlock = false

const (
	stageDescribeDone = iota + 1
	stageSetupDone
	stageWaitCodecData
	stageCodecDataDone
)

type Client struct {
	DebugRtsp bool
	DebugRtp  bool
	Headers   []string

	SkipErrRtpBlock bool

	RtspTimeout          time.Duration
	RtpTimeout           time.Duration
	RtpKeepAliveTimeout  time.Duration
	rtpKeepaliveTimer    time.Time
	rtpKeepaliveEnterCnt int

	stage int

	setupIdx []int
	setupMap []int

	authHeaders func(method string) []string

	url         *url.URL
	conn        *connWithTimeout
	brconn      *bufio.Reader
	requestUri  string
	cseq        uint
	streams     []*Stream
	streamsintf []av.CodecData
	session     string
	body        io.Reader
}

type Request struct {
	Header []string
	Uri    string
	Method string
}

type Response struct {
	StatusCode    int
	Headers       textproto.MIMEHeader
	ContentLength int
	Body          []byte

	Block []byte
}

func DialTimeout(uri string, timeout time.Duration) (self *Client, err error) {
	var URL *url.URL
	if URL, err = url.Parse(uri); err != nil {
		return
	}

	if _, _, err := net.SplitHostPort(URL.Host); err != nil {
		URL.Host = URL.Host + ":554"
	}

	dailer := net.Dialer{Timeout: timeout}
	var conn net.Conn
	if conn, err = dailer.Dial("tcp", URL.Host); err != nil {
		return
	}

	u2 := *URL
	u2.User = nil

	connt := &connWithTimeout{Conn: conn}

	self = &Client{
		conn:            connt,
		brconn:          bufio.NewReaderSize(connt, 256),
		url:             URL,
		requestUri:      u2.String(),
		DebugRtp:        DebugRtp,
		DebugRtsp:       DebugRtsp,
		SkipErrRtpBlock: SkipErrRtpBlock,
	}
	return
}

func Dial(uri string) (self *Client, err error) {
	return DialTimeout(uri, 0)
}

func (self *Client) allCodecDataReady() bool {
	for _, si := range self.setupIdx {
		stream := self.streams[si]
		if stream.CodecData == nil {
			return false
		}
	}
	return true
}

func (self *Client) probe() (err error) {
	for {
		if self.allCodecDataReady() {
			break
		}
		if _, err = self.readPacket(); err != nil {
			return
		}
	}
	self.stage = stageCodecDataDone
	return
}

func (self *Client) prepare(stage int) (err error) {
	for self.stage < stage {
		switch self.stage {
		case 0:
			if _, err = self.Describe(); err != nil {
				return
			}

		case stageDescribeDone:
			if err = self.SetupAll(); err != nil {
				return
			}

		case stageSetupDone:
			if err = self.Play(); err != nil {
				return
			}

		case stageWaitCodecData:
			if err = self.probe(); err != nil {
				return
			}
		}
	}
	return
}

func (self *Client) Streams() (streams []av.CodecData, err error) {
	if err = self.prepare(stageCodecDataDone); err != nil {
		return
	}
	for _, si := range self.setupIdx {
		stream := self.streams[si]
		streams = append(streams, stream.CodecData)
	}
	return
}

func (self *Client) SendRtpKeepalive() (err error) {
	if self.RtpKeepAliveTimeout > 0 {
		if self.rtpKeepaliveTimer.IsZero() {
			self.rtpKeepaliveTimer = time.Now()
		} else if time.Now().Sub(self.rtpKeepaliveTimer) > self.RtpKeepAliveTimeout {
			self.rtpKeepaliveTimer = time.Now()
			if self.DebugRtsp {
				fmt.Println("rtp: keep alive")
			}
			req := Request{
				Method: "OPTIONS",
				Uri:    self.requestUri,
			}
			if err = self.WriteRequest(req); err != nil {
				return
			}
		}
	}
	return
}

func (self *Client) WriteRequest(req Request) (err error) {
	self.conn.Timeout = self.RtspTimeout
	self.cseq++

	buf := &bytes.Buffer{}

	fmt.Fprintf(buf, "%s %s RTSP/1.0\r\n", req.Method, req.Uri)
	fmt.Fprintf(buf, "CSeq: %d\r\n", self.cseq)

	if self.authHeaders != nil {
		headers := self.authHeaders(req.Method)
		for _, s := range headers {
			io.WriteString(buf, s)
			io.WriteString(buf, "\r\n")
		}
	}
	for _, s := range req.Header {
		io.WriteString(buf, s)
		io.WriteString(buf, "\r\n")
	}
	for _, s := range self.Headers {
		io.WriteString(buf, s)
		io.WriteString(buf, "\r\n")
	}
	io.WriteString(buf, "\r\n")

	bufout := buf.Bytes()

	if self.DebugRtsp {
		fmt.Print("> ", string(bufout))
	}

	if _, err = self.conn.Write(bufout); err != nil {
		return
	}

	return
}

func (self *Client) parseBlockHeader(h []byte) (length int, no int, valid bool) {
	length = int(h[2])<<8 + int(h[3])
	no = int(h[1])
	if no/2 >= len(self.streams) {
		return
	}

	if no%2 == 0 { // rtp
		if length < 8 {
			return
		}

		// V=2
		if h[4]&0xc0 != 0x80 {
			return
		}

		stream := self.streams[no/2]
		if int(h[5]&0x7f) != stream.Sdp.PayloadType {
			return
		}

		timestamp := binary.BigEndian.Uint32(h[8:12])
		if stream.firsttimestamp != 0 {
			timestamp -= stream.firsttimestamp
			if timestamp < stream.timestamp {
				return
			} else if timestamp-stream.timestamp > uint32(stream.timeScale()*60*60) {
				return
			}
		}
	} else { // rtcp
	}

	valid = true
	return
}

func (self *Client) parseHeaders(b []byte) (statusCode int, headers textproto.MIMEHeader, err error) {
	var line string
	r := textproto.NewReader(bufio.NewReader(bytes.NewReader(b)))
	if line, err = r.ReadLine(); err != nil {
		err = fmt.Errorf("rtsp: header invalid")
		return
	}

	if codes := strings.Split(line, " "); len(codes) >= 2 {
		if statusCode, err = strconv.Atoi(codes[1]); err != nil {
			err = fmt.Errorf("rtsp: header invalid: %s", err)
			return
		}
	}

	headers, _ = r.ReadMIMEHeader()
	return
}

func (self *Client) handleResp(res *Response) (err error) {
	if sess := res.Headers.Get("Session"); sess != "" && self.session == "" {
		if fields := strings.Split(sess, ";"); len(fields) > 0 {
			self.session = fields[0]
		}
	}
	if res.StatusCode == 401 {
		if err = self.handle401(res); err != nil {
			return
		}
	}
	return
}

func (self *Client) handle401(res *Response) (err error) {
	/*
		RTSP/1.0 401 Unauthorized
		CSeq: 2
		Date: Wed, May 04 2016 10:10:51 GMT
		WWW-Authenticate: Digest realm="LIVE555 Streaming Media", nonce="c633aaf8b83127633cbe98fac1d20d87"
	*/
	authval := res.Headers.Get("WWW-Authenticate")
	hdrval := strings.SplitN(authval, " ", 2)
	var realm, nonce string

	if len(hdrval) == 2 {
		for _, field := range strings.Split(hdrval[1], ",") {
			field = strings.Trim(field, ", ")
			if keyval := strings.Split(field, "="); len(keyval) == 2 {
				key := keyval[0]
				val := strings.Trim(keyval[1], `"`)
				switch key {
				case "realm":
					realm = val
				case "nonce":
					nonce = val
				}
			}
		}

		if realm != "" {
			var username string
			var password string

			if self.url.User == nil {
				err = fmt.Errorf("rtsp: no username")
				return
			}
			username = self.url.User.Username()
			password, _ = self.url.User.Password()

			self.authHeaders = func(method string) []string {
				var headers []string
				if nonce == "" {
					headers = []string{
						fmt.Sprintf(`Authorization: Basic %s`, base64.StdEncoding.EncodeToString([]byte(username+":"+password))),
					}
				} else {
					hs1 := md5hash(username + ":" + realm + ":" + password)
					hs2 := md5hash(method + ":" + self.requestUri)
					response := md5hash(hs1 + ":" + nonce + ":" + hs2)
					headers = []string{fmt.Sprintf(
						`Authorization: Digest username="%s", realm="%s", nonce="%s", uri="%s", response="%s"`,
						username, realm, nonce, self.requestUri, response)}
				}
				return headers
			}
		}
	}

	return
}

func (self *Client) findRTSP() (block []byte, data []byte, err error) {
	const (
		R = iota + 1
		T
		S
		Header
		Dollar
	)
	var _peek [8]byte
	peek := _peek[0:0]
	stat := 0

	for i := 0; ; i++ {
		var b byte
		if b, err = self.brconn.ReadByte(); err != nil {
			return
		}
		switch b {
		case 'R':
			if stat == 0 {
				stat = R
			}
		case 'T':
			if stat == R {
				stat = T
			}
		case 'S':
			if stat == T {
				stat = S
			}
		case 'P':
			if stat == S {
				stat = Header
			}
		case '$':
			if stat != Dollar {
				stat = Dollar
				peek = _peek[0:0]
			}
		default:
			if stat != Dollar {
				stat = 0
				peek = _peek[0:0]
			}
		}

		if false && self.DebugRtp {
			fmt.Println("rtsp: findRTSP", i, b)
		}

		if stat != 0 {
			peek = append(peek, b)
		}
		if stat == Header {
			data = peek
			return
		}

		if stat == Dollar && len(peek) >= 12 {
			if self.DebugRtp {
				fmt.Println("rtsp: dollar at", i, len(peek))
			}
			if blocklen, _, ok := self.parseBlockHeader(peek); ok {
				left := blocklen + 4 - len(peek)
				block = append(peek, make([]byte, left)...)
				if _, err = io.ReadFull(self.brconn, block[len(peek):]); err != nil {
					return
				}
				return
			}
			stat = 0
			peek = _peek[0:0]
		}
	}

	return
}

func (self *Client) readLFLF() (block []byte, data []byte, err error) {
	const (
		LF = iota + 1
		LFLF
	)
	peek := []byte{}
	stat := 0
	dollarpos := -1
	lpos := 0
	pos := 0

	for {
		var b byte
		if b, err = self.brconn.ReadByte(); err != nil {
			return
		}
		switch b {
		case '\n':
			if stat == 0 {
				stat = LF
				lpos = pos
			} else if stat == LF {
				if pos-lpos <= 2 {
					stat = LFLF
				} else {
					lpos = pos
				}
			}
		case '$':
			dollarpos = pos
		}
		peek = append(peek, b)

		if stat == LFLF {
			data = peek
			return
		} else if dollarpos != -1 && dollarpos-pos >= 12 {
			hdrlen := dollarpos - pos
			start := len(peek) - hdrlen
			if blocklen, _, ok := self.parseBlockHeader(peek[start:]); ok {
				block = append(peek[start:], make([]byte, blocklen+4-hdrlen)...)
				if _, err = io.ReadFull(self.brconn, block[hdrlen:]); err != nil {
					return
				}
				return
			}
			dollarpos = -1
		}

		pos++
	}

	return
}

func (self *Client) readResp(b []byte) (res Response, err error) {
	if res.StatusCode, res.Headers, err = self.parseHeaders(b); err != nil {
		return
	}
	res.ContentLength, _ = strconv.Atoi(res.Headers.Get("Content-Length"))
	if res.ContentLength > 0 {
		res.Body = make([]byte, res.ContentLength)
		if _, err = io.ReadFull(self.brconn, res.Body); err != nil {
			return
		}
	}
	if err = self.handleResp(&res); err != nil {
		return
	}
	return
}

func (self *Client) poll() (res Response, err error) {
	var block []byte
	var rtsp []byte
	var headers []byte

	self.conn.Timeout = self.RtspTimeout
	for {
		if block, rtsp, err = self.findRTSP(); err != nil {
			return
		}
		if len(block) > 0 {
			res.Block = block
			return
		} else {
			if block, headers, err = self.readLFLF(); err != nil {
				return
			}
			if len(block) > 0 {
				res.Block = block
				return
			}
			if res, err = self.readResp(append(rtsp, headers...)); err != nil {
				return
			}
		}
		return
	}

	return
}

func (self *Client) ReadResponse() (res Response, err error) {
	for {
		if res, err = self.poll(); err != nil {
			return
		}
		if res.StatusCode > 0 {
			return
		}
	}
	return
}

func (self *Client) SetupAll() (err error) {
	idx := []int{}
	for i := range self.streams {
		idx = append(idx, i)
	}
	return self.Setup(idx)
}

func (self *Client) Setup(idx []int) (err error) {
	if err = self.prepare(stageDescribeDone); err != nil {
		return
	}

	self.setupMap = make([]int, len(self.streams))
	for i := range self.setupMap {
		self.setupMap[i] = -1
	}
	self.setupIdx = idx

	for i, si := range idx {
		self.setupMap[si] = i

		uri := ""
		control := self.streams[si].Sdp.Control
		if strings.HasPrefix(control, "rtsp://") {
			uri = control
		} else {
			uri = self.requestUri + "/" + control
		}
		req := Request{Method: "SETUP", Uri: uri}
		req.Header = append(req.Header, fmt.Sprintf("Transport: RTP/AVP/TCP;unicast;interleaved=%d-%d", si*2, si*2+1))
		if self.session != "" {
			req.Header = append(req.Header, "Session: "+self.session)
		}
		if err = self.WriteRequest(req); err != nil {
			return
		}
		if _, err = self.ReadResponse(); err != nil {
			return
		}
	}

	if self.stage == stageDescribeDone {
		self.stage = stageSetupDone
	}
	return
}

func md5hash(s string) string {
	h := md5.Sum([]byte(s))
	return hex.EncodeToString(h[:])
}

func (self *Client) Describe() (streams []sdp.Media, err error) {
	var res Response

	for i := 0; i < 2; i++ {
		req := Request{
			Method: "DESCRIBE",
			Uri:    self.requestUri,
			Header: []string{"Accept: application/sdp"},
		}
		if err = self.WriteRequest(req); err != nil {
			return
		}
		if res, err = self.ReadResponse(); err != nil {
			return
		}
		if res.StatusCode == 200 {
			break
		}
	}
	if res.ContentLength == 0 {
		err = fmt.Errorf("rtsp: Describe failed, StatusCode=%d", res.StatusCode)
		return
	}

	body := string(res.Body)

	if self.DebugRtsp {
		fmt.Println("<", body)
	}

	_, medias := sdp.Parse(body)

	self.streams = []*Stream{}
	for _, media := range medias {
		stream := &Stream{Sdp: media, client: self}
		stream.makeCodecData()
		self.streams = append(self.streams, stream)
		streams = append(streams, media)
	}

	if self.stage == 0 {
		self.stage = stageDescribeDone
	}
	return
}

func (self *Client) Options() (err error) {
	req := Request{
		Method: "OPTIONS",
		Uri:    self.requestUri,
	}
	if self.session != "" {
		req.Header = append(req.Header, "Session: "+self.session)
	}
	if err = self.WriteRequest(req); err != nil {
		return
	}
	if _, err = self.ReadResponse(); err != nil {
		return
	}
	return
}

func (self *Client) HandleCodecDataChange() (_newcli *Client, err error) {
	newcli := &Client{}
	*newcli = *self

	newcli.streams = []*Stream{}
	for _, stream := range self.streams {
		newstream := &Stream{}
		*newstream = *stream
		newstream.client = newcli

		if newstream.isCodecDataChange() {
			if err = newstream.makeCodecData(); err != nil {
				return
			}
			newstream.clearCodecDataChange()
		}
		newcli.streams = append(newcli.streams, newstream)
	}

	_newcli = newcli
	return
}

func (self *Stream) clearCodecDataChange() {
	self.spsChanged = false
	self.ppsChanged = false
}

func (self *Stream) isCodecDataChange() bool {
	if self.spsChanged && self.ppsChanged {
		return true
	}
	return false
}

func (self *Stream) timeScale() int {
	t := self.Sdp.TimeScale
	if t == 0 {
		// https://tools.ietf.org/html/rfc5391
		t = 8000
	}
	return t
}

func (self *Stream) makeCodecData() (err error) {
	media := self.Sdp

	if media.PayloadType >= 96 && media.PayloadType <= 127 {
		switch media.Type {
		case av.H264:
			for _, nalu := range media.SpropParameterSets {
				if len(nalu) > 0 {
					self.handleH264Payload(0, nalu)
				}
			}

			if len(self.sps) == 0 || len(self.pps) == 0 {
				if nalus, typ := h264parser.SplitNALUs(media.Config); typ != h264parser.NALU_RAW {
					for _, nalu := range nalus {
						if len(nalu) > 0 {
							self.handleH264Payload(0, nalu)
						}
					}
				}
			}

			if len(self.sps) > 0 && len(self.pps) > 0 {
				if self.CodecData, err = h264parser.NewCodecDataFromSPSAndPPS(self.sps, self.pps); err != nil {
					err = fmt.Errorf("rtsp: h264 sps/pps invalid: %s", err)
					return
				}
			} else {
				err = fmt.Errorf("rtsp: missing h264 sps or pps")
				return
			}

		case av.AAC:
			if len(media.Config) == 0 {
				err = fmt.Errorf("rtsp: aac sdp config missing")
				return
			}
			if self.CodecData, err = aacparser.NewCodecDataFromMPEG4AudioConfigBytes(media.Config); err != nil {
				err = fmt.Errorf("rtsp: aac sdp config invalid: %s", err)
				return
			}
		}
	} else {
		switch media.PayloadType {
		case 0:
			self.CodecData = codec.NewPCMMulawCodecData()

		case 8:
			self.CodecData = codec.NewPCMAlawCodecData()

		default:
			err = fmt.Errorf("rtsp: PayloadType=%d unsupported", media.PayloadType)
			return
		}
	}

	return
}

func (self *Stream) handleBuggyAnnexbH264Packet(timestamp uint32, packet []byte) (isBuggy bool, err error) {
	if len(packet) >= 4 && packet[0] == 0 && packet[1] == 0 && packet[2] == 0 && packet[3] == 1 {
		isBuggy = true
		if nalus, typ := h264parser.SplitNALUs(packet); typ != h264parser.NALU_RAW {
			for _, nalu := range nalus {
				if len(nalu) > 0 {
					if err = self.handleH264Payload(timestamp, nalu); err != nil {
						return
					}
				}
			}
		}
	}
	return
}

func (self *Stream) handleH264Payload(timestamp uint32, packet []byte) (err error) {
	if len(packet) < 2 {
		err = fmt.Errorf("rtp: h264 packet too short")
		return
	}

	var isBuggy bool
	if isBuggy, err = self.handleBuggyAnnexbH264Packet(timestamp, packet); isBuggy {
		return
	}

	naluType := packet[0] & 0x1f

	/*
		Table 7-1 – NAL unit type codes
		1   ￼Coded slice of a non-IDR picture
		5    Coded slice of an IDR picture
		6    Supplemental enhancement information (SEI)
		7    Sequence parameter set
		8    Picture parameter set
		1-23     NAL unit  Single NAL unit packet             5.6
		24       STAP-A    Single-time aggregation packet     5.7.1
		25       STAP-B    Single-time aggregation packet     5.7.1
		26       MTAP16    Multi-time aggregation packet      5.7.2
		27       MTAP24    Multi-time aggregation packet      5.7.2
		28       FU-A      Fragmentation unit                 5.8
		29       FU-B      Fragmentation unit                 5.8
		30-31    reserved                                     -
	*/
	switch {
	case naluType >= 1 && naluType <= 5:
		if naluType == 5 {
			self.pkt.IsKeyFrame = true
		}
		self.gotpkt = true
		// raw nalu to avcc
		b := make([]byte, 4+len(packet))
		pio.PutU32BE(b[0:4], uint32(len(packet)))
		copy(b[4:], packet)
		self.pkt.Data = b
		self.timestamp = timestamp

	case naluType == 7: // sps
		if self.client != nil && self.client.DebugRtp {
			fmt.Println("rtsp: got sps")
		}
		if len(self.sps) == 0 {
			self.sps = packet
			self.makeCodecData()
		} else if bytes.Compare(self.sps, packet) != 0 {
			self.spsChanged = true
			self.sps = packet
			if self.client != nil && self.client.DebugRtp {
				fmt.Println("rtsp: sps changed")
			}
		}

	case naluType == 8: // pps
		if self.client != nil && self.client.DebugRtp {
			fmt.Println("rtsp: got pps")
		}
		if len(self.pps) == 0 {
			self.pps = packet
			self.makeCodecData()
		} else if bytes.Compare(self.pps, packet) != 0 {
			self.ppsChanged = true
			self.pps = packet
			if self.client != nil && self.client.DebugRtp {
				fmt.Println("rtsp: pps changed")
			}
		}

	case naluType == 28: // FU-A
		/*
			0                   1                   2                   3
			0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
			+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
			| FU indicator  |   FU header   |                               |
			+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+                               |
			|                                                               |
			|                         FU payload                            |
			|                                                               |
			|                               +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
			|                               :...OPTIONAL RTP padding        |
			+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
			Figure 14.  RTP payload format for FU-A

			The FU indicator octet has the following format:
			+---------------+
			|0|1|2|3|4|5|6|7|
			+-+-+-+-+-+-+-+-+
			|F|NRI|  Type   |
			+---------------+


			The FU header has the following format:
			+---------------+
			|0|1|2|3|4|5|6|7|
			+-+-+-+-+-+-+-+-+
			|S|E|R|  Type   |
			+---------------+

			S: 1 bit
			When set to one, the Start bit indicates the start of a fragmented
			NAL unit.  When the following FU payload is not the start of a
			fragmented NAL unit payload, the Start bit is set to zero.

			E: 1 bit
			When set to one, the End bit indicates the end of a fragmented NAL
			unit, i.e., the last byte of the payload is also the last byte of
			the fragmented NAL unit.  When the following FU payload is not the
			last fragment of a fragmented NAL unit, the End bit is set to
			zero.

			R: 1 bit
			The Reserved bit MUST be equal to 0 and MUST be ignored by the
			receiver.

			Type: 5 bits
			The NAL unit payload type as defined in table 7-1 of [1].
		*/
		fuIndicator := packet[0]
		fuHeader := packet[1]
		isStart := fuHeader&0x80 != 0
		isEnd := fuHeader&0x40 != 0
		if isStart {
			self.fuStarted = true
			self.fuBuffer = []byte{fuIndicator&0xe0 | fuHeader&0x1f}
		}
		if self.fuStarted {
			self.fuBuffer = append(self.fuBuffer, packet[2:]...)
			if isEnd {
				self.fuStarted = false
				if err = self.handleH264Payload(timestamp, self.fuBuffer); err != nil {
					return
				}
			}
		}

	case naluType == 24: // STAP-A
		/*
			0                   1                   2                   3
			0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
			+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
			|                          RTP Header                           |
			+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
			|STAP-A NAL HDR |         NALU 1 Size           | NALU 1 HDR    |
			+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
			|                         NALU 1 Data                           |
			:                                                               :
			+               +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
			|               | NALU 2 Size                   | NALU 2 HDR    |
			+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
			|                         NALU 2 Data                           |
			:                                                               :
			|                               +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
			|                               :...OPTIONAL RTP padding        |
			+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

			Figure 7.  An example of an RTP packet including an STAP-A
			containing two single-time aggregation units
		*/
		packet = packet[1:]
		for len(packet) >= 2 {
			size := int(packet[0])<<8 | int(packet[1])
			if size+2 > len(packet) {
				break
			}
			if err = self.handleH264Payload(timestamp, packet[2:size+2]); err != nil {
				return
			}
			packet = packet[size+2:]
		}
		return

	case naluType >= 6 && naluType <= 23: // other single NALU packet
	case naluType == 25: // STAB-B
	case naluType == 26: // MTAP-16
	case naluType == 27: // MTAP-24
	case naluType == 28: // FU-B

	default:
		err = fmt.Errorf("rtsp: unsupported H264 naluType=%d", naluType)
		return
	}

	return
}

func (self *Stream) handleRtpPacket(packet []byte) (err error) {
	if self.isCodecDataChange() {
		err = ErrCodecDataChange
		return
	}

	if self.client != nil && self.client.DebugRtp {
		fmt.Println("rtp: packet", self.CodecData.Type(), "len", len(packet))
		dumpsize := len(packet)
		if dumpsize > 32 {
			dumpsize = 32
		}
		fmt.Print(hex.Dump(packet[:dumpsize]))
	}

	/*
		0                   1                   2                   3
		0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
		+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
		|V=2|P|X|  CC   |M|     PT      |       sequence number         |
		+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
		|                           timestamp                           |
		+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
		|           synchronization source (SSRC) identifier            |
		+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+
		|            contributing source (CSRC) identifiers             |
		|                             ....                              |
		+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	*/
	if len(packet) < 8 {
		err = fmt.Errorf("rtp: packet too short")
		return
	}
	payloadOffset := 12 + int(packet[0]&0xf)*4
	if payloadOffset > len(packet) {
		err = fmt.Errorf("rtp: packet too short")
		return
	}
	timestamp := binary.BigEndian.Uint32(packet[4:8])
	payload := packet[payloadOffset:]

	/*
		PT 	Encoding Name 	Audio/Video (A/V) 	Clock Rate (Hz) 	Channels 	Reference
		0	PCMU	A	8000	1	[RFC3551]
		1	Reserved
		2	Reserved
		3	GSM	A	8000	1	[RFC3551]
		4	G723	A	8000	1	[Vineet_Kumar][RFC3551]
		5	DVI4	A	8000	1	[RFC3551]
		6	DVI4	A	16000	1	[RFC3551]
		7	LPC	A	8000	1	[RFC3551]
		8	PCMA	A	8000	1	[RFC3551]
		9	G722	A	8000	1	[RFC3551]
		10	L16	A	44100	2	[RFC3551]
		11	L16	A	44100	1	[RFC3551]
		12	QCELP	A	8000	1	[RFC3551]
		13	CN	A	8000	1	[RFC3389]
		14	MPA	A	90000		[RFC3551][RFC2250]
		15	G728	A	8000	1	[RFC3551]
		16	DVI4	A	11025	1	[Joseph_Di_Pol]
		17	DVI4	A	22050	1	[Joseph_Di_Pol]
		18	G729	A	8000	1	[RFC3551]
		19	Reserved	A
		20	Unassigned	A
		21	Unassigned	A
		22	Unassigned	A
		23	Unassigned	A
		24	Unassigned	V
		25	CelB	V	90000		[RFC2029]
		26	JPEG	V	90000		[RFC2435]
		27	Unassigned	V
		28	nv	V	90000		[RFC3551]
		29	Unassigned	V
		30	Unassigned	V
		31	H261	V	90000		[RFC4587]
		32	MPV	V	90000		[RFC2250]
		33	MP2T	AV	90000		[RFC2250]
		34	H263	V	90000		[Chunrong_Zhu]
		35-71	Unassigned	?
		72-76	Reserved for RTCP conflict avoidance				[RFC3551]
		77-95	Unassigned	?
		96-127	dynamic	?			[RFC3551]
	*/
	//payloadType := packet[1]&0x7f

	switch self.Sdp.Type {
	case av.H264:
		if err = self.handleH264Payload(timestamp, payload); err != nil {
			return
		}

	case av.AAC:
		if len(payload) < 4 {
			err = fmt.Errorf("rtp: aac packet too short")
			return
		}
		payload = payload[4:] // TODO: remove this hack
		self.gotpkt = true
		self.pkt.Data = payload
		self.timestamp = timestamp

	default:
		self.gotpkt = true
		self.pkt.Data = payload
		self.timestamp = timestamp
	}

	return
}

func (self *Client) Play() (err error) {
	req := Request{
		Method: "PLAY",
		Uri:    self.requestUri,
	}
	req.Header = append(req.Header, "Session: "+self.session)
	if err = self.WriteRequest(req); err != nil {
		return
	}

	if self.allCodecDataReady() {
		self.stage = stageCodecDataDone
	} else {
		self.stage = stageWaitCodecData
	}
	return
}

func (self *Client) Teardown() (err error) {
	req := Request{
		Method: "TEARDOWN",
		Uri:    self.requestUri,
	}
	req.Header = append(req.Header, "Session: "+self.session)
	if err = self.WriteRequest(req); err != nil {
		return
	}
	return
}

func (self *Client) Close() (err error) {
	return self.conn.Conn.Close()
}

func (self *Client) handleBlock(block []byte) (pkt av.Packet, ok bool, err error) {
	_, blockno, _ := self.parseBlockHeader(block)
	if blockno%2 != 0 {
		if self.DebugRtp {
			fmt.Println("rtsp: rtcp block len", len(block)-4)
		}
		return
	}

	i := blockno / 2
	if i >= len(self.streams) {
		err = fmt.Errorf("rtsp: block no=%d invalid", blockno)
		return
	}
	stream := self.streams[i]

	herr := stream.handleRtpPacket(block[4:])
	if herr != nil {
		if !self.SkipErrRtpBlock {
			err = herr
			return
		}
	}

	if stream.gotpkt {
		/*
			TODO: sync AV by rtcp NTP timestamp
			TODO: handle timestamp overflow
			https://tools.ietf.org/html/rfc3550
			A receiver can then synchronize presentation of the audio and video packets by relating
			their RTP timestamps using the timestamp pairs in RTCP SR packets.
		*/
		if stream.firsttimestamp == 0 {
			stream.firsttimestamp = stream.timestamp
		}
		stream.timestamp -= stream.firsttimestamp

		ok = true
		pkt = stream.pkt
		pkt.Time = time.Duration(stream.timestamp) * time.Second / time.Duration(stream.timeScale())
		pkt.Idx = int8(self.setupMap[i])

		if pkt.Time < stream.lasttime || pkt.Time-stream.lasttime > time.Minute*30 {
			err = fmt.Errorf("rtp: time invalid stream#%d time=%v lasttime=%v", pkt.Idx, pkt.Time, stream.lasttime)
			return
		}
		stream.lasttime = pkt.Time

		if self.DebugRtp {
			fmt.Println("rtp: pktout", pkt.Idx, pkt.Time, len(pkt.Data))
		}

		stream.pkt = av.Packet{}
		stream.gotpkt = false
	}

	return
}

func (self *Client) readPacket() (pkt av.Packet, err error) {
	if err = self.SendRtpKeepalive(); err != nil {
		return
	}

	for {
		var res Response
		for {
			if res, err = self.poll(); err != nil {
				return
			}
			if len(res.Block) > 0 {
				break
			}
		}

		var ok bool
		if pkt, ok, err = self.handleBlock(res.Block); err != nil {
			return
		}
		if ok {
			return
		}
	}

	return
}

func (self *Client) ReadPacket() (pkt av.Packet, err error) {
	if err = self.prepare(stageCodecDataDone); err != nil {
		return
	}
	return self.readPacket()
}

func Handler(h *avutil.RegisterHandler) {
	h.UrlDemuxer = func(uri string) (ok bool, demuxer av.DemuxCloser, err error) {
		if !strings.HasPrefix(uri, "rtsp://") {
			return
		}
		ok = true
		demuxer, err = Dial(uri)
		return
	}
}
