package cmd

import (
	"os"
	"strings"
	"yes4all/ads-noti-api/pkg/config"

	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	rootCmd = &cobra.Command{
		Use:     "Example",
		Short:   "Example Short",
		Long:    `Example Long`,
		Version: "1.0.0",
	}
	customViper = viper.New()
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.AddCommand(apiCmd)

	defaultFlags()
}

func initConfig() {
	viper.SetConfigFile(".env")
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
}

// nolint
func defaultFlags() {
	pflags := rootCmd.PersistentFlags()

	// app port
	pflags.Int("port", 8081, "app binding port")
	_ = customViper.BindPFlag("app.port", pflags.Lookup("port"))

	// TODO: Add more default config here
}
