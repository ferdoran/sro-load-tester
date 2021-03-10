package main

import (
	"github.com/ferdoran/sro-load-tester/flows"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var AvailableFlows = make(map[string]flows.Flow)

func main() {
	// 1. Load config
	initConfig()
	//logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:               true,
		EnvironmentOverrideColors: true,
		FullTimestamp:             true,
		TimestampFormat:           "2006-01-01 15:04:05.000",
	})

	// 2. Initialize flows
	AvailableFlows["login"] = flows.NewLoginFlow()

	// 3. Play flows
	flows := viper.GetStringSlice("flows.active")

	logrus.Infof("found flows: %v", flows)

	runner := NewRunner(flows)
	runner.Start()
}

func initConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")

	viper.SetDefault("gateway.host", "127.0.0.1")
	viper.SetDefault("gateway.port", 15779)
	viper.SetDefault("agent.host", "127.0.0.1")
	viper.SetDefault("agent.port", 15882)
	viper.SetDefault("config.duration", "60s")
	viper.SetDefault("config.reschedule-timeout", "50ms")
	viper.SetDefault("config.concurrent-clients", 10)

	err := viper.ReadInConfig()

	if err != nil {
		logrus.Panicf("failed to read in config file: %s", err)
	}
}
