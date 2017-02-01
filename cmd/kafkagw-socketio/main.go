package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/googollee/go-socket.io"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/moul/kafka-gateway/gen/pb"
)

var (
	client kafkapb.KafkaServiceClient
	conn   *grpc.ClientConn
	ctx    context.Context
	server *socketio.Server
)

func main() {
	var err error
	conn, err = grpc.Dial(
		"127.0.0.1:9000",
		grpc.WithInsecure(),
		grpc.WithTimeout(time.Second),
	)
	if err != nil {
		panic(err)
	}
	ctx = context.Background()
	client = kafkapb.NewKafkaServiceClient(conn)

	server, err = socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}

	server.On("connection", func(so socketio.Socket) {
		log.Println("on connection")
		so.Join("test")
		so.On("produce", func(msg string) string {
			log.Println("produce", msg)
			ret, err := client.Producer(ctx, &kafkapb.ProducerRequest{
				Topic: "test",
				Key:   "test",
				Value: msg,
			})
			if err != nil {
				return fmt.Sprintf("Error: %v", err)
			}
			return fmt.Sprintf("%v", ret)
		})
		so.On("disconnection", func() {
			log.Println("on disconnect")
		})
	})

	server.On("error", func(so socketio.Socket, err error) {
		log.Println("error:", err)
	})

	go kafkaConsumer()

	http.Handle("/socket.io/", server)
	// http.Handle("/", http.FileServer(http.Dir("./asset")))
	log.Println("Serving at localhost:5000...")
	log.Fatal(http.ListenAndServe(":5000", nil))

}

func kafkaConsumer() {
	stream, err := client.ConsumerStream(ctx, &kafkapb.ConsumerStreamRequest{
		Topics:   []string{"test"},
		ClientId: "slackbot",
	})
	if err != nil {
		panic(err)
	}

	for {
		message, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		server.BroadcastTo("test", "message", message.Value)
	}
}
