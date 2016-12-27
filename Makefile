DOCKER_IMAGE ?=	moul/kafkagw

.PHONY: build
build: kafkagw kafkagw-client kafkagw-slackbot

kafkagw: gen/pb/kafka.pb.go cmd/kafkagw/main.go service/service.go
	go build -o $@ ./cmd/$@

kafkagw-client: gen/pb/kafka.pb.go cmd/kafkagw-client/main.go
	go build -o $@ ./cmd/$@

kafkagw-slackbot: gen/pb/kafka.pb.go cmd/kafkagw-slackbot/main.go
	go build -o $@ ./cmd/$@

gen/pb/kafka.pb.go:	pb/kafka.proto
	@mkdir -p gen/pb
	cd pb; protoc --gotemplate_out=destination_dir=../gen,template_dir=../vendor/github.com/moul/protoc-gen-gotemplate/examples/go-kit/templates/{{.File.Package}}/gen:../gen ./kafka.proto
	gofmt -w gen
	cd pb; protoc --gogo_out=plugins=grpc:../gen/pb ./kafka.proto

.PHONY: stats
stats:
	wc -l service/service.go cmd/*/*.go pb/*.proto
	wc -l $(shell find gen -name "*.go")


.PHONY: test
test:
	go test -v $(shell go list ./... | grep -v /vendor/)

.PHONY: install
install:
	go install ./cmd/kafkagw
	go install ./cmd/kafkagw-client
	go install ./cmd/kafkagw-slackbot

.PHONY: docker.build
docker.build:
	docker build -t $(DOCKER_IMAGE) .

.PHONY: docker.run
docker.run:
	docker run -p 8000:8000 -p 9000:9000 $(DOCKER_IMAGE)

.PHONY: docker.test
docker.test: docker.build
	docker run $(DOCKER_IMAGE) make test

.PHONY: clean
clean:
	rm -rf gen

.PHONY: re
re: clean build
