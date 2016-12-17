DOCKER_IMAGE ?=	moul/kafkagw
SOURCES := cmd/kafkagw/main.go service/service.go

.PHONY: build
build: kafkagw

kafkagw: gen/pb/kafka.pb.go gen/.generated $(SOURCES)
	go build -o kafkagw ./cmd/kafkagw

gen/pb/kafka.pb.go: pb/kafka.proto
	@mkdir -p gen/pb
	cd pb; protoc --gogo_out=plugins=grpc:../gen/pb ./kafka.proto

gen/.generated:	pb/kafka.proto
	@mkdir -p gen
	cd pb; protoc --gotemplate_out=destination_dir=../gen,template_dir=../vendor/github.com/moul/protoc-gen-gotemplate/examples/go-kit/templates/{{.File.Package}}/gen:../gen ./kafka.proto
	@touch gen/.generated

.PHONY: install
install:
	go install ./cmd/kafkagw

.PHONY: docker.build
docker.build:
	docker build -t $(DOCKER_IMAGE) .

.PHONY: docker.run
docker.run:
	docker run -p 8000:8000 -p 9000:9000 $(DOCKER_IMAGE)
