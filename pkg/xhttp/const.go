package xhttp

import (
	"yes4all/ads-noti-api/pkg/utils/constants"
)

const (
	GgChatWebhooksErrorEvent  = "https://chat.googleapis.com/v1/spaces/AAAAIgLtFS8/messages?key=AIzaSyDdI0hCZtE6vySjMm-WEfRq3CPzqKqqsHI&token=YMiHvU43mxiwip93CaLp8Dy5o6yEIS-g0RsncyqgCho"
	ggChatWebhooksHttpTimeout = "https://chat.googleapis.com/v1/spaces/AAAADwKN0wA/messages?key=AIzaSyDdI0hCZtE6vySjMm-WEfRq3CPzqKqqsHI&token=Ph97Xg_Yh9nYIitBrgzk86U7SQD9vUEnLVOdUr3_7Jc"
)

func GetEnviromentHttpTimeout(env string) (weedHook string) {
	switch env {
	case constants.ProductionEnv:
		return GgChatWebhooksErrorEvent
	default:
		return ggChatWebhooksHttpTimeout
	}
}

type ggResponse struct {
	Name string `json:"name"`
	Text string `json:"text"`
}

type dataResponse struct {
	URL       string `json:"url"`
	Method    string `json:"method"`
	RequestId string `json:"y4aMess"`
}
