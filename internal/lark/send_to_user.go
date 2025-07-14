package lark

import (
	"github.com/Yostardev/requests"
	"github.com/pkg/errors"
)

type sendToUserResponse struct {
	Code int `json:"code"`
	Data struct {
		Body struct {
			Content string `json:"content"`
		} `json:"body"`
		ChatID     string `json:"chat_id"`
		CreateTime string `json:"create_time"`
		Deleted    bool   `json:"deleted"`
		MessageID  string `json:"message_id"`
		MsgType    string `json:"msg_type"`
		Sender     struct {
			ID         string `json:"id"`
			IDType     string `json:"id_type"`
			SenderType string `json:"sender_type"`
			TenantKey  string `json:"tenant_key"`
		} `json:"sender"`
		UpdateTime string `json:"update_time"`
		Updated    bool   `json:"updated"`
	} `json:"data"`
	Msg string `json:"msg"`
}

func sendToUser(userID, content string) (string, error) {
	token, err := getToken()
	if err != nil {
		return "", errors.WithStack(err)
	}

	res, err := requests.New().SetUrl("https://open.feishu.cn/open-apis/im/v1/messages?receive_id_type=user_id").
		SetJsonBody(map[string]string{
			"receive_id": userID,
			"msg_type":   "interactive",
			"content":    content,
		}).
		AddJsonHeader().
		AddHeader("Authorization", "Bearer "+token).
		Post()
	if err != nil {
		return "", errors.WithStack(err)
	}
	if res.StatusCode != 200 {
		return "", errors.New("feishu send message error: " + res.Body.String())
	}
	var resp sendToUserResponse
	err = res.Body.JsonBind(&resp)
	if err != nil {
		return "", errors.WithStack(err)
	}
	return resp.Data.MessageID, nil
}
