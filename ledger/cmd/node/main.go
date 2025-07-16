package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"

	"github.com/evrblk/monstera"
	"github.com/evrblk/monstera-example/ledger"
	"google.golang.org/grpc"
)

var (
	monsteraPort       = flag.Int("port", 0, "Monstera server port")
	dataDir            = flag.String("data-dir", ".", "Base directory for data")
	monsteraConfigPath = flag.String("monstera-config", "", "Monstera cluster config path")
	nodeId             = flag.String("node-id", "", "Monstera node id")
)

func main() {
	log.Println("Initializing Node server...")

	flag.Parse()

	// Load monstera cluster config
	data, err := os.ReadFile(*monsteraConfigPath)
	if err != nil {
		log.Fatal(err)
	}

	clusterConfig, err := monstera.LoadConfigFromProto(data)
	if err != nil {
		log.Fatal(err)
	}

	dataStore := monstera.NewBadgerStore(filepath.Join(*dataDir, "data"))
	raftStore := monstera.NewBadgerStore(filepath.Join(*dataDir, "raft"))

	monsteraNode := monstera.NewNode(*dataDir, *nodeId, clusterConfig, raftStore, monstera.DefaultMonsteraNodeConfig)

	monsteraNode.RegisterApplicationCore(&monstera.ApplicationCoreDescriptor{
		Name:                   "Accounts",
		RestoreSnapshotOnStart: false,
		CoreFactoryFunc: func(application *monstera.Application, shard *monstera.Shard, replica *monstera.Replica) monstera.ApplicationCore {
			return ledger.NewAccountsCoreAdapter(ledger.NewAccountsCore(dataStore, shard.LowerBound, shard.UpperBound))
		},
	})

	// Starting Monstera node
	monsteraNode.Start()

	// Starting Monstera gRPC server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *monsteraPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	monsteraServer := monstera.NewMonsteraServer(monsteraNode)

	grpcServer := grpc.NewServer(
		grpc.MaxRecvMsgSize(50 * 1024 * 1024),
	)
	monstera.RegisterMonsteraApiServer(grpcServer, monsteraServer)

	cleanupDone := &sync.WaitGroup{}
	cleanupDone.Add(1)

	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	go func() {
		select {
		case <-c:
			log.Println("Received SIGINT. Shutting down...")
			cancel()
			monsteraNode.Stop()
			grpcServer.GracefulStop()
			dataStore.Close()
			raftStore.Close()
		case <-ctx.Done():
		}
		cleanupDone.Done()
		log.Printf("Cleanup done")
	}()
	defer func() {
		signal.Stop(c)
		cancel()
	}()

	err = grpcServer.Serve(lis)
	if err != nil {
		log.Printf("Monstera gRPC server stopped: %s", err)
	} else {
		log.Printf("Monstera gRPC server stopped")
	}

	cleanupDone.Wait()

	log.Printf("Exiting...")
}
