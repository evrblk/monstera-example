package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/evrblk/monstera-example/gatewaypb"
	"log"
	"net"

	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/evrblk/monstera"
	"github.com/evrblk/monstera-example"
	"google.golang.org/grpc"
)

var (
	dataDir = flag.String("data-dir", ".", "Base directory for data")
	port    = flag.Int("port", 0, "The server port")
)

func main() {
	log.Println("Initializing standalone application...")

	flag.Parse()

	dataStore := monstera.NewBadgerStore(filepath.Join(*dataDir, "data"))

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	locksCore := monsteraexample.NewLocksCore(dataStore, []byte{0x00, 0x00, 0x00, 0x00}, []byte{0x00, 0x00, 0x00, 0x00})
	namespacesCore := monsteraexample.NewNamespacesCore(dataStore, []byte{0x00, 0x00, 0x00, 0x00}, []byte{0x00, 0x00, 0x00, 0x00})
	accountsCore := monsteraexample.NewAccountsCore(dataStore)

	// ExampleService client
	exampleServiceCoreApiClient := monsteraexample.NewExampleServiceCoreApiStandaloneStub(accountsCore, namespacesCore, locksCore)

	authenticationMiddleware := monsteraexample.AuthenticationMiddleware{}
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(authenticationMiddleware.Unary),
	)

	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	go func() {
		select {
		case <-c:
			log.Println("Received SIGINT. Shutting down...")
			cancel()
			grpcServer.GracefulStop()
		case <-ctx.Done():
		}
	}()
	defer func() {
		signal.Stop(c)
		cancel()
	}()

	// Create and register Gateway server
	exampleServiceApiGatewayServer := monsteraexample.NewExampleServiceApiServer(exampleServiceCoreApiClient)
	defer exampleServiceApiGatewayServer.Close()
	gatewaypb.RegisterExampleServiceApiServer(grpcServer, exampleServiceApiGatewayServer)

	log.Println("Starting API Gateway Server...")
	grpcServer.Serve(lis)

}
