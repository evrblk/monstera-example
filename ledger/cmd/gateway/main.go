package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/evrblk/monstera"
	"github.com/evrblk/monstera-example/ledger"
	"github.com/evrblk/monstera-example/ledger/gatewaypb"
	"google.golang.org/grpc"
)

var (
	port               = flag.Int("port", 0, "The server port")
	monsteraConfigPath = flag.String("monstera-config", "", "Monstera cluster config path")
)

func main() {
	log.Println("Initializing API Gateway Server...")

	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Load monstera cluster config
	data, err := os.ReadFile(*monsteraConfigPath)
	if err != nil {
		log.Fatal(err)
	}

	clusterConfig, err := monstera.LoadConfigFromProto(data)
	if err != nil {
		log.Fatal(err)
	}

	// Create Monstera client
	monsteraClient := monstera.NewMonsteraClient(clusterConfig)
	monsteraClient.Start()

	// LedgerService client
	ledgerServiceCoreApiClient := ledger.NewLedgerServiceCoreApiMonsteraStub(monsteraClient, &ledger.ShardKeyCalculator{})

	grpcServer := grpc.NewServer()

	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	go func() {
		select {
		case <-c:
			log.Println("Received SIGINT. Shutting down...")
			cancel()
			grpcServer.GracefulStop()
			monsteraClient.Stop()
		case <-ctx.Done():
		}
	}()
	defer func() {
		signal.Stop(c)
		cancel()
	}()

	// Create and register Gateway server
	ledgerServiceApiGatewayServer := ledger.NewLedgerServiceApiServer(ledgerServiceCoreApiClient)
	defer ledgerServiceApiGatewayServer.Close()
	gatewaypb.RegisterLedgerServiceApiServer(grpcServer, ledgerServiceApiGatewayServer)

	log.Println("Starting API Gateway Server...")
	grpcServer.Serve(lis)
}
