package lark

import (
	"github.com/Yostardev/requests"
	"github.com/pkg/errors"
)

type tokenResp struct {
	Code              int    `json:"code"`
	Msg               string `json:"msg"`
	TenantAccessToken string `json:"tenant_access_token"`
	Expire            int    `json:"expire"`
}

func getToken() (string, error) {
	res, err := requests.New().SetUrl("https://open.feishu.cn/open-apis/auth/v3/tenant_access_token/internal").
		SetJsonBody(map[string]string{
			"app_id":     appID,
			"app_secret": appSecret,
		}).Post()
	if err != nil {
		return "", errors.WithStack(err)
	}
	if res.StatusCode != 200 {
		return "", errors.New("get feishu token error: " + res.Body.String())
	}

	var resp tokenResp
	err = res.Body.JsonBind(&resp)
	if err != nil {
		return "", errors.WithStack(err)
	}

	return resp.TenantAccessToken, nil
}
