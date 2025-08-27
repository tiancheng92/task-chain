package lark

import (
	"github.com/pkg/errors"
)

func updateMsg(messageID, content string) error {
	token, err := getToken()
	if err != nil {
		return err
	}

	res, err := larkClient().R().
		SetBody(map[string]string{
			"content": content,
		}).
		SetAuthToken(token).
		SetPathParam("message_id", messageID).
		Post("/open-apis/im/v1/messages/{message_id}")
	if err != nil {
		return errors.WithStack(err)
	}
	if res.StatusCode() != 200 {
		return errors.New("feishu send message error: " + res.String())
	}
	return nil
}
