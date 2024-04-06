package cloudwalkservice

func CreateRequestBody(question string, imageData string) map[string]interface{} {
	data := map[string]interface{}{
		"id":          "b31c92228a264137-90bb01a8a0e217a7",
		"command":     "lmm",
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
