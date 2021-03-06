package main

import (
	"testing"
)

func TestChatRequestCatchUp(t *testing.T) {
	messaging := NewChat()
	go messaging.Run()

	for i := 0; i < 2; i++ {
		messaging.Incoming <- chatUpdate{
			Message: &ChatMessage{
				SessionID: "session",
				Message:   "hello",
			},
		}
	}

	receiver := make(chan *chatUpdate)
	messaging.UpdateRequest <- ChatUpdateRequest{
		FirstMessage: 0,
		Receiver:     receiver,
	}

	var count int
	for m := range receiver {
		if m.Message.SessionID != "session" {
			t.Error("SessionId not as expected:", m.Message.SessionID)
		}
		if m.Message.Message != "hello" {
			t.Error("Message not as expected:", m.Message.Message)
		}
		if m.UpdateID != count {
			t.Error("UpdateID not as expected:", m.UpdateID)
		}
		count++
	}
	if count != 2 {
		t.Error("Expected two messages")
	}
}

func TestChatRequestPending(t *testing.T) {
	messaging := NewChat()
	go messaging.Run()

	receiver := make(chan *chatUpdate)
	messaging.UpdateRequest <- ChatUpdateRequest{
		FirstMessage: 0,
		Receiver:     receiver,
	}

	messaging.Incoming <- chatUpdate{
		Message: &ChatMessage{
			SessionID: "session",
			Message:   "hello",
		},
	}

	var count int
	for m := range receiver {
		if m.Message.SessionID != "session" {
			t.Error("SessionId not as expected:", m.Message.SessionID)
		}
		if m.Message.Message != "hello" {
			t.Error("Message not as expected:", m.Message.Message)
		}
		if m.UpdateID != count {
			t.Error("UpdateID not as expected:", m.UpdateID)
		}
		count++
	}
	if count != 1 {
		t.Error("Expected one message")
	}
}
