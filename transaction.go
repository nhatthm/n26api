package n26api

import (
	"context"
	"time"

	"github.com/bool64/ctxd"

	"github.com/nhatthm/n26api/internal/api"
	"github.com/nhatthm/n26api/pkg/transaction"
	"github.com/nhatthm/n26api/pkg/util"
)

var _ transaction.Finder = (*Client)(nil)

// WithTransactionsPageSize sets page size limit for finding transactions.
func WithTransactionsPageSize(limit int64) Option {
	return func(c *Client) {
		c.config.transactionsPageSize = limit
	}
}

func (c *Client) findTransactions(ctx context.Context, req api.GetAPISmrtTransactionsRequest) ([]transaction.Transaction, error) {
	res, err := c.api.GetAPISmrtTransactions(ctx, req)
	if err != nil {
		return nil, err
	}

	if res.ValueUnauthorized != nil {
		return nil, ctxd.NewError(ctx, "invalid token", "response", res)
	}

	if res.ValueOK == nil {
		return nil, ctxd.NewError(ctx, "unexpected response", "response", res)
	}

	return res.ValueOK, nil
}

// FindAllTransactionsInRange finds all transactions in a time period.
func (c *Client) FindAllTransactionsInRange(ctx context.Context, from time.Time, to time.Time) ([]transaction.Transaction, error) {
	page := 1
	limit := c.config.transactionsPageSize
	count := limit
	last := (*string)(nil)
	result := make([]transaction.Transaction, 0, limit)

	for {
		if count < limit {
			break
		}

		trans, err := c.findTransactions(ctx, api.GetAPISmrtTransactionsRequest{
			From:   util.Int64Ptr(util.UnixTimestampMS(from)),
			To:     util.Int64Ptr(util.UnixTimestampMS(to)),
			Limit:  util.Int64Ptr(limit),
			LastID: last,
		})
		if err != nil {
			return nil, ctxd.WrapError(ctx, err, "could not find transactions",
				"from", from,
				"to", to,
				"limit", limit,
				"page", page,
			)
		}

		count = int64(len(trans))
		page++

		if count > 0 {
			last = util.StringPtr(trans[count-1].ID.String())
			result = append(result, trans...)
		}
	}

	return result, nil
}
