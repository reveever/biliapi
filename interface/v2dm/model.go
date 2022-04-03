package v2dm

type WebSegOpt struct {
	Type         int   `url:"type"`          // 弹幕类: 1: 视频弹幕
	Oid          int64 `url:"oid"`           // 视频 cid
	Pid          int   `url:"pid,omitempty"` // 稿件 avid
	SegmentIndex int   `url:"segment_index"` // 分包, 6分钟一包
}
