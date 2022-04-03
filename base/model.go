package base

import "encoding/json"

type APIResponse struct {
	Code    int              `json:"code"`
	Message string           `json:"message"`
	TTL     int              `json:"ttl"`
	Data    *json.RawMessage `json:"data"`
}
