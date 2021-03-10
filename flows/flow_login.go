package flows

import (
	"github.com/ferdoran/go-sro-framework/client"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	ConfigKeyShardName = "flows.config.login.shard-name"
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
	logrus.Debugf("%s connecting to gateway server", f.logPrefix())
	gatewayClient.Connect()
	logrus.Debugf("%s why don't you come here?", f.logPrefix())
	for _, flow := range f.subFlows {
		logrus.Debugf("%s playing sub flow %s", f.logPrefix(), flow.Name())
		flow.Play(gatewayClient, agentClient, globalState)
	}

	f.isPlaying = false
}

func (f *Login) Clone() Flow {
	return &Login{
		shardName: f.shardName,
		user:      f.user,
		password:  f.password,
		subFlows: [2]Flow{
			NewShardListPingFlow(),
			NewShardListRequestFlow(),
		},
		isPlaying: false,
	}
}
