package functionsProcess

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"sparkai/model"
	"sparkai/model/constant"
	"sparkai/model/mem"

	"github.com/gorilla/websocket"
)

func BigModelFunc(userId string) (string, error) {
	// sessionId 处理历史数据

	d := websocket.Dialer{
		HandshakeTimeout: 5 * time.Second,
	}
	//握手并建立websocket 连接
	conn, resp, err := d.Dial(assembleAuthUrl1(constant.WssConfig.HostUrl, constant.WssConfig.ApiKey, constant.WssConfig.ApiSecret), nil)

	if err != nil {
		return "", err
	} else if resp.StatusCode != 101 {
		log.Println(resp)
		return "", errors.New("AI错误的状态码")
	}

	defer conn.Close()

	// goroutine 方法调用，可以维持wss通道不失效
	// c := make(chan int)
	// defer func() {
	// 	c <- 0
	// }()

	// go heartbeat(c, conn)

	go waitUserInput(conn, constant.WssConfig.Appid, userId)

	return WaitSparkaiOutput(conn, userId)
}

func WaitSparkaiOutput(conn *websocket.Conn, userId string) (string, error) {
	var answer = ""

	v, ok := mem.WSConnContainers[userId]
	if !ok {
		panic("Id为" + userId + "的用户不在线")
	}

	//获取返回的数据
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("read message error:", err)
			return "", err
		}

		var data map[string]interface{}
		err1 := json.Unmarshal(msg, &data)
		if err1 != nil {
			fmt.Println("Error parsing JSON:", err)
			return "", err1
		}
		log.Println(string(msg))
		//解析数据
		payload := data["payload"].(map[string]interface{})
		choices := payload["choices"].(map[string]interface{})
		header := data["header"].(map[string]interface{})
		code := header["code"].(float64)

		if code != 0 {
			log.Println(data["payload"])
			return "", errors.New("sparkai response err")
		}
		status := choices["status"].(float64)
		text := choices["text"].([]interface{})
		content := text[0].(map[string]interface{})["content"].(string)

		var wsResponse model.WSBodyResponse
		wsResponse.Code = int(code)

		wsResponse.ContentType = "text"
		wsResponse.Status = responseConvert("text", int(status))
		wsResponse.Content = content

		if wsResponse.Status != 9 {
			answer += content
		} else {
			answer += content
			usage := payload["usage"].(map[string]interface{})
			temp := usage["text"].(map[string]interface{})
			totalTokens := temp["total_tokens"].(float64)
			fmt.Println("total_tokens:", totalTokens)
			conn.Close()
		}

		textByteData, err := json.Marshal(wsResponse)
		if err == nil {
			if e := v.Send(textByteData); e != nil {
				return "", e
			}
		} else {
			return "", err
		}

		if wsResponse.Status == 9 {
			break
		}
	}
	return answer, nil
}

func waitUserInput(conn *websocket.Conn, appid string, userId string) {
	if v, ok := mem.WSConnContainers[userId]; ok {
		message := v.Messages[len(v.Messages)-1]
		data := genParams(appid, userId, v.ChatId, message.Content)

		byteData, err := json.Marshal(data)
		if err != nil {
			log.Println("map cast byte error ", err)
			return
		}

		log.Println("BigModel发送数据：" + string(byteData))
		conn.WriteMessage(websocket.TextMessage, byteData)
	} else {
		panic("Id为" + userId + "的用户不在线")
	}

	// conn.WriteJSON(data)
}

func heartbeat(c chan int, conn *websocket.Conn) {
	for {
		time.Sleep(30 * time.Second) // 每隔 30 秒发送一次心跳消息
		select {
		case <-c:
			log.Println("心跳检测关闭")
			return
		default:
			err := conn.WriteMessage(websocket.PingMessage, []byte{})
			if err != nil {
				log.Println("发送心跳消息时发生错误:", err)
				return
			}
			log.Println("发送心跳.....")
		}
	}
}

// 创建鉴权url  apikey 即 hmac username
func assembleAuthUrl1(hosturl string, apiKey, apiSecret string) string {
	ul, err := url.Parse(hosturl)
	if err != nil {
		fmt.Println(err)
	}
	date := time.Now().UTC().Format(time.RFC1123)
	signString := []string{"host: " + ul.Host, "date: " + date, "GET " + ul.Path + " HTTP/1.1"}
	sgin := strings.Join(signString, "\n")
	// fmt.Println(sgin)
	sha := HmacWithShaTobase64("hmac-sha256", sgin, apiSecret)
	// fmt.Println(sha)
	//构建请求参数
	authUrl := fmt.Sprintf("hmac username=\"%s\", algorithm=\"%s\", headers=\"%s\", signature=\"%s\"", apiKey,
		"hmac-sha256", "host date request-line", sha)
	//将请求参数使用base64编码
	authorization := base64.StdEncoding.EncodeToString([]byte(authUrl))

	v := url.Values{}
	v.Add("host", ul.Host)
	v.Add("date", date)
	v.Add("authorization", authorization)
	//将编码后的字符串url encode后添加到url后面
	callurl := hosturl + "?" + v.Encode()
	return callurl
}

func HmacWithShaTobase64(algorithm, data, key string) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(data))
	encodeData := mac.Sum(nil)
	return base64.StdEncoding.EncodeToString(encodeData)
}

func readResp(resp *http.Response) string {
	if resp == nil {
		return ""
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("code=%d,body=%s", resp.StatusCode, string(b))
}

func responseConvert(contentType string, code int) int {
	var res int = 9
	if contentType == "text" {
		switch code {
		case 0:
			res = 0
		case 1:
			res = 1
		case 2:
			res = 9
		}
	} else if contentType == "function" {
		return code
	}
	return res
}

// 生成参数
func genParams(appid string, uid string, chat_id string, text string) map[string]interface{} {
	messages := []model.Message{
		{
			Role:    constant.SYSTEM,
			Content: constant.GeneralPromptConfig,
		},
		{
			Role:    constant.USER,
			Content: text,
		},
	}

	data := map[string]interface{}{
		"header": map[string]interface{}{
			"app_id": appid,
			"uid":    uid,
		},
		"parameter": map[string]interface{}{
			"chat": map[string]interface{}{
				"domain":      "generalv3.5",
				"temperature": float64(0.5),
				"top_k":       int64(4),
				"max_tokens":  int64(8192),
				"chat_id":     chat_id,
			},
		},
		"payload": map[string]interface{}{
			"message": map[string]interface{}{
				"text": messages,
			},
		},
	}

	return data
}
