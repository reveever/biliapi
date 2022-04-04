package v2dm

import (
	"net/http"

	"github.com/reveever/biliapi/base"
	"github.com/reveever/biliapi/proto/dm"
)

var xV2Dm = base.APIEndpoint + "/x/v2/dm"

type API struct {
	base *base.Base
}

func NewAPI(base *base.Base) *API {
	return &API{
		base: base,
	}
}

// 获取实时弹幕(网页)新
// https://github.com/SocialSisterYi/bilibili-API-collect/blob/master/danmaku/danmaku_proto.md#获取实时弹幕
func (a *API) WebSeg(opt WebSegOpt) (*dm.DmSegMobileReply, error) {
	resp := dm.DmSegMobileReply{}
	err := a.base.GetPb(xV2Dm+"/web/seg.so", opt, &resp)
	if base.IsHttpNOK(err, http.StatusNotModified) {
		return nil, nil
	}
	return &resp, err
}
