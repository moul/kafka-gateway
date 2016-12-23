package kafkasvc

import (
	"fmt"

	sarama "github.com/Shopify/sarama"
	sarama_cluster "github.com/bsm/sarama-cluster"
	"golang.org/x/net/context"

	"github.com/moul/kafka-gateway/gen/pb"
)

type Service struct {
	producer sarama.SyncProducer
	brokers  []string
}

func New(producer sarama.SyncProducer, brokers []string) kafkapb.KafkaServiceServer {
	return &Service{
		producer: producer,
		brokers:  brokers,
	}
}

func (s *Service) Consumer(ctx context.Context, input *kafkapb.ConsumerRequest) (*kafkapb.ConsumerResponse, error) {
	config := sarama_cluster.NewConfig()
	config.ClientID = input.ClientId
	config.Consumer.Return.Errors = true
	config.Group.Return.Notifications = true
	consumer, err := sarama_cluster.NewConsumer(s.brokers, "groupTest", input.Topics, config)
	if err != nil {
		return nil, err
	}

	defer consumer.Close()

	for msg := range consumer.Messages() {
		return &kafkapb.ConsumerResponse{
			Value: string(msg.Value),
		}, nil
	}

	return nil, fmt.Errorf("not implemented")
}

func (s *Service) ConsumerStream(input *kafkapb.ConsumerStreamRequest, server kafkapb.KafkaService_ConsumerStreamServer) (err error) {
	return fmt.Errorf("not implemented")
}

func (s *Service) Producer(ctx context.Context, input *kafkapb.ProducerRequest) (*kafkapb.ProducerResponse, error) {
	packet := &sarama.ProducerMessage{
		Topic: input.Topic,
		Value: sarama.ByteEncoder(input.Value),
	}
	if input.Key != "" {
		packet.Key = sarama.StringEncoder(input.Key)
	}
	partition, offset, err := s.producer.SendMessage(packet)
	if err != nil {
		return nil, err
	}
	return &kafkapb.ProducerResponse{
		Partition: partition,
		Offset:    offset,
	}, nil
}
