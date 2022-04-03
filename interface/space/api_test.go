package space

import (
	"encoding/json"
	"testing"

	"github.com/reveever/biliapi/interface/test"
)

func newTestAPI(t *testing.T) *API {
	return NewAPI(test.NewTestBase(t))
}

func TestArcList(t *testing.T) {
	resp, err := newTestAPI(t).ArcList(ArcListOpt{
		Mid: test.Mid,
		Pn:  1,
		Ps:  5,
	})
	if err != nil {
		t.Fatal(err)
	}
	buf, _ := json.Marshal(resp)
	t.Log(string(buf))
}

func TestArcSearch(t *testing.T) {
	resp, err := newTestAPI(t).ArcSearch(ArcSearchOpt{
		Mid:     test.Mid,
		Tid:     1,
		Order:   "click",
		Keyword: "C",
		Pn:      1,
		Ps:      5,
		// CheckType: "channel", // recv -1200
		// CheckID:   1,
		Jsonp: "jsonp",
	})
	if err != nil {
		t.Fatal(err)
	}
	buf, _ := json.Marshal(resp)
	t.Log(string(buf))
}
