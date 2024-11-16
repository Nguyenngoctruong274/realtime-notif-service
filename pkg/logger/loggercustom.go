package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/spf13/cast"
	"net/http"
	"os"
	"time"
	"yes4all/ads-noti-api/pkg/utils/constants"
	"yes4all/ads-noti-api/pkg/utils/timeutils"
)

const (
	GroupPathHeader          = "X-Group-Path"
	GgChatWebhooksErrorDB    = "https://chat.googleapis.com/v1/spaces/AAAADwKN0wA/messages?key=AIzaSyDdI0hCZtE6vySjMm-WEfRq3CPzqKqqsHI&token=Ph97Xg_Yh9nYIitBrgzk86U7SQD9vUEnLVOdUr3_7Jc"
	GgChatWebhooksErrorEvent = "https://chat.googleapis.com/v1/spaces/AAAAIgLtFS8/messages?key=AIzaSyDdI0hCZtE6vySjMm-WEfRq3CPzqKqqsHI&token=YMiHvU43mxiwip93CaLp8Dy5o6yEIS-g0RsncyqgCho"
)

var (
	pathRepositoryError = "repository"
	errorRecordNotFound = "not found"
)

func WriteLogger(ctx context.Context, dataLog map[string]interface{}, name string, startTime time.Time) {
	allTime := cast.ToFloat64(timeutils.Since(startTime).Milliseconds())
	dataLog["allTime"] = allTime
	logger := NewLogger()
	if v, ok := dataLog[constants.ErrorKey]; ok {
		logger.WithKeyword(ctx, name+"_error").
			WithOutput(dataLog).
			WithResponseTime(allTime).
			WithErrorStr(fmt.Sprintf("%+v", v)).
			Error()
	} else if v, ok := dataLog[constants.WarnKey]; ok { //nolint: gocritic
		logger.WithKeyword(ctx, name+"_warn").
			WithOutput(dataLog).
			WithResponseTime(allTime).
			WithErrorStr(fmt.Sprintf("%+v", v)).
			Warn()
	} else if allTime > constants.WarnTime {
		logger.WithKeyword(ctx, name+"_warn").
			WithOutput(dataLog).
			WithResponseTime(allTime).
			Warn()
	} else {
		logger.WithKeyword(ctx, name+"_info").
			WithOutput(dataLog).
			WithResponseTime(allTime).
			Info()
	}
}

func sendGGchatCustom(ctx context.Context, messages string, webhookURL string) (err error) {

	messageData := struct {
		Text string `json:"text"`
	}{
		Text: messages,
	}

	messageBytes, _ := json.Marshal(messageData)
	//if err != nil {
	//	return fmt.Errorf("unable to marshal message to JSON: %v", err)
	//}

	req, _ := http.NewRequestWithContext(ctx, "POST", webhookURL, bytes.NewReader(messageBytes))
	//if err != nil {
	//	return fmt.Errorf("unable to create HTTP request: %v", err)
	//}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(GroupPathHeader, "/chat.googleapis.com/v1/spaces")
	client := &http.Client{}
	resp, _ := client.Do(req)
	//if err != nil {
	//	return fmt.Errorf("unable to send HTTP request: %v", err)
	//}
	defer func() {
		_ = resp.Body.Close()
	}()

	//body, err := io.ReadAll(resp.Body)
	//if err != nil {
	//	return fmt.Errorf("unable to read response body: %v", err)
	//}
	//
	//if resp.StatusCode != http.StatusOK {
	//	return fmt.Errorf("received non-OK response status: %s, response body: %s", resp.Status, string(body))
	//}

	return nil
}

func GetEnviromentHttpTimeout(env string) (weedHook string) {
	switch env {
	case constants.ProductionEnv:
		return GgChatWebhooksErrorEvent
	default:
		return GgChatWebhooksErrorDB
	}
}

func sendGGchatErrorEvent(ctx context.Context, errMessage, fromSource, trackId, topic string) {
	// send notification to google chat after begin calculation
	env := os.Getenv(constants.EnvironmentsEnv)
	path := GetEnviromentHttpTimeout(env)
	if env == "" {
		return
	}
	message := fmt.Sprintf("```Log data started at: %v \n \ntopic: %v \nenv: %v \n"+
		"from_source: %v \nerror: %v \ntrack_id: %v```",
		timeutils.NowInGMT07().Format(timeutils.CustomDateTimeVN), topic, env, fromSource, errMessage, trackId)
	go sendGGchatCustom(ctx,
		message,
		path)
}
