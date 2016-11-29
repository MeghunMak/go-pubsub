package natss

import (
	"github.com/nats-io/go-nats-streaming"
	"github.com/nats-io/go-nats-streaming/pb"
	"github.com/utilitywarehouse/go-pubsub"
)

var _ pubsub.MessageSource = (*messageSource)(nil)

type messageSource struct {
	topic string
	conn  stan.Conn
}

func NewMessageSource(clusterID, topic, consumerID, natsURL string) (pubsub.MessageSource, error) {

	conn, err := stan.Connect(clusterID, consumerID, stan.NatsURL(natsURL))
	if err != nil {
		return nil, err
	}

	return &messageSource{
		topic: topic,
		conn:  conn,
	}, nil
}

func (mq *messageSource) ConsumeMessages(handler pubsub.ConsumerMessageHandler, onError pubsub.ConsumerErrorHandler) error {

	f := func(msg *stan.Msg) {
		m := pubsub.ConsumerMessage{msg.Data}
		err := handler(m)
		if err != nil {
			if err := onError(m, err); err != nil {
				panic("write the error handling")
			}
		}
	}

	startOpt := stan.StartAt(pb.StartPosition_NewOnly)

	groupName := "groupName"

	_, err := mq.conn.QueueSubscribe("demo-topic", groupName, f, startOpt, stan.DurableName(groupName))
	if err != nil {
		return err
	}

	return nil
}

func (mq *messageSource) Close() error {
	return mq.conn.Close()
}