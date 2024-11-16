package xhttp

import (
	"testing"
)

func TestPostJSONWithCustomHeader_TimeOut(t *testing.T) {
	//Set up environment
	//
	//err := os.Setenv("AWS_URL", "https://adv-gateway.yes4all.com/amz_y4a")
	//err = os.Setenv("ENVIRONMENT", "test")
	//viper.SetConfigFile("../../../.env")
	//viper.AutomaticEnv()
	//
	//if err := viper.ReadInConfig(); err != nil {
	//	fmt.Printf("error while reading config file: %s", err.Error())
	//}
	//for _, env := range viper.AllKeys() {
	//	if viper.GetString(env) != "" {
	//		_ = os.Setenv(env, viper.GetString(env))
	//		_ = os.Setenv(strings.ToUpper(env), viper.GetString(env))
	//	}
	//}
	//config.InitConfig()
	//logger.InitLogger()
	//if err != nil {
	//	return
	//}
	//c := NewClient()
	//
	//ctx := tracking.InitContextWithTrackID()
	//ctx = xcontext.SetFromSource(ctx, "backend")
	////Set up data test
	//dataTests := []struct {
	//	name           string
	//	url            string
	//	data           interface{}
	//	target         interface{}
	//	customHeader   http.Header
	//	method         string
	//	reqOpts        RequestOption
	//	expectedStatus int
	//	expectedError  error
	//	sleepDuration  time.Duration
	//}{
	//	{
	//		name:           "Request exceeds timeout",
	//		url:            os.Getenv(constants.AWS_URL),
	//		data:           nil,
	//		target:         nil,
	//		customHeader:   http.Header{"Content-Type": []string{"application/json"}},
	//		method:         "POST",
	//		reqOpts:        RequestOption{GroupPath: "amz_y4a/sp/campaigns/list"},
	//		expectedStatus: http.StatusRequestTimeout,
	//		expectedError:  context.DeadlineExceeded,
	//		sleepDuration:  1 * time.Second,
	//	},
	//}
	//
	//for _, test := range dataTests {
	//	t.Run(test.name, func(t *testing.T) {
	//		//// Set Timeout into context
	//
	//		req, err := NewRequestBuilderWithCtx(ctx).
	//			WithMethod(test.method).
	//			WithURL(test.url).
	//			WithBody(MIMEJSON, test.data).
	//			Build()
	//		if err != nil {
	//			return
	//		}
	//		decodeNumber := false
	//		// call PostJSONWithCustomHeader
	//		statusCode, err := c.Do(ctx, req, test.target, decodeNumber)
	//		// Assert the results
	//		log.Error(statusCode, err)
	//	})
	//}
}
