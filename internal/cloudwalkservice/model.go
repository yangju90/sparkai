package cloudwalkservice

import (
	"strings"

	"github.com/google/uuid"
)

func CreateRequestBody(question string, imageData string, modelType string) map[string]interface{} {
	data := map[string]interface{}{
		"id":          strings.ReplaceAll(uuid.New().String(), "-", ""),
		"command":     modelType,
		"model":       "",
		"max_tokens":  1024,
		"temperature": 1,
		"stream":      false,
		"stops":       make([]string, 0),
		"messages": []map[string]interface{}{
			{
				"role": "user",
				"content": []map[string]interface{}{
					{
						"type": "text",
						"text": question,
					},
					{
						"type": "image_url",
						"image_url": map[string]string{
							"url": imageData,
						},
					},
				},
			},
		},
		"params": map[string]interface{}{
			"general_ocr_server_addr": "",
			"llm_server_addr":         "",
		},
	}

	return data
}
