package rtmp

import (
	"bufio"
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/tls"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net"
	"net/url"
	"strings"
	"time"

	"github.com/datarhei/joy4/av"
	"github.com/datarhei/joy4/av/avutil"
	"github.com/datarhei/joy4/format/flv"
	"github.com/datarhei/joy4/format/flv/flvio"
	"github.com/datarhei/joy4/utils/bits/pio"
)

var Debug bool

func ParseURL(uri string) (u *url.URL, err error) {
	if u, err = url.Parse(uri); err != nil {
		return
	}
	if _, _, serr := net.SplitHostPort(u.Host); serr != nil {
		u.Host += ":1935"
	}
	return
}

type DialOptions struct {
	MaxProbePacketCount int
	DebugChunks         func(conn net.Conn) bool
}

func Dial(uri string, options DialOptions) (conn *Conn, err error) {
	return DialTimeout(uri, 0, options)
}

func DialTimeout(uri string, timeout time.Duration, options DialOptions) (conn *Conn, err error) {
	var u *url.URL
	if u, err = ParseURL(uri); err != nil {
		return
	}

	dailer := net.Dialer{Timeout: timeout}
	var netconn net.Conn
	if netconn, err = dailer.Dial("tcp", u.Host); err != nil {
		return
	}

	conn = NewConn(netconn)
	conn.URL = u
	conn.prober = flv.NewProber(options.MaxProbePacketCount)
	if options.DebugChunks != nil {
		conn.debugChunks = options.DebugChunks(netconn)
	}

	return
}

var ErrServerClosed = errors.New("server closed")

type Server struct {
	Addr          string
	TLSConfig     *tls.Config
	HandlePublish func(*Conn)
	HandlePlay    func(*Conn)
	HandleConn    func(*Conn)

	MaxProbePacketCount   int
	SkipInvalidMessages   bool
	DebugChunks           func(conn net.Conn) bool
	ConnectionIdleTimeout time.Duration

	listener net.Listener
	doneChan chan struct{}
}

func (s *Server) HandleNetConn(netconn net.Conn) (err error) {
	conn := NewConn(netconn)
	conn.prober = flv.NewProber(s.MaxProbePacketCount)
	conn.skipInvalidMessages = s.SkipInvalidMessages
	if s.DebugChunks != nil {
		conn.debugChunks = s.DebugChunks(netconn)
	}
	conn.netconn.SetIdleTimeout(s.ConnectionIdleTimeout)
	conn.isserver = true

	err = s.handleConn(conn)
	if Debug {
		fmt.Println("rtmp: server: client closed err:", err)
	}

	conn.Close()

	return
}

func (s *Server) handleConn(conn *Conn) (err error) {
	if s.HandleConn != nil {
		s.HandleConn(conn)
	} else {
		if err = conn.prepare(stageCommandDone, 0); err != nil {
			return
		}

		if conn.playing {
			conn.netconn.SetReadIdleTimeout(0)
			if s.HandlePlay != nil {
				s.HandlePlay(conn)
			}
		} else if conn.publishing {
			conn.netconn.SetWriteIdleTimeout(0)
			if s.HandlePublish != nil {
				s.HandlePublish(conn)
			}
		}
	}

	return
}

func (s *Server) ListenAndServe() error {
	addr := s.Addr
	if addr == "" {
		addr = ":1935"
	}

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	return s.Serve(listener)
}

func (s *Server) ListenAndServeTLS(certFile, keyFile string) error {
	addr := s.Addr
	if addr == "" {
		addr = ":1935"
	}

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	return s.ServeTLS(listener, certFile, keyFile)
}

func (s *Server) ServeTLS(listener net.Listener, certFile, keyFile string) error {
	var config *tls.Config

	if s.TLSConfig != nil {
		config = s.TLSConfig.Clone()
	} else {
		config = &tls.Config{}
	}

	configHasCert := len(config.Certificates) > 0 || config.GetCertificate != nil

	if !configHasCert || certFile != "" || keyFile != "" {
		config.Certificates = make([]tls.Certificate, 1)

		var err error
		config.Certificates[0], err = tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			return err
		}
	}

	listener = tls.NewListener(listener, config)

	return s.Serve(listener)
}

func (s *Server) Serve(listener net.Listener) error {
	s.doneChan = make(chan struct{})

	s.listener = listener
	defer s.listener.Close()

	if Debug {
		fmt.Println("rtmp: server: listening on", s.listener.Addr().String())
	}

	for {
		netconn, err := s.listener.Accept()
		if err != nil {
			select {
			case <-s.doneChan:
				return ErrServerClosed
			default:
			}

			return err
		}

		if Debug {
			fmt.Println("rtmp: server: accepted")
		}

		go func(conn net.Conn) {
			err := s.HandleNetConn(conn)
			if Debug {
				fmt.Println("rtmp: server: client closed err:", err)
			}
		}(netconn)
	}
}

func (s *Server) Close() {
	if s.listener != nil {
		s.listener.Close()
		s.listener = nil

		close(s.doneChan)
	}
}

const (
	stageHandshakeDone = iota + 1
	stageCommandDone
	stageCodecDataDone
)

const (
	prepareReading = iota + 1
	prepareWriting
)

type Conn struct {
	URL             *url.URL
	OnPlayOrPublish func(string, flvio.AMFMap) error

	prober  *flv.Prober
	streams []av.CodecData

	bufr *bufio.Reader
	bufw *bufio.Writer
	ackn uint32

	writebuf []byte
	readbuf  []byte

	netconn   *idleConn
	txrxcount *txrxcount

	writeMaxChunkSize int
	readMaxChunkSize  int
	readAckSize       uint32
	readcsmap         map[uint32]*chunkStream

	isserver            bool
	publishing, playing bool
	reading, writing    bool
	stage               int

	avmsgsid uint32

	gotcommand     bool
	commandname    string
	commandtransid float64
	commandobj     flvio.AMFMap
	commandparams  []interface{}

	gotmsg      bool
	timestamp   uint32
	msgdata     []byte
	msgtypeid   uint8
	datamsgvals []interface{}
	avtag       flvio.Tag

	metadata flvio.AMFMap

	eventtype uint16

	start time.Time

	skipInvalidMessages bool

	debugChunks bool
}

type idleConn struct {
	net.Conn
	ReadIdleTimeout  time.Duration
	WriteIdleTimeout time.Duration
}

func (t *idleConn) Read(p []byte) (int, error) {
	if t.ReadIdleTimeout > 0 {
		t.Conn.SetReadDeadline(time.Now().Add(t.ReadIdleTimeout))
	}

	n, err := t.Conn.Read(p)
	return n, err
}

func (t *idleConn) Write(p []byte) (int, error) {
	if t.WriteIdleTimeout > 0 {
		t.Conn.SetWriteDeadline(time.Now().Add(t.WriteIdleTimeout))
	}

	n, err := t.Conn.Write(p)
	return n, err
}

func (t *idleConn) SetReadIdleTimeout(d time.Duration) error {
	t.ReadIdleTimeout = d

	deadline := time.Time{}
	if t.ReadIdleTimeout > 0 {
		deadline = time.Now().Add(d)
	}

	return t.Conn.SetReadDeadline(deadline)
}

func (t *idleConn) SetWriteIdleTimeout(d time.Duration) error {
	t.WriteIdleTimeout = d

	deadline := time.Time{}
	if t.WriteIdleTimeout > 0 {
		deadline = time.Now().Add(d)
	}

	return t.Conn.SetWriteDeadline(deadline)
}

func (t *idleConn) SetIdleTimeout(d time.Duration) error {
	err := t.SetReadIdleTimeout(d)
	if err != nil {
		return err
	}

	return t.SetWriteIdleTimeout(d)
}

type txrxcount struct {
	io.ReadWriteCloser
	txbytes uint64
	rxbytes uint64
}

func (t *txrxcount) Read(p []byte) (int, error) {
	n, err := t.ReadWriteCloser.Read(p)
	t.rxbytes += uint64(n)
	return n, err
}

func (t *txrxcount) Write(p []byte) (int, error) {
	n, err := t.ReadWriteCloser.Write(p)
	t.txbytes += uint64(n)
	return n, err
}

func NewConn(netconn net.Conn) *Conn {
	conn := &Conn{}
	conn.prober = &flv.Prober{}
	conn.netconn = &idleConn{
		Conn: netconn,
	}
	conn.readcsmap = make(map[uint32]*chunkStream)
	conn.readMaxChunkSize = 128
	conn.writeMaxChunkSize = 128
	conn.readAckSize = 1048576
	conn.txrxcount = &txrxcount{ReadWriteCloser: conn.netconn}
	conn.bufr = bufio.NewReaderSize(conn.txrxcount, pio.RecommendBufioSize)
	conn.bufw = bufio.NewWriterSize(conn.txrxcount, pio.RecommendBufioSize)
	conn.writebuf = make([]byte, 4096)
	conn.readbuf = make([]byte, 4096)
	conn.start = time.Now()
	return conn
}

type chunkStream struct {
	timenow          uint32
	prevtimenow      uint32
	sametscount      int
	genwallclocktime bool
	timedelta        uint32
	hastimeext       bool
	timeext          uint32
	msgsid           uint32
	msgtypeid        uint8
	msgdatalen       uint32
	msgdataleft      uint32
	msghdrtype       uint8
	msgdata          []byte
	msgcount         uint64
}

func (cs *chunkStream) Start() {
	cs.msgdataleft = cs.msgdatalen
	cs.msgdata = make([]byte, cs.msgdatalen)
}

const (
	msgtypeidUserControl      = 4
	msgtypeidAck              = 3
	msgtypeidWindowAckSize    = 5
	msgtypeidSetPeerBandwidth = 6
	msgtypeidSetChunkSize     = 1
	msgtypeidCommandMsgAMF0   = 20
	msgtypeidCommandMsgAMF3   = 17
	msgtypeidDataMsgAMF0      = 18
	msgtypeidDataMsgAMF3      = 15
	msgtypeidVideoMsg         = 9
	msgtypeidAudioMsg         = 8
)

const (
	eventtypeStreamBegin      = 0
	eventtypeSetBufferLength  = 3
	eventtypeStreamIsRecorded = 4
	eventtypePingRequest      = 6
	eventtypePingResponse     = 7
)

func (conn *Conn) NetConn() net.Conn {
	return conn.netconn.Conn
}

func (conn *Conn) LocalAddr() net.Addr {
	return conn.netconn.LocalAddr()
}

func (conn *Conn) RemoteAddr() net.Addr {
	return conn.netconn.RemoteAddr()
}

func (conn *Conn) SetReadIdleTimeout(d time.Duration) error {
	return conn.netconn.SetReadIdleTimeout(d)
}

func (conn *Conn) SetWriteIdleTimeout(d time.Duration) error {
	return conn.netconn.SetWriteIdleTimeout(d)
}

func (conn *Conn) SetIdleTimeout(d time.Duration) error {
	return conn.netconn.SetIdleTimeout(d)
}

func (conn *Conn) TxBytes() uint64 {
	return conn.txrxcount.txbytes
}

func (conn *Conn) RxBytes() uint64 {
	return conn.txrxcount.rxbytes
}

func (conn *Conn) Close() (err error) {
	return conn.netconn.Close()
}

func (conn *Conn) pollCommand() (err error) {
	for {
		if err = conn.pollMsg(); err != nil {
			return
		}
		if conn.gotcommand {
			return
		}
	}
}

func (conn *Conn) pollAVTag() (tag flvio.Tag, err error) {
	for {
		if err = conn.pollMsg(); err != nil {
			return
		}
		switch conn.msgtypeid {
		case msgtypeidVideoMsg, msgtypeidAudioMsg:
			tag = conn.avtag
			return
		}
	}
}

func (conn *Conn) pollMsg() (err error) {
	conn.gotmsg = false
	conn.gotcommand = false
	conn.datamsgvals = nil
	conn.avtag = flvio.Tag{}
	for {
		if err = conn.readChunk(); err != nil {
			return
		}

		if conn.readAckSize != 0 && conn.ackn > conn.readAckSize {
			if err = conn.writeAck(conn.ackn, false); err != nil {
				return fmt.Errorf("writeACK: %w", err)
			}
			conn.flushWrite()
			conn.ackn = 0
		}

		if conn.gotmsg {
			return
		}
	}
}

func SplitPath(u *url.URL) (app, stream string) {
	pathsegs := strings.SplitN(u.RequestURI(), "/", 3)
	if len(pathsegs) > 1 {
		app = pathsegs[1]
	}
	if len(pathsegs) > 2 {
		stream = pathsegs[2]
	}
	return
}

func getTcUrl(u *url.URL) string {
	app, _ := SplitPath(u)
	nu := *u
	nu.Path = "/" + app
	return nu.String()
}

func createURL(tcurl, app, play string) (*url.URL, error) {
	ps := strings.Split(app+"/"+play, "/")
	out := []string{""}
	for _, s := range ps {
		if len(s) > 0 {
			out = append(out, s)
		}
	}
	if len(out) < 2 {
		out = append(out, "")
	}
	path := strings.Join(out, "/")

	u, err := url.ParseRequestURI(path)
	if err != nil {
		return nil, err
	}

	if tcurl != "" {
		tu, _ := url.Parse(tcurl)
		if tu != nil {
			u.Host = tu.Host
			u.Scheme = tu.Scheme
		}
	}

	return u, nil
}

var CodecTypes = flv.CodecTypes

func (conn *Conn) writeBasicConf() (err error) {
	// > WindowAckSize
	if err = conn.writeWindowAckSize(2500000, false); err != nil {
		return
	}
	// > SetPeerBandwidth
	if err = conn.writeSetPeerBandwidth(2500000, 2, true); err != nil {
		return
	}
	if conn.isserver {
		// > StreamBegin
		if err = conn.writeStreamBegin(0, true); err != nil {
			return
		}
	}
	// > SetChunkSize
	if err = conn.writeSetChunkSize(1024*64, true); err != nil {
		return
	}
	return
}

func (conn *Conn) readConnect() (err error) {
	var connectpath string

	// < connect("app")
	if err = conn.pollCommand(); err != nil {
		return
	}
	if conn.commandname != "connect" {
		err = fmt.Errorf("first command is not connect")
		return
	}
	if conn.commandobj == nil {
		err = fmt.Errorf("connect command params invalid")
		return
	}

	//fmt.Printf("readConnect: %+v\n", conn.commandobj)

	var ok bool
	var _app, _tcurl interface{}
	if _app, ok = conn.commandobj["app"]; !ok {
		err = fmt.Errorf("the `connect` params missing `app`")
		return
	}
	connectpath, _ = _app.(string)

	var tcurl string
	if _tcurl, ok = conn.commandobj["tcUrl"]; !ok {
		_tcurl, ok = conn.commandobj["tcurl"]
	}
	if ok {
		tcurl, _ = _tcurl.(string)
	}

	connectparams := conn.commandobj

	if err = conn.writeBasicConf(); err != nil {
		return
	}

	objectEncoding, ok := connectparams["objectEncoding"]
	if !ok {
		objectEncoding = 3
	}

	// > _result("NetConnection.Connect.Success")
	if err = conn.writeCommandMsg(false, 3, 0, "_result", conn.commandtransid,
		flvio.AMFMap{
			"fmsVer":       "FMS/3,0,1,123",
			"capabilities": 31,
		},
		flvio.AMFMap{
			"level":          "status",
			"code":           "NetConnection.Connect.Success",
			"description":    "Connection succeeded.",
			"objectEncoding": objectEncoding,
		},
	); err != nil {
		return
	}

	// > onBWDone()
	if err = conn.writeCommandMsg(true, 3, 0, "onBWDone", 0, nil, 8192); err != nil {
		return
	}

	if err = conn.flushWrite(); err != nil {
		return
	}

	for {
		if err = conn.pollMsg(); err != nil {
			return
		}
		if conn.gotcommand {
			switch conn.commandname {

			// < createStream
			case "createStream":
				conn.avmsgsid = uint32(1)
				// > _result(streamid)
				if err = conn.writeCommandMsg(false, 3, 0, "_result", conn.commandtransid, nil, conn.avmsgsid); err != nil {
					return
				}
				if err = conn.flushWrite(); err != nil {
					return
				}

			// < publish("path")
			case "publish":
				if Debug {
					fmt.Println("rtmp: < publish")
				}

				if len(conn.commandparams) < 1 {
					err = fmt.Errorf("publish params invalid")
					return
				}
				publishpath, _ := conn.commandparams[0].(string)

				var cberr error
				if conn.OnPlayOrPublish != nil {
					cberr = conn.OnPlayOrPublish(conn.commandname, connectparams)
				}

				// > onStatus()
				if err = conn.writeCommandMsg(false, 5, conn.avmsgsid,
					"onStatus", conn.commandtransid, nil,
					flvio.AMFMap{
						"level":       "status",
						"code":        "NetStream.Publish.Start",
						"description": "Start publishing",
					},
				); err != nil {
					return
				}
				if err = conn.flushWrite(); err != nil {
					return
				}

				if cberr != nil {
					err = fmt.Errorf("OnPlayOrPublish check failed")
					return
				}

				u, uerr := createURL(tcurl, connectpath, publishpath)
				if uerr != nil {
					err = fmt.Errorf("invalid URL: %w", uerr)
					return
				}

				conn.URL = u
				conn.publishing = true
				conn.reading = true
				conn.stage++
				return

			// < play("path")
			case "play":
				if Debug {
					fmt.Println("rtmp: < play")
				}

				if len(conn.commandparams) < 1 {
					err = fmt.Errorf("command play params invalid")
					return
				}
				playpath, _ := conn.commandparams[0].(string)

				// > streamBegin(streamid)
				if err = conn.writeStreamBegin(conn.avmsgsid, false); err != nil {
					return
				}

				// > onStatus()
				if err = conn.writeCommandMsg(false, 5, conn.avmsgsid,
					"onStatus", conn.commandtransid, nil,
					flvio.AMFMap{
						"level":       "status",
						"code":        "NetStream.Play.Start",
						"description": "Start live",
					},
				); err != nil {
					return
				}

				// > |RtmpSampleAccess()
				//if err = self.writeDataMsg(5, self.avmsgsid,
				//	"|RtmpSampleAccess", true, true,
				//); err != nil {
				//	return
				//}

				if err = conn.flushWrite(); err != nil {
					return
				}

				u, uerr := createURL(tcurl, connectpath, playpath)
				if uerr != nil {
					err = fmt.Errorf("invalid URL: %w", uerr)
					return
				}

				conn.URL = u
				conn.playing = true
				conn.writing = true
				conn.stage++
				return
			}

		}
	}
}

func (conn *Conn) checkConnectResult() (ok bool, errmsg string) {
	if len(conn.commandparams) < 1 {
		errmsg = "params length < 1"
		return
	}

	obj, _ := conn.commandparams[0].(flvio.AMFMap)
	if obj == nil {
		errmsg = "params[0] not object"
		return
	}

	_code := obj["code"]
	if _code == nil {
		errmsg = "code invalid"
		return
	}

	code, _ := _code.(string)
	if code != "NetConnection.Connect.Success" {
		errmsg = "code != NetConnection.Connect.Success"
		return
	}

	ok = true
	return
}

func (conn *Conn) checkCreateStreamResult() (ok bool, avmsgsid uint32) {
	if len(conn.commandparams) < 1 {
		return
	}

	ok = true
	_avmsgsid, _ := conn.commandparams[0].(float64)
	avmsgsid = uint32(_avmsgsid)
	return
}

func (conn *Conn) probe() (err error) {
	for !conn.prober.Probed() {
		var tag flvio.Tag
		if tag, err = conn.pollAVTag(); err != nil {
			return
		}
		if err = conn.prober.PushTag(tag, int32(conn.timestamp)); err != nil {
			if Debug {
				fmt.Printf("rtmp: error probing tag: %s\n", err.Error())
			}
		}
	}

	conn.streams = conn.prober.Streams
	conn.stage++
	return
}

func (conn *Conn) writeConnect(path string) (err error) {
	if err = conn.writeBasicConf(); err != nil {
		return
	}

	// > connect("app")
	if Debug {
		fmt.Printf("rtmp: > connect('%s') host=%s\n", path, conn.URL.Host)
	}
	if err = conn.writeCommandMsg(false, 3, 0, "connect", 1,
		flvio.AMFMap{
			"app":           path,
			"flashVer":      "MAC 22,0,0,192",
			"tcUrl":         getTcUrl(conn.URL),
			"fpad":          false,
			"capabilities":  15,
			"audioCodecs":   4071,
			"videoCodecs":   252,
			"videoFunction": 1,
			"fourCcList":    flvio.AMFArray{"av01", "vp09", "hvc1"},
		},
	); err != nil {
		return
	}

	if err = conn.flushWrite(); err != nil {
		return
	}

	for {
		if err = conn.pollMsg(); err != nil {
			return
		}
		if conn.gotcommand {
			// < _result("NetConnection.Connect.Success")
			if conn.commandname == "_result" {
				var ok bool
				var errmsg string
				if ok, errmsg = conn.checkConnectResult(); !ok {
					err = fmt.Errorf("command connect failed: %s", errmsg)
					return
				}
				if Debug {
					fmt.Printf("rtmp: < _result() of connect\n")
				}
				break
			}
		} else {
			if conn.msgtypeid == msgtypeidWindowAckSize {
				if len(conn.msgdata) == 4 {
					conn.readAckSize = pio.U32BE(conn.msgdata) >> 1
				}
				//if err = self.writeWindowAckSize(0xffffffff); err != nil {
				//	return
				//}
			}
		}
	}

	return
}

func (conn *Conn) connectPublish() (err error) {
	connectpath, publishpath := SplitPath(conn.URL)

	if err = conn.writeConnect(connectpath); err != nil {
		return
	}

	transid := 2

	// > createStream()
	if Debug {
		fmt.Printf("rtmp: > createStream()\n")
	}
	if err = conn.writeCommandMsg(false, 3, 0, "createStream", transid, nil); err != nil {
		return
	}
	transid++

	if err = conn.flushWrite(); err != nil {
		return
	}

	for {
		if err = conn.pollMsg(); err != nil {
			return
		}
		if conn.gotcommand {
			// < _result(avmsgsid) of createStream
			if conn.commandname == "_result" {
				var ok bool
				if ok, conn.avmsgsid = conn.checkCreateStreamResult(); !ok {
					err = fmt.Errorf("createStream command failed")
					return
				}
				break
			}
		}
	}

	// > publish('app')
	if Debug {
		fmt.Printf("rtmp: > publish('%s')\n", publishpath)
	}
	if err = conn.writeCommandMsg(false, 8, conn.avmsgsid, "publish", transid, nil, publishpath); err != nil {
		return
	}
	transid++

	if err = conn.flushWrite(); err != nil {
		return
	}

	conn.writing = true
	conn.publishing = true
	conn.stage++
	return
}

func (conn *Conn) connectPlay() (err error) {
	connectpath, playpath := SplitPath(conn.URL)

	if err = conn.writeConnect(connectpath); err != nil {
		return
	}

	// > createStream()
	if Debug {
		fmt.Printf("rtmp: > createStream()\n")
	}
	if err = conn.writeCommandMsg(false, 3, 0, "createStream", 2, nil); err != nil {
		return
	}

	// > SetBufferLength 0,100ms
	if err = conn.writeSetBufferLength(0, 100, false); err != nil {
		return
	}

	if err = conn.flushWrite(); err != nil {
		return
	}

	for {
		if err = conn.pollMsg(); err != nil {
			return
		}
		if conn.gotcommand {
			// < _result(avmsgsid) of createStream
			if conn.commandname == "_result" {
				var ok bool
				if ok, conn.avmsgsid = conn.checkCreateStreamResult(); !ok {
					err = fmt.Errorf("createStream command failed")
					return
				}
				break
			}
		}
	}

	// > play('app')
	if Debug {
		fmt.Printf("rtmp: > play('%s')\n", playpath)
	}
	if err = conn.writeCommandMsg(false, 8, conn.avmsgsid, "play", 0, nil, playpath); err != nil {
		return
	}
	if err = conn.flushWrite(); err != nil {
		return
	}

	conn.reading = true
	conn.playing = true
	conn.stage++
	return
}

func (conn *Conn) ReadPacket() (pkt av.Packet, err error) {
	if err = conn.prepare(stageCodecDataDone, prepareReading); err != nil {
		return
	}

	if !conn.prober.Empty() {
		pkt = conn.prober.PopPacket()
		return
	}

	for {
		var tag flvio.Tag
		if tag, err = conn.pollAVTag(); err != nil {
			return
		}

		var ok bool
		if pkt, ok = conn.prober.TagToPacket(tag, int32(conn.timestamp)); ok {
			return pkt, nil
		}
	}
}

func (conn *Conn) Prepare() (err error) {
	return conn.prepare(stageCommandDone, 0)
}

func (conn *Conn) prepare(stage int, flags int) (err error) {
	for conn.stage < stage {
		switch conn.stage {
		case 0:
			if conn.isserver {
				if err = conn.handshakeServer(); err != nil {
					return
				}
			} else {
				if err = conn.handshakeClient(); err != nil {
					return
				}
			}

		case stageHandshakeDone:
			if conn.isserver {
				if err = conn.readConnect(); err != nil {
					return
				}
			} else {
				if flags == prepareReading {
					if err = conn.connectPlay(); err != nil {
						return
					}
				} else {
					if err = conn.connectPublish(); err != nil {
						return
					}
				}
			}

		case stageCommandDone:
			if flags == prepareReading {
				if err = conn.probe(); err != nil {
					return
				}
			} else {
				err = fmt.Errorf("call WriteHeader() before WritePacket()")
				return
			}
		}
	}
	return
}

func (conn *Conn) Streams() (streams []av.CodecData, err error) {
	if err = conn.prepare(stageCodecDataDone, prepareReading); err != nil {
		return
	}
	streams = conn.streams
	return
}

func (conn *Conn) WritePacket(pkt av.Packet) (err error) {
	if err = conn.prepare(stageCodecDataDone, prepareWriting); err != nil {
		return
	}

	stream := conn.streams[pkt.Idx]
	tag, timestamp := flv.PacketToTag(pkt, stream)

	if Debug {
		fmt.Println("rtmp: WritePacket", pkt.Idx, pkt.Time, pkt.CompositionTime)
	}

	if err = conn.writeAVTag(tag, int32(timestamp)); err != nil {
		return
	}

	return
}

func (conn *Conn) WriteTrailer() (err error) {
	if err = conn.flushWrite(); err != nil {
		return
	}
	return
}

func (conn *Conn) SetMetaData(data flvio.AMFMap) {
	conn.metadata = data
}

func (conn *Conn) GetMetaData() flvio.AMFMap {
	return conn.metadata
}

func (conn *Conn) WriteHeader(streams []av.CodecData) (err error) {
	if err = conn.prepare(stageCommandDone, prepareWriting); err != nil {
		return
	}

	var metadata flvio.AMFMap

	if metadata, err = flv.NewMetadataByStreams(streams); err != nil {
		return
	}

	if conn.metadata != nil {
		for key, value := range conn.metadata {
			if _, hasKey := metadata[key]; hasKey {
				continue
			}

			metadata[key] = value
		}
	}

	// > onMetaData()
	if err = conn.writeDataMsg(5, conn.avmsgsid, "onMetaData", metadata); err != nil {
		return
	}

	// > Videodata(decoder config)
	// > Audiodata(decoder config)
	for _, stream := range streams {
		var ok bool
		var tag flvio.Tag
		if tag, ok, err = flv.CodecDataToTag(stream); err != nil {
			return
		}
		if ok {
			if err = conn.writeAVTag(tag, 0); err != nil {
				return
			}
		}
	}

	conn.streams = streams
	conn.stage++
	return
}

func (conn *Conn) tmpwbuf(n int) []byte {
	if len(conn.writebuf) < n {
		conn.writebuf = make([]byte, n)
	}
	return conn.writebuf
}

func (conn *Conn) writeSetChunkSize(size int, append bool) (err error) {
	conn.writeMaxChunkSize = size
	var b []byte
	var n int

	b = conn.tmpwbuf(chunkHeaderLength + 4)
	n = conn.fillChunkHeader(append, b, 2, 0, msgtypeidSetChunkSize, 0, 4)
	pio.PutU32BE(b[n:], uint32(size))
	n += 4
	_, err = conn.bufw.Write(b[:n])
	return
}

func (conn *Conn) writeAck(seqnum uint32, append bool) (err error) {
	var b []byte
	var n int

	b = conn.tmpwbuf(chunkHeaderLength + 4)
	n = conn.fillChunkHeader(append, b, 2, 0, msgtypeidAck, 0, 4)
	pio.PutU32BE(b[n:], seqnum)
	n += 4
	_, err = conn.bufw.Write(b[:n])
	return
}

func (conn *Conn) writeWindowAckSize(size uint32, append bool) (err error) {
	var b []byte
	var n int

	b = conn.tmpwbuf(chunkHeaderLength + 4)
	n = conn.fillChunkHeader(append, b, 2, 0, msgtypeidWindowAckSize, 0, 4)
	pio.PutU32BE(b[n:], size)
	n += 4
	_, err = conn.bufw.Write(b[:n])
	return
}

func (conn *Conn) writeSetPeerBandwidth(acksize uint32, limittype uint8, append bool) (err error) {
	var b []byte
	var n int

	b = conn.tmpwbuf(chunkHeaderLength + 5)
	n = conn.fillChunkHeader(append, b, 2, 0, msgtypeidSetPeerBandwidth, 0, 5)
	pio.PutU32BE(b[n:], acksize)
	n += 4
	b[n] = limittype
	n++
	_, err = conn.bufw.Write(b[:n])
	return
}

func (conn *Conn) writeCommandMsg(append bool, csid, msgsid uint32, args ...interface{}) (err error) {
	return conn.writeAMF0Msg(append, msgtypeidCommandMsgAMF0, csid, msgsid, args...)
}

func (conn *Conn) writeDataMsg(csid, msgsid uint32, args ...interface{}) (err error) {
	return conn.writeAMF0Msg(false, msgtypeidDataMsgAMF0, csid, msgsid, args...)
}

func (conn *Conn) writeAMF0Msg(append bool, msgtypeid uint8, csid, msgsid uint32, args ...interface{}) (err error) {
	size := 0
	for _, arg := range args {
		size += flvio.LenAMF0Val(arg)
	}

	b := conn.tmpwbuf(chunkHeaderLength + size)
	n := conn.fillChunkHeader(append, b, csid, 0, msgtypeid, msgsid, size)
	for _, arg := range args {
		n += flvio.FillAMF0Val(b[n:], arg)
	}

	_, err = conn.bufw.Write(b[:n])
	return
}

func (conn *Conn) writeAVTag(tag flvio.Tag, ts int32) (err error) {
	var msgtypeid uint8
	var csid uint32
	var data []byte

	switch tag.Type {
	case flvio.TAG_AUDIO:
		msgtypeid = msgtypeidAudioMsg
		csid = 6
		data = tag.Data

	case flvio.TAG_VIDEO:
		msgtypeid = msgtypeidVideoMsg
		csid = 7
		data = tag.Data
	}

	actualChunkHeaderLength := chunkHeaderLength
	if uint32(ts) > FlvTimestampMax {
		actualChunkHeaderLength += 4
	}

	b := conn.tmpwbuf(actualChunkHeaderLength + flvio.MaxTagSubHeaderLength)
	hdrlen := tag.FillHeader(b[actualChunkHeaderLength:])
	conn.fillChunkHeader(false, b, csid, ts, msgtypeid, conn.avmsgsid, hdrlen+len(data))
	n := actualChunkHeaderLength + hdrlen

	if _, err = conn.bufw.Write(b[:n]); err != nil {
		return
	}

	chunksize := conn.writeMaxChunkSize
	msgdataleft := len(data)

	n = msgdataleft
	if msgdataleft > chunksize-hdrlen {
		n = chunksize - hdrlen
	}

	if _, err = conn.bufw.Write(data[:n]); err != nil {
		return
	}

	data = data[n:]
	msgdataleft -= n

	for {
		if msgdataleft == 0 {
			break
		}

		n = conn.fillChunkHeader3(b, csid, ts)

		if _, err = conn.bufw.Write(b[:n]); err != nil {
			return
		}

		n = msgdataleft
		if msgdataleft > chunksize {
			n = chunksize
		}

		if _, err = conn.bufw.Write(data[:n]); err != nil {
			return
		}

		data = data[n:]
		msgdataleft -= n
	}

	err = conn.bufw.Flush()

	return
}

func (conn *Conn) writeStreamBegin(msgsid uint32, append bool) (err error) {
	b := conn.tmpwbuf(chunkHeaderLength + 6)
	n := conn.fillChunkHeader(append, b, 2, 0, msgtypeidUserControl, 0, 6)
	pio.PutU16BE(b[n:], eventtypeStreamBegin)
	n += 2
	pio.PutU32BE(b[n:], msgsid)
	n += 4
	_, err = conn.bufw.Write(b[:n])
	return
}

func (conn *Conn) writeSetBufferLength(msgsid uint32, timestamp uint32, append bool) (err error) {
	b := conn.tmpwbuf(chunkHeaderLength + 10)
	n := conn.fillChunkHeader(append, b, 2, 0, msgtypeidUserControl, 0, 10)
	pio.PutU16BE(b[n:], eventtypeSetBufferLength)
	n += 2
	pio.PutU32BE(b[n:], msgsid)
	n += 4
	pio.PutU32BE(b[n:], timestamp)
	n += 4
	_, err = conn.bufw.Write(b[:n])
	return
}

func (conn *Conn) writePingResponse(timestamp uint32, append bool) (err error) {
	b := conn.tmpwbuf(chunkHeaderLength + 10)
	n := conn.fillChunkHeader(append, b, 2, 0, msgtypeidUserControl, 0, 6)
	pio.PutU16BE(b[n:], eventtypePingResponse)
	n += 2
	pio.PutU32BE(b[n:], timestamp)
	n += 4
	_, err = conn.bufw.Write(b[:n])
	return
}

const chunkHeaderLength = 12
const FlvTimestampMax = 0xFFFFFF

func (conn *Conn) fillChunkHeader(append bool, b []byte, csid uint32, timestamp int32, msgtypeid uint8, msgsid uint32, msgdatalen int) (n int) {
	if !append {
		//  0                   1                   2                   3
		//  0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
		// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
		// |                   timestamp                   |message length |
		// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
		// |     message length (cont)     |message type id| msg stream id |
		// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
		// |           message stream id (cont)            |
		// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
		//
		//       Figure 9 Chunk Message Header – Type 0

		b[n] = byte(csid) & 0x3f
		n++
		if uint32(timestamp) <= FlvTimestampMax {
			pio.PutU24BE(b[n:], uint32(timestamp))
		} else {
			pio.PutU24BE(b[n:], FlvTimestampMax)
		}
		n += 3
		pio.PutU24BE(b[n:], uint32(msgdatalen))
		n += 3
		b[n] = msgtypeid
		n++
		pio.PutU32LE(b[n:], msgsid)
		n += 4
		if uint32(timestamp) > FlvTimestampMax {
			pio.PutU32BE(b[n:], uint32(timestamp))
			n += 4
		}
	} else {
		//  0                   1                   2                   3
		//  0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
		// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
		// |                timestamp delta                |message length |
		// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
		// |     message length (cont)     |message type id|
		// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
		//
		//       Figure 10 Chunk Message Header – Type 1

		b[n] = 1 << 6
		b[n] += byte(csid) & 0x3f
		n++
		pio.PutU24BE(b[n:], 0)
		n += 3
		pio.PutU24BE(b[n:], uint32(msgdatalen))
		n += 3
		b[n] = msgtypeid
		n++
	}

	if Debug {
		fmt.Printf("rtmp: write chunk msgdatalen=%d msgsid=%d\n", msgdatalen, msgsid)
		fmt.Print(hex.Dump(b[:n]))
	}

	return
}

func (c *Conn) fillChunkHeader3(b []byte, csid uint32, timestamp int32) (n int) {
	pio.PutU8(b, (uint8(csid)&0x3f)|3<<6)
	n++
	if uint32(timestamp) >= FlvTimestampMax {
		pio.PutU32BE(b[n:], uint32(timestamp))
		n += 4
	}

	return
}

func (conn *Conn) flushWrite() (err error) {
	if err = conn.bufw.Flush(); err != nil {
		return
	}
	return
}

func (conn *Conn) readChunk() (err error) {
	b := conn.readbuf
	n := 0
	if _, err = io.ReadFull(conn.bufr, b[:1]); err != nil {
		return
	}
	header := b[0]
	n += 1

	var msghdrtype uint8
	var csid uint32

	msghdrtype = header >> 6

	csid = uint32(header) & 0x3f
	switch csid {
	default: // Chunk basic header 1
	case 0: // Chunk basic header 2
		if _, err = io.ReadFull(conn.bufr, b[:1]); err != nil {
			return fmt.Errorf("chunk basic header 2: %w", err)
		}
		n += 1
		csid = uint32(b[0]) + 64
	case 1: // Chunk basic header 3
		if _, err = io.ReadFull(conn.bufr, b[:2]); err != nil {
			return fmt.Errorf("chunk basic header 3: %w", err)
		}
		n += 2
		csid = uint32(pio.U16BE(b)) + 64
	}

	newcs := false
	cs := conn.readcsmap[csid]
	if cs == nil {
		cs = &chunkStream{}
		conn.readcsmap[csid] = cs
		newcs = true
	}

	if len(conn.readcsmap) > 16 {
		err = fmt.Errorf("too many chunk stream ids")
		return
	}

	var timestamp uint32

	switch msghdrtype {
	case 0:
		//  0                   1                   2                   3
		//  0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
		// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
		// |                   timestamp                   |message length |
		// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
		// |     message length (cont)     |message type id| msg stream id |
		// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
		// |           message stream id (cont)            |
		// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
		//
		//       Figure 9 Chunk Message Header – Type 0
		if cs.msgdataleft != 0 {
			if !conn.skipInvalidMessages {
				err = fmt.Errorf("chunk msgdataleft=%d invalid", cs.msgdataleft)
				return
			}
		}
		h := b[:11]
		if _, err = io.ReadFull(conn.bufr, h); err != nil {
			return
		}
		n += len(h)
		timestamp = pio.U24BE(h[0:3])
		cs.msghdrtype = msghdrtype
		cs.msgdatalen = pio.U24BE(h[3:6])
		cs.msgtypeid = h[6]
		cs.msgsid = pio.U32LE(h[7:11])
		if timestamp == FlvTimestampMax {
			if _, err = io.ReadFull(conn.bufr, b[:4]); err != nil {
				return
			}
			n += 4
			timestamp = pio.U32BE(b)
			cs.hastimeext = true
			cs.timeext = timestamp
		} else {
			cs.hastimeext = false
		}
		cs.timenow = timestamp
		cs.Start()

	case 1:
		//  0                   1                   2                   3
		//  0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
		// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
		// |                timestamp delta                |message length |
		// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
		// |     message length (cont)     |message type id|
		// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
		//
		//       Figure 10 Chunk Message Header – Type 1
		if newcs {
			err = fmt.Errorf("chunk message type 1 without previous chunk")
			return
		}
		if cs.msgdataleft != 0 {
			if !conn.skipInvalidMessages {
				err = fmt.Errorf("chunk msgdataleft=%d invalid", cs.msgdataleft)
				return
			}
		}
		h := b[:7]
		if _, err = io.ReadFull(conn.bufr, h); err != nil {
			return
		}
		n += len(h)
		timestamp = pio.U24BE(h[0:3])
		cs.msghdrtype = msghdrtype
		cs.msgdatalen = pio.U24BE(h[3:6])
		cs.msgtypeid = h[6]
		if timestamp == FlvTimestampMax {
			if _, err = io.ReadFull(conn.bufr, b[:4]); err != nil {
				return
			}
			n += 4
			timestamp = pio.U32BE(b)
			cs.hastimeext = true
			cs.timeext = timestamp
		} else {
			cs.hastimeext = false
		}
		cs.timedelta = timestamp
		cs.timenow += timestamp
		cs.Start()

	case 2:
		//  0                   1                   2
		//  0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3
		// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
		// |                timestamp delta                |
		// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
		//
		//       Figure 11 Chunk Message Header – Type 2
		if newcs {
			err = fmt.Errorf("chunk message type 2 without previous chunk")
			return
		}
		if cs.msgdataleft != 0 {
			if !conn.skipInvalidMessages {
				err = fmt.Errorf("chunk msgdataleft=%d invalid", cs.msgdataleft)
				return
			}
		}
		h := b[:3]
		if _, err = io.ReadFull(conn.bufr, h); err != nil {
			return
		}
		n += len(h)
		cs.msghdrtype = msghdrtype
		timestamp = pio.U24BE(h[0:3])
		if timestamp == FlvTimestampMax {
			if _, err = io.ReadFull(conn.bufr, b[:4]); err != nil {
				return
			}
			n += 4
			timestamp = pio.U32BE(b)
			cs.hastimeext = true
			cs.timeext = timestamp
		} else {
			cs.hastimeext = false
		}
		cs.timedelta = timestamp
		cs.timenow += timestamp
		cs.Start()

	case 3:
		if newcs {
			err = fmt.Errorf("chunk message type 3 without previous chunk")
			return
		}

		if cs.msgdataleft == 0 {
			switch cs.msghdrtype {
			case 0:
				if cs.hastimeext {
					if _, err = io.ReadFull(conn.bufr, b[:4]); err != nil {
						return
					}
					n += 4
					timestamp = pio.U32BE(b)
					cs.timenow = timestamp
					cs.timeext = timestamp
				}
			case 1, 2:
				if cs.hastimeext {
					if _, err = io.ReadFull(conn.bufr, b[:4]); err != nil {
						return
					}
					n += 4
					timestamp = pio.U32BE(b)
				} else {
					timestamp = cs.timedelta
				}
				cs.timenow += timestamp
			}
			cs.Start()
		} else {
			if cs.hastimeext {
				var b []byte
				if b, err = conn.bufr.Peek(4); err != nil {
					return
				}
				if pio.U32BE(b) == cs.timeext {
					if _, err = io.ReadFull(conn.bufr, b[:4]); err != nil {
						return
					}
				}
			}
		}

	default:
		err = fmt.Errorf("invalid chunk msg header type=%d", msghdrtype)
		return
	}

	size := int(cs.msgdataleft)
	if size > conn.readMaxChunkSize {
		size = conn.readMaxChunkSize
	}
	off := cs.msgdatalen - cs.msgdataleft
	buf := cs.msgdata[off : int(off)+size]
	if _, err = io.ReadFull(conn.bufr, buf); err != nil {
		return
	}
	n += len(buf)
	cs.msgdataleft -= uint32(size)

	if Debug {
		fmt.Printf("rtmp: chunk msgsid=%d msgtypeid=%d msghdrtype=%d len=%d left=%d max=%d",
			cs.msgsid, cs.msgtypeid, msghdrtype, cs.msgdatalen, cs.msgdataleft, conn.readMaxChunkSize)
	}

	if conn.debugChunks {
		path := "***"
		if conn.URL != nil {
			path = conn.URL.Path
		}
		data := fmt.Sprintf("%s: chunk id=%d msgsid=%d msgtypeid=%d msghdrtype=%d timestamp=%d ext=%v wallclock=%v len=%d left=%d max=%d",
			path, csid, cs.msgsid, cs.msgtypeid, msghdrtype, cs.timenow, cs.hastimeext, cs.genwallclocktime, cs.msgdatalen, cs.msgdataleft, conn.readMaxChunkSize)

		if cs.msgtypeid != msgtypeidVideoMsg && cs.msgtypeid != msgtypeidAudioMsg {
			if len(cs.msgdata) > 1024 {
				data += " data=" + hex.EncodeToString(cs.msgdata[:1024]) + "... "
			} else {
				data += " data=" + hex.EncodeToString(cs.msgdata) + " "
			}
		} else {
			data += " data= "
		}

		data += fmt.Sprintf("(%d bytes)", len(cs.msgdata))

		fmt.Printf("%s\n", data)
	}

	if cs.msgdataleft == 0 {
		if Debug {
			fmt.Println("rtmp: chunk data")
			fmt.Print(hex.Dump(cs.msgdata))
		}

		timestamp = cs.timenow

		if cs.msgtypeid == msgtypeidVideoMsg || cs.msgtypeid == msgtypeidAudioMsg {
			if !cs.genwallclocktime {
				if cs.prevtimenow >= cs.timenow {
					cs.sametscount++
				} else {
					cs.sametscount = 0
				}

				// if the previous timestamp is the same as the current for too often in a row, assume defect timestamps
				if cs.sametscount > 10 {
					if cs.msgcount < 20 { // only consider the first video and audio messages
						cs.genwallclocktime = true
					} else { // otherwise bail out
						err = fmt.Errorf("detected sequence of non-changing timestamps: %d (msgtypeid %d)", cs.timenow, cs.msgtypeid)
						return
					}
				}

				cs.prevtimenow = cs.timenow
				cs.msgcount++
			}

			if cs.genwallclocktime {
				timestamp = uint32(time.Since(conn.start).Milliseconds() % 0xFFFFFFFF)
			}
		}

		if err = conn.handleMsg(timestamp, cs.msgsid, cs.msgtypeid, cs.msgdata); err != nil {
			return fmt.Errorf("handleMsg: %w", err)
		}

		cs.msgdata = nil
	}

	conn.ackn += uint32(n)

	return
}

func (conn *Conn) handleCommandMsgAMF0(b []byte) (n int, err error) {
	var name, transid, obj interface{}
	var size int

	if name, size, err = flvio.ParseAMF0Val(b[n:]); err != nil {
		return
	}
	n += size
	if transid, size, err = flvio.ParseAMF0Val(b[n:]); err != nil {
		return
	}
	n += size
	if obj, size, err = flvio.ParseAMF0Val(b[n:]); err != nil {
		return
	}
	n += size

	var ok bool
	if conn.commandname, ok = name.(string); !ok {
		err = fmt.Errorf("CommandMsgAMF0 command is not string")
		return
	}
	conn.commandtransid, _ = transid.(float64)
	conn.commandobj, _ = obj.(flvio.AMFMap)
	conn.commandparams = []interface{}{}

	for n < len(b) {
		if obj, size, err = flvio.ParseAMF0Val(b[n:]); err != nil {
			return
		}
		n += size
		conn.commandparams = append(conn.commandparams, obj)
	}
	if n < len(b) {
		err = fmt.Errorf("CommandMsgAMF0 left bytes=%d", len(b)-n)
		return
	}

	conn.gotcommand = true
	return
}

func (conn *Conn) handleMsg(timestamp uint32, msgsid uint32, msgtypeid uint8, msgdata []byte) (err error) {
	conn.msgdata = msgdata
	conn.msgtypeid = msgtypeid
	conn.timestamp = timestamp

	switch msgtypeid {
	case msgtypeidCommandMsgAMF0:
		if _, err = conn.handleCommandMsgAMF0(msgdata); err != nil {
			err = fmt.Errorf("AMF0: %w", err)
			return
		}

	case msgtypeidCommandMsgAMF3:
		if len(msgdata) < 1 {
			err = fmt.Errorf("short packet of CommandMsgAMF3")
			return
		}
		// skip first byte
		if _, err = conn.handleCommandMsgAMF0(msgdata[1:]); err != nil {
			err = fmt.Errorf("AMF3: %w", err)
			return
		}

	case msgtypeidUserControl:
		if len(msgdata) < 2 {
			err = fmt.Errorf("short packet of UserControl")
			return
		}
		conn.eventtype = pio.U16BE(msgdata)

		if conn.eventtype == eventtypePingRequest {
			if len(msgdata) != 6 {
				err = fmt.Errorf("wrong length for UserControl.PingRequest")
				return
			}
			pingtimestamp := pio.U32BE(msgdata[2:])
			conn.writePingResponse(pingtimestamp, false)
			conn.flushWrite()
		}

	case msgtypeidDataMsgAMF0:
		b := msgdata
		n := 0
		for n < len(b) {
			var obj interface{}
			var size int
			if obj, size, err = flvio.ParseAMF0Val(b[n:]); err != nil {
				return
			}
			n += size
			conn.datamsgvals = append(conn.datamsgvals, obj)
		}
		if n < len(b) {
			err = fmt.Errorf("DataMsgAMF0 left bytes=%d", len(b)-n)
			return
		}

		// store metadata

		metaindex := -1

		for i, x := range conn.datamsgvals {
			switch x := x.(type) {
			case string:
				if x == "onMetaData" {
					metaindex = i + 1
				}
			}
		}

		if metaindex != -1 && metaindex < len(conn.datamsgvals) {
			metadata, ok := conn.datamsgvals[metaindex].(flvio.AMFMap)
			if ok {
				conn.metadata = metadata

				if _, hasVideo := metadata["videocodecid"]; hasVideo {
					conn.prober.HasVideo = true
				}

				if _, hasAudio := metadata["audiocodecid"]; hasAudio {
					conn.prober.HasAudio = true
				}
			}
		}

	case msgtypeidVideoMsg:
		if len(msgdata) == 0 {
			return
		}
		//fmt.Printf("msgdata: %#08x\n", msgdata[:5])
		tag := flvio.Tag{Type: flvio.TAG_VIDEO}
		var n int
		if n, err = (&tag).ParseHeader(msgdata); err != nil {
			return
		}
		//fmt.Printf("tag: %+v\n", tag)
		if !(tag.FrameType == flvio.FRAME_INTER || tag.FrameType == flvio.FRAME_KEY) {
			return
		}
		tag.Data = msgdata[n:]
		conn.avtag = tag

	case msgtypeidAudioMsg:
		if len(msgdata) == 0 {
			return
		}
		tag := flvio.Tag{Type: flvio.TAG_AUDIO}
		var n int
		if n, err = (&tag).ParseHeader(msgdata); err != nil {
			return
		}
		tag.Data = msgdata[n:]
		conn.avtag = tag

	case msgtypeidSetChunkSize:
		if len(msgdata) < 4 {
			err = fmt.Errorf("short packet of SetChunkSize")
			return
		}
		conn.readMaxChunkSize = int(pio.U32BE(msgdata))
	case msgtypeidWindowAckSize:
		if len(conn.msgdata) != 4 {
			return fmt.Errorf("invalid packet of WindowAckSize")
		}
		conn.readAckSize = pio.U32BE(conn.msgdata) >> 1
	default:
		if Debug {
			fmt.Printf("rtmp: unhandled msg: %d\n", msgtypeid)
		}
	}

	conn.gotmsg = true
	return
}

var (
	hsClientFullKey = []byte{
		'G', 'e', 'n', 'u', 'i', 'n', 'e', ' ', 'A', 'd', 'o', 'b', 'e', ' ',
		'F', 'l', 'a', 's', 'h', ' ', 'P', 'l', 'a', 'y', 'e', 'r', ' ',
		'0', '0', '1',
		0xF0, 0xEE, 0xC2, 0x4A, 0x80, 0x68, 0xBE, 0xE8, 0x2E, 0x00, 0xD0, 0xD1,
		0x02, 0x9E, 0x7E, 0x57, 0x6E, 0xEC, 0x5D, 0x2D, 0x29, 0x80, 0x6F, 0xAB,
		0x93, 0xB8, 0xE6, 0x36, 0xCF, 0xEB, 0x31, 0xAE,
	}
	hsServerFullKey = []byte{
		'G', 'e', 'n', 'u', 'i', 'n', 'e', ' ', 'A', 'd', 'o', 'b', 'e', ' ',
		'F', 'l', 'a', 's', 'h', ' ', 'M', 'e', 'd', 'i', 'a', ' ',
		'S', 'e', 'r', 'v', 'e', 'r', ' ',
		'0', '0', '1',
		0xF0, 0xEE, 0xC2, 0x4A, 0x80, 0x68, 0xBE, 0xE8, 0x2E, 0x00, 0xD0, 0xD1,
		0x02, 0x9E, 0x7E, 0x57, 0x6E, 0xEC, 0x5D, 0x2D, 0x29, 0x80, 0x6F, 0xAB,
		0x93, 0xB8, 0xE6, 0x36, 0xCF, 0xEB, 0x31, 0xAE,
	}
	hsClientPartialKey = hsClientFullKey[:30]
	hsServerPartialKey = hsServerFullKey[:36]
)

func hsMakeDigest(key []byte, src []byte, gap int) (dst []byte) {
	h := hmac.New(sha256.New, key)
	if gap <= 0 {
		h.Write(src)
	} else {
		h.Write(src[:gap])
		h.Write(src[gap+32:])
	}
	return h.Sum(nil)
}

func hsCalcDigestPos(p []byte, base int) (pos int) {
	for i := 0; i < 4; i++ {
		pos += int(p[base+i])
	}
	pos = (pos % 728) + base + 4
	return
}

func hsFindDigest(p []byte, key []byte, base int) int {
	gap := hsCalcDigestPos(p, base)
	digest := hsMakeDigest(key, p, gap)
	if !bytes.Equal(p[gap:gap+32], digest) {
		return -1
	}
	return gap
}

func hsParse1(p []byte, peerkey []byte, key []byte) (ok bool, digest []byte) {
	var pos int
	if pos = hsFindDigest(p, peerkey, 772); pos == -1 {
		if pos = hsFindDigest(p, peerkey, 8); pos == -1 {
			return
		}
	}
	ok = true
	digest = hsMakeDigest(key, p[pos:pos+32], -1)
	return
}

func hsCreate01(p []byte, time uint32, ver uint32, key []byte) {
	p[0] = 3
	p1 := p[1:]
	rand.Read(p1[8:])
	pio.PutU32BE(p1[0:4], time)
	pio.PutU32BE(p1[4:8], ver)
	gap := hsCalcDigestPos(p1, 8)
	digest := hsMakeDigest(key, p1, gap)
	copy(p1[gap:], digest)
}

func hsCreate2(p []byte, key []byte) {
	rand.Read(p)
	gap := len(p) - 32
	digest := hsMakeDigest(key, p, gap)
	copy(p[gap:], digest)
}

func (conn *Conn) handshakeClient() (err error) {
	var random [(1 + 1536*2) * 2]byte

	C0C1C2 := random[:1536*2+1]
	C0 := C0C1C2[:1]
	//C1 := C0C1C2[1:1536+1]
	C0C1 := C0C1C2[:1536+1]
	var C2 []byte

	S0S1S2 := random[1536*2+1:]
	//S0 := S0S1S2[:1]
	S1 := S0S1S2[1 : 1536+1]
	//S0S1 := S0S1S2[:1536+1]
	//S2 := S0S1S2[1536+1:]

	C0[0] = 3
	//hsCreate01(C0C1, hsClientFullKey)

	// > C0C1
	if _, err = conn.bufw.Write(C0C1); err != nil {
		return
	}
	if err = conn.bufw.Flush(); err != nil {
		return
	}

	// < S0S1S2
	if _, err = io.ReadFull(conn.bufr, S0S1S2); err != nil {
		return
	}

	if Debug {
		fmt.Println("handshakeClient: server version", S1[4], S1[5], S1[6], S1[7])
	}

	if ver := pio.U32BE(S1[4:8]); ver != 0 {
		C2 = S1
	} else {
		C2 = S1
	}

	// > C2
	if _, err = conn.bufw.Write(C2); err != nil {
		return
	}

	conn.stage++
	return
}

func (conn *Conn) handshakeServer() (err error) {
	var random [(1 + 1536*2) * 2]byte

	C0C1C2 := random[:1536*2+1]
	C0 := C0C1C2[:1]
	C1 := C0C1C2[1 : 1536+1]
	C0C1 := C0C1C2[:1536+1]
	C2 := C0C1C2[1536+1:]

	S0S1S2 := random[1536*2+1:]
	S0 := S0S1S2[:1]
	S1 := S0S1S2[1 : 1536+1]
	S0S1 := S0S1S2[:1536+1]
	S2 := S0S1S2[1536+1:]

	// < C0C1
	if _, err = io.ReadFull(conn.bufr, C0C1); err != nil {
		return
	}
	if C0[0] != 3 {
		err = fmt.Errorf("handshake version=%d invalid", C0[0])
		return
	}

	S0[0] = 3

	clitime := pio.U32BE(C1[0:4])
	srvtime := clitime
	srvver := uint32(0x0d0e0a0d)
	cliver := pio.U32BE(C1[4:8])

	if cliver != 0 {
		var ok bool
		var digest []byte
		if ok, digest = hsParse1(C1, hsClientPartialKey, hsServerFullKey); !ok {
			err = fmt.Errorf("handshake server: C1 invalid")
			return
		}
		hsCreate01(S0S1, srvtime, srvver, hsServerPartialKey)
		hsCreate2(S2, digest)
	} else {
		copy(S1, C1)
		copy(S2, C2)
	}

	// > S0S1S2
	if _, err = conn.bufw.Write(S0S1S2); err != nil {
		return
	}
	if err = conn.bufw.Flush(); err != nil {
		return
	}

	// < C2
	if _, err = io.ReadFull(conn.bufr, C2); err != nil {
		return
	}

	conn.stage++
	return
}

type closeConn struct {
	*Conn
	waitclose chan bool
}

func (cc closeConn) Close() error {
	cc.waitclose <- true
	return nil
}

func Handler(h *avutil.RegisterHandler) {
	h.UrlDemuxer = func(uri string) (ok bool, demuxer av.DemuxCloser, err error) {
		if !strings.HasPrefix(uri, "rtmp://") {
			return
		}
		ok = true
		demuxer, err = Dial(uri, DialOptions{})
		return
	}

	h.UrlMuxer = func(uri string) (ok bool, muxer av.MuxCloser, err error) {
		if !strings.HasPrefix(uri, "rtmp://") {
			return
		}
		ok = true
		muxer, err = Dial(uri, DialOptions{})
		return
	}

	h.ServerMuxer = func(uri string) (ok bool, muxer av.MuxCloser, err error) {
		if !strings.HasPrefix(uri, "rtmp://") {
			return
		}
		ok = true

		var u *url.URL
		if u, err = ParseURL(uri); err != nil {
			return
		}
		server := &Server{
			Addr: u.Host,
		}

		waitstart := make(chan error)
		waitconn := make(chan *Conn)
		waitclose := make(chan bool)

		server.HandlePlay = func(conn *Conn) {
			waitconn <- conn
			<-waitclose
		}

		go func() {
			waitstart <- server.ListenAndServe()
		}()

		select {
		case err = <-waitstart:
			if err != nil {
				return
			}

		case conn := <-waitconn:
			muxer = closeConn{Conn: conn, waitclose: waitclose}
			return
		}

		return
	}

	h.ServerDemuxer = func(uri string) (ok bool, demuxer av.DemuxCloser, err error) {
		if !strings.HasPrefix(uri, "rtmp://") {
			return
		}
		ok = true

		var u *url.URL
		if u, err = ParseURL(uri); err != nil {
			return
		}
		server := &Server{
			Addr: u.Host,
		}

		waitstart := make(chan error)
		waitconn := make(chan *Conn)
		waitclose := make(chan bool)

		server.HandlePublish = func(conn *Conn) {
			waitconn <- conn
			<-waitclose
		}

		go func() {
			waitstart <- server.ListenAndServe()
		}()

		select {
		case err = <-waitstart:
			if err != nil {
				return
			}

		case conn := <-waitconn:
			demuxer = closeConn{Conn: conn, waitclose: waitclose}
			return
		}

		return
	}

	h.CodecTypes = CodecTypes
}
