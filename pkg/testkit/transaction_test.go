package testkit

import (
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/nhatthm/httpmock"
	plannerMock "github.com/nhatthm/httpmock/mock/planner"
	"github.com/nhatthm/httpmock/request"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

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
		scenario         string
		transactions     []transaction.Transaction
		expectedRequests []*Request
	}{
		{
			scenario:     "first page is empty",
			transactions: []transaction.Transaction{},
			expectedRequests: []*Request{
				request.NewRequest(&sync.Mutex{}, httpmock.MethodGet, "/api/smrt/transactions?from=1577934245000&limit=2&to=1580612645000").
					ReturnJSON([]transaction.Transaction{}),
			},
		},
		{
			scenario:     "first page size is less than the limit",
			transactions: []transaction.Transaction{{ID: id1}},
			expectedRequests: []*Request{
				request.NewRequest(&sync.Mutex{}, httpmock.MethodGet, "/api/smrt/transactions?from=1577934245000&limit=2&to=1580612645000").
					ReturnJSON([]transaction.Transaction{{ID: id1}}),
			},
		},
		{
			scenario:     "first page size is same as the limit",
			transactions: []transaction.Transaction{{ID: id1}, {ID: id2}},
			expectedRequests: []*Request{
				request.NewRequest(&sync.Mutex{}, httpmock.MethodGet, "/api/smrt/transactions?from=1577934245000&limit=2&to=1580612645000").
					ReturnJSON([]transaction.Transaction{{ID: id1}, {ID: id2}}),
				request.NewRequest(&sync.Mutex{}, httpmock.MethodGet, httpmock.Exactf("/api/smrt/transactions?from=1577934245000&lastId=%s&limit=2&to=1580612645000", id2.String())).
					ReturnJSON([]transaction.Transaction{}),
			},
		},
		{
			scenario:     "two pages",
			transactions: []transaction.Transaction{{ID: id1}, {ID: id2}, {ID: id3}},
			expectedRequests: []*Request{
				request.NewRequest(&sync.Mutex{}, httpmock.MethodGet, "/api/smrt/transactions?from=1577934245000&limit=2&to=1580612645000").
					ReturnJSON([]transaction.Transaction{{ID: id1}, {ID: id2}}),
				request.NewRequest(&sync.Mutex{}, httpmock.MethodGet, httpmock.Exactf("/api/smrt/transactions?from=1577934245000&lastId=%s&limit=2&to=1580612645000", id2.String())).
					ReturnJSON([]transaction.Transaction{{ID: id3}}),
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			actualRequests := make([]*request.Request, 0)

			p := plannerMock.Mock(func(p *plannerMock.Planner) {
				p.On("Expect", mock.Anything).
					Run(func(args mock.Arguments) {
						actualRequests = append(actualRequests, args[0].(*request.Request))
					})

				p.On("IsEmpty").Return(true)
			})(t)

			MockEmptyServer(func(s *Server) {
				s.WithPlanner(p)
			}, WithFindAllTransactionsInRange(from, to, pageSize, tc.transactions))(t)

			assert.Equal(t, len(tc.expectedRequests), len(actualRequests))

			for i, expected := range tc.expectedRequests {
				actual := actualRequests[i]

				expectedBody, err := handleRequestSuccess(t, expected)
				assert.NoError(t, err)

				actualBody, err := handleRequestSuccess(t, expected)
				assert.NoError(t, err)

				assert.Equal(t, request.URIMatcher(expected), request.URIMatcher(actual))
				assert.Equal(t, expectedBody, actualBody)
			}
		})
	}
}
