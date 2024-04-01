package io

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func WaitCommandInput() (string, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("请输入文本: ")
	text, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("读取输入时发生错误: %v", err)
	}

	// 去除末尾的换行符
	text = strings.TrimSpace(text)

	return text, nil
}
