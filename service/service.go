package kafkasvc

import (
	"encoding/json"
	"fmt"

	sarama "github.com/Shopify/sarama"
	"golang.org/x/net/context"

	"github.com/moul/kafka-gateway/gen/pb"
)

type Service struct {
	producer sarama.SyncProducer
}

func New(producer sarama.SyncProducer) kafkapb.KafkaServiceServer {
	return &Service{
		producer: producer,
	}
}

func (s *Service) Consumer(ctx context.Context, input *kafkapb.ConsumerRequest) (*kafkapb.ConsumerResponse, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *Service) Producer(ctx context.Context, input *kafkapb.ProducerRequest) (*kafkapb.ProducerResponse, error) {
	val, err := json.Marshal(input.Value)
	if err != nil {
		return nil, err
	}
	packet := &sarama.ProducerMessage{
		Topic: input.Topic,
		Value: sarama.ByteEncoder(val),
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
