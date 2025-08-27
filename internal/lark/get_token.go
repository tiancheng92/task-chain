package lark

import (
	"github.com/pkg/errors"
)

type tokenResp struct {
	Code              int    `json:"code"`
	Msg               string `json:"msg"`
	TenantAccessToken string `json:"tenant_access_token"`
	Expire            int    `json:"expire"`
}

func getToken() (string, error) {
	var resp tokenResp
	res, err := larkClient().R().SetResult(&resp).SetBody(map[string]string{
		"app_id":     appID,
		"app_secret": appSecret,
	}).Post("/open-apis/auth/v3/tenant_access_token/internal")
	if err != nil {
		return "", errors.WithStack(err)
	}
	if res.StatusCode() != 200 {
		return "", errors.New("get feishu token error: " + res.String())
	}

	return resp.TenantAccessToken, nil
}
