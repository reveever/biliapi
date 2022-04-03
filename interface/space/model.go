package space

import (
	"github.com/reveever/biliapi/proto/archive"
)

type ArcListOpt struct {
	Mid int64 `url:"mid"` // 目标用户 mid
	Pn  int   `url:"pn"`  // 页码 >=1
	Ps  int   `url:"ps"`  // 每页项数 <=50
}

type ArcListResp struct {
	Page     ArcPage        `json:"page"`
	Archives []*ArcListItem `json:"archives"`
}

type ArcPage struct {
	Pn    int32 `json:"pn"`
	Ps    int32 `json:"ps"`
	Count int64 `json:"count"`
}

type ArcListItem struct {
	Aid      int64  `json:"aid"`
	Bvid     string `json:"bvid"`
	Pic      string `json:"pic"`
	Title    string `json:"title"`
	Duration int64  `json:"duration"`
	Author   struct {
		Mid  int64  `json:"mid"`
		Name string `json:"name"`
		Face string `json:"face"`
	} `json:"author"`
	Stat struct {
		View     int `json:"view"`
		Danmaku  int `json:"danmaku"`
		Reply    int `json:"reply"`
		Favorite int `json:"favorite"`
		Coin     int `json:"coin"`
		Share    int `json:"share"`
		Like     int `json:"like"`
	} `json:"stat"`
	Rights  archive.Rights `json:"rights"`
	Pubdate int64          `json:"pubdate"`
}

type ArcSearchOpt struct {
	Mid       int64  `url:"mid"`                  // 目标用户 mid
	Tid       int64  `url:"tid,omitempty"`        // 筛选目标分区: 不进行分区筛选(默认): 0, 所筛选的分区
	Order     string `url:"order,omitempty"`      // 排序方式: 最新发布(默认): pubdate, 最多播放: click, 最多收藏: stow
	Keyword   string `url:"keyword,omitempty"`    // 用于使用关键词搜索该 UP 主视频稿件
	Pn        int    `url:"pn" `                  // 页码 >=0
	Ps        int    `url:"ps"`                   // 每页项数 <=50
	CheckType string `url:"check_type,omitempty"` // ?校验类型: channel
	CheckID   int64  `url:"check_id,omitempty"`   // ?校验 CID >0
	Jsonp     string `url:"jsonp,omitempty"`      // ?Jsonp: jsonp
}

type ArcSearchResp struct {
	Page           ArcPage           `json:"page"`
	List           ArcSearchRes      `json:"list"`
	EpisodicButton ArcEpisodicButton `json:"episodic_button"`
}

type ArcSearchRes struct {
	Tlist map[string]ArcSearchTList `json:"tlist"`
	Vlist []ArcSearchVList          `json:"vlist"`
}

type ArcSearchTList struct {
	Tid   int64  `json:"tid"`
	Count int64  `json:"count"`
	Name  string `json:"name"`
}

type ArcSearchVList struct {
	Comment        int64  `json:"comment"`
	TypeID         int64  `json:"typeid"`
	Play           int64  `json:"play"`
	Pic            string `json:"pic"`
	Subtitle       string `json:"subtitle"`
	Description    string `json:"description"`
	Copyright      string `json:"copyright"`
	Title          string `json:"title"`
	Review         int64  `json:"review"`
	Author         string `json:"author"`
	Mid            int64  `json:"mid"`
	Created        int64  `json:"created"`
	Length         string `json:"length"`
	VideoReview    int64  `json:"video_review"`
	Aid            int64  `json:"aid"`
	Bvid           string `json:"bvid"`
	HideClick      bool   `json:"hide_click"`
	IsPay          int    `json:"is_pay"`
	IsUnionVideo   int    `json:"is_union_video"`
	IsSteinsGate   int    `json:"is_steins_gate"`
	IsLivePlayback int    `json:"is_live_playback"`
}

type ArcEpisodicButton struct {
	Text string `json:"text"`
	URI  string `json:"uri"`
}
