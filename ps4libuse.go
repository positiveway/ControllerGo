package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/mrasband/ps4"
)

func mainPS() {
	inputs, err := ps4.Discover()
	if err != nil {
		fmt.Printf("Error discovering controller: %s\n", err)
		os.Exit(1)
	}

	var device *ps4.Input
	for _, input := range inputs {
		if input.Type == ps4.Controller {
			device = input
			break
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	events, _ := ps4.Watch(ctx, device)
	for e := range events {
		fmt.Printf("%+v\n", e)
		time.Sleep(100 * time.Millisecond)
	}
}
