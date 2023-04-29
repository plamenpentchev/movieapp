package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"movieexample.com/pkg/discovery"
	"movieexample.com/pkg/discovery/consul"
	"movieexample.com/rating/internal/controller/rating"
	httpHandler "movieexample.com/rating/internal/handler/http"
	"movieexample.com/rating/internal/repository/memory"
)

const serviceName = "rating"

func main() {
	var port int
	flag.IntVar(&port, "port", 8082, "API handler port")
	flag.Parse()

	log.Printf("starting rating service on port %d ...", port)

	registry, err := consul.NewRegistry("localhost:8500")
	if err != nil {
		panic(err)
	}

	instanceID := discovery.GenerateInstanceID(serviceName)
	ctx := context.Background()
	if err = registry.Register(ctx, instanceID, serviceName, fmt.Sprintf("localhost:%d", port)); err != nil {
		panic("Failed to register rating service: " + err.Error())
	}

	go func() {
		for {
			if err := registry.ReportHealthyStatus(instanceID, serviceName); err != nil {
				log.Printf("")
			}
			time.Sleep(1 * time.Second)
		}
	}()
	defer registry.Deregister(ctx, instanceID, serviceName)

	repo := memory.New()
	ctrl := rating.New(repo)
	h := httpHandler.New(ctrl)
	http.Handle("/rating", http.HandlerFunc(h.Handle))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		log.Fatal(err)
	}

}
