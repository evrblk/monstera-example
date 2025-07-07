package commands

import (
	"context"
	"fmt"
	"github.com/evrblk/monstera-example/tinyurl/gatewaypb"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

var userId string

var scenario1Cmd = &cobra.Command{
	Use:   "scenario-1",
	Short: "scenario-1",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("user id: %s\n", userId)

		// Connect to grpc server
		conn, err := grpc.NewClient("localhost:8000", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		client := gatewaypb.NewTinyUrlServiceApiClient(conn)

		ctx := context.Background()

		// create short url
		resp1, err := client.CreateShortUrl(ctx, &gatewaypb.CreateShortUrlRequest{
			UserId:  userId,
			FullUrl: "https://everblack.dev/docs/monstera/overview/",
		})
		if err != nil {
			log.Fatalf("could not create a short url: %v", err)
		}
		fmt.Printf("Created Short URL: %s\n", resp1.ShortUrl.Id)

		// GetShortUrl can read from followers and there might be replication lag, so we need to wait a bit
		time.Sleep(1 * time.Second)

		// get short url
		resp2, err := client.GetShortUrl(ctx, &gatewaypb.GetShortUrlRequest{
			ShortUrl: resp1.ShortUrl.Id,
		})
		if err != nil {
			log.Fatalf("could not get a short url: %v", err)
		}
		fmt.Printf("Short URL: %s -> %s\n", resp2.ShortUrl.Id, resp2.ShortUrl.FullUrl)

	},
}

func init() {
	rootCmd.AddCommand(scenario1Cmd)

	scenario1Cmd.PersistentFlags().StringVarP(&userId, "user-id", "", "", "User id")
	err := scenario1Cmd.MarkPersistentFlagRequired("user-id")
	if err != nil {
		panic(err)
	}
}
