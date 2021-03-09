package flows

import (
	"github.com/ferdoran/go-sro-framework/client"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	ConfigKeyShardName = "flows.login.shard_name"
	ConfigKeyUser      = "flows.login.user"
	ConfigKeyPassword  = "flows.login.password"
)

type LoginFlow struct {
	CurrentStep int
	shardName   string
	user        string
	password    string
	subFlows    [2]Flow
	isPlaying   bool
}

func NewLoginFlow() *LoginFlow {
	shardName := viper.GetString(ConfigKeyShardName)
	user := viper.GetString(ConfigKeyUser)
	password := viper.GetString(ConfigKeyPassword)

	if shardName == "" || user == "" || password == "" {
		logrus.Panicf("at least one property is empty: %s, %s, %s", ConfigKeyShardName, ConfigKeyUser, ConfigKeyShardName)
	}

	// Init Steps
	steps := [2]Flow{
		&ShardListPing{},
		&ShardListRequest{WantedShardName: shardName},
	}

	return &LoginFlow{
		CurrentStep: 0,
		shardName:   shardName,
		user:        user,
		password:    password,
		subFlows:    steps,
		isPlaying:   false,
	}

}

func (f LoginFlow) Name() string {
	return "login"
}

func (f LoginFlow) logPrefix() string {
	return "[LOGIN FLOW]::"
}

func (f *LoginFlow) IsPlaying() bool {
	return f.isPlaying
}

func (f *LoginFlow) Play(gatewayClient, agentClient *client.Client, globalState map[string]interface{}) {
	f.isPlaying = true
	logrus.Infof("%s connecting to gateway server", f.logPrefix())
	gatewayClient.Connect()

	for _, flow := range f.subFlows {
		logrus.Infof("%s playing sub flow %s", f.logPrefix(), flow.Name())
		flow.Play(gatewayClient, agentClient, globalState)
	}

	f.isPlaying = false
}
