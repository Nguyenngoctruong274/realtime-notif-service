package ggchat

import (
	"context"
	"fmt"
	"os"
	"time"
	"yes4all/ads-noti-api/pkg/utils/constants"
	"yes4all/ads-noti-api/pkg/xhttp"

	"go.elastic.co/apm"
)

type IGGchat interface {
	SenGGchatCustom(ctx context.Context, messages string, path string, threadName string) (err error)
	SendGGchatErrorEvent(ctx context.Context, errMessage string, topic string)
	SendGGChatAuto(ctx context.Context, errMessage interface{}, topic string, isStart bool)
}

type ggChat struct {
	client xhttp.Client
}

func NewClient(client xhttp.Client) IGGchat {
	return &ggChat{
		client: client,
	}
}

const (
	GgChatWebhooksErrorBudget   = "https://chat.googleapis.com/v1/spaces/AAAA7NMyrAI/messages?key=AIzaSyDdI0hCZtE6vySjMm-WEfRq3CPzqKqqsHI&token=GuuEisM9oXeUNdbCVvcb3mfkrfBugUgKB3kmHRc7Dwg"
	GgChatErrorBudgetThreadName = "spaces/AAAA7NMyrAI/threads/SJMDrx-Re5U"

	GgChatWebhooksErrorEvent = "https://chat.googleapis.com/v1/spaces/AAAA7IjrSf0/messages?key=AIzaSyDdI0hCZtE6vySjMm-WEfRq3CPzqKqqsHI&token=c2OwxZMEV2l6cSV4ALINxamTpFnNfLBW7ETgmzoVSyI"

	GgChatWebhooksErrorAuto     = "https://chat.googleapis.com/v1/spaces/AAAABDPZt1Y/messages?key=AIzaSyDdI0hCZtE6vySjMm-WEfRq3CPzqKqqsHI&token=9PCIetIh2KZZf6snm74viQ1e3WnZbNe3vBd8HASKOu0"
	GgChatWebhooksErrorAutoTest = "https://chat.googleapis.com/v1/spaces/AAAAW60qHPk/messages?key=AIzaSyDdI0hCZtE6vySjMm-WEfRq3CPzqKqqsHI&token=GpDl7PNVlyIxy2gXO9gCCl6EYVSrSkpIAT1qYohY0qg"
	//https://chat.googleapis.com/v1/spaces/AAAAW60qHPk/messages?key=AIzaSyDdI0hCZtE6vySjMm-WEfRq3CPzqKqqsHI&token=GpDl7PNVlyIxy2gXO9gCCl6EYVSrSkpIAT1qYohY0qg
)

const (
	GgChatWebhooksErrorBudgetTest   = "https://chat.googleapis.com/v1/spaces/AAAAW60qHPk/messages?key=AIzaSyDdI0hCZtE6vySjMm-WEfRq3CPzqKqqsHI&token=GpDl7PNVlyIxy2gXO9gCCl6EYVSrSkpIAT1qYohY0qg"
	GgChatErrorBudgetThreadNameTest = "spaces/AAAAW60qHPk/threads/Q2ncJNrBRz8"
)

func GetThreadByENV(env string) (weedHook, thread string) {
	switch env {
	case constants.ProductionEnv:
		return GgChatWebhooksErrorBudget, GgChatErrorBudgetThreadName
	default:
		return GgChatWebhooksErrorBudgetTest, GgChatErrorBudgetThreadNameTest
	}
}

func GetThreadErrorEvent(env string) (weedHook string) {
	switch env {
	case constants.ProductionEnv:
		return GgChatWebhooksErrorEvent
	default:
		return GgChatWebhooksErrorEvent
		return ""
	}
}

func GetThreadErrorAuto(env string) (weedHook string) {
	switch env {
	case constants.ProductionEnv:
		return GgChatWebhooksErrorAuto
	default:
		return GgChatWebhooksErrorAutoTest
	}
}

type ggResponse struct {
	Name string `json:"name"`
	Text string `json:"text"`
}

func (c *ggChat) SenGGchatCustom(
	ctx context.Context,
	message string,
	path string,
	threadName string,
) (err error) {
	type thread struct {
		Name string `json:"name"`
	}
	messageData := struct {
		Text   string  `json:"text"`
		Thread *thread `json:"thread"`
	}{
		Text: message,
	}

	if threadName == "" {
		messageData.Thread = &thread{
			Name: threadName,
		}
	}
	var returnValue ggResponse
	xopt := xhttp.RequestOption{GroupPath: "/chat.googleapis.com/v1/spaces"}
	if _, err = c.client.PostJSON(ctx, path, &messageData, &returnValue, xopt); err != nil {
		apm.CaptureError(context.Background(), err).Send()
		return
	}

	return
}

func (c *ggChat) SendGGchatErrorEvent(ctx context.Context, errMessage string, topic string) {
	// send notification to google chat after begin calculation
	env := os.Getenv(constants.EnvironmentsEnv)
	path := GetThreadErrorEvent(env)
	if env == "" {
		return
	}
	c.SenGGchatCustom(ctx,
		fmt.Sprintf("```Log data started at: %v \n \n error: %v \n topic: %v \n env: %v```", time.Now(), errMessage, topic, env),
		path, "",
	)
}

func (c *ggChat) SendGGChatAuto(ctx context.Context, errMessage interface{}, topic string, isStart bool) {
	// send enum_notification to google chat after begin calculation
	env := os.Getenv(constants.EnvironmentsEnv)
	path := GetThreadErrorAuto(env)
	if env == "" {
		return
	}
	message := fmt.Sprintf("Log data ended at: %v \n \n error: %+v \n topic: %v \n env: %v",
		time.Now(), errMessage, topic, env,
	)
	if isStart {
		message = fmt.Sprintf("Log data started at: %v \n \n topic: %v \n env: %v",
			time.Now(), topic, env,
		)
	}

	messageTemp := ""
	limit := 4000
	for _, v := range message {
		messageTemp += string(v)
		if len(messageTemp) >= limit {
			messageTemp = "```" + messageTemp + "```"
			c.SenGGchatCustom(ctx,
				messageTemp,
				path, "",
			)
			messageTemp = ""
		}
	}
	if len(messageTemp) != 0 {
		messageTemp = "```" + messageTemp + "```"
		c.SenGGchatCustom(ctx,
			messageTemp,
			path, "",
		)
	}

}
