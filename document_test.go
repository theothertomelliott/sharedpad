package main

import (
	"strconv"
	"testing"
	"time"
)

func TestGetSession(t *testing.T) {
	messaging := NewDocument()
	go messaging.Run()
	for i := 0; i < 10; i++ {
		c := make(chan Session)
		messaging.SessionRequest <- c
		select {
		case session := <-c:
			if session.Id != strconv.Itoa(i) {
				t.Errorf("Wrong id received")
			}
		case <-time.After(1 * time.Second):
			t.Errorf("Timed out waiting for session")
		}
	}
}

func TestMessageRequestCatchUp(t *testing.T) {
	messaging := NewDocument()
	go messaging.Run()

	for i := 0; i < 2; i++ {
		messaging.Incoming <- update{
			SessionID: "session",
		}
	}

	receiver := make(chan *update)
	messaging.UpdateRequest <- UpdateRequest{
		FirstMessage: 0,
		Receiver:     receiver,
	}

	var count int
	for m := range receiver {
		if m.SessionID != "session" {
			t.Error("SessionId not as expected:", m.SessionID)
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

func TestMessageRequestPending(t *testing.T) {
	messaging := NewDocument()
	go messaging.Run()

	receiver := make(chan *update)
	messaging.UpdateRequest <- UpdateRequest{
		FirstMessage: 0,
		Receiver:     receiver,
	}

	messaging.Incoming <- update{
		SessionID: "session",
	}

	var count int
	for m := range receiver {
		if m.SessionID != "session" {
			t.Error("SessionId not as expected:", m.SessionID)
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
