package main

import (
	"fmt"
	"time"

	"google.golang.org/protobuf/types/known/anypb"
)

//go:generate protoc --go_out=. heartbeat.proto

func main() {
	ping := &WebsocketPing{
		Timestamp: time.Now().UnixMilli(),
	}

	// https://stackoverflow.com/questions/64055785/how-to-use-protobuf-any-in-golang
	message1, err := anypb.New(ping)
	if err != nil {
		panic(err)
	}

	payload := &LeafPayload{
		LeafId:   1,
		Message1: message1,
		// Message2: message2,
	}

	fmt.Println(payload)
}
