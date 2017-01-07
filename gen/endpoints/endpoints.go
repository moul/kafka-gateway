package kafka_endpoints

import (
	"fmt"

	"github.com/go-kit/kit/endpoint"
	pb "github.com/moul/kafka-gateway/gen/pb"
	context "golang.org/x/net/context"
)

var _ = endpoint.Chain
var _ = fmt.Errorf
var _ = context.Background

type StreamEndpoint func(server interface{}, req interface{}) (err error)

type Endpoints struct {
	ConsumerEndpoint endpoint.Endpoint

	ProducerEndpoint endpoint.Endpoint

	ConsumerStreamEndpoint StreamEndpoint
}

func (e *Endpoints) Consumer(ctx context.Context, in *pb.ConsumerRequest) (*pb.ConsumerResponse, error) {
	out, err := e.ConsumerEndpoint(ctx, in)
	if err != nil {
		return &pb.ConsumerResponse{ErrMsg: err.Error()}, err
	}
	return out.(*pb.ConsumerResponse), err
}

func (e *Endpoints) Producer(ctx context.Context, in *pb.ProducerRequest) (*pb.ProducerResponse, error) {
	out, err := e.ProducerEndpoint(ctx, in)
	if err != nil {
		return &pb.ProducerResponse{ErrMsg: err.Error()}, err
	}
	return out.(*pb.ProducerResponse), err
}

func (e *Endpoints) ConsumerStream(in *pb.ConsumerStreamRequest, server pb.KafkaService_ConsumerStreamServer) error {
	return fmt.Errorf("not implemented")
}

func MakeConsumerEndpoint(svc pb.KafkaServiceServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*pb.ConsumerRequest)
		rep, err := svc.Consumer(ctx, req)
		if err != nil {
			return &pb.ConsumerResponse{ErrMsg: err.Error()}, err
		}
		return rep, nil
	}
}

func MakeProducerEndpoint(svc pb.KafkaServiceServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*pb.ProducerRequest)
		rep, err := svc.Producer(ctx, req)
		if err != nil {
			return &pb.ProducerResponse{ErrMsg: err.Error()}, err
		}
		return rep, nil
	}
}

func MakeConsumerStreamEndpoint(svc pb.KafkaServiceServer) StreamEndpoint {
	return func(server interface{}, request interface{}) error {

		return svc.ConsumerStream(request.(*pb.ConsumerStreamRequest), server.(pb.KafkaService_ConsumerStreamServer))

	}
}

func MakeEndpoints(svc pb.KafkaServiceServer) Endpoints {
	return Endpoints{

		ConsumerEndpoint: MakeConsumerEndpoint(svc),

		ProducerEndpoint: MakeProducerEndpoint(svc),

		ConsumerStreamEndpoint: MakeConsumerStreamEndpoint(svc),
	}
}
