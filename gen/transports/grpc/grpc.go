package kafka_grpctransport

import (
	"fmt"

	grpctransport "github.com/go-kit/kit/transport/grpc"
	context "golang.org/x/net/context"

	endpoints "github.com/moul/kafka-gateway/gen/endpoints"
	pb "github.com/moul/kafka-gateway/gen/pb"
)

// avoid import errors
var _ = fmt.Errorf

func MakeGRPCServer(ctx context.Context, endpoints endpoints.Endpoints) pb.KafkaServiceServer {
	options := []grpctransport.ServerOption{}
	return &grpcServer{

		consumer: grpctransport.NewServer(
			ctx,
			endpoints.ConsumerEndpoint,
			decodeRequest,
			encodeConsumerResponse,
			options...,
		),

		producer: grpctransport.NewServer(
			ctx,
			endpoints.ProducerEndpoint,
			decodeRequest,
			encodeProducerResponse,
			options...,
		),

		consumerstream: &server{
			e: endpoints.ConsumerStreamEndpoint,
		},
	}
}

type grpcServer struct {
	consumer grpctransport.Handler

	producer grpctransport.Handler

	consumerstream streamHandler
}

func (s *grpcServer) Consumer(ctx context.Context, req *pb.ConsumerRequest) (*pb.ConsumerResponse, error) {
	_, rep, err := s.consumer.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.ConsumerResponse), nil
}

func encodeConsumerResponse(ctx context.Context, response interface{}) (interface{}, error) {
	resp := response.(*pb.ConsumerResponse)
	return resp, nil
}

func (s *grpcServer) Producer(ctx context.Context, req *pb.ProducerRequest) (*pb.ProducerResponse, error) {
	_, rep, err := s.producer.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.ProducerResponse), nil
}

func encodeProducerResponse(ctx context.Context, response interface{}) (interface{}, error) {
	resp := response.(*pb.ProducerResponse)
	return resp, nil
}

func (s *grpcServer) ConsumerStream(req *pb.ConsumerStreamRequest, server pb.KafkaService_ConsumerStreamServer) error {
	return s.consumerstream.Do(server, req)
}

func decodeRequest(ctx context.Context, grpcReq interface{}) (interface{}, error) {
	return grpcReq, nil
}

type streamHandler interface {
	Do(server interface{}, req interface{}) (err error)
}

type server struct {
	e endpoints.StreamEndpoint
}

func (s server) Do(server interface{}, req interface{}) (err error) {
	if err := s.e(server, req); err != nil {
		return err
	}
	return nil
}
