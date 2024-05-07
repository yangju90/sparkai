package qwen

import (
	"sparkai/model"
	"sparkai/model/constant"
)

type OllamaReqBody struct {
	Model    string          `json:"model"`
	Stream   bool            `json:"stream"`
	Messages []model.Message `json:"messages"`
}

func NewOllamaReqBody() OllamaReqBody {
	request := OllamaReqBody{
		Model:  "qwen:7b",
		Stream: true,
		Messages: []model.Message{
			{
				Role:    constant.USER,
				Content: "你是谁？",
			},
		},
	}

	return request
}
