package kafka

import (
	"net"
	"strconv"

	"github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"
)

const (
	Topic = "banktransfer"
)

func EnsureTransactionTopic() {
	if err := ensureTopic(Topic, 10); err != nil {
		panic(err.Error())
	}
}

func ensureTopic(topic string, numPartitions int) error {
	conn, err := kafka.Dial("tcp", connect)
	if err != nil {
		return err
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return err
	}
	var leaderConn *kafka.Conn
	leaderConn, err = kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		return err
	}
	defer leaderConn.Close()

	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             topic,
			NumPartitions:     numPartitions,
			ReplicationFactor: 1,
		},
	}

	err = leaderConn.CreateTopics(topicConfigs...)
	if err != nil {
		return err
	}

	log.Infof(
		"Topic '%s' with %d partitions successfully created\n",
		topic,
		numPartitions,
	)
	return nil
}
