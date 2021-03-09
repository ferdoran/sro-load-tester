package flows

import (
	"github.com/ferdoran/go-sro-framework/client"
	"github.com/ferdoran/go-sro-framework/network"
	"github.com/ferdoran/go-sro-framework/network/opcode"
	"github.com/sirupsen/logrus"
)

type ShardListRequest struct {
	WantedShardName string
	isPlaying       bool
}

func (s ShardListRequest) logPrefix() string {
	return "[SHARD LIST REQUEST]::"
}

func (s *ShardListRequest) IsPlaying() bool {
	return s.isPlaying
}

func (s *ShardListRequest) Play(gatewayClient, agentClient *client.Client, globalState map[string]interface{}) {
	s.isPlaying = true
	logrus.Infof("%s, sending shard list request", s.logPrefix())
	p := network.EmptyClientPacket()
	p.Encrypted = true
	p.MessageID = opcode.ShardlistRequest
	gatewayClient.OutgoingPacketChannel <- p

	for response := range gatewayClient.IncomingPacketChannel {
		if response.MessageID == opcode.ShardlistResponse {

			for {
				hasFarmEntries, err := response.ReadBool()

				if err != nil {
					logrus.Errorf("%s failed to read if farm entries exist", s.logPrefix())
					return
				}

				if !hasFarmEntries {
					break
				}

				farmID, err := response.ReadByte()
				if err != nil {
					logrus.Errorf("%s failed to read farm id", s.logPrefix())
					return
				}

				farmName, err := response.ReadString()
				if err != nil {
					logrus.Errorf("%s failed to read farm name", s.logPrefix())
				}

				logrus.Infof("%s received farm with id = %d, name = %s", s.logPrefix(), farmID, farmName)
			}

			for {
				hasShardEntries, err := response.ReadBool()

				if err != nil {
					logrus.Errorf("%s failed to read if shard entries exist", s.logPrefix())
					return
				}

				if !hasShardEntries {
					break
				}

				shardID, err := response.ReadUInt16()
				if err != nil {
					logrus.Errorf("%s failed to read shard id", s.logPrefix())
					return
				}

				shardName, err := response.ReadString()
				if err != nil {
					logrus.Errorf("%s failed to read shard name", s.logPrefix())
				}

				shardOnlineCount, err := response.ReadUInt16()
				if err != nil {
					logrus.Errorf("%s failed to read shard online count", s.logPrefix())
				}

				capacity, err := response.ReadUInt16()
				if err != nil {
					logrus.Errorf("%s failed to read shard capacity", s.logPrefix())
				}

				shardIsOperating, err := response.ReadBool()
				if err != nil {
					logrus.Errorf("%s failed to read if is shard operating", s.logPrefix())
				}

				farmId, err := response.ReadByte()
				if err != nil {
					logrus.Errorf("%s failed to read shard's farm id", s.logPrefix())
				}

				logrus.Infof("%s received shard with id = %d, name = %s, online count = %d, capacity = %d, is operating = %t, farm id = %d", s.logPrefix(), shardID, shardName, shardOnlineCount, capacity, shardIsOperating, farmId)

				if s.WantedShardName == shardName {
					logrus.Infof("%s found wanted shard: %s", s.logPrefix(), s.WantedShardName)
					globalState["shard.id"] = shardID
					globalState["shard.name"] = shardName
				}
			}

			break
		}
	}

	s.isPlaying = false
}

func (s ShardListRequest) Name() string {
	return "shard-list-request"
}
