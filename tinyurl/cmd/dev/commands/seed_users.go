package commands

import (
	"context"
	"fmt"
	"log"
	"math/rand/v2"
	"os"
	"time"

	"github.com/evrblk/monstera"
	"github.com/evrblk/monstera-example/tinyurl"
	"github.com/evrblk/monstera-example/tinyurl/corepb"
	"github.com/spf13/cobra"
)

var seedUsersCmd = &cobra.Command{
	Use:   "seed-users",
	Short: "seed-users",
	Run: func(cmd *cobra.Command, args []string) {
		// Monstera cluster config
		data, err := os.ReadFile("./cluster_config.pb")
		if err != nil {
			log.Fatal(err)
		}

		clusterConfig, err := monstera.LoadConfigFromProto(data)
		if err != nil {
			log.Fatal(err)
		}

		// Monstera client
		monsteraClient := monstera.NewMonsteraClient(clusterConfig)

		// TinyUrlService client
		tinyUrlServiceCoreApiClient := tinyurl.NewTinyUrlServiceCoreApiMonsteraStub(monsteraClient, &tinyurl.ShardKeyCalculator{})

		numberOfAccounts := 100

		for i := 0; i < numberOfAccounts; i++ {
			now := time.Now()
			userId := rand.Uint32()

			// Create a user
			_, err := tinyUrlServiceCoreApiClient.CreateUser(context.Background(), &corepb.CreateUserRequest{
				UserId: userId,
				Now:    now.UnixNano(),
			})
			if err != nil {
				log.Fatalf("could not create account: %v", err)
			}

			fmt.Printf("created user %s\n", tinyurl.EncodeUserId(userId))
		}
	},
}

func init() {
	rootCmd.AddCommand(seedUsersCmd)
}
