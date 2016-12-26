package kafka_clientgrpc

import (
	jwt "github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	context "golang.org/x/net/context"
	"google.golang.org/grpc"

	endpoints "github.com/moul/kafka-gateway/gen/endpoints"
	pb "github.com/moul/kafka-gateway/gen/pb"
)

func New(conn *grpc.ClientConn, logger log.Logger) pb.KafkaServiceServer {

	var consumerEndpoint endpoint.Endpoint
	{
		consumerEndpoint = grpctransport.NewClient(
			conn,
			"kafka.KafkaService",
			"Consumer",
			EncodeConsumerRequest,
			DecodeConsumerResponse,
			pb.ConsumerResponse{},
			append([]grpctransport.ClientOption{}, grpctransport.ClientBefore(jwt.FromGRPCContext()))...,
		).Endpoint()
	}

	var producerEndpoint endpoint.Endpoint
	{
		producerEndpoint = grpctransport.NewClient(
			conn,
			"kafka.KafkaService",
			"Producer",
			EncodeProducerRequest,
			DecodeProducerResponse,
			pb.ProducerResponse{},
			append([]grpctransport.ClientOption{}, grpctransport.ClientBefore(jwt.FromGRPCContext()))...,
		).Endpoint()
	}

	return &endpoints.Endpoints{

		ConsumerEndpoint: consumerEndpoint,

		ProducerEndpoint: producerEndpoint,
	}
}

func EncodeConsumerRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.ConsumerRequest)
	return req, nil
}

func DecodeConsumerResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(*pb.ConsumerResponse)
	return response, nil
}

func EncodeProducerRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.ProducerRequest)
	return req, nil
}

func DecodeProducerResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(*pb.ProducerResponse)
	return response, nil
}
