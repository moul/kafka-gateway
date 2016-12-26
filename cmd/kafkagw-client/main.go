package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/urfave/cli"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/moul/kafka-gateway/gen/pb"
)

var (
	conn   *grpc.ClientConn
	ctx    context.Context
	logger log.Logger
	client kafkapb.KafkaServiceClient
)

func main() {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "grpc-server",
			Value:  "127.0.0.1:9000",
			EnvVar: "GRPC_SERVER",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:      "producer",
			Action:    producer,
			ArgsUsage: "[message]",
			Usage:     "Produce a message",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "kafka-topic",
					Value:  "test",
					EnvVar: "KAFKA_TOPIC",
				},
				cli.StringFlag{
					Name:   "kafka-key",
					Value:  "test",
					EnvVar: "KAFKA_KEY",
				},
			},
		},
		{
			Name:   "consumer",
			Action: consumer,
			Usage:  "Consume a message",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "kafka-topics",
					Value:  "test",
					EnvVar: "KAFKA_TOPICS",
				},
				cli.StringFlag{
					Name:   "kafka-client-id",
					Value:  "test",
					EnvVar: "KAFKA_CLIENT_ID",
				},
			},
		},
		{
			Name:   "consumer-stream",
			Action: consumerStream,
			Usage:  "Consume multiple messages",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "kafka-topics",
					Value:  "test",
					EnvVar: "KAFKA_TOPICS",
				},
				cli.StringFlag{
					Name:   "kafka-client-id",
					Value:  "test",
					EnvVar: "KAFKA_CLIENT_ID",
				},
			},
		},
	}
	app.Before = before
	defer after()
	app.Run(os.Args)
}

func before(c *cli.Context) (err error) {
	conn, err = grpc.Dial(
		c.String("grpc-server"),
		grpc.WithInsecure(),
		grpc.WithTimeout(time.Second),
	)
	ctx = context.Background()
	logger = log.NewNopLogger()
	client = kafkapb.NewKafkaServiceClient(conn)
	return
}

func after() {
	if conn != nil {
		conn.Close()
	}
}

func producer(c *cli.Context) error {
	var value string
	if len(c.Args()) > 0 {
		value = c.Args()[0]
	}
	ret, err := client.Producer(ctx, &kafkapb.ProducerRequest{
		Topic: c.String("kafka-topic"),
		Key:   c.String("kafka-key"),
		Value: value,
	})
	fmt.Println(ret, err)
	return nil
}

func consumer(c *cli.Context) error {
	var topics []string
	if len(c.String("kafka-topics")) > 0 {
		topics = strings.Split(c.String("kafka-topics"), ",")
	}
	ret, err := client.Consumer(ctx, &kafkapb.ConsumerRequest{
		Topics:   topics,
		ClientId: c.String("kafka-client-id"),
	})
	fmt.Println(ret, err)
	return nil
}

func consumerStream(c *cli.Context) error {
	var topics []string
	if len(c.String("kafka-topics")) > 0 {
		topics = strings.Split(c.String("kafka-topics"), ",")
	}
	stream, err := client.ConsumerStream(ctx, &kafkapb.ConsumerStreamRequest{
		Topics:   topics,
		ClientId: c.String("kafka-client-id"),
	})
	if err != nil {
		return err
	}

	for {
		message, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		fmt.Println(message)
	}

	return nil
}
