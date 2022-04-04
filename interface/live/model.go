package live

type WebRoomInfoOpt struct {
	ID   int64 `url:"id"`             // 直播间真实 id
	Type int   `url:"type,omitempry"` // ?类型, 默认为 0
}

type WebRoomInfoResp struct {
	Group            string         `json:"group"`
	BusinessID       int            `json:"business_id"`
	RefreshRowFactor float64        `json:"refresh_row_factor"`
	RefreshRate      int            `json:"refresh_rate"`
	MaxDelay         int            `json:"max_delay"`
	Token            string         `json:"token"`
	HostList         []LiveHostList `json:"host_list"`
}

type LiveHostList struct {
	Host    string `json:"host"`
	Port    int    `json:"port"`
	WssPort int    `json:"wss_port"`
	WsPort  int    `json:"ws_port"`
}
