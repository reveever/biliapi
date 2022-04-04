package space

import (
	"github.com/reveever/biliapi/base"
)

var xSpace = base.APIEndpoint + "/x/space"

type API struct {
	base *base.Base
}

func NewAPI(base *base.Base) *API {
	return &API{
		base: base,
	}
}

// 列出用户投稿视频明细, 按照时间倒序
// https://-
func (a *API) ArcList(opt ArcListOpt) (*ArcListResp, error) {
	resp := ArcListResp{}
	return &resp, a.base.GetJson(xSpace+"/arc/list", opt, &resp)
}

// 查询用户投稿视频明细
// https://github.com/SocialSisterYi/bilibili-API-collect/blob/master/user/space.md#投稿
func (a *API) ArcSearch(opt ArcSearchOpt) (*ArcSearchResp, error) {
	resp := ArcSearchResp{}
	return &resp, a.base.GetJson(xSpace+"/arc/search", opt, &resp)
}
