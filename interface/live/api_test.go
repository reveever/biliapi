package live

import (
	"encoding/json"
	"testing"

	"github.com/reveever/biliapi/interface/test"
)

func newTestAPI(t *testing.T) *API {
	return NewAPI(test.NewTestBase(t))
}

func TestWebRoomInfo(t *testing.T) {
	resp, err := newTestAPI(t).WebRoomInfo(WebRoomInfoOpt{
		ID:   test.RoomID,
		Type: 0,
	})
	if err != nil {
		t.Fatal(err)
	}
	buf, _ := json.Marshal(resp)
	t.Log(string(buf))
}
