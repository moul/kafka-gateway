FROM    golang:1.7.3
COPY    . /go/src/github.com/moul/kafka
WORKDIR /go/src/github.com/moul/kafka
CMD     ["kafkagw"]
EXPOSE  8000 9000
RUN     make install
