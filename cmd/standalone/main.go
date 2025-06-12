package main

import (
	"flag"
	"fmt"
	"github.com/evrblk/monstera-example/gatewaypb"
	"log"
	"net"

	"github.com/evrblk/monstera"
	"github.com/evrblk/monstera-example"
	"google.golang.org/grpc"
	"path/filepath"
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

	// Create and register Gateway server
	exampleServiceApiGatewayServer := monsteraexample.NewExampleServiceApiServer(exampleServiceCoreApiClient)
	defer exampleServiceApiGatewayServer.Close()
	gatewaypb.RegisterExampleServiceApiServer(grpcServer, exampleServiceApiGatewayServer)

	log.Println("Starting API Gateway Server...")
	grpcServer.Serve(lis)
}
