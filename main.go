package main

import (
	"github.com/ferdoran/go-sro-framework/client"
	"github.com/ferdoran/sro-load-tester/flows"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	// 1. Load config
	initConfig()
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:               true,
		EnvironmentOverrideColors: true,
		FullTimestamp:             true,
		TimestampFormat:           "2006-01-01 15:04:05.000",
	})

	gatewayClient := client.NewClient(viper.GetString("gateway.host"), viper.GetInt("gateway.port"), "SR_Client")
	agentClient := client.NewClient(viper.GetString("agent.host"), viper.GetInt("agent.port"), "SR_Client")
	sharedState := make(map[string]interface{})

	// 2. Initialize flows
	availableFlows := make(map[string]flows.Flow)
	availableFlows["login"] = flows.NewLoginFlow()
	// 3. Play flows
	flows := viper.GetStringSlice("playFlows")

	logrus.Infof("found flows: %v", flows)

	for _, flow := range flows {
		if f, exists := availableFlows[flow]; exists {
			logrus.Infof("playing flow %s", f.Name())
			f.Play(gatewayClient, agentClient, sharedState)
		}
	}
}

func initConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")

	viper.SetDefault("gateway.host", "127.0.0.1")
	viper.SetDefault("gateway.port", 15779)
	viper.SetDefault("agent.host", "127.0.0.1")
	viper.SetDefault("agent.port", 15882)

	err := viper.ReadInConfig()

	if err != nil {
		logrus.Panicf("failed to read in config file: %s", err)
	}
}
