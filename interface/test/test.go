package test

import (
	"testing"

	"github.com/reveever/biliapi/base"
)

const (
	RoomID = 213       // https://live.bilibili.com/213
	Mid    = 67141     // https://space.bilibili.com/67141
	Oid    = 343721233 // https://www.bilibili.com/video/BV1bf4y1h7bK
	Pid    = 290815541 // https://www.bilibili.com/video/av5006123
)

type Logger struct {
	t *testing.T
}

func (l Logger) Println(v ...interface{}) {
	l.t.Log(v...)
}

func (l Logger) Printf(format string, v ...interface{}) {
	l.t.Logf(format, v...)
}

func NewTestBase(t *testing.T) *base.Base {
	base := &base.Base{
		Log: Logger{t},
	}
	err := base.Init()
	if err != nil {
		t.Fatal(err)
	}
	return base
}
