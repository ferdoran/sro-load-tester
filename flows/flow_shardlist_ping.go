package flows

import (
	"github.com/ferdoran/go-sro-framework/client"
	"github.com/ferdoran/go-sro-framework/network"
	"github.com/ferdoran/go-sro-framework/network/opcode"
	"github.com/sirupsen/logrus"
)

type ShardListPing struct {
	isPlaying bool
}

func NewShardListPingFlow() *ShardListPing {
	return &ShardListPing{isPlaying: false}
}

func (s ShardListPing) logPrefix() string {
	return "[SHARD LIST PING]::"
}

func (s ShardListPing) IsPlaying() bool {
	return s.isPlaying
}

func (s *ShardListPing) Play(gatewayClient, agentClient *client.Client, globalState map[string]interface{}) {
	s.isPlaying = true
	logrus.Debugf("%s sending shard list ping", s.logPrefix())
	p := network.EmptyClientPacket()
	p.Encrypted = true
	p.MessageID = opcode.ShardlistPing
	gatewayClient.OutgoingPacketChannel <- p

	for response := range gatewayClient.IncomingPacketChannel {
		if response.MessageID == opcode.ShardlistPong {
			wasSuccessful, err := response.ReadBool()

			if err != nil {
				logrus.Errorf("%s failed to read response status", s.logPrefix())
			}

			if wasSuccessful {
				farmId, err := response.ReadByte()
				if err != nil {
					logrus.Errorf("%s failed to read farm id", s.logPrefix())
				}

				farmIp, err := response.ReadBytes(4)
				if err != nil {
					logrus.Errorf("%s failed to read farm ip", s.logPrefix())
				}

				logrus.Debugf("%s, received farm id = %d, ip = %s", s.logPrefix(), farmId, farmIp)
			} else {
				errCode, err := response.ReadByte()
				if err != nil {
					logrus.Errorf("%s failed to read error code", s.logPrefix())
				}

				logrus.Errorf("%s received error code %d", s.logPrefix(), errCode)
			}

			break
		}
	}
	s.isPlaying = false
}

func (s ShardListPing) Name() string {
	return "shard-list-ping"
}

func (s *ShardListPing) Clone() Flow {
	return &ShardListPing{isPlaying: false}
}
