package biliapi

import (
	"log"
	"net/http"
	"os"

	"github.com/reveever/biliapi/base"
	"github.com/reveever/biliapi/interface/live"
	"github.com/reveever/biliapi/interface/room"
	"github.com/reveever/biliapi/interface/space"
	"github.com/reveever/biliapi/interface/v2dm"
)

type BiliAPI struct {
	Base  *base.Base
	Space *space.API
	V2Dm  *v2dm.API
	Live  *live.API
	Room  *room.API
}

func NewBiliApi(options ...func(*BiliAPI)) (*BiliAPI, error) {
	base := new(base.Base)
	biliAPI := &BiliAPI{
		Base:  base,
		Space: space.NewAPI(base),
		V2Dm:  v2dm.NewAPI(base),
		Live:  live.NewAPI(base),
		Room:  room.NewAPI(base),
	}

	for _, o := range options {
		o(biliAPI)
	}

	if err := base.Init(); err != nil {
		return nil, err
	}

	return biliAPI, nil
}

func EnableDebugLogger() func(*BiliAPI) {
	return func(b *BiliAPI) {
		b.Base.Log = log.New(os.Stderr, "", log.LstdFlags)
	}
}

func WithDebugLogger(l base.BaseLogger) func(*BiliAPI) {
	return func(b *BiliAPI) {
		b.Base.Log = l
	}
}

func WithHttpClient(c *http.Client) func(*BiliAPI) {
	return func(b *BiliAPI) {
		b.Base.Client = c
	}
}

func WithUserAgent(ua string) func(*BiliAPI) {
	return func(b *BiliAPI) {
		b.Base.UserAgent = ua
	}
}
