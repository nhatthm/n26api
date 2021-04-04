package transaction_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"

	transactionMock "github.com/nhatthm/n26api/pkg/testkit/transaction"
	"github.com/nhatthm/n26api/pkg/transaction"
)

func TestFinder_FindAllTransactionsInRange(t *testing.T) {
	t.Parallel()

	from := time.Date(2020, 1, 2, 3, 4, 5, 0, time.UTC)
	to := time.Date(2020, 2, 2, 3, 4, 5, 0, time.UTC)
	id := uuid.New()

	testCases := []struct {
		scenario       string
		mockFinder     transactionMock.FinderMocker
		expectedResult []transaction.Transaction
		expectedError  string
	}{
		{
			scenario: "result is nil",
			mockFinder: transactionMock.MockFinder(func(f *transactionMock.Finder) {
				f.On("FindAllTransactionsInRange", context.Background(), from, to).
					Return(nil, nil)
			}),
		},
		{
			scenario: "result is not nil",
			mockFinder: transactionMock.MockFinder(func(f *transactionMock.Finder) {
				f.On("FindAllTransactionsInRange", context.Background(), from, to).
					Return([]transaction.Transaction{{ID: id}}, nil)
			}),
			expectedResult: []transaction.Transaction{{ID: id}},
		},
		{
			scenario: "error",
			mockFinder: transactionMock.MockFinder(func(f *transactionMock.Finder) {
				f.On("FindAllTransactionsInRange", context.Background(), from, to).
					Return(nil, errors.New("find error"))
			}),
			expectedError: "find error",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

		})
	}
}
