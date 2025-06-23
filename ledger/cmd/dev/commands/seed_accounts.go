package commands

import (
	"context"
	"fmt"
	"log"
	"math/rand/v2"
	"os"
	"time"

	"github.com/evrblk/monstera"
	"github.com/evrblk/monstera-example/ledger"
	"github.com/evrblk/monstera-example/ledger/corepb"
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

		// LedgerService client
		ledgerServiceCoreApiClient := ledger.NewLedgerServiceCoreApiMonsteraStub(monsteraClient, &ledger.ShardKeyCalculator{})

		numberOfAccounts := 100

		for i := 0; i < numberOfAccounts; i++ {
			now := time.Now()
			accountId := rand.Uint64()

			// Create an account
			_, err := ledgerServiceCoreApiClient.CreateAccount(context.Background(), &corepb.CreateAccountRequest{
				AccountId: accountId,
				Now:       now.UnixNano(),
			})
			if err != nil {
				log.Fatalf("could not create account: %v", err)
			}

			fmt.Printf("created account %s\n", ledger.EncodeAccountId(accountId))
		}
	},
}

func init() {
	rootCmd.AddCommand(seedAccountsCmd)
}
