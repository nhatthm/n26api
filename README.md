# N26 API Client

[![GitHub Releases](https://img.shields.io/github/v/release/nhatthm/n26api)](https://github.com/nhatthm/n26api/releases/latest)
[![Build Status](https://github.com/nhatthm/n26api/actions/workflows/test.yaml/badge.svg)](https://github.com/nhatthm/n26api/actions/workflows/test.yaml)
[![codecov](https://codecov.io/gh/nhatthm/n26api/branch/master/graph/badge.svg?token=eTdAgDE2vR)](https://codecov.io/gh/nhatthm/n26api)
[![Go Report Card](https://goreportcard.com/badge/nhatthm/n26api)](https://goreportcard.com/report/nhatthm/n26api)
[![GoDevDoc](https://img.shields.io/badge/dev-doc-00ADD8?logo=go)](https://pkg.go.dev/github.com/nhatthm/n26api)
[![Donate](https://img.shields.io/badge/Donate-PayPal-green.svg)](https://www.paypal.com/donate/?hosted_button_id=PJZSGJN57TDJY)

Unofficial N26 API Client for Golang.

## Disclaimer

This project is NOT sponsored or funded by N26 nor any of its competitors.

The client is built on a collection of observed API calls and methods provided by [Rots/n26-api](https://github.com/Rots/n26-api).

## Prerequisites

- `Go >= 1.17`

## Install

```bash
go get github.com/nhatthm/n26api
```

## Development

### API Service

TBA

### API Client

TBD

### Testkit

#### Unit Test

The services can be easily mocked by using the `testkit`, all the mocked interface is provided
by [stretchr/testify/mock](https://github.com/stretchr/testify#mock-package)

For example: Mocking `transaction.Finder` service

```go
package test

import (
    "context"
    "testing"
    "time"

    "github.com/google/uuid"
    transactionMock "github.com/nhatthm/n26api/pkg/testkit/transaction"
    "github.com/nhatthm/n26api/pkg/transaction"
    "github.com/stretchr/testify/assert"
)

func TestTransactions(t *testing.T) {
    t.Parallel()

    from := time.Now()
    to := from.Add(time.Hour)
    id := uuid.New()

    f := transactionMock.MockFinder(func(f *transactionMock.Finder) {
        f.On("FindAllTransactionsInRange", context.Background(), from, to).
            Return(
                []transaction.Transaction{
                    {ID: id, OriginalAmount: 3.5},
                },
                nil,
            )
    })(t)

    expected := []transaction.Transaction{
        {ID: id, OriginalAmount: 3.5},
    }

    result, err := f.FindAllTransactionsInRange(context.Background(), from, to)

    assert.Equal(t, expected, result)
    assert.NoError(t, err)
}
```

#### Integration Test

The `testkit` provides a mocked API server for testing.

For example:

```go
package test

import (
    "context"
    "testing"
    "time"

    "github.com/google/uuid"
    "github.com/nhatthm/n26api"
    "github.com/nhatthm/n26api/pkg/testkit"
    "github.com/nhatthm/n26api/pkg/transaction"
)

func TestFindTransactions(t *testing.T) {
    t.Parallel()

    username := "username"
    password := "password"
    deviceID := uuid.New()
    from := time.Now()
    to := from.Add(time.Hour)
    pageSize := int64(1)
    id1 := uuid.New()
    id2 := uuid.New()

    s := testkit.MockServer(username, password, deviceID,
        testkit.WithFindAllTransactionsInRange(
            from, to, pageSize,
            []transaction.Transaction{{ID: id1}, {ID: id2}},
        ),
    )(t)

    c := n26api.NewClient(
        n26api.WithBaseURL(s.URL()),
        n26api.WithDeviceID(deviceID),
        n26api.WithCredentials(username, password),
        n26api.WithMFAWait(5*time.Millisecond),
        n26api.WithMFATimeout(time.Second),
        n26api.WithTransactionsPageSize(pageSize),
    )

    result, err := c.FindAllTransactionsInRange(context.Background(), from, to)

    // assertions
}
```

##### Server Options

###### `WithFindAllTransactionsInRange`

TBD

## Test

### Unit Test

Run

```bash
make test
```

or

```bash
make test-unit
```

### Integration Test

TBD

## GDPR

No information is used or store by this module.

## References

- https://github.com/guitmz/n26 (for MFA)
- https://github.com/Rots/n26-api (for OpenAPI doc)

## Donation

If this project help you reduce time to develop, you can give me a cup of coffee :)

### Paypal donation

[![paypal](https://www.paypalobjects.com/en_US/i/btn/btn_donateCC_LG.gif)](https://www.paypal.com/donate/?hosted_button_id=PJZSGJN57TDJY)

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;or scan this

<img src="https://user-images.githubusercontent.com/1154587/113494222-ad8cb200-94e6-11eb-9ef3-eb883ada222a.png" width="147px" />
