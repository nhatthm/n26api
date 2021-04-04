package transaction

import (
	"context"
	"time"
)

// Finder is a service to find n26 transactions.
type Finder interface {
	// FindAllTransactionsInRange finds all transactions in a time period.
	FindAllTransactionsInRange(ctx context.Context, from time.Time, to time.Time) ([]Transaction, error)
}
