package kafka_httptransport



import (
       "log"
	"net/http"
	"encoding/json"
	context "golang.org/x/net/context"

        pb "github.com/moul/kafka-gateway/gen/pb"
        gokit_endpoint "github.com/go-kit/kit/endpoint"
        httptransport "github.com/go-kit/kit/transport/http"
        endpoints "github.com/moul/kafka-gateway/gen/endpoints"
)


func MakeConsumerHandler(ctx context.Context, svc pb.KafkaServiceServer, endpoint gokit_endpoint.Endpoint) *httptransport.Server {
	return httptransport.NewServer(
		ctx,
		endpoint,
		decodeConsumerRequest,
		encodeConsumerResponse,
                []httptransport.ServerOption{}...,
	)
}

func decodeConsumerRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req pb.ConsumerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return &req, nil
}

func encodeConsumerResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func MakeProducerHandler(ctx context.Context, svc pb.KafkaServiceServer, endpoint gokit_endpoint.Endpoint) *httptransport.Server {
	return httptransport.NewServer(
		ctx,
		endpoint,
		decodeProducerRequest,
		encodeProducerResponse,
                []httptransport.ServerOption{}...,
	)
}

func decodeProducerRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req pb.ProducerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return &req, nil
}

func encodeProducerResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}


func RegisterHandlers(ctx context.Context, svc pb.KafkaServiceServer, mux *http.ServeMux, endpoints endpoints.Endpoints) error {
	
        log.Println("new HTTP endpoint: \"/Consumer\" (service=Kafka)")
	mux.Handle("/Consumer", MakeConsumerHandler(ctx, svc, endpoints.ConsumerEndpoint))
	
        log.Println("new HTTP endpoint: \"/Producer\" (service=Kafka)")
	mux.Handle("/Producer", MakeProducerHandler(ctx, svc, endpoints.ProducerEndpoint))
	
	return nil
}
