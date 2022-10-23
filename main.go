package main

import (
	"context"
	"fmt"

	"github.com/hazelcast/hazelcast-go-client"
)

func main() {
	ctx := context.TODO()
	client, err := hazelcast.StartNewClient(ctx)
	if err != nil {
		panic(err)
	}
	m, err := client.GetMap(ctx, "my-map")
	if err != nil {
		panic(err)
	}
	if err := m.Set(ctx, "some-key", "some-value"); err != nil {
		panic(err)
	}
	v, err := m.Get(ctx, "some-key")
	if err != nil {
		panic(err)
	}
	fmt.Println("Received value:", v)
	if err := client.Shutdown(ctx); err != nil {
		panic(err)
	}
}
