package pubsub_test

import (
	"testing"

	"github.com/aicacia/pubsub"
)

type Message struct {
	name string
}

func TestPubSub(t *testing.T) {
	pubsub := pubsub.NewPubSub[Message]()
	subscriber0 := pubsub.Subscribe()
	subscriber1 := pubsub.Subscribe()

	if len(pubsub.Subscribers()) != 2 {
		t.Fatalf("Invalid subscriber count %d != 2", len(pubsub.Subscribers()))
	}

	pubsub.Publish(Message{name: "Hello, world!"})

	msg0, ok0 := <-subscriber0.C
	if !ok0 {
		t.Fatalf("Failed to read message from subscriber0")
	}
	if msg0.name != "Hello, world!" {
		t.Fatalf("Invalid message subscriber0: %s", msg0)
	}
	subscriber0.Close()

	msg1, ok1 := <-subscriber1.C
	if !ok1 {
		t.Fatalf("Failed to read message from subscriber1")
	}
	if msg1.name != "Hello, world!" {
		t.Fatalf("Invalid message for subscriber1: %s", msg1)
	}
	subscriber1.Close()

	if len(pubsub.Subscribers()) != 0 {
		t.Fatalf("Invalid subscriber count %d != 0", len(pubsub.Subscribers()))
	}
}
