package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/evrblk/monstera-example/ledger/gatewaypb"

	"path/filepath"

	"github.com/evrblk/monstera"
	"github.com/evrblk/monstera-example/ledger"
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

	accountsCore := ledger.NewAccountsCore(dataStore, []byte{0x00, 0x00, 0x00, 0x00}, []byte{0x00, 0x00, 0x00, 0x00})

	// LedgerService client
	ledgerServiceCoreApiClient := ledger.NewLedgerServiceCoreApiStandaloneStub(accountsCore)

	grpcServer := grpc.NewServer()

	// Create and register Gateway server
	ledgerServiceApiGatewayServer := ledger.NewLedgerServiceApiServer(ledgerServiceCoreApiClient)
	defer ledgerServiceApiGatewayServer.Close()
	gatewaypb.RegisterLedgerServiceApiServer(grpcServer, ledgerServiceApiGatewayServer)

	log.Println("Starting API Gateway Server...")
	grpcServer.Serve(lis)
}
