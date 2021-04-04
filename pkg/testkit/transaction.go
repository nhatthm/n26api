package testkit

import (
	"net/url"
	"strconv"
	"time"

	"github.com/nhatthm/n26api/internal/api"
	"github.com/nhatthm/n26api/pkg/transaction"
	"github.com/nhatthm/n26api/pkg/util"
)

func buildSmrtTransactionsURI(request api.GetAPISmrtTransactionsRequest) string {
	requestURI := "/api/smrt/transactions"

	query := make(url.Values, 7)

	if request.From != nil {
		query.Set("from", strconv.FormatInt(*request.From, 10))
	}

	if request.To != nil {
		query.Set("to", strconv.FormatInt(*request.To, 10))
	}

	if request.Limit != nil {
		query.Set("limit", strconv.FormatInt(*request.Limit, 10))
	}

	if request.Pending != nil {
		query.Set("pending", strconv.FormatBool(*request.Pending))
	}

	if request.Categories != nil {
		query.Set("categories", *request.Categories)
	}

	if request.TextFilter != nil {
		query.Set("textFilter", *request.TextFilter)
	}

	if request.LastID != nil {
		query.Set("lastId", *request.LastID)
	}

	if len(query) > 0 {
		requestURI += "?" + query.Encode()
	}

	return requestURI
}

// WithFindAllTransactionsInRange sets expectations for finding all transactions in a range.
func WithFindAllTransactionsInRange(
	from time.Time,
	to time.Time,
	pageSize int64,
	result []transaction.Transaction,
) ServerOption {
	return func(s *Server) {
		last := (*string)(nil)

		for {
			requestURI := buildSmrtTransactionsURI(api.GetAPISmrtTransactionsRequest{
				From:   util.Int64Ptr(util.UnixTimestampMS(from)),
				To:     util.Int64Ptr(util.UnixTimestampMS(to)),
				Limit:  util.Int64Ptr(pageSize),
				LastID: last,
			})
			count := int64(len(result))
			end := pageSize

			if count < pageSize {
				end = count
			}

			json := result[0:end]

			s.ExpectGet(requestURI).ReturnJSON(json)

			// If this is the last page.
			if count < pageSize {
				break
			}

			last = util.StringPtr(json[pageSize-1].ID.String())
			result = result[end:]
		}
	}
}
