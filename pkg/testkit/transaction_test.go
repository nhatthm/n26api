package testkit

import (
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

	testCases := []struct {
		scenario     string
		transactions []transaction.Transaction
		expect       func() []*Request
	}{
		{
			scenario:     "first page is empty",
			transactions: []transaction.Transaction{},
			expect: func() []*Request {
				s := NewServer(t)

				s.ExpectGet("/api/smrt/transactions?from=1577934245000&limit=2&to=1580612645000").
					ReturnJSON([]transaction.Transaction{})

				return s.ExpectedRequests
			},
		},
		{
			scenario:     "first page size is less than the limit",
			transactions: []transaction.Transaction{{ID: id1}},
			expect: func() []*Request {
				s := NewServer(t)

				s.ExpectGet("/api/smrt/transactions?from=1577934245000&limit=2&to=1580612645000").
					ReturnJSON([]transaction.Transaction{{ID: id1}})

				return s.ExpectedRequests
			},
		},
		{
			scenario:     "first page size is same as the limit",
			transactions: []transaction.Transaction{{ID: id1}, {ID: id2}},
			expect: func() []*Request {
				s := NewServer(t)

				s.ExpectGet("/api/smrt/transactions?from=1577934245000&limit=2&to=1580612645000").
					ReturnJSON([]transaction.Transaction{{ID: id1}, {ID: id2}})

				s.ExpectGet(httpmock.Exactf("/api/smrt/transactions?from=1577934245000&lastId=%s&limit=2&to=1580612645000", id2.String())).
					ReturnJSON([]transaction.Transaction{})

				return s.ExpectedRequests
			},
		},
		{
			scenario:     "two pages",
			transactions: []transaction.Transaction{{ID: id1}, {ID: id2}, {ID: id3}},
			expect: func() []*Request {
				s := NewServer(t)

				s.ExpectGet("/api/smrt/transactions?from=1577934245000&limit=2&to=1580612645000").
					ReturnJSON([]transaction.Transaction{{ID: id1}, {ID: id2}})

				s.ExpectGet(httpmock.Exactf("/api/smrt/transactions?from=1577934245000&lastId=%s&limit=2&to=1580612645000", id2.String())).
					ReturnJSON([]transaction.Transaction{{ID: id3}})

				return s.ExpectedRequests
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
			expected := tc.expect()

			assert.Equal(t, len(expected), len(s.ExpectedRequests))

			for i, expected := range expected {
				actual := s.ExpectedRequests[i]

				expectedBody, err := expected.Handle(nil)
				assert.NoError(t, err)

				actualBody, err := actual.Handle(nil)
				assert.NoError(t, err)

				assert.Equal(t, expected.RequestURI, actual.RequestURI)
				assert.Equal(t, expectedBody, actualBody)
			}
		})
	}
}
