package functionsProcess

import (
	"encoding/json"
	"log"
	"sparkai/image/subImage"
	"sparkai/internal/cloudwalkservice"
	"sparkai/model"
	"sparkai/model/mem"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

var colors = []string{"红色", "绿色", "蓝色", "黄色", "品红", "青色"}

func ChoiceFuntionCall(name string, userId string) error {
	log.Println("调用function call : " + name)
	var resErr error
	switch name {
	case AUDIO_DETECTION:
		sendCallFunctionSignal(AUDIO_DETECTION, userId)
		url := audioChannelChoice(userId)
		var text string
		if len(url) == 0 {
			text = "调用的设备不在线，请重新选择！"
		} else {
			text = AUDIO_DETECTION + "已调用！"
		}
		resErr = sendMsg(userId, text, url, 4)
	case IMAGE_UNDERSTANDING:
		b := perImageProcess(userId)
		if b {
			sendCallFunctionSignal(IMAGE_UNDERSTANDING, userId)
			imageUnderstanding, err := cloudwalkservice.Service(userId)
			if err != nil {
				return err
			} else {
				resErr = sendMsg(userId, imageUnderstanding, "", 2)
			}
		} else {
			_, resErr = BigModelFunc(userId)
		}

	case OBJECT_RECOGNITION:
		b := perImageProcess(userId)
		if b {
			sendCallFunctionSignal(OBJECT_RECOGNITION, userId)
			objects, err := cloudwalkservice.ServiceRecognition(userId)
			if err != nil {
				return err
			} else {
				text, path, e := fiterResultAndDrawPic(userId, objects)
				if e != nil {
					resErr = sendMsg(userId, OBJECT_RECOGNITION+"调用图像生成错误！", "", 5)
				} else {
					resErr = sendMsg(userId, text, path, 5)
				}
			}
		} else {
			_, resErr = BigModelFunc(userId)
		}

	case ROUTER_PLANING:
		sendCallFunctionSignal(ROUTER_PLANING, userId)
		resErr = sendMsg(userId, ROUTER_PLANING+"已调用！", "", 5)
	default:
		_, resErr = BigModelFunc(userId)
	}

	if v, ok := mem.WSConnContainers[userId]; ok {
		v.Messages = v.Messages[:len(v.Messages)-1]
	}
	return resErr
}

func sendCallFunctionSignal(funcName string, userId string) error {
	if v, ok := mem.WSConnContainers[userId]; ok {
		wsResponse := model.WSBodyResponse{
			Code:        int(0),
			Status:      3,
			Content:     "【功能调度 " + funcName + "】\n",
			ContentType: "function",
		}
		ccc, _ := json.Marshal(wsResponse)
		if e := v.Send(ccc); e != nil {
			return e
		}
	} else {
		panic("Id为" + userId + "的用户不在线")
	}

	return nil
}

func sendMsg(userId string, text string, url string, status int) error {
	if v, ok := mem.WSConnContainers[userId]; ok {
		wsResponse := model.WSBodyResponse{
			Code:        int(0),
			Status:      status,
			Content:     text,
			ContentType: CodeToType(status),
			Url:         url,
		}
		ccc, _ := json.Marshal(wsResponse)
		if e := v.Send(ccc); e != nil {
			return e
		}
	} else {
		panic("Id为" + userId + "的用户不在线")
	}

	return nil
}

func perImageProcess(userId string) bool {
	if v, ok := mem.WSConnContainers[userId]; ok {
		if len(v.ImageData) == 0 {
			return false
		}
	} else {
		return false
	}

	return true
}

func fiterResultAndDrawPic(userId string, res []interface{}) (text string, path string, e error) {
	// 过滤有效的识别
	drawObject := make(map[string][]interface{})
	for _, v := range res {
		r := v.(map[string]interface{})
		cls := r["cls"].(string)
		box := r["box"]
		score := r["score"].(float64)
		if score > 0.6 {
			if v, ok := drawObject[cls]; ok {
				drawObject[cls] = append(v, box)
			} else {
				t := []interface{}{box}
				drawObject[cls] = t
			}
		}
	}

	// 识别画图，生成
	id := uuid.New().String()
	path = "http://localhost:8090/image/object/" + id + ".png"
	if v, ok := mem.WSConnContainers[userId]; ok {
		subImage.DrawRectanglePic(v.ImageData, id, drawObject)
	} else {
		panic("Id为" + userId + "的用户不在线")
	}

	text = countNumText(drawObject)

	return
}

func countNumText(m map[string][]interface{}) string {
	var res string
	colorIndex := 0
	res += "<table><thead><tr><th>类别</th><th>颜色</th><th>数量</th></tr></thead><tbody>"
	res += "<tr><td></td><td></td>"
	for k, v := range m {
		res += "<tr><td>" + k + "</td><td>" + colors[colorIndex] + "</td><td>" + strconv.Itoa(len(v)) + "</td><tr>"
		colorIndex++
		if colorIndex == 6 {
			break
		}
	}
	res += "</tbody></table>"
	return res
}

func audioChannelChoice(userId string) string {
	if v, ok := mem.WSConnContainers[userId]; ok {
		message := v.Messages[len(v.Messages)-1]
		content := message.Content

		if strings.Contains(content, "无人机") {
			return "1"
		} else if strings.Contains(content, "摄像头") {
			return "2"
		}
	} else {
		panic("Id为" + userId + "的用户不在线")
	}
	return ""
}
