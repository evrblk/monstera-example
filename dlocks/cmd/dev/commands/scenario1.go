package commands

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/evrblk/monstera-example/dlocks/gatewaypb"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var accountId string

var scenario1Cmd = &cobra.Command{
	Use:   "scenario-1",
	Short: "scenario-1",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("account id: %s\n", accountId)

		// Connect to grpc server
		conn, err := grpc.NewClient("localhost:8000", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		client := gatewaypb.NewLocksServiceApiClient(conn)

		// Add header into request context
		ctx := metadata.AppendToOutgoingContext(context.Background(),
			"account-id", accountId)

		// Get namespace
		_, err = client.GetNamespace(ctx, &gatewaypb.GetNamespaceRequest{
			NamespaceName: "my-namespace",
		})
		if err != nil {
			if errors.Is(err, status.Error(codes.NotFound, "namespace not found")) {
				fmt.Printf("namespace not found, creating namespace\n")
				// Create namespace
				_, err = client.CreateNamespace(ctx, &gatewaypb.CreateNamespaceRequest{
					Name:        "my-namespace",
					Description: "my-namespace-description",
				})
				if err != nil {
					log.Fatalf("could not create namespace: %v", err)
				}
			} else {
				log.Fatalf("could not get namespace: %v", err)
			}
		}

		// Acquire lock my-lock-1 with write lock
		resp1, err := client.AcquireLock(ctx, &gatewaypb.AcquireLockRequest{
			NamespaceName: "my-namespace",
			LockName:      "my-lock-1",
			ExpiresAt:     time.Now().Add(1 * time.Minute).UnixNano(),
			WriteLock:     true,
			ProcessId:     "process-id-1",
		})
		if err != nil {
			log.Fatalf("could not acquire lock: %v", err)
		}

		fmt.Printf("acquired lock my-lock-1 with write lock: %v\n", resp1)
	},
}

func init() {
	rootCmd.AddCommand(scenario1Cmd)

	scenario1Cmd.PersistentFlags().StringVarP(&accountId, "account-id", "", "", "Account id")
	err := scenario1Cmd.MarkPersistentFlagRequired("account-id")
	if err != nil {
		panic(err)
	}
}
