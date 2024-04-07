package functionsProcess

const (
	BID_MODEL           string = "大模型"
	AUDIO_DETECTION     string = "视频检测"
	IMAGE_UNDERSTANDING string = "图片问答"
	OBJECT_RECOGNITION  string = "目标识别"
	ROUTER_PLANING      string = "路径规划"
)

const (
	RTSPUrl string = "rtsp://127.0.0.1/road_2k_fusion.sdp"
)

func CodeToType(code int) string {
	var res string
	switch code {
	case 0:
		res = "text"
	case 1:
		res = "text"
	case 2:
		res = "text"
	case 3:
		res = "function"
	case 4:
		res = "audio"
	case 5:
		res = "image"
	case 6:
		res = "table"
	case 7:

	case 8:

	case 9:
		res = "text"
	}

	return res
}
