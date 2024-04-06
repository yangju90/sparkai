package functionsProcess

import (
	"log"
	"sparkai/internal/cloudwalkservice"
)

func ChoiceFuntionCall(name string, userId string) (string, bool, error) {
	log.Println("调用function call : " + name)

	var res string
	next := true
	var resErr error
	switch name {
	case BID_MODEL:
		res = BID_MODEL + "，调用完成！"
	case AUDIO_DETECTION:
		res = AUDIO_DETECTION + "，调用完成！"
	case IMAGE_UNDERSTANDING:
		imageUnderstanding, err := cloudwalkservice.Service(userId)
		if err != nil {
			bigModelRes, err1 := BigModelFunc(userId)
			if err1 != nil {
				resErr = err1
			} else {
				res = bigModelRes
				next = false
			}
		} else {
			res = imageUnderstanding
		}
	case OBJECT_RECOGNITION:
		res = OBJECT_RECOGNITION + "，调用完成！"
	}

	return res, next, resErr
}
