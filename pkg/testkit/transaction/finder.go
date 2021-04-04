package transaction

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/nhatthm/n26api/pkg/transaction"
)

// FinderMocker is Finder mocker.
type FinderMocker func(t testing.TB) *Finder

// NoMockFinder is no mock Finder.
var NoMockFinder = MockFinder()

var _ transaction.Finder = (*Finder)(nil)

// Finder is a transaction.Finder.Finder.
type Finder struct {
	mock.Mock
}

// FindAllTransactionsInRange satisfies transaction.Finder.
func (f *Finder) FindAllTransactionsInRange(ctx context.Context, from time.Time, to time.Time) ([]transaction.Transaction, error) {
	ret := f.Called(ctx, from, to)

	ret1 := ret.Get(0)
	ret2 := ret.Error(1)

	if ret1 == nil {
		return nil, ret2
	}

	return ret1.([]transaction.Transaction), ret2
}

// mockFinder mocks transaction.Finder.Finder interface.
func mockFinder(mocks ...func(f *Finder)) *Finder {
	f := &Finder{}

	for _, m := range mocks {
		m(f)
	}

	return f
}

// MockFinder creates Finder mock with cleanup to ensure all the expectations are met.
func MockFinder(mocks ...func(f *Finder)) FinderMocker {
	return func(t testing.TB) *Finder {
		f := mockFinder(mocks...)

		t.Cleanup(func() {
			assert.True(t, f.Mock.AssertExpectations(t))
		})

		return f
	}
}
