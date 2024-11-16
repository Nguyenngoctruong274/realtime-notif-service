package mail

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"yes4all/ads-noti-apipkg/config"
	"yes4all/ads-noti-apipkg/logger"
	"yes4all/ads-noti-apipkg/xhttp"

	"github.com/spf13/viper"
)

func Test_SendEmailAlert(t *testing.T) {
	viper.SetConfigFile("../../../.env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("error while reading config file: %s", err.Error())
	}
	for _, env := range viper.AllKeys() {
		if viper.GetString(env) != "" {
			_ = os.Setenv(env, viper.GetString(env))
			_ = os.Setenv(strings.ToUpper(env), viper.GetString(env))
		}
	}
	config.InitConfig()
	logger.InitLogger()
	testInterface := NewClient(xhttp.NewClient())
	err, response := testInterface.SendEmailAlert(context.Background(), WarningEmail{
		TicketID:              43,
		TeamName:              "US Operation",
		BudgetSpendPercentage: 10,
		BudgetSpend:           10,
	})
	fmt.Println(response)
	if err != nil {
		t.Fail()
	}
}
