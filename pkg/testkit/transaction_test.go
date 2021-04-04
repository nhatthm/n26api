package testkit

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/nhatthm/httpmock"
	"github.com/stretchr/testify/assert"

	"github.com/nhatthm/n26api/pkg/transaction"
)

func TestWithFindAllTransactionsInRange(t *testing.T) {
	t.Parallel()

	pageSize := int64(2)
	from := time.Date(2020, 1, 2, 3, 4, 5, 0, time.UTC)
	to := time.Date(2020, 2, 2, 3, 4, 5, 0, time.UTC)
	id1 := uuid.New()
	id2 := uuid.New()
	id3 := uuid.New()

	expectRequest := func(requestURI string, transactions []transaction.Transaction) *Request {
		return &Request{
			RequestURI: requestURI,
			Do: func(r *http.Request) ([]byte, error) {
				return json.Marshal(transactions)
			},
		}
	}

	testCases := []struct {
		scenario     string
		transactions []transaction.Transaction
		expected     []*Request
	}{
		{
			scenario:     "first page is empty",
			transactions: []transaction.Transaction{},
			expected: []*Request{
				expectRequest(
					"/api/smrt/transactions?from=1577934245000&limit=2&to=1580612645000",
					[]transaction.Transaction{},
				),
			},
		},
		{
			scenario:     "first page size is less than the limit",
			transactions: []transaction.Transaction{{ID: id1}},
			expected: []*Request{
				expectRequest(
					"/api/smrt/transactions?from=1577934245000&limit=2&to=1580612645000",
					[]transaction.Transaction{{ID: id1}},
				),
			},
		},
		{
			scenario:     "first page size is same as the limit",
			transactions: []transaction.Transaction{{ID: id1}, {ID: id2}},
			expected: []*Request{
				expectRequest(
					"/api/smrt/transactions?from=1577934245000&limit=2&to=1580612645000",
					[]transaction.Transaction{{ID: id1}, {ID: id2}},
				),
				expectRequest(
					fmt.Sprintf("/api/smrt/transactions?from=1577934245000&lastId=%s&limit=2&to=1580612645000", id2.String()),
					[]transaction.Transaction{},
				),
			},
		},
		{
			scenario:     "two pages",
			transactions: []transaction.Transaction{{ID: id1}, {ID: id2}, {ID: id3}},
			expected: []*Request{
				expectRequest(
					"/api/smrt/transactions?from=1577934245000&limit=2&to=1580612645000",
					[]transaction.Transaction{{ID: id1}, {ID: id2}},
				),
				expectRequest(
					fmt.Sprintf("/api/smrt/transactions?from=1577934245000&lastId=%s&limit=2&to=1580612645000", id2.String()),
					[]transaction.Transaction{{ID: id3}},
				),
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			var mu sync.Mutex

			mu.Lock()
			defer mu.Unlock()

			s := &Server{
				Server: &httpmock.Server{},
			}

			WithFindAllTransactionsInRange(from, to, pageSize, tc.transactions)(s)

			assert.Equal(t, len(tc.expected), len(s.ExpectedRequests))

			for i, expected := range tc.expected {
				actual := s.ExpectedRequests[i]

				expectedBody, err := expected.Do(nil)
				assert.NoError(t, err)

				actualBody, err := actual.Do(nil)
				assert.NoError(t, err)

				assert.Equal(t, expected.RequestURI, actual.RequestURI)
				assert.Equal(t, expectedBody, actualBody)
			}
		})
	}
}
