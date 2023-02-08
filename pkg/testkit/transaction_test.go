package testkit

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.nhat.io/httpmock"
	plannerMock "go.nhat.io/httpmock/mock/planner"
	"go.nhat.io/httpmock/planner"
	"go.nhat.io/matcher/v2"

	"github.com/nhatthm/n26api/pkg/transaction"
)

func TestWithFindAllTransactionsInRange(t *testing.T) {
	t.Parallel()

	type result struct {
		URL  string
		Body string
	}

	type expectation interface {
		httpmock.ExpectationHandler
		planner.Expectation
	}

	createTxsResult := func(url string, txs []transaction.Transaction) result {
		data, err := json.Marshal(txs)
		if err != nil {
			panic(err)
		}

		return result{
			URL:  url,
			Body: string(data),
		}
	}

	pageSize := int64(2)
	from := time.Date(2020, 1, 2, 3, 4, 5, 0, time.UTC)
	to := time.Date(2020, 2, 2, 3, 4, 5, 0, time.UTC)
	id1 := uuid.New()
	id2 := uuid.New()
	id3 := uuid.New()

	testCases := []struct {
		scenario        string
		transactions    []transaction.Transaction
		expectedResults []result
	}{
		{
			scenario:     "first page is empty",
			transactions: []transaction.Transaction{},
			expectedResults: []result{
				createTxsResult("/api/smrt/transactions?from=1577934245000&limit=2&to=1580612645000", []transaction.Transaction{}),
			},
		},
		{
			scenario:     "first page size is less than the limit",
			transactions: []transaction.Transaction{{ID: id1}},
			expectedResults: []result{
				createTxsResult("/api/smrt/transactions?from=1577934245000&limit=2&to=1580612645000", []transaction.Transaction{{ID: id1}}),
			},
		},
		{
			scenario:     "first page size is same as the limit",
			transactions: []transaction.Transaction{{ID: id1}, {ID: id2}},
			expectedResults: []result{
				createTxsResult("/api/smrt/transactions?from=1577934245000&limit=2&to=1580612645000", []transaction.Transaction{{ID: id1}, {ID: id2}}),
				createTxsResult(fmt.Sprintf("/api/smrt/transactions?from=1577934245000&lastId=%s&limit=2&to=1580612645000", id2.String()), []transaction.Transaction{}),
			},
		},
		{
			scenario:     "two pages",
			transactions: []transaction.Transaction{{ID: id1}, {ID: id2}, {ID: id3}},
			expectedResults: []result{
				createTxsResult("/api/smrt/transactions?from=1577934245000&limit=2&to=1580612645000", []transaction.Transaction{{ID: id1}, {ID: id2}}),
				createTxsResult(fmt.Sprintf("/api/smrt/transactions?from=1577934245000&lastId=%s&limit=2&to=1580612645000", id2.String()), []transaction.Transaction{{ID: id3}}),
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			actualRequests := make([]expectation, 0)

			p := plannerMock.Mock(func(p *plannerMock.Planner) {
				p.On("Expect", mock.Anything).
					Run(func(args mock.Arguments) {
						actualRequests = append(actualRequests, args[0].(expectation))
					})

				p.On("IsEmpty").Return(true)
			})(t)

			MockEmptyServer(func(s *Server) {
				s.WithPlanner(p)
			}, WithFindAllTransactionsInRange(from, to, pageSize, tc.transactions))(t)

			require.Equal(t, len(tc.expectedResults), len(actualRequests))

			for i, expected := range tc.expectedResults {
				actual := actualRequests[i]

				actualBody, err := handleRequestSuccess(t, actual)
				assert.NoError(t, err)

				assert.Equal(t, matcher.Match(expected.URL), actual.URIMatcher())
				assert.Equal(t, expected.Body, string(actualBody))
			}
		})
	}
}
