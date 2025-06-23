package commands

import (
	"context"
	"fmt"
	"github.com/evrblk/monstera-example/ledger/gatewaypb"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
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
		client := gatewaypb.NewLedgerServiceApiClient(conn)

		ctx := context.Background()

		// get account balance
		resp1, err := client.GetAccount(ctx, &gatewaypb.GetAccountRequest{
			AccountId: accountId,
		})
		if err != nil {
			log.Fatalf("could not get account: %v", err)
		}
		fmt.Printf("Account %s balance: %d available, %d settled\n", resp1.Account.Id, resp1.Account.AvailableBalance, resp1.Account.SettledBalance)

		// create an instant topup transaction +100
		resp2, err := client.CreateTransaction(ctx, &gatewaypb.CreateTransactionRequest{
			AccountId:   accountId,
			Amount:      100,
			Description: "Instant Topup 1",
			Settled:     true,
		})
		if err != nil {
			log.Fatalf("could not create transaction: %v", err)
		}
		fmt.Printf("Created transaction %s for amount %d: %s\n", resp2.Transaction.Id, resp2.Transaction.Amount, resp2.Transaction.Status)

		// create a pending purchase transaction -10
		resp3, err := client.CreateTransaction(ctx, &gatewaypb.CreateTransactionRequest{
			AccountId:   accountId,
			Amount:      -10,
			Description: "Purchase 1",
			Settled:     false,
		})
		if err != nil {
			log.Fatalf("could not create transaction: %v", err)
		}
		fmt.Printf("Created transaction %s for amount %d: %s\n", resp3.Transaction.Id, resp3.Transaction.Amount, resp3.Transaction.Status)

		// get account balance
		resp4, err := client.GetAccount(ctx, &gatewaypb.GetAccountRequest{
			AccountId: accountId,
		})
		if err != nil {
			log.Fatalf("could not get account: %v", err)
		}
		fmt.Printf("Account %s balance: %d available, %d settled\n", resp4.Account.Id, resp4.Account.AvailableBalance, resp4.Account.SettledBalance)

		// settle pending transaction
		resp5, err := client.SettleTransaction(ctx, &gatewaypb.SettleTransactionRequest{
			TransactionId: resp3.Transaction.Id,
		})
		if err != nil {
			log.Fatalf("could not settle transaction: %v", err)
		}
		fmt.Printf("Settled transaction %s for amount %d: %s\n", resp5.Transaction.Id, resp5.Transaction.Amount, resp5.Transaction.Status)

		// get account balance
		resp6, err := client.GetAccount(ctx, &gatewaypb.GetAccountRequest{
			AccountId: accountId,
		})
		if err != nil {
			log.Fatalf("could not get account: %v", err)
		}
		fmt.Printf("Account %s balance: %d available, %d settled\n", resp6.Account.Id, resp6.Account.AvailableBalance, resp6.Account.SettledBalance)

		// exceed available funds with -1000 transaction
		resp7, err := client.CreateTransaction(ctx, &gatewaypb.CreateTransactionRequest{
			AccountId:   accountId,
			Amount:      -1000,
			Description: "Purchase 2",
			Settled:     true,
		})
		if err != nil {
			log.Fatalf("could not create transaction: %v", err)
		}
		fmt.Printf("Created transaction %s for amount %d: %s\n", resp7.Transaction.Id, resp7.Transaction.Amount, resp7.Transaction.Status)
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
