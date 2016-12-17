package kafkasvc

import (
	"fmt"

	"golang.org/x/net/context"

	"github.com/moul/kafka-gateway/gen/pb"
)

type Service struct{}

func New() kafkapb.KafkaServiceServer {
	return &Service{}
}

func (s *Service) Consumer(ctx context.Context, input *kafkapb.ConsumerRequest) (*kafkapb.ConsumerResponse, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *Service) Producer(ctx context.Context, input *kafkapb.ProducerRequest) (*kafkapb.ProducerResponse, error) {
	return nil, fmt.Errorf("not implemented")
}
