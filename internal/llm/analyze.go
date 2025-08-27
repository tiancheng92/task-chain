package llm

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/schema"
	"github.com/pkg/errors"
)

func AnalyzeError(errorInfo string) (string, error) {
	if errorInfo == "" {
		return "", nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()
	model, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		BaseURL: baseUrl,
		APIKey:  apiKey,
		Model:   modelName,
	})
	if err != nil {
		return "", errors.WithStack(err)
	}
	prompt := fmt.Sprintf(`你是一位资深多云平台SRE专家，请分析以下错误信息:
---
%s
---

要求：
1. 请用简体中文言简意赅的语言简述错误内容，而不要简单的翻译错误信息
2. 信息不存在时使用无而非猜测值
`, errorInfo)
	messages := []*schema.Message{
		schema.SystemMessage("你是企业级云平台错误处理专家，擅长分析各个云接口返回的各种错误。"),
		schema.UserMessage(prompt),
	}

	response, err := model.Generate(ctx, messages)
	if err != nil {
		return "", errors.WithStack(err)
	}

	return response.Content, nil
}
