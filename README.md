# biliapi
[![Go Reference](https://pkg.go.dev/badge/github.com/reveever/biliapi.svg)](https://pkg.go.dev/github.com/reveever/biliapi)

部分 bilibili API SDK，大致定了个框架，挑着用得到的 API 先实现了

API 接口信息大多来自 [bilibili-API-collect](https://github.com/SocialSisterYi/bilibili-API-collect) 项目

~~B站官方代码结构和接口也是可以参考一下的~~

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
## Doc & More examples
[godoc](https://pkg.go.dev/github.com/reveever/biliapi)
