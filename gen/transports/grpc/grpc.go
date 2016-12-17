package kafka_grpctransport



import (
        "fmt"

	context "golang.org/x/net/context"
        pb "github.com/moul/kafka-gateway/gen/pb"
        endpoints "github.com/moul/kafka-gateway/gen/endpoints"
	grpctransport "github.com/go-kit/kit/transport/grpc"
)

// avoid import errors
var _ = fmt.Errorf

func MakeGRPCServer(ctx context.Context, endpoints endpoints.Endpoints) pb.KafkaServiceServer {
	options := []grpctransport.ServerOption{}
	return &grpcServer{
		
                
                
		consumer: grpctransport.NewServer(
			ctx,
			endpoints.ConsumerEndpoint,
			decodeConsumerRequest,
			encodeConsumerResponse,
			options...,
		),
                
		
                
                
                
		producer: grpctransport.NewServer(
			ctx,
			endpoints.ProducerEndpoint,
			decodeProducerRequest,
			encodeProducerResponse,
			options...,
		),
                
		
                
	}
}

type grpcServer struct {
	
	consumer grpctransport.Handler
	
	producer grpctransport.Handler
	
}


func (s *grpcServer) Consumer(ctx context.Context, req *pb.ConsumerRequest) (*pb.ConsumerResponse, error) {
	_, rep, err := s.consumer.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.ConsumerResponse), nil
}

func decodeConsumerRequest(ctx context.Context, grpcReq interface{}) (interface{}, error) {
	return grpcReq, nil
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

func decodeProducerRequest(ctx context.Context, grpcReq interface{}) (interface{}, error) {
	return grpcReq, nil
}

func encodeProducerResponse(ctx context.Context, response interface{}) (interface{}, error) {
	resp := response.(*pb.ProducerResponse)
	return resp, nil
}

