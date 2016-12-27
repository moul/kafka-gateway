package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"google.golang.org/grpc"

	slackbot "github.com/BeepBoopHQ/go-slackbot"
	"github.com/nlopes/slack"
	"golang.org/x/net/context"

	"github.com/moul/kafka-gateway/gen/pb"
)

var (
	client kafkapb.KafkaServiceClient
	conn   *grpc.ClientConn
	ctx    context.Context
)

func init() {
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
}

func main() {
	bot := slackbot.New(os.Getenv("SLACK_TOKEN"))
	bot.Hear("(?)kafka produce (.*)").MessageHandler(kafkaProduceHandler)
	bot.Run()
}

func kafkaProduceHandler(ctx context.Context, bot *slackbot.Bot, msg *slack.MessageEvent) {
	text := strings.Join(strings.Split(slackbot.StripDirectMention(msg.Text), " ")[2:], " ")
	ret, err := client.Producer(ctx, &kafkapb.ProducerRequest{
		Topic: "test",
		Key:   "test",
		Value: text,
	})
	if err != nil {
		bot.Reply(msg, fmt.Sprintf("Failed: %v", err), slackbot.WithTyping)
	} else {
		bot.Reply(msg, fmt.Sprintf("Done: %v", ret), slackbot.WithTyping)
	}
}
