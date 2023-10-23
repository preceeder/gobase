//   File Name:  run.go
//    Description:
//    Author:      Chenghu
//    Date:       2023/10/23 13:38
//    Change Activity:

package consumer

import (
	"github.com/preceeder/gobase/utils"
	"github.com/spf13/viper"
	"time"
)

var nsqConfig *NsqConfig

type NsqConfig struct {
	NsqLookupd   []string `json:"lookupds"`
	PollInterval int64    `json:"pollInterval"`
	MaxInFlight  int      `json:"maxInFlight"`
}

func InitNsqConsumerConfig(config viper.Viper) {
	nsqConfig = &NsqConfig{}
	utils.ReadViperConfig(config, "nsq-consumer", nsqConfig)
}

// 启动服务
type Server struct {
	cfg *NsqConfig
	nsq *NsqConsumer
}

func NewServer(cfg ...*NsqConfig) *Server {
	if len(cfg) > 0 {
		nsqConfig = cfg[0]
	}
	return &Server{
		cfg: nsqConfig,
	}
}

func (s *Server) Run() {
	// 启动 nsq consumber    10, 2
	s.nsq = NewNsqConsumer(s.cfg.NsqLookupd, time.Duration(s.cfg.PollInterval)*time.Second, s.cfg.MaxInFlight)

	if err := s.nsq.Start(); err != nil {
		panic(err)
	}
}

// 关闭服务
func (s *Server) Close() {
	if s.nsq != nil {
		s.nsq.Close()
	}
}

func Start(cfg ...*NsqConfig) *Server {
	// 启动服务
	if len(cfg) > 0 {
		nsqConfig = cfg[0]
	}
	srv := NewServer(nsqConfig)
	srv.Run()
	return srv
}
