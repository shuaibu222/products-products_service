package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/shuaibu222/p-products/products"
	"google.golang.org/grpc"
)

const (
	webPort  = "80"
	gRpcPort = "50002"
)

type Config struct{}

func main() {
	app := Config{}

	go app.gRPCListen()

	// start web server
	log.Println("Starting service on port", webPort)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic()
	}

}

func (app *Config) gRPCListen() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", gRpcPort))
	if err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}

	s := grpc.NewServer()

	products.RegisterProductsServiceServer(s, &ProductsServer{})

	log.Printf("gRPC Server started on port %s", gRpcPort)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}
}
