# N26 API Client

[![Build Status](https://github.com/nhatthm/n26api/actions/workflows/test.yaml/badge.svg)](https://github.com/nhatthm/n26api/actions/workflows/test.yaml)
[![codecov](https://codecov.io/gh/nhatthm/n26api/branch/master/graph/badge.svg?token=eTdAgDE2vR)](https://codecov.io/gh/nhatthm/n26api)
[![Go Report Card](https://goreportcard.com/badge/github.com/nhatthm/httpmock)](https://goreportcard.com/report/github.com/nhatthm/httpmock)
[![GoDevDoc](https://img.shields.io/badge/dev-doc-00ADD8?logo=go)](https://pkg.go.dev/github.com/nhatthm/n26api)
[![Donate](https://img.shields.io/badge/Donate-PayPal-green.svg)](https://www.paypal.com/donate/?hosted_button_id=PJZSGJN57TDJY)

Unofficial N26 API Client for Golang.

## Disclaimer

This project is NOT sponsored or funded by N26 nor any of its competitors.

The client is built on a collection of observed API calls and methods provided by [Rots/n26-api](https://github.com/Rots/n26-api).

## Prerequisites

- `Go >= 1.14`

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

The service can be easily mocked by using the `testkit`, all the mocked interface is provided
by [stretchr/testify/mock](https://github.com/stretchr/testify#mock-package)

For example

```go
package test

import (
    "context"
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
