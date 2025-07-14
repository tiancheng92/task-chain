package lark

import (
	"github.com/Yostardev/requests"
	"github.com/pkg/errors"
)

func updateMsg(messageID, content string) error {
	token, err := getToken()
	if err != nil {
		return err
	}

	res, err := requests.New().SetUrl("https://open.feishu.cn/open-apis/im/v1/messages/"+messageID).
		SetJsonBody(map[string]string{
			"content": content,
		}).
		AddJsonHeader().
		AddHeader("Authorization", "Bearer "+token).
		Patch()
	if err != nil {
		return errors.WithStack(err)
	}
	if res.StatusCode != 200 {
		return errors.New("feishu send message error: " + res.Body.String())
	}
	return nil
}
