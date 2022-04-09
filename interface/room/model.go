package room

type RoomInitInfoOpt struct {
	ID   int64  `url:"id"`             // 直播间短 ID / 直播间真实 ID
	Lang string `url:"lang,omitempty"` // ?语言: hant: 国际版繁体中文, hans: 国际版简体中文
}

type RoomInitInfoResp struct {
	RoomID      int64 `json:"room_id"`
	ShortID     int64 `json:"short_id"`
	Uid         int64 `json:"uid"`
	NeedP2P     int   `json:"need_p2p"`
	IsHidden    bool  `json:"is_hidden"`
	IsLocked    bool  `json:"is_locked"`
	IsPortrait  bool  `json:"is_portrait"`
	LiveStatus  int   `json:"live_status"`
	HiddenTill  int   `json:"hidden_till"`
	LockTill    int   `json:"lock_till"`
	Encrypted   bool  `json:"encrypted"`
	PwdVerified bool  `json:"pwd_verified"`
	LiveTime    int   `json:"live_time"`
	RoomShield  int   `json:"room_shield"`
	IsSP        int   `json:"is_sp"`
	SpecialType int   `json:"special_type"`
}
