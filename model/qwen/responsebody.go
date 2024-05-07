package qwen

import "sparkai/model"

type OllamaRespBody struct {
	Model              string `json:"model"`
	CreatedAt          bool   `json:"created_at"`
	Done               bool   `json:"done"`
	TotalDuration      int64  `json:"total_duration"`
	LoadDuration       int64  `json:"load_duration"`
	PromptEvalCount    int32  `json:"prompt_eval_count"`
	PromptEvalDuration int32  `json:"prompt_eval_duration"`
	EvalCount          int16  `json:"eval_count"`
	EvalDuration       int64  `json:"eval_duration"`

	Messages []model.Message `json:"messages"`
}
