package lark

import (
	"sync"
	"time"

	"github.com/bytedance/sonic"
	"github.com/go-resty/resty/v2"
)

var (
	cc *resty.Client
	o  sync.Once
)

func larkClient() *resty.Client {
	o.Do(func() {
		cc = resty.New().
			SetJSONMarshaler(sonic.Marshal).
			SetJSONUnmarshaler(sonic.Unmarshal).
			SetTimeout(10*time.Second).
			SetHeader("Content-Type", "application/json").SetBaseURL("https://open.feishu.cn")
	})
	return cc
}
