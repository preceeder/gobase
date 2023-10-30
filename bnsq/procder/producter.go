//   File Name:  producter.go
//    Description:
//    Author:      Chenghu
//    Date:       2023/10/20 16:18
//    Change Activity:

package producer

import (
	"fmt"
	"github.com/nsqio/go-nsq"
	"github.com/preceeder/gobase/utils"
	"github.com/spf13/viper"
	"log/slog"
	"os"
)

var NsqProducer = map[string]*nsq.Producer{}

type NsqProducerConfig struct {
	Name string `json:"name" default:"default"`
	Addr string `json:"nsqdAddr" default:"127.0.0.1:8009"`
}

func InitNsqProducer(config viper.Viper) {
	nsqConfig := []NsqProducerConfig{}
	utils.ReadViperConfig(config, "nsq-producer", &nsqConfig)
	for _, nsqc := range nsqConfig {
		err := NewProduct(nsqc.Addr, nsqc.Name)
		if err != nil {
			slog.Error("InitNsqProducer error ", "error", err.Error())
			panic("InitNsqProducer error: " + err.Error())
		}
	}
	//开启信号监听
	signl := utils.StartSignalLister()

	//开启信号处理
	go utils.SignalHandler(signl, func() {
		//平滑关闭
		for _, v := range NsqProducer {
			slog.Info("stop nsq producer", "addr", v.String())
			v.Stop()
		}
		os.Exit(1)
	})
}

func NewProduct(nsqdAddr, name string) error {
	config := nsq.NewConfig()
	producer, err := nsq.NewProducer(nsqdAddr, config)
	if err != nil {
		fmt.Printf("create producer failed, err:%v\n", err.Error())
		return err
	}
	err = producer.Ping()
	if err != nil {
		fmt.Printf("nsq producer ping error: %v\n", err.Error())
	}
	NsqProducer[name] = producer
	slog.Info("开启 nsq producer", "addr", producer.String())
	return nil
}
