//   File Name:  consumer.go
//    Description:
//    Author:      Chenghu
//    Date:       2023/10/23 13:37
//    Change Activity:

package nsq_f

import (
	"github.com/nsqio/go-nsq"
	"golang.org/x/net/context"
	"time"
)

type SubRouter struct {
	Topic   string
	Channel string
	Handler MqHandlerFunc
}

// 保存路由
var Routers []SubRouter = make([]SubRouter, 0)

// 添加路由到 队列
func RegisterRouter(topic, channel string, handler MqHandlerFunc) {
	Routers = append(Routers, SubRouter{topic, channel, handler})
}

// 定义 新的 NsqConsumer 对象
type NsqConsumer struct {
	// nsq lookupds 地址
	lookupds []string
	// nsq consumer config
	lookupdPollInterval time.Duration
	// nsq consumer config
	maxInFlight int
	// 注册的 consumers 对象
	consumers []*nsq.Consumer
}

// 消息 handler 函数定义
type MqHandlerFunc func(ctx context.Context, msg *nsq.Message) error

// 初始化 Nsq consumer 对象
func NewNsqConsumer(lookupds []string, pollInterval time.Duration, maxInFlight int) *NsqConsumer {
	return &NsqConsumer{
		lookupds:            lookupds,
		lookupdPollInterval: pollInterval,
		maxInFlight:         maxInFlight,
	}
}

// 将 MqHandlerFunc 转换成 nsq.Consumer 内接受的 nsq.HandlerFunc
func (n *NsqConsumer) toNsqHandler(handlerFunc MqHandlerFunc) nsq.HandlerFunc {
	return func(msg *nsq.Message) error {
		ctx := context.TODO()
		return handlerFunc(ctx, msg)
	}
}

// 注册 topic 和 handler Func
func (n *NsqConsumer) registerHandler(topic, channel string, handler MqHandlerFunc) error {
	cfg := nsq.NewConfig()
	cfg.LookupdPollInterval = n.lookupdPollInterval
	c, err := nsq.NewConsumer(topic, channel, cfg)
	if err != nil {
		return err
	}
	c.ChangeMaxInFlight(n.maxInFlight)
	c.AddHandler(n.toNsqHandler(handler))
	n.consumers = append(n.consumers, c)
	return nil
}

func (n *NsqConsumer) Start() error {
	for _, router := range Routers {
		if err := n.registerHandler(router.Topic, router.Channel, router.Handler); err != nil {
			return err
		}
	}

	for _, h := range n.consumers {
		if err := h.ConnectToNSQLookupds(n.lookupds); err != nil {
			return err
		}
	}
	return nil
}

// close consumer
func (n *NsqConsumer) Close() {
	for _, h := range n.consumers {
		h.Stop()
	}
}
