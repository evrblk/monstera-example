package ledger

import (
	"github.com/evrblk/monstera-example/ledger/corepb"
	"github.com/evrblk/monstera-example/ledger/gatewaypb"
)

func accountToFront(account *corepb.Account) *gatewaypb.Account {
	if account == nil {
		return nil
	}

	return &gatewaypb.Account{
		Id:               EncodeAccountId(account.Id),
		AvailableBalance: account.AvailableBalance,
		SettledBalance:   account.SettledBalance,
		CreatedAt:        account.CreatedAt,
		UpdatedAt:        account.UpdatedAt,
	}
}

func transactionToFront(transaction *corepb.Transaction) *gatewaypb.Transaction {
	if transaction == nil {
		return nil
	}

	return &gatewaypb.Transaction{
		Id:          EncodeTransactionId(transaction.Id),
		Amount:      transaction.Amount,
		Description: transaction.Description,
		Status:      transactionStatusToFront(transaction.Status),
		CreatedAt:   transaction.CreatedAt,
		UpdatedAt:   transaction.UpdatedAt,
	}
}

func transactionStatusToFront(status corepb.TransactionStatus) gatewaypb.TransactionStatus {
	switch status {
	case corepb.TransactionStatus_TRANSACTION_STATUS_SETTLED:
		return gatewaypb.TransactionStatus_TRANSACTION_STATUS_SETTLED
	case corepb.TransactionStatus_TRANSACTION_STATUS_CANCELLED:
		return gatewaypb.TransactionStatus_TRANSACTION_STATUS_CANCELLED
	case corepb.TransactionStatus_TRANSACTION_STATUS_INSUFFICIENT_FUNDS:
		return gatewaypb.TransactionStatus_TRANSACTION_STATUS_INSUFFICIENT_FUNDS
	case corepb.TransactionStatus_TRANSACTION_STATUS_PENDING:
		return gatewaypb.TransactionStatus_TRANSACTION_STATUS_PENDING
	default:
		return gatewaypb.TransactionStatus_TRANSACTION_STATUS_INVALID
	}
}
