package pub_sub

import (
	"cloud.google.com/go/pubsub"
	"context"
	"log"
	"os"
)

type PubSub struct {
	Topic *pubsub.Topic
}

func NewPubSub() *PubSub {
	ctx := context.Background()

	client, err := pubsub.NewClient(ctx, "logservice-311815")
	if err != nil {
		log.Println("Could not create pubsub Client: %v", err)
		return nil
	}

	topic := os.Getenv("TOPIC")
	if topic == "" {
		log.Println("There is no topic set in env variable")
		return nil
	}

	return &PubSub{
		Topic: client.Topic(topic),
	}
}
