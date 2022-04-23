package live

// https://github.com/lovelyyoshino/Bilibili-Live-API/blob/master/API.WebSocket.md
// https://github.com/SocialSisterYi/bilibili-API-collect/blob/master/live/message_stream.md

import (
	"bytes"
	"compress/zlib"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/andybalholm/brotli"
	"github.com/gorilla/websocket"
	"github.com/reveever/biliapi/base"
)

const (
	wsSubURLFormat  = "ws://%s:%d/sub"
	wssSubURLFormat = "wss://%s:%d/sub"
)

var defaultWsSubHost = LiveHost{
	Host:    "broadcastlv.chat.bilibili.com",
	Port:    2243,
	WssPort: 443,
	WsPort:  2244,
}

type Config struct {
	NoSecure           bool
	ReadHeartbeatReply bool
	HeartBeatInterval  time.Duration
	MsgCap             int
	Log                base.BaseLogger
	AuthInfo           UserAuth
	HostList           []LiveHost
	CookieJar          http.CookieJar
	UserAgent          string
	Dialer             *websocket.Dialer
	Header             http.Header
}

type Live struct {
	heartbeatReply bool
	closed         *int32
	subMsgs        chan wsSubMsgWithErr
	log            base.BaseLogger
	conn           *websocket.Conn
	ctx            context.Context
	cancel         context.CancelFunc
}

type UserAuth struct {
	UID      int64  `json:"uid"`
	RoomID   int64  `json:"roomid"`
	ProtoVer int    `json:"protover"`
	Platform string `json:"platform"`
	Type     int    `json:"type"`
	Key      string `json:"key"`
}

type wsSubMsgWithErr struct {
	msg WsSubMsg
	err error
}

type WsSubMsg struct {
	Ver  uint16
	Op   uint32
	Cmd  string
	Body []byte
}

type wsSubMsgCmd struct {
	Cmd string `json:"cmd"`
}

func NewLive(ctx context.Context, cfg Config) (*Live, error) {
	if cfg.MsgCap == 0 {
		cfg.MsgCap = 64
	}

	if cfg.HeartBeatInterval == 0 {
		cfg.HeartBeatInterval = time.Second * 30
	}

	if cfg.Log == nil {
		cfg.Log = base.NoLogger{}
	}

	if cfg.AuthInfo.RoomID == 0 {
		return nil, fmt.Errorf("empty roomID")
	}

	if cfg.AuthInfo.ProtoVer == 0 {
		cfg.AuthInfo.ProtoVer = WsBodyProtocolVersionBrotli
	} else if cfg.AuthInfo.ProtoVer != WsBodyProtocolVersionZlib && cfg.AuthInfo.ProtoVer != WsBodyProtocolVersionBrotli {
		return nil, fmt.Errorf("invalid ws body protocal version")
	}

	if cfg.AuthInfo.Platform == "" {
		cfg.AuthInfo.Platform = "web"
	}

	if cfg.AuthInfo.Type == 0 {
		cfg.AuthInfo.Type = 2
	}

	if cfg.AuthInfo.Key == "" {
		return nil, fmt.Errorf("empty token")
	}

	if len(cfg.HostList) == 0 {
		cfg.HostList = []LiveHost{defaultWsSubHost}
	}

	if cfg.Dialer == nil {
		cfg.Dialer = websocket.DefaultDialer
	}

	if cfg.UserAgent != "" {
		if cfg.Header == nil {
			cfg.Header = make(http.Header, 1)
		}
		cfg.Header.Set("User-Agent", cfg.UserAgent)
	}

	if cfg.CookieJar != nil {
		cfg.Dialer.Jar = cfg.CookieJar
	}

	var (
		url            string
		err            error
		closed         = int32(0)
		subCtx, cancel = context.WithCancel(ctx)
	)
	l := &Live{
		heartbeatReply: cfg.ReadHeartbeatReply,
		closed:         &closed,
		subMsgs:        make(chan wsSubMsgWithErr, cfg.MsgCap),
		log:            cfg.Log,
		conn:           nil,
		ctx:            subCtx,
		cancel:         cancel,
	}

	for _, host := range cfg.HostList {
		if cfg.NoSecure {
			url = fmt.Sprintf(wsSubURLFormat, host.Host, host.WsPort)
		} else {
			url = fmt.Sprintf(wssSubURLFormat, host.Host, host.WssPort)
		}
		l.log.Printf("dial %s", url)
		l.conn, _, err = cfg.Dialer.DialContext(subCtx, url, cfg.Header)
		if err != nil {
			l.log.Println(err)
			continue
		}

		go l.run()
		if err := l.auth(cfg.AuthInfo); err != nil {
			return nil, fmt.Errorf("user auth: %v", err)
		}
		go l.heartbeat(cfg.HeartBeatInterval)
		return l, nil
	}
	cancel()
	return nil, err
}

func (l *Live) auth(authInfo UserAuth) error {
	buf, err := json.Marshal(authInfo)
	if err != nil {
		return err
	}
	err = l.conn.WriteMessage(websocket.BinaryMessage,
		NewWsSubPkt(WsBodyProtocolVersionNormal, WsOpUserAuthentication, buf))
	if err != nil {
		return err
	}

	authResp, err := l.ReadWsMsg(l.ctx)
	if err != nil {
		return err
	}
	if authResp.Op != WsOpConnectSuccess {
		return fmt.Errorf("unexpect op %d", authResp.Op)
	}
	if !bytes.Equal(authResp.Body, []byte(`{"code":0}`)) {
		return fmt.Errorf(string(authResp.Body))
	}
	l.log.Printf("auth succ: %s", string(authResp.Body))
	return nil
}

func (l *Live) heartbeat(d time.Duration) {
	l.log.Println("start heartbeat")
	defer l.log.Println("stop heartbeat")
	defer l.Close()
	defer l.recover()

	f := func() {
		l.log.Println("send WsOpHeartbeat")
		err := l.conn.WriteMessage(websocket.BinaryMessage,
			NewWsSubPkt(WsBodyProtocolVersionNormal, WsOpHeartbeat, nil))
		if err != nil {
			l.error(fmt.Errorf("send WsOpHeartbeat: %x", err))
		}
	}

	f()
	ticker := time.NewTicker(d)
	defer ticker.Stop()

	for {
		select {
		case <-l.ctx.Done():
			return
		case <-ticker.C:
			f()
		}
	}
}

func (l *Live) run() {
	l.log.Println("start live run")
	defer l.log.Println("stop live run")
	defer l.Close()
	defer l.recover()

	var (
		typ int
		msg []byte
		err error
	)

	for {
		typ, msg, err = l.conn.ReadMessage()
		if err != nil {
			l.error(fmt.Errorf("read ws: %v", err))
			return
		}
		switch typ {
		case websocket.BinaryMessage:
			err := l.handleWsMsg(msg)
			if err == io.EOF || err == nil {
				continue
			}
			l.error(fmt.Errorf("handle ws msg: %v %x", err, msg))
		case websocket.TextMessage:
			l.error(fmt.Errorf("unexpected text msg: %s", string(msg)))
		default:
			l.error(fmt.Errorf("unexpected msg type %d: %x", typ, msg))
		}
	}
}

func (l *Live) handleWsMsg(pkt WsSubPkt) error {
	ver, op := pkt.Version(), pkt.Operation()

	l.log.Printf("handle WsSubPkt: pktLen:%d, hdrLen:%d, ver:%d, op:%d, seq:%d",
		pkt.PktLen(), pkt.HdrLen(), ver, op, pkt.Sequence())

	switch ver {
	case WsBodyProtocolVersionBrotli:
		r := brotli.NewReader(bytes.NewReader(pkt.Body()))
		for {
			msg, err := ReadWsSubPkt(r)
			if err != nil {
				return err
			}
			return l.handleWsMsg(msg) //lint:ignore SA4004 Ignore the deprecation warnings
		}

	case WsBodyProtocolVersionZlib:
		r, err := zlib.NewReader(bytes.NewReader(pkt.Body()))
		if err != nil {
			return fmt.Errorf("zlib reader: %v", err)
		}
		for {
			msg, err := ReadWsSubPkt(r)
			if err != nil {
				return err
			}
			return l.handleWsMsg(msg) //lint:ignore SA4004 Ignore the deprecation warnings
		}

	case WsBodyProtocolVersionNormal:
		switch op {
		case WsOpMessage:
			var wsMsgCmd wsSubMsgCmd
			body := pkt.Body()
			err := json.Unmarshal(body, &wsMsgCmd)
			if err != nil {
				return fmt.Errorf("read WsOpMessage: %v", err)
			}
			l.log.Printf("recv WsOpMessage: %s", string(body))
			l.write(ver, op, wsMsgCmd.Cmd, body)

		case WsOpHeartbeatReply:
			body := pkt.Body()
			l.log.Printf("recv WsOpHeartbeatReply: %x", body)
			if l.heartbeatReply {
				l.write(ver, op, "", body)
			}

		case WsOpConnectSuccess:
			l.write(ver, op, "", pkt.Body())
			return nil

		case WsOpHeartbeat, WsOpUserAuthentication:
			return fmt.Errorf("unexpected op: %d", op)

		default:
			return fmt.Errorf("unknown op: %d", op)
		}

	default:
		return fmt.Errorf("unknown protocol version: %d", ver)
	}
	return nil
}

func (l *Live) ReadWsMsg(ctx context.Context) (WsSubMsg, error) {
	select {
	case <-ctx.Done():
		return WsSubMsg{}, l.ctx.Err()
	case m, ok := <-l.subMsgs:
		if !ok {
			return WsSubMsg{}, io.EOF
		}
		return m.msg, m.err
	}
}

func (l *Live) Close() {
	if !atomic.CompareAndSwapInt32(l.closed, 0, 1) {
		l.log.Println("live closed")
		return
	}
	l.log.Println("close live, send close message")

	l.cancel()
	err := l.conn.WriteMessage(websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		l.log.Printf("failed to send close message: %v", err)
	}
	close(l.subMsgs)
}

func (l *Live) write(ver uint16, op uint32, cmd string, body []byte) {
	if atomic.LoadInt32(l.closed) == 1 {
		l.log.Println("drop msg: ver: %d op: %d cmd: %s bodyLen: %d", ver, op, cmd, len(body))
		return
	}
	select {
	case <-l.ctx.Done():
		l.log.Println("drop msg: ver: %d op: %d cmd: %s bodyLen: %d", ver, op, cmd, len(body))
	case l.subMsgs <- wsSubMsgWithErr{
		msg: WsSubMsg{Ver: ver, Op: op, Cmd: cmd, Body: body}}:
	}
}

func (l *Live) error(err error) {
	l.log.Printf("error: %v", err)
	if atomic.LoadInt32(l.closed) == 1 {
		l.log.Printf("drop error: %v", err)
		return
	}
	select {
	case <-l.ctx.Done():
		l.log.Printf("drop error: %v", err)
	case l.subMsgs <- wsSubMsgWithErr{
		err: err}:
	}
}

func (l *Live) recover() {
	if r := recover(); r != nil {
		l.error(fmt.Errorf("panic: %v", r))
	}
}
