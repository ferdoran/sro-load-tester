package main

import (
	"github.com/ferdoran/go-sro-framework/client"
	"github.com/panjf2000/ants"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"time"
)

type Runner struct {
	pool              *ants.Pool
	concurrentUsers   int
	duration          time.Duration
	rescheduleTimeout time.Duration
	running           bool
	flows             []string
}

func (r *Runner) IsRunning() bool {
	return r.running
}

func NewRunner(flows []string) *Runner {
	concurrentUsers := viper.GetInt("config.concurrent-clients")
	duration := viper.GetDuration("config.duration")
	rescheduleTimeout := viper.GetDuration("config.reschedule-timeout")

	pool, err := ants.NewPool(concurrentUsers, ants.WithNonblocking(true))

	if err != nil {
		logrus.Panic("failed to initialize ants pool. ", err)
	}

	return &Runner{
		pool:              pool,
		concurrentUsers:   concurrentUsers,
		duration:          duration,
		rescheduleTimeout: rescheduleTimeout,
		running:           false,
		flows:             flows,
	}
}

func (r *Runner) Start() {
	logrus.Infof("starting load test with %d concurrent users over %s", r.concurrentUsers, r.duration)
	r.running = true
	defer ants.Release()
	defer r.pool.Release()
	timer := time.After(r.duration)

	for {
		select {
		case <-timer:
			r.pool.Release()
			logrus.Infof("finished load test")
			r.running = false
			return
		default:
			r.pool.Submit(func() {
				gatewayClient := client.NewClient(viper.GetString("gateway.host"), viper.GetInt("gateway.port"), "SR_Client")
				agentClient := client.NewClient(viper.GetString("agent.host"), viper.GetInt("agent.port"), "SR_Client")
				gatewayClient.AutoReconnect = false
				agentClient.AutoReconnect = false
				sharedState := make(map[string]interface{})

				for _, flow := range r.flows {
					if f, exists := AvailableFlows[flow]; exists {
						fl := f.Clone()
						logrus.Debugf("playing flow %s", f.Name())
						fl.Play(gatewayClient, agentClient, sharedState)
					}
				}

				gatewayClient.Conn.Close()
				//agentClient.Conn.Close()
				time.Sleep(r.rescheduleTimeout)
			})
		}
	}
}
