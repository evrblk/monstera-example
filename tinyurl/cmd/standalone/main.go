package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/evrblk/monstera-example/tinyurl/gatewaypb"

	"path/filepath"

	"github.com/evrblk/monstera"
	"github.com/evrblk/monstera-example/tinyurl"
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

	usersCore := tinyurl.NewUsersCore(dataStore, []byte{0x00, 0x00, 0x00, 0x00}, []byte{0x00, 0x00, 0x00, 0x00})
	shortUrlsCore := tinyurl.NewShortUrlsCore(dataStore, []byte{0x00, 0x00, 0x00, 0x00}, []byte{0x00, 0x00, 0x00, 0x00})

	// TinyUrlService client
	tinyUrlServiceCoreApiClient := tinyurl.NewTinyUrlServiceCoreApiStandaloneStub(shortUrlsCore, usersCore)

	grpcServer := grpc.NewServer()

	// Create and register Gateway server
	tinyUrlServiceApiGatewayServer := tinyurl.NewTinyUrlServiceApiServer(tinyUrlServiceCoreApiClient)
	defer tinyUrlServiceApiGatewayServer.Close()
	gatewaypb.RegisterTinyUrlServiceApiServer(grpcServer, tinyUrlServiceApiGatewayServer)

	log.Println("Starting API Gateway Server...")
	grpcServer.Serve(lis)
}
