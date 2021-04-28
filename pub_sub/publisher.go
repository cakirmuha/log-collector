package pub_sub

import (
	"context"
	"fmt"
	"cloud.google.com/go/pubsub"
)

func (p *PubSub) Publish(msgData []byte) error {
	ctx := context.Background()

	result := p.Topic.Publish(ctx, &pubsub.Message{
		Data: msgData,
	})
	// Block until the result is returned and a server-generated
	// ID is returned for the published message
	id, err := result.Get(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("Published a message; msg ID: %v\n", id)
	return nil
}
