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
		Model:  "qwen:14b",
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

func CreateHistoryOllamaReqBody(messages []model.Message) OllamaReqBody {
	request := OllamaReqBody{
		Model:    "qwen:14b",
		Stream:   true,
		Messages: messages,
	}

	return request
}

func CreateNewOllamaReqBody(messages []model.Message) OllamaReqBody {

	current := []model.Message{
		messages[0],
		messages[len(messages)-1],
	}

	request := OllamaReqBody{
		Model:    "qwen:14b",
		Stream:   true,
		Messages: current,
	}

	return request
}

func CreateNewOllamaReqBodyWithGeneralPrompt(messages []model.Message) OllamaReqBody {
	current := []model.Message{
		{
			Role:    constant.SYSTEM,
			Content: constant.FuncPromptConfig,
		},
		messages[len(messages)-1],
	}

	request := OllamaReqBody{
		Model:    "qwen:14b",
		Stream:   true,
		Messages: current,
	}

	return request
}
