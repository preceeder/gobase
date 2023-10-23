//   File Name:  producter.go
//    Description:
//    Author:      Chenghu
//    Date:       2023/10/20 16:18
//    Change Activity:

package procder

import (
	"fmt"
	"github.com/nsqio/go-nsq"
	"github.com/preceeder/gobase/utils"
	"log/slog"
	"os"
)

var NsqProduct = map[string]*nsq.Producer{}

func InitNsqProducer() {
	//开启信号监听
	signl := utils.StartSignalLister()

	//开启信号处理
	go utils.SignalHandler(signl, func() {
		//平滑关闭
		for _, v := range NsqProduct {
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
	NsqProduct[name] = producer
	slog.Info("开启 nsq producer", "addr", producer.String())
	return nil
}
