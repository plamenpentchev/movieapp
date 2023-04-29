package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"movieexample.com/movie/internal/controller/movie"
	metadatgateway "movieexample.com/movie/internal/gateway/metadata/with_discovery/http"
	ratinggateway "movieexample.com/movie/internal/gateway/rating/with_discovery/http"
	httpHandler "movieexample.com/movie/internal/handler/http"
	"movieexample.com/pkg/discovery"
	"movieexample.com/pkg/discovery/consul"
)

const serviceName = "movie"
const serviceNameMetadata = "metadata"
const serviceNameRating = "rating"

func main() {
	var port int
	flag.IntVar(&port, "port", 8083, "API handler port")
	flag.Parse()

	log.Printf("starting movie service on port %d ...", port)

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

	ratingGateway := ratinggateway.New(registry)
	metadataGateway := metadatgateway.New(registry)
	ctrl := movie.New(ratingGateway, metadataGateway)
	h := httpHandler.New(ctrl)
	http.Handle("/movie", http.HandlerFunc(h.GetMovieDetails))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		log.Fatal(err)
	}

}
