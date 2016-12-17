# kafka-gateway
:microphone: Kafka Gateway (gRPC+http) micro-service

A kafka gateway / proxy, developed as a micro-service using Sarama, Protobuf and Go-Kit, with more that 75% boilerplate code generated automatically using protoc-gen-gotemplate.

## Usage

```console
$> curl -X POST localhost:8000/Producer -d '{"topic":"test","value":{"test":42}}'
{"offset":14}
$> curl -X POST localhost:8000/Producer -d '{"topic":"test","value":{"test":42}}'
{"offset":15}
```
