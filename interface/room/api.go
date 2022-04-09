package room

import (
	"github.com/reveever/biliapi/base"
)

var roomV1 = base.APILiveEndpoint + "/room/v1"

type API struct {
	base *base.Base
}

func NewAPI(base *base.Base) *API {
	return &API{
		base: base,
	}
}

// 获取房间页初始化信息
// https://github.com/SocialSisterYi/bilibili-API-collect/blob/master/live/info.md#获取房间页初始化信息
func (a *API) RoomInitInfo(opt RoomInitInfoOpt) (*RoomInitInfoResp, error) {
	resp := RoomInitInfoResp{}
	return &resp, a.base.GetJson(roomV1+"/Room/room_init", opt, &resp)
}
