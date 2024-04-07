package functionsProcess

import (
	"encoding/json"
	"log"
	"sparkai/internal/cloudwalkservice"
	"sparkai/model"
	"sparkai/model/mem"
	"strings"
)

func ChoiceFuntionCall(name string, userId string) error {
	log.Println("调用function call : " + name)
	var resErr error
	switch name {
	case AUDIO_DETECTION:
		url := audioChannelChoice(userId)
		var text string
		if len(url) == 0 {
			text = "调用的设备不在线，请重新选择！"
		} else {
			text = AUDIO_DETECTION + "已调用！"
		}
		resErr = sendMsg(userId, text, url, 4)
	case IMAGE_UNDERSTANDING:
		imageUnderstanding, err := cloudwalkservice.Service(userId)
		if err != nil {
			_, resErr = BigModelFunc(userId)
		} else {
			resErr = sendMsg(userId, imageUnderstanding, "", 2)
		}
	case OBJECT_RECOGNITION:
		resErr = sendMsg(userId, OBJECT_RECOGNITION+"已调用！", "", 5)
	case ROUTER_PLANING:
		resErr = sendMsg(userId, ROUTER_PLANING+"已调用！", "", 5)
	default:
		_, resErr = BigModelFunc(userId)
	}

	if v, ok := mem.WSConnContainers[userId]; ok {
		v.Messages = v.Messages[:len(v.Messages)-1]
	}
	return resErr
}

func sendMsg(userId string, text string, rtspUrl string, status int) error {
	if v, ok := mem.WSConnContainers[userId]; ok {
		wsResponse := model.WSBodyResponse{
			Code:        int(0),
			Status:      status,
			Content:     text,
			ContentType: CodeToType(status),
			RTSPUrl:     rtspUrl,
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
