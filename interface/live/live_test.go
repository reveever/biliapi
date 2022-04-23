package live

import (
	"context"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/reveever/biliapi/interface/test"
)

func TestLive(t *testing.T) {
	api := newTestAPI(t)

	liveInfo, err := api.WebRoomInfo(WebRoomInfoOpt{
		ID:   test.RoomID,
		Type: 0,
	})
	if err != nil {
		t.Fatal(err)
	}

	cfg := Config{
		NoSecure:           false,
		ReadHeartbeatReply: true,
		HeartBeatInterval:  time.Second * 30,
		MsgCap:             1,
		Log:                test.NewLogger(t),
		AuthInfo: UserAuth{
			UID:      0,
			RoomID:   test.RoomID,
			ProtoVer: WsBodyProtocolVersionBrotli,
			Platform: "web",
			Type:     2,
			Key:      liveInfo.Token,
		},
		HostList: liveInfo.HostList,
		Dialer: &websocket.Dialer{
			Proxy:            http.ProxyFromEnvironment,
			HandshakeTimeout: 45 * time.Second,
			Jar:              api.base.Client.Jar,
		},
		Header: map[string][]string{
			"User-Agent": {api.base.UserAgent},
		},
	}

	l, err := NewLive(context.TODO(), cfg)
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 4; i++ {
		msg, err := l.ReadWsMsg(context.TODO())
		if err != nil {
			if err == io.EOF {
				t.Logf("EOF")
				break
			}
			t.Logf("err: %v", err)
		} else {
			if msg.Op == WsOpHeartbeatReply {
				t.Logf("%d %d %s: %x\n", msg.Ver, msg.Op, msg.Cmd, msg.Body)
			} else {
				t.Logf("%d %d %s: %s\n", msg.Ver, msg.Op, msg.Cmd, string(msg.Body))
			}
		}
	}
	l.Close()
	time.Sleep(time.Second)
}
