package live

import (
	"github.com/reveever/biliapi/base"
)

var xLiveWebRoomV1 = base.APILiveEndpoint + "/xlive/web-room/v1"

type API struct {
	base *base.Base
}

func NewAPI(base *base.Base) *API {
	return &API{
		base: base,
	}
}

// 获取信息流认证秘钥
// https://github.com/SocialSisterYi/bilibili-API-collect/blob/master/live/message_stream.md#获取信息流认证秘钥
func (a *API) WebRoomInfo(opt WebRoomInfoOpt) (*WebRoomInfoResp, error) {
	resp := WebRoomInfoResp{}
	return &resp, a.base.GetJson(xLiveWebRoomV1+"/index/getDanmuInfo", opt, &resp)
}
