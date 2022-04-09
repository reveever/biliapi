package room

import (
	"encoding/json"
	"testing"

	"github.com/reveever/biliapi/interface/test"
)

func newTestAPI(t *testing.T) *API {
	return NewAPI(test.NewTestBase(t))
}

func TestRoomInitInfo(t *testing.T) {
	resp, err := newTestAPI(t).RoomInitInfo(RoomInitInfoOpt{
		ID: test.ShortID,
	})
	if err != nil {
		t.Fatal(err)
	}
	buf, _ := json.Marshal(resp)
	t.Log(string(buf))
}
