package flows

import "github.com/ferdoran/go-sro-framework/client"

type Flow interface {
	Name() string
	Play(gatewayClient *client.Client, agentClient *client.Client, globalState map[string]interface{})
	IsPlaying() bool
}
