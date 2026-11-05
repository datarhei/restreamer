package srt

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"sync"
	"time"

	"github.com/datarhei/gosrt/crypto"
	srtnet "github.com/datarhei/gosrt/net"
	"github.com/datarhei/gosrt/packet"
)

// ConnType represents the kind of connection as returned
// from the AcceptFunc. It is one of REJECT, PUBLISH, or SUBSCRIBE.
type ConnType int

// String returns a string representation of the ConnType.
func (c ConnType) String() string {
	switch c {
	case REJECT:
		return "REJECT"
	case PUBLISH:
		return "PUBLISH"
	case SUBSCRIBE:
		return "SUBSCRIBE"
	default:
		return ""
	}
}

const (
	REJECT    ConnType = ConnType(1 << iota) // Reject a connection
	PUBLISH                                  // This connection is meant to write data to the server
	SUBSCRIBE                                // This connection is meant to read data from a PUBLISHed stream
)

// RejectionReason are the rejection reasons that can be returned from the AcceptFunc in order to send
// another reason than the default one (REJ_PEER) to the client.
type RejectionReason uint32

// Table 7: Handshake Rejection Reason Codes
const (
	REJ_UNKNOWN    RejectionReason = 1000 // unknown reason
	REJ_SYSTEM     RejectionReason = 1001 // system function error
	REJ_PEER       RejectionReason = 1002 // rejected by peer
	REJ_RESOURCE   RejectionReason = 1003 // resource allocation problem
	REJ_ROGUE      RejectionReason = 1004 // incorrect data in handshake
	REJ_BACKLOG    RejectionReason = 1005 // listener's backlog exceeded
	REJ_IPE        RejectionReason = 1006 // internal program error
	REJ_CLOSE      RejectionReason = 1007 // socket is closing
	REJ_VERSION    RejectionReason = 1008 // peer is older version than agent's min
	REJ_RDVCOOKIE  RejectionReason = 1009 // rendezvous cookie collision
	REJ_BADSECRET  RejectionReason = 1010 // wrong password
	REJ_UNSECURE   RejectionReason = 1011 // password required or unexpected
	REJ_MESSAGEAPI RejectionReason = 1012 // stream flag collision
	REJ_CONGESTION RejectionReason = 1013 // incompatible congestion-controller type
	REJ_FILTER     RejectionReason = 1014 // incompatible packet filter
	REJ_GROUP      RejectionReason = 1015 // incompatible group
)

// These are the extended rejection reasons that may be less well supported
// Codes & their meanings taken from https://github.com/Haivision/srt/blob/f477af533562505abf5295f059cf2156b17be740/srtcore/access_control.h
const (
	REJX_BAD_REQUEST   RejectionReason = 1400 // General syntax error in the SocketID specification (also a fallback code for undefined cases)
	REJX_UNAUTHORIZED  RejectionReason = 1401 // Authentication failed, provided that the user was correctly identified and access to the required resource would be granted
	REJX_OVERLOAD      RejectionReason = 1402 // The server is too heavily loaded, or you have exceeded credits for accessing the service and the resource.
	REJX_FORBIDDEN     RejectionReason = 1403 // Access denied to the resource by any kind of reason.
	REJX_NOTFOUND      RejectionReason = 1404 // Resource not found at this time.
	REJX_BAD_MODE      RejectionReason = 1405 // The mode specified in `m` key in StreamID is not supported for this request.
	REJX_UNACCEPTABLE  RejectionReason = 1406 // The requested parameters specified in SocketID cannot be satisfied for the requested resource. Also when m=publish and the data format is not acceptable.
	REJX_CONFLICT      RejectionReason = 1407 // The resource being accessed is already locked for modification. This is in case of m=publish and the specified resource is currently read-only.
	REJX_NOTSUP_MEDIA  RejectionReason = 1415 // The media type is not supported by the application. This is the `t` key that specifies the media type as stream, file and auth, possibly extended by the application.
	REJX_LOCKED        RejectionReason = 1423 // The resource being accessed is locked for any access.
	REJX_FAILED_DEPEND RejectionReason = 1424 // The request failed because it specified a dependent session ID that has been disconnected.
	REJX_ISE           RejectionReason = 1500 // Unexpected internal server error
	REJX_UNIMPLEMENTED RejectionReason = 1501 // The request was recognized, but the current version doesn't support it.
	REJX_GW            RejectionReason = 1502 // The server acts as a gateway and the target endpoint rejected the connection.
	REJX_DOWN          RejectionReason = 1503 // The service has been temporarily taken over by a stub reporting this error. The real service can be down for maintenance or crashed.
	REJX_VERSION       RejectionReason = 1505 // SRT version not supported. This might be either unsupported backward compatibility, or an upper value of a version.
	REJX_NOROOM        RejectionReason = 1507 // The data stream cannot be archived due to lacking storage space. This is in case when the request type was to send a file or the live stream to be archived.
)

// ConnRequest is an incoming connection request
type ConnRequest interface {
	// RemoteAddr returns the address of the peer. The returned net.Addr
	// is a copy and can be used at will.
	RemoteAddr() net.Addr

	// Version returns the handshake version of the incoming request. Currently
	// known versions are 4 and 5. With version 4 the StreamId will always be
	// empty and IsEncrypted will always return false. An incoming version 4
	// connection will always be publishing.
	Version() uint32

	// StreamId returns the streamid of the requesting connection. Use this
	// to decide what to do with the connection.
	StreamId() string

	// IsEncrypted returns whether the connection is encrypted. If it is
	// encrypted, use SetPassphrase to set the passphrase for decrypting.
	IsEncrypted() bool

	// SetPassphrase sets the passphrase in order to decrypt the incoming
	// data. Returns an error if the passphrase did not work or the connection
	// is not encrypted.
	SetPassphrase(p string) error

	// SetRejectionReason sets the rejection reason for the connection. If
	// no set, REJ_PEER will be used.
	SetRejectionReason(r RejectionReason)
}

// connRequest implements the ConnRequest interface
type connRequest struct {
	addr      net.Addr
	start     time.Time
	socketId  uint32
	timestamp uint32

	config          Config
	handshake       *packet.CIFHandshake
	crypto          crypto.Crypto
	passphrase      string
	rejectionReason RejectionReason
}

func (req *connRequest) RemoteAddr() net.Addr {
	addr, _ := net.ResolveUDPAddr("udp", req.addr.String())
	return addr
}

func (req *connRequest) Version() uint32 {
	return req.handshake.Version
}

func (req *connRequest) StreamId() string {
	return req.handshake.StreamId
}

func (req *connRequest) IsEncrypted() bool {
	return req.crypto != nil
}

func (req *connRequest) SetPassphrase(passphrase string) error {
	if req.handshake.Version == 5 {
		if req.crypto == nil {
			return fmt.Errorf("listen: request without encryption")
		}

		if err := req.crypto.UnmarshalKM(req.handshake.SRTKM, passphrase); err != nil {
			return err
		}
	}

	req.passphrase = passphrase

	return nil
}

func (req *connRequest) SetRejectionReason(reason RejectionReason) {
	req.rejectionReason = reason
}

// ErrListenerClosed is returned when the listener is about to shutdown.
var ErrListenerClosed = errors.New("srt: listener closed")

// AcceptFunc receives a connection request and returns the type of connection
// and is required by the Listener for each Accept of a new connection.
type AcceptFunc func(req ConnRequest) ConnType

// Listener waits for new connections
type Listener interface {
	// Accept waits for new connections. For each new connection the AcceptFunc
	// gets called. Conn is a new connection if AcceptFunc is PUBLISH or SUBSCRIBE.
	// If AcceptFunc returns REJECT, Conn is nil. In case of failure error is not
	// nil, Conn is nil and ConnType is REJECT. On closing the listener err will
	// be ErrListenerClosed and ConnType is REJECT.
	Accept(AcceptFunc) (Conn, ConnType, error)

	// Close closes the listener. It will stop accepting new connections and
	// close all currently established connections.
	Close()

	// Addr returns the address of the listener.
	Addr() net.Addr
}

// listener implements the Listener interface.
type listener struct {
	pc   *net.UDPConn
	addr net.Addr

	config Config

	backlog chan connRequest
	conns   map[uint32]*srtConn
	lock    sync.RWMutex

	start time.Time

	rcvQueue chan packet.Packet

	sndMutex sync.Mutex
	sndData  bytes.Buffer

	syncookie *srtnet.SYNCookie

	shutdown     bool
	shutdownLock sync.RWMutex
	shutdownOnce sync.Once

	stopReader context.CancelFunc

	doneChan chan struct{}
	doneErr  error
	doneOnce sync.Once
}

// Listen returns a new listener on the SRT protocol on the address with
// the provided config. The network parameter needs to be "srt".
//
// The address has the form "host:port".
//
// Examples:
//
//	Listen("srt", "127.0.0.1:3000", DefaultConfig())
//
// In case of an error, the returned Listener is nil and the error is non-nil.
func Listen(network, address string, config Config) (Listener, error) {
	if network != "srt" {
		return nil, fmt.Errorf("listen: the network must be 'srt'")
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("listen: invalid config: %w", err)
	}

	if config.Logger == nil {
		config.Logger = NewLogger(nil)
	}

	ln := &listener{
		config: config,
	}

	lc := net.ListenConfig{
		Control: ListenControl(config),
	}

	lp, err := lc.ListenPacket(context.Background(), "udp", address)
	if err != nil {
		return nil, fmt.Errorf("listen: %w", err)
	}

	pc := lp.(*net.UDPConn)

	ln.pc = pc
	ln.addr = pc.LocalAddr()
	if ln.addr == nil {
		return nil, fmt.Errorf("listen: no local address")
	}

	ln.conns = make(map[uint32]*srtConn)

	ln.backlog = make(chan connRequest, 128)

	ln.rcvQueue = make(chan packet.Packet, 2048)

	syncookie, err := srtnet.NewSYNCookie(ln.addr.String(), nil)
	if err != nil {
		ln.Close()
		return nil, err
	}
	ln.syncookie = syncookie

	ln.doneChan = make(chan struct{})

	ln.start = time.Now()

	var readerCtx context.Context
	readerCtx, ln.stopReader = context.WithCancel(context.Background())
	go ln.reader(readerCtx)

	go func() {
		buffer := make([]byte, config.MSS) // MTU size

		for {
			if ln.isShutdown() {
				ln.markDone(ErrListenerClosed)
				return
			}

			ln.pc.SetReadDeadline(time.Now().Add(3 * time.Second))
			n, addr, err := ln.pc.ReadFrom(buffer)
			if err != nil {
				if errors.Is(err, os.ErrDeadlineExceeded) {
					continue
				}

				if ln.isShutdown() {
					ln.markDone(ErrListenerClosed)
					return
				}

				ln.markDone(err)
				return
			}

			p, err := packet.NewPacketFromData(addr, buffer[:n])
			if err != nil {
				continue
			}

			// non-blocking
			select {
			case ln.rcvQueue <- p:
			default:
				ln.log("listen", func() string { return "receive queue is full" })
			}
		}
	}()

	return ln, nil
}

func (ln *listener) Accept(acceptFn AcceptFunc) (Conn, ConnType, error) {
	if ln.isShutdown() {
		return nil, REJECT, ErrListenerClosed
	}

	select {
	case <-ln.doneChan:
		return nil, REJECT, ln.error()
	case request := <-ln.backlog:
		if acceptFn == nil {
			ln.reject(request, REJ_PEER)
			break
		}

		mode := acceptFn(&request)
		if mode != PUBLISH && mode != SUBSCRIBE {
			// Figure out the reason
			reason := REJ_PEER
			if request.rejectionReason > 0 {
				reason = request.rejectionReason
			}
			ln.reject(request, reason)
			break
		}

		if request.crypto != nil && len(request.passphrase) == 0 {
			ln.reject(request, REJ_BADSECRET)
			break
		}

		// Create a new socket ID
		socketId := uint32(time.Since(ln.start).Microseconds())

		// Select the largest TSBPD delay advertised by the caller, but at least 120ms
		recvTsbpdDelay := uint16(request.config.ReceiverLatency.Milliseconds())
		sendTsbpdDelay := uint16(request.config.PeerLatency.Milliseconds())

		if request.handshake.Version == 5 {
			if request.handshake.SRTHS.SendTSBPDDelay > recvTsbpdDelay {
				recvTsbpdDelay = request.handshake.SRTHS.SendTSBPDDelay
			}

			if request.handshake.SRTHS.RecvTSBPDDelay > sendTsbpdDelay {
				sendTsbpdDelay = request.handshake.SRTHS.RecvTSBPDDelay
			}

			request.config.StreamId = request.handshake.StreamId
		}

		request.config.Passphrase = request.passphrase

		// Create a new connection
		conn := newSRTConn(srtConnConfig{
			version:                     request.handshake.Version,
			localAddr:                   ln.addr,
			remoteAddr:                  request.addr,
			config:                      request.config,
			start:                       request.start,
			socketId:                    socketId,
			peerSocketId:                request.handshake.SRTSocketId,
			tsbpdTimeBase:               uint64(request.timestamp),
			tsbpdDelay:                  uint64(recvTsbpdDelay) * 1000,
			peerTsbpdDelay:              uint64(sendTsbpdDelay) * 1000,
			initialPacketSequenceNumber: request.handshake.InitialPacketSequenceNumber,
			crypto:                      request.crypto,
			keyBaseEncryption:           packet.EvenKeyEncrypted,
			onSend:                      ln.send,
			onShutdown:                  ln.handleShutdown,
			logger:                      request.config.Logger,
		})

		ln.log("connection:new", func() string { return fmt.Sprintf("%#08x (%s) %s", conn.SocketId(), conn.StreamId(), mode) })

		request.handshake.SRTSocketId = socketId
		request.handshake.SynCookie = 0

		if request.handshake.Version == 5 {
			//  3.2.1.1.1.  Handshake Extension Message Flags
			request.handshake.SRTHS.SRTVersion = SRT_VERSION
			request.handshake.SRTHS.SRTFlags.TSBPDSND = true
			request.handshake.SRTHS.SRTFlags.TSBPDRCV = true
			request.handshake.SRTHS.SRTFlags.CRYPT = true
			request.handshake.SRTHS.SRTFlags.TLPKTDROP = true
			request.handshake.SRTHS.SRTFlags.PERIODICNAK = true
			request.handshake.SRTHS.SRTFlags.REXMITFLG = true
			request.handshake.SRTHS.SRTFlags.STREAM = false
			request.handshake.SRTHS.SRTFlags.PACKET_FILTER = false
			request.handshake.SRTHS.RecvTSBPDDelay = recvTsbpdDelay
			request.handshake.SRTHS.SendTSBPDDelay = sendTsbpdDelay
		}

		ln.accept(request)

		// Add the connection to the list of known connections
		ln.lock.Lock()
		ln.conns[socketId] = conn
		ln.lock.Unlock()

		return conn, mode, nil
	}

	return nil, REJECT, nil
}

// markDone marks the listener as done by closing
// the done channel & sets the error
func (ln *listener) markDone(err error) {
	ln.doneOnce.Do(func() {
		ln.lock.Lock()
		defer ln.lock.Unlock()
		ln.doneErr = err
		close(ln.doneChan)
	})
}

// error returns the error that caused the listener to be done
// if it's nil then the listener is not done
func (ln *listener) error() error {
	ln.lock.Lock()
	defer ln.lock.Unlock()
	return ln.doneErr
}

func (ln *listener) handleShutdown(socketId uint32) {
	ln.lock.Lock()
	delete(ln.conns, socketId)
	ln.lock.Unlock()
}

func (ln *listener) reject(request connRequest, reason RejectionReason) {
	p := packet.NewPacket(request.addr)
	p.Header().IsControlPacket = true

	p.Header().ControlType = packet.CTRLTYPE_HANDSHAKE
	p.Header().SubType = 0
	p.Header().TypeSpecific = 0

	p.Header().Timestamp = uint32(time.Since(ln.start).Microseconds())
	p.Header().DestinationSocketId = request.socketId

	request.handshake.HandshakeType = packet.HandshakeType(reason)

	p.MarshalCIF(request.handshake)

	ln.log("handshake:send:dump", func() string { return p.Dump() })
	ln.log("handshake:send:cif", func() string { return request.handshake.String() })

	ln.send(p)
}

func (ln *listener) accept(request connRequest) {
	p := packet.NewPacket(request.addr)

	p.Header().IsControlPacket = true

	p.Header().ControlType = packet.CTRLTYPE_HANDSHAKE
	p.Header().SubType = 0
	p.Header().TypeSpecific = 0

	p.Header().Timestamp = uint32(time.Since(request.start).Microseconds())
	p.Header().DestinationSocketId = request.socketId

	p.MarshalCIF(request.handshake)

	ln.log("handshake:send:dump", func() string { return p.Dump() })
	ln.log("handshake:send:cif", func() string { return request.handshake.String() })

	ln.send(p)
}

func (ln *listener) isShutdown() bool {
	ln.shutdownLock.RLock()
	defer ln.shutdownLock.RUnlock()

	return ln.shutdown
}

func (ln *listener) Close() {
	ln.shutdownOnce.Do(func() {
		ln.shutdownLock.Lock()
		ln.shutdown = true
		ln.shutdownLock.Unlock()

		ln.lock.RLock()
		for _, conn := range ln.conns {
			conn.close()
		}
		ln.lock.RUnlock()

		ln.stopReader()

		ln.log("listen", func() string { return "closing socket" })

		ln.pc.Close()
	})
}

func (ln *listener) Addr() net.Addr {
	addrString := "0.0.0.0:0"
	if ln.addr != nil {
		addrString = ln.addr.String()
	}

	addr, _ := net.ResolveUDPAddr("udp", addrString)
	return addr
}

func (ln *listener) reader(ctx context.Context) {
	defer func() {
		ln.log("listen", func() string { return "left reader loop" })
	}()

	ln.log("listen", func() string { return "reader loop started" })

	for {
		select {
		case <-ctx.Done():
			return
		case p := <-ln.rcvQueue:
			if ln.isShutdown() {
				break
			}

			ln.log("packet:recv:dump", func() string { return p.Dump() })

			if p.Header().DestinationSocketId == 0 {
				if p.Header().IsControlPacket && p.Header().ControlType == packet.CTRLTYPE_HANDSHAKE {
					ln.handleHandshake(p)
				}

				break
			}

			ln.lock.RLock()
			conn, ok := ln.conns[p.Header().DestinationSocketId]
			ln.lock.RUnlock()

			if !ok {
				// ignore the packet, we don't know the destination
				break
			}

			conn.push(p)
		}
	}
}

// Send a packet to the wire. This function must be synchronous in order to allow to safely call Packet.Decommission() afterward.
func (ln *listener) send(p packet.Packet) {
	ln.sndMutex.Lock()
	defer ln.sndMutex.Unlock()

	ln.sndData.Reset()

	if err := p.Marshal(&ln.sndData); err != nil {
		p.Decommission()
		ln.log("packet:send:error", func() string { return "marshalling packet failed" })
		return
	}

	buffer := ln.sndData.Bytes()

	ln.log("packet:send:dump", func() string { return p.Dump() })

	// Write the packet's contents to the wire
	ln.pc.WriteTo(buffer, p.Header().Addr)

	if p.Header().IsControlPacket {
		// Control packets can be decommissioned because they will not be sent again (data packets might be retransferred)
		p.Decommission()
	}
}

func (ln *listener) handleHandshake(p packet.Packet) {
	cif := &packet.CIFHandshake{}

	err := p.UnmarshalCIF(cif)

	ln.log("handshake:recv:dump", func() string { return p.Dump() })
	ln.log("handshake:recv:cif", func() string { return cif.String() })

	if err != nil {
		ln.log("handshake:recv:error", func() string { return err.Error() })
		return
	}

	// Assemble the response (4.3.1.  Caller-Listener Handshake)

	p.Header().ControlType = packet.CTRLTYPE_HANDSHAKE
	p.Header().SubType = 0
	p.Header().TypeSpecific = 0
	p.Header().Timestamp = uint32(time.Since(ln.start).Microseconds())
	p.Header().DestinationSocketId = cif.SRTSocketId

	cif.PeerIP.FromNetAddr(ln.addr)

	// Create a copy of the configuration for the connection
	config := ln.config

	if cif.HandshakeType == packet.HSTYPE_INDUCTION {
		// cif
		cif.Version = 5
		cif.EncryptionField = 0 // Don't advertise any specific encryption method
		cif.ExtensionField = 0x4A17
		//cif.initialPacketSequenceNumber = newCircular(0, MAX_SEQUENCENUMBER)
		//cif.maxTransmissionUnitSize = 0
		//cif.maxFlowWindowSize = 0
		//cif.SRTSocketId = 0
		cif.SynCookie = ln.syncookie.Get(p.Header().Addr.String())

		p.MarshalCIF(cif)

		ln.log("handshake:send:dump", func() string { return p.Dump() })
		ln.log("handshake:send:cif", func() string { return cif.String() })

		ln.send(p)
	} else if cif.HandshakeType == packet.HSTYPE_CONCLUSION {
		// Verify the SYN cookie
		if !ln.syncookie.Verify(cif.SynCookie, p.Header().Addr.String()) {
			cif.HandshakeType = packet.HandshakeType(REJ_ROGUE)
			ln.log("handshake:recv:error", func() string { return "invalid SYN cookie" })
			p.MarshalCIF(cif)
			ln.log("handshake:send:dump", func() string { return p.Dump() })
			ln.log("handshake:send:cif", func() string { return cif.String() })
			ln.send(p)

			return
		}

		// Peer is advertising a too big MSS
		if cif.MaxTransmissionUnitSize > MAX_MSS_SIZE {
			cif.HandshakeType = packet.HandshakeType(REJ_ROGUE)
			ln.log("handshake:recv:error", func() string { return fmt.Sprintf("MTU is too big (%d bytes)", cif.MaxTransmissionUnitSize) })
			p.MarshalCIF(cif)
			ln.log("handshake:send:dump", func() string { return p.Dump() })
			ln.log("handshake:send:cif", func() string { return cif.String() })
			ln.send(p)

			return
		}

		// If the peer has a smaller MTU size, adjust to it
		if cif.MaxTransmissionUnitSize < config.MSS {
			config.MSS = cif.MaxTransmissionUnitSize
			config.PayloadSize = config.MSS - SRT_HEADER_SIZE - UDP_HEADER_SIZE

			if config.PayloadSize < MIN_PAYLOAD_SIZE {
				cif.HandshakeType = packet.HandshakeType(REJ_ROGUE)
				ln.log("handshake:recv:error", func() string { return fmt.Sprintf("payload size is too small (%d bytes)", config.PayloadSize) })
				p.MarshalCIF(cif)
				ln.log("handshake:send:dump", func() string { return p.Dump() })
				ln.log("handshake:send:cif", func() string { return cif.String() })
				ln.send(p)
			}
		}

		// We only support HSv4 and HSv5
		if cif.Version == 4 {
			// Check if the type (encryption field + extension field) has the value 2
			if cif.EncryptionField != 0 || cif.ExtensionField != 2 {
				cif.HandshakeType = packet.HandshakeType(REJ_ROGUE)
				ln.log("handshake:recv:error", func() string { return "invalid type, expecting a value of 2 (UDT_DGRAM)" })
				p.MarshalCIF(cif)
				ln.log("handshake:send:dump", func() string { return p.Dump() })
				ln.log("handshake:send:cif", func() string { return cif.String() })
				ln.send(p)

				return
			}
		} else if cif.Version == 5 {
			// Check if the peer version is sufficient
			if cif.SRTHS.SRTVersion < config.MinVersion {
				cif.HandshakeType = packet.HandshakeType(REJ_VERSION)
				ln.log("handshake:recv:error", func() string {
					return fmt.Sprintf("peer version insufficient (%#06x), expecting at least %#06x", cif.SRTHS.SRTVersion, config.MinVersion)
				})
				p.MarshalCIF(cif)
				ln.log("handshake:send:dump", func() string { return p.Dump() })
				ln.log("handshake:send:cif", func() string { return cif.String() })
				ln.send(p)

				return
			}

			// Check the required SRT flags
			if !cif.SRTHS.SRTFlags.TSBPDSND || !cif.SRTHS.SRTFlags.TSBPDRCV || !cif.SRTHS.SRTFlags.TLPKTDROP || !cif.SRTHS.SRTFlags.PERIODICNAK || !cif.SRTHS.SRTFlags.REXMITFLG {
				cif.HandshakeType = packet.HandshakeType(REJ_ROGUE)
				ln.log("handshake:recv:error", func() string { return "not all required flags are set" })
				p.MarshalCIF(cif)
				ln.log("handshake:send:dump", func() string { return p.Dump() })
				ln.log("handshake:send:cif", func() string { return cif.String() })
				ln.send(p)

				return
			}

			// We only support live streaming
			if cif.SRTHS.SRTFlags.STREAM {
				cif.HandshakeType = packet.HandshakeType(REJ_MESSAGEAPI)
				ln.log("handshake:recv:error", func() string { return "only live streaming is supported" })
				p.MarshalCIF(cif)
				ln.log("handshake:send:dump", func() string { return p.Dump() })
				ln.log("handshake:send:cif", func() string { return cif.String() })
				ln.send(p)

				return
			}
		} else {
			cif.HandshakeType = packet.HandshakeType(REJ_ROGUE)
			ln.log("handshake:recv:error", func() string { return fmt.Sprintf("only HSv4 and HSv5 are supported (got HSv%d)", cif.Version) })
			p.MarshalCIF(cif)
			ln.log("handshake:send:dump", func() string { return p.Dump() })
			ln.log("handshake:send:cif", func() string { return cif.String() })
			ln.send(p)

			return
		}

		// Fill up a connection request with all relevant data and put it into the backlog

		c := connRequest{
			addr:      p.Header().Addr,
			start:     time.Now(),
			socketId:  cif.SRTSocketId,
			timestamp: p.Header().Timestamp,
			config:    config,

			handshake: cif,
		}

		if cif.SRTKM != nil {
			cr, err := crypto.New(int(cif.SRTKM.KLen))
			if err != nil {
				cif.HandshakeType = packet.HandshakeType(REJ_ROGUE)
				ln.log("handshake:recv:error", func() string { return fmt.Sprintf("crypto: %s", err) })
				p.MarshalCIF(cif)
				ln.log("handshake:send:dump", func() string { return p.Dump() })
				ln.log("handshake:send:cif", func() string { return cif.String() })
				ln.send(p)

				return
			}

			c.crypto = cr
		}

		// If the backlog is full, reject the connection
		select {
		case ln.backlog <- c:
		default:
			cif.HandshakeType = packet.HandshakeType(REJ_BACKLOG)
			ln.log("handshake:recv:error", func() string { return "backlog is full" })
			p.MarshalCIF(cif)
			ln.log("handshake:send:dump", func() string { return p.Dump() })
			ln.log("handshake:send:cif", func() string { return cif.String() })
			ln.send(p)
		}
	} else {
		if cif.HandshakeType.IsRejection() {
			ln.log("handshake:recv:error", func() string { return fmt.Sprintf("connection rejected: %s", cif.HandshakeType.String()) })
		} else {
			ln.log("handshake:recv:error", func() string { return fmt.Sprintf("unsupported handshake: %s", cif.HandshakeType.String()) })
		}
	}
}

func (ln *listener) log(topic string, message func() string) {
	if ln.config.Logger == nil {
		return
	}

	ln.config.Logger.Print(topic, 0, 2, message)
}
