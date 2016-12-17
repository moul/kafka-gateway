# kafka-gateway
:microphone: Kafka Gateway (gRPC+http) micro-service

A [kafka](https://kafka.apache.org/) gateway / proxy, developed as a micro-service using [Sarama](https://github.com/Shopify/sarama), [Protobuf](https://github.com/google/protobuf) and [Go-Kit](https://github.com/go-kit/kit), with more that 75% boilerplate code generated automatically using [protoc-gen-gotemplate](https://github.com/moul/protoc-gen-gotemplate).

## Usage

```console
$> curl -X POST localhost:8000/Producer -d '{"topic":"test","value":{"test":42}}'
{"offset":14}
$> curl -X POST localhost:8000/Producer -d '{"topic":"test","value":{"test":42}}'
{"offset":15}
```

## Code generation

```console
# custom code
$> wc -l service/service.go pb/*.proto cmd/*/main.go
      47 service/service.go
      31 pb/kafka.proto
      85 cmd/kafkagw/main.go
     163 total

# generated code
$> wc -l $(find gen -name "*.go")
      73 gen/endpoints/endpoints.go
     263 gen/pb/kafka.pb.go
      90 gen/transports/grpc/grpc.go
      72 gen/transports/http/http.go
     498 total
```
