package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/evrblk/monstera-example/dlocks/gatewaypb"

	"path/filepath"

	"github.com/evrblk/monstera"
	"github.com/evrblk/monstera-example/dlocks"
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

	locksCore := dlocks.NewLocksCore(dataStore, []byte{0x00, 0x00, 0x00, 0x00}, []byte{0x00, 0x00, 0x00, 0x00})
	namespacesCore := dlocks.NewNamespacesCore(dataStore, []byte{0x00, 0x00, 0x00, 0x00}, []byte{0x00, 0x00, 0x00, 0x00})
	accountsCore := dlocks.NewAccountsCore(dataStore)

	// LocksService client
	locksServiceCoreApiClient := dlocks.NewLocksServiceCoreApiStandaloneStub(accountsCore, namespacesCore, locksCore)

	authenticationMiddleware := dlocks.AuthenticationMiddleware{}
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(authenticationMiddleware.Unary),
	)

	// Create and register Gateway server
	locksServiceApiGatewayServer := dlocks.NewLocksServiceApiServer(locksServiceCoreApiClient)
	defer locksServiceApiGatewayServer.Close()
	gatewaypb.RegisterLocksServiceApiServer(grpcServer, locksServiceApiGatewayServer)

	log.Println("Starting API Gateway Server...")
	grpcServer.Serve(lis)
}
