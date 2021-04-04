package n26api_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/nhatthm/n26api"
	"github.com/nhatthm/n26api/internal/api"
	"github.com/nhatthm/n26api/pkg/testkit"
	"github.com/nhatthm/n26api/pkg/transaction"
)

func TestClient_FindAllTransactionsInRange(t *testing.T) {
	t.Parallel()

	pageSize := int64(2)
	from := time.Date(2020, 1, 2, 3, 4, 5, 0, time.UTC)
	to := time.Date(2020, 2, 2, 3, 4, 5, 0, time.UTC)
	deviceID := uuid.New()
	id1 := uuid.New()
	id2 := uuid.New()
	id3 := uuid.New()

	transactionURL := fmt.Sprintf("/api/smrt/transactions?from=1577934245000&limit=%d&to=1580612645000", pageSize)

	withFindAllTransactionsInRange := func(result []transaction.Transaction) testkit.ServerOption {
		return testkit.WithFindAllTransactionsInRange(from, to, pageSize, result)
	}

	testCases := []struct {
		scenario             string
		mockServer           testkit.ServerMocker
		expectedTransactions []transaction.Transaction
		expectedError        string
	}{
		{
			scenario: "invalid token",
			mockServer: mockServer(deviceID, func(s *testkit.Server) {
				s.ExpectGet(transactionURL).
					ReturnCode(http.StatusUnauthorized).
					ReturnJSON(api.InvalidTokenError{
						Status: http.StatusUnauthorized,
						Detail: "Invalid token",
						Type:   "error",
						UserMessage: api.UserMessage{
							Title:  "Login attempt expired",
							Detail: "That took too long, please try again.",
						},
						Error:            "invalid_token",
						ErrorDescription: "Invalid token",
					})
			}),
			expectedError: "could not find transactions: invalid token",
		},
		{
			scenario: "server error",
			mockServer: mockServer(deviceID, func(s *testkit.Server) {
				s.ExpectGet(transactionURL).
					ReturnCode(http.StatusInternalServerError)
			}),
			expectedError: "could not find transactions: unexpected response status: 500 Internal Server Error",
		},
		{
			scenario: "success with an empty list",
			mockServer: mockServer(deviceID, withFindAllTransactionsInRange(
				[]transaction.Transaction{},
			)),
			expectedTransactions: []transaction.Transaction{},
		},
		{
			scenario: "success with first page size is less than the limit",
			mockServer: mockServer(deviceID, withFindAllTransactionsInRange(
				[]transaction.Transaction{{ID: id1}},
			)),
			expectedTransactions: []transaction.Transaction{{ID: id1}},
		},
		{
			scenario: "success with first page size is same as the limit and second page is empty",
			mockServer: mockServer(deviceID, withFindAllTransactionsInRange(
				[]transaction.Transaction{{ID: id1}, {ID: id2}},
			)),
			expectedTransactions: []transaction.Transaction{{ID: id1}, {ID: id2}},
		},
		{
			scenario: "success with 2 pages",
			mockServer: mockServer(deviceID, withFindAllTransactionsInRange(
				[]transaction.Transaction{{ID: id1}, {ID: id2}, {ID: id3}},
			)),
			expectedTransactions: []transaction.Transaction{{ID: id1}, {ID: id2}, {ID: id3}},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			s := tc.mockServer(t)
			c := n26api.NewClient(
				n26api.WithBaseURL(s.URL()),
				n26api.WithDeviceID(deviceID),
				n26api.WithCredentials(n26Username, n26Password),
				n26api.WithMFAWait(5*time.Millisecond),
				n26api.WithMFATimeout(time.Second),
				n26api.WithTransactionsPageSize(2),
			)

			result, err := c.FindAllTransactionsInRange(context.Background(), from, to)

			assert.Equal(t, tc.expectedTransactions, result)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}
