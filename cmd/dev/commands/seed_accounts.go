package commands

import (
	"context"
	"fmt"
	"log"
	"math/rand/v2"
	"os"
	"time"

	"github.com/evrblk/monstera"
	monsteraexample "github.com/evrblk/monstera-example"
	"github.com/evrblk/monstera-example/corepb"
	"github.com/spf13/cobra"
)

var seedAccountsCmd = &cobra.Command{
	Use:   "seed-accounts",
	Short: "seed-accounts",
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

		// ExampleService client
		exampleServiceCoreApiClient := monsteraexample.NewExampleServiceCoreApiMonsteraStub(monsteraClient, &monsteraexample.ShardKeyCalculator{})

		numberOfAccounts := 100

		for i := 0; i < numberOfAccounts; i++ {
			email := fmt.Sprintf("account+%04d@example.com", i+1)
			now := time.Now()
			accountId := rand.Uint64()

			// Create account
			_, err := exampleServiceCoreApiClient.CreateAccount(context.Background(), &corepb.CreateAccountRequest{
				AccountId:             accountId,
				FullName:              "Doogie Howser",
				Email:                 email,
				Now:                   now.UnixNano(),
				MaxNumberOfNamespaces: 25,
			})
			if err != nil {
				log.Fatalf("could not create account: %v", err)
			}

			fmt.Printf("created account %s\n", monsteraexample.EncodeAccountId(accountId))
		}
	},
}

func init() {
	rootCmd.AddCommand(seedAccountsCmd)
}
