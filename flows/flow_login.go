package flows

import (
	"github.com/ferdoran/go-sro-framework/client"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	ConfigKeyShardName = "flows.config.login.shard_name"
	ConfigKeyUser      = "flows.config.login.user"
	ConfigKeyPassword  = "flows.config.login.password"
)

type Login struct {
	shardName string
	user      string
	password  string
	subFlows  [2]Flow
	isPlaying bool
}

func NewLoginFlow() *Login {
	shardName := viper.GetString(ConfigKeyShardName)
	user := viper.GetString(ConfigKeyUser)
	password := viper.GetString(ConfigKeyPassword)

	if shardName == "" || user == "" || password == "" {
		logrus.Panicf("at least one property is empty: %s, %s, %s", ConfigKeyShardName, ConfigKeyUser, ConfigKeyShardName)
	}

	// Init Steps
	steps := [2]Flow{
		NewShardListPingFlow(),
		NewShardListRequestFlow(),
	}

	return &Login{
		shardName: shardName,
		user:      user,
		password:  password,
		subFlows:  steps,
		isPlaying: false,
	}

}

func (f Login) Name() string {
	return "login"
}

func (f Login) logPrefix() string {
	return "[LOGIN FLOW]::"
}

func (f *Login) IsPlaying() bool {
	return f.isPlaying
}

func (f *Login) Play(gatewayClient, agentClient *client.Client, globalState map[string]interface{}) {
	f.isPlaying = true
	logrus.Infof("%s connecting to gateway server", f.logPrefix())
	gatewayClient.Connect()

	for _, flow := range f.subFlows {
		logrus.Infof("%s playing sub flow %s", f.logPrefix(), flow.Name())
		flow.Play(gatewayClient, agentClient, globalState)
	}

	f.isPlaying = false
}
