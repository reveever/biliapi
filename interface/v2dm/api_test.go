package v2dm

import (
	"encoding/json"
	"testing"

	"github.com/reveever/biliapi/interface/test"
)

func newTestAPI(t *testing.T) *API {
	return NewAPI(test.NewTestBase(t))
}

func TestArcList(t *testing.T) {
	resp, err := newTestAPI(t).WebSeg(WebSegOpt{
		Type:         1,
		Oid:          test.Oid,
		Pid:          test.Pid,
		SegmentIndex: 2,
	})
	if err != nil {
		t.Fatal(err)
	}
	buf, _ := json.Marshal(resp)
	t.Log(string(buf))
}
