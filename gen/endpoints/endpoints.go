package kafka_endpoints



import (
	"fmt"

	context "golang.org/x/net/context"
        pb "github.com/moul/kafka-gateway/gen/pb"
	"github.com/go-kit/kit/endpoint"
)

var _ = fmt.Errorf

type Endpoints struct {
	
	ConsumerEndpoint endpoint.Endpoint
	
	ProducerEndpoint endpoint.Endpoint
	
}


func (e *Endpoints)Consumer(ctx context.Context, in *pb.ConsumerRequest) (*pb.ConsumerResponse, error) {
	out, err := e.ConsumerEndpoint(ctx, in)
	if err != nil {
		return &pb.ConsumerResponse{ErrMsg: err.Error()}, err
	}
	return out.(*pb.ConsumerResponse), err
}

func (e *Endpoints)Producer(ctx context.Context, in *pb.ProducerRequest) (*pb.ProducerResponse, error) {
	out, err := e.ProducerEndpoint(ctx, in)
	if err != nil {
		return &pb.ProducerResponse{ErrMsg: err.Error()}, err
	}
	return out.(*pb.ProducerResponse), err
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


func MakeEndpoints(svc pb.KafkaServiceServer) Endpoints {
	return Endpoints{
		
		ConsumerEndpoint: MakeConsumerEndpoint(svc),
		
		ProducerEndpoint: MakeProducerEndpoint(svc),
		
	}
}
