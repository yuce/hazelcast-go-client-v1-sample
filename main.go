package main

import (
	"context"
	"fmt"
	"github.com/hazelcast/hazelcast-go-client"
	"log"
)

func main() {
	ctx := context.TODO()
	client, err := hazelcast.StartNewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	m, err := client.GetMap(ctx, "my-map")
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Set(ctx, "some-key", "some-value"); err != nil {
		log.Fatal(err)
	}
	v, err := m.Get(ctx, "some-key")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Received value:", v)
	client.Shutdown(ctx)
}
