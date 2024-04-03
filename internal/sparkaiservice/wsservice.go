package sparkaiservice

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"sparkai/internal/io"
	"sparkai/model/constant"

	"github.com/gorilla/websocket"
)

func Wsservice(sessionId string, text string) {
	// sessionId 处理历史数据

	d := websocket.Dialer{
		HandshakeTimeout: 5 * time.Second,
	}
	//握手并建立websocket 连接
	conn, resp, err := d.Dial(assembleAuthUrl1(constant.WssConfig.HostUrl, constant.WssConfig.ApiKey, constant.WssConfig.ApiSecret), nil)
	if err != nil {
		panic(readResp(resp) + err.Error())
		return
	} else if resp.StatusCode != 101 {
		panic(readResp(resp) + err.Error())
	}

	defer conn.Close()

	// goroutine 方法调用，可以维持wss通道不失效
	// c := make(chan int)
	// defer func() {
	// 	c <- 0
	// }()

	// go heartbeat(c, conn)

	go io.WaitUserInput(conn, constant.WssConfig.Appid, text)

	io.WaitSparkaiOutput(conn)
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
