package biliapi

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"github.com/reveever/biliapi/base"
	"github.com/reveever/biliapi/interface/live"
)

type BiliLive struct {
	*live.Live
}

func NewBiliLive(ctx context.Context, roomID int64, token string, options ...func(*live.Config)) (*BiliLive, error) {
	var cfg live.Config
	for _, o := range options {
		o(&cfg)
	}
	if roomID != 0 {
		cfg.AuthInfo.RoomID = roomID
	}
	if token != "" {
		cfg.AuthInfo.Key = token
	}

	l, err := live.NewLive(ctx, cfg)
	if err != nil {
		return nil, nil
	}
	return &BiliLive{l}, nil
}

func LiveNoSecure() func(*live.Config) {
	return func(c *live.Config) {
		c.NoSecure = true
	}
}

func LiveReadHeartbeatReply() func(*live.Config) {
	return func(c *live.Config) {
		c.ReadHeartbeatReply = true
	}
}

func LiveSetHeartBeatInterval(d time.Duration) func(*live.Config) {
	return func(c *live.Config) {
		c.HeartBeatInterval = d
	}
}

func LiveSetQueueCapacity(l int) func(*live.Config) {
	return func(c *live.Config) {
		c.MsgCap = l
	}
}

func LiveWithDebugLogger(l base.BaseLogger) func(*live.Config) {
	return func(c *live.Config) {
		c.Log = l
	}
}

func LiveEnableDebugLogger() func(*live.Config) {
	return func(c *live.Config) {
		if c.Log == nil {
			c.Log = log.New(os.Stderr, "", log.LstdFlags)
		}
	}
}

func LiveWithAuthInfo(v live.UserAuth) func(*live.Config) {
	return func(c *live.Config) {
		c.AuthInfo = v
	}
}

func LiveWithHostList(v []live.LiveHost) func(*live.Config) {
	return func(c *live.Config) {
		c.HostList = v
	}
}

func LiveWithCookieJar(j http.CookieJar) func(*live.Config) {
	return func(c *live.Config) {
		c.CookieJar = j
	}
}

func LiveWithBase(b *base.Base) func(*live.Config) {
	return func(c *live.Config) {
		c.UserAgent = b.UserAgent
		if c.CookieJar == nil {
			c.CookieJar = b.Client.Jar
		}
		if c.Log == nil {
			c.Log = b.Log
		}
	}
}

func LiveWithWsDialer(d *websocket.Dialer) func(*live.Config) {
	return func(c *live.Config) {
		c.Dialer = d
	}
}

func LiveSetHeader(h http.Header) func(*live.Config) {
	return func(c *live.Config) {
		c.Header = h
	}
}
