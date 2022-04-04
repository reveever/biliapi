package base

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"github.com/google/go-querystring/query"
	"golang.org/x/net/publicsuffix"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var (
	APIEndpoint     = "https://api.bilibili.com"
	APILiveEndpoint = "https://api.live.bilibili.com"
	DefaultUA       = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.102 Safari/537.36 Edg/98.0.1108.62"
)

type Base struct {
	Client    *http.Client
	Log       BaseLogger
	UserAgent string
}

type BaseLogger interface {
	Println(v ...interface{})
	Printf(format string, v ...interface{})
}

func (b *Base) Init() error {
	if b == nil {
		return errors.New("nil base")
	}

	if b.Client == nil {
		jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
		if err != nil {
			return err
		}
		b.Client = &http.Client{Jar: jar}
	}

	if b.UserAgent == "" {
		b.UserAgent = DefaultUA
	}

	_, err := b.MakeRequest("GET", APIEndpoint+"/x/web-interface/nav", nil, nil)
	if err != nil {
		return fmt.Errorf("api test failed: %v", err)
	}

	return nil
}

func (b *Base) GetJson(url string, opt interface{}, result interface{}) error {
	resp, err := b.MakeRequest("GET", url, opt, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	contentTypes := resp.Header.Get("Content-Type")
	if !strings.Contains(contentTypes, "application/json") {
		buf, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("unexpected content type %s: %s", contentTypes, string(buf))
	}

	var apiResp APIResponse
	err = json.NewDecoder(resp.Body).Decode(&apiResp)
	if err != nil {
		return err
	}

	if apiResp.Code != 0 {
		if apiResp.Message == "" {
			apiResp.Message = ErrorMmsg(apiResp.Code)
		}
		var data string
		if apiResp.Data != nil {
			data = ": " + string(*apiResp.Data)
		}
		return fmt.Errorf("[%d] %s%s", apiResp.Code, apiResp.Message, data)
	}

	if apiResp.Data != nil {
		return json.Unmarshal(*apiResp.Data, result)
	}
	return nil
}

func (b *Base) GetPb(url string, opt interface{}, result protoreflect.ProtoMessage) error {
	resp, err := b.MakeRequest("GET", url, opt, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	contentTypes := resp.Header.Get("Content-Type")
	if !strings.Contains(contentTypes, "application/octet-stream") {
		buf, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("unexpected content type %s: %s", contentTypes, string(buf))
	}

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return proto.Unmarshal(buf, result)
}

func (b *Base) Post(path string, values url.Values, body []byte, result interface{}) error {
	return errors.New("TODO")
}

func (b *Base) MakeRequest(method, url string, opt interface{}, body []byte) (*http.Response, error) {
	if opt != nil {
		values, err := query.Values(opt)
		if err != nil {
			return nil, err
		}

		url += "?" + values.Encode()
	}

	if b.Log != nil {
		b.Log.Printf("%s: %s", method, url)
	}

	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	if b.UserAgent != "" {
		req.Header.Set("User-Agent", b.UserAgent)
	}

	resp, err := b.Client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		return nil, NewHttpNOK(resp)
	}

	return resp, err
}

type HttpNOK struct {
	StatusCode int
	Status     string
	Body       []byte
}

func NewHttpNOK(resp *http.Response) *HttpNOK {
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		buf = []byte(err.Error())
	}
	return &HttpNOK{
		StatusCode: resp.StatusCode,
		Status:     resp.Status,
		Body:       buf,
	}
}

func (h *HttpNOK) Error() string {
	return h.Status
}

func GetHttpNOK(err error) *HttpNOK {
	if v, ok := err.(*HttpNOK); ok {
		return v
	}
	return nil
}

func IsHttpNOK(err error, code int) bool {
	if err == nil {
		return false
	}
	v, ok := err.(*HttpNOK)
	if !ok || v.StatusCode != code {
		return false
	}
	return true
}
