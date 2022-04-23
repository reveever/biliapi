# biliapi
[![Go Reference](https://pkg.go.dev/badge/github.com/reveever/biliapi.svg)](https://pkg.go.dev/github.com/reveever/biliapi)

部分 bilibili API SDK，大致定了个框架，挑着用得到的 API 先实现了

API 接口信息大多来自 [bilibili-API-collect](https://github.com/SocialSisterYi/bilibili-API-collect) 项目，~~B站官方代码结构和接口也是可以参考一下的~~

已实现的接口信息可以参考 godoc:

- [用户空间相关](https://pkg.go.dev/github.com/reveever/biliapi/interface/space)
- [视频弹幕相关](https://pkg.go.dev/github.com/reveever/biliapi/interface/v2dm)
- [直播间相关](https://pkg.go.dev/github.com/reveever/biliapi/interface/live)

## Example

引用：`go get -u github.com/reveever/biliapi`.

获取用户投稿
```go
package main

import (
	"encoding/json"
	"fmt"

	"github.com/reveever/biliapi"
	"github.com/reveever/biliapi/interface/space"
)

func main() {
	api, err := biliapi.NewBiliApi(biliapi.WithUserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:98.0) Gecko/20100101 Firefox/98.0"))
	if err != nil {
		panic(err)
	}

	resp, err := api.Space.ArcList(space.ArcListOpt{
		Mid: 67141, Pn: 1, Ps: 5,
	})
	if err != nil {
		panic(err)
	}

	buf, _ := json.Marshal(resp)
	fmt.Println(string(buf))
}
```

获取视频弹幕分片
```go
package main

import (
	"encoding/json"
	"fmt"

	"github.com/reveever/biliapi"
	"github.com/reveever/biliapi/interface/v2dm"
)

func main() {
	api, err := biliapi.NewBiliApi(biliapi.EnableDebugLogger())
	if err != nil {
		panic(err)
	}

	resp, err := api.V2Dm.WebSeg(v2dm.WebSegOpt{
		Type: 1, Oid: 343721233, SegmentIndex: 2,
	})
	if err != nil {
		panic(err)
	}

	buf, _ := json.Marshal(resp)
	fmt.Println(string(buf))
}
```

直播间信息流获取
```go
package main

import (
	"context"
	"io"
	"log"

	"github.com/reveever/biliapi"
	"github.com/reveever/biliapi/interface/live"
	"github.com/reveever/biliapi/interface/room"
)

func main() {
	api, err := biliapi.NewBiliApi(biliapi.EnableDebugLogger())
	if err != nil {
		log.Fatal(err)
	}

	roomInfo, err := api.Room.RoomInitInfo(room.RoomInitInfoOpt{ID: 213})
	if err != nil {
		log.Fatal(err)
	}

	liveInfo, err := api.Live.WebRoomInfo(live.WebRoomInfoOpt{ID: roomInfo.RoomID})
	if err != nil {
		log.Fatal(err)
	}

	l, err := biliapi.NewBiliLive(context.Background(), roomInfo.RoomID, liveInfo.Token,
		biliapi.LiveWithBase(api.Base),
		biliapi.LiveWithHostList(liveInfo.HostList),
	)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < 10; i++ {
		msg, err := l.ReadWsMsg(context.TODO())
		if err != nil {
			if err == io.EOF {
				log.Println("EOF")
				break
			}
			log.Printf("[%02d] err: %v", i, err)
			continue
		}
		if msg.Op == live.WsOpHeartbeatReply {
			log.Printf("[%02d] %d %d %s: %x", i, msg.Ver, msg.Op, msg.Cmd, msg.Body)
		} else {
			log.Printf("[%02d] %d %d %s: %s", i, msg.Ver, msg.Op, msg.Cmd, string(msg.Body))
		}
	}
	l.Close()
}
```
## Doc & More examples
[godoc](https://pkg.go.dev/github.com/reveever/biliapi)
