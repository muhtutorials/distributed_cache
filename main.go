package main

import (
	"distributed_cache/cache"
	"log"
)

func main() {
	opts := ServerOpts{
		ListenAddr: ":4000",
		IsLeader:   false,
		LeaderAddr: ":3000",
	}

	server := NewServer(opts, cache.New())
	err := server.Start()
	if err != nil {
		log.Fatal(err)
	}
}
