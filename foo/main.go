package main

import (
	"context"
	"distributed_cache/client"
	"fmt"
	"log"
)

func main() {
	cl, err := client.New(":3000", client.Options{})
	if err != nil {
		log.Fatal(err)
	}

	err = cl.Set(context.Background(), []byte("foo"), []byte("bar"), 0)
	if err != nil {
		log.Fatal(err)
	}

	value, err := cl.Get(context.Background(), []byte("foo"))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Client response: %s\n", value)

	cl.Close()
}
