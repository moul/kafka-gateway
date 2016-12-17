FROM    golang:1.7.3
COPY    . /go/src/github.com/moul/kafka-gateway
WORKDIR /go/src/github.com/moul/kafka-gateway
CMD     ["kafkagw"]
EXPOSE  8000 9000
RUN     make install
