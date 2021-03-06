// Code generated by github.com/swaggest/swac v0.1.19, DO NOT EDIT.

package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/nhatthm/n26api/pkg/transaction"
)

// GetAPISmrtTransactionsRequest is operation request value.
type GetAPISmrtTransactionsRequest struct {
	// From is an optional `from` parameter in query.
	// Timestamp - milliseconds since 1970 in CET
	From *int64
	To   *int64 // To is an optional `to` parameter in query.
	// Limit is an optional `limit` parameter in query.
	// Limit the number of transactions to return
	Limit   *int64
	Pending *bool // Pending is an optional `pending` parameter in query.
	// Categories is an optional `categories` parameter in query.
	// Comma separated list of category IDs
	Categories *string
	// TextFilter is an optional `textFilter` parameter in query.
	// Query string to search for
	TextFilter *string
	LastID     *string // LastID is an optional `lastId` parameter in query.
}

// encode creates *http.Request for request data.
func (request *GetAPISmrtTransactionsRequest) encode(ctx context.Context, baseURL string) (*http.Request, error) {
	requestURI := baseURL + "/api/smrt/transactions"

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

	req, err := http.NewRequest(http.MethodGet, requestURI, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")

	req = req.WithContext(ctx)

	return req, err
}

// GetAPISmrtTransactionsResponse is operation response value.
type GetAPISmrtTransactionsResponse struct {
	StatusCode        int
	ValueOK           []transaction.Transaction // ValueOK is a value of 200 OK response.
	ValueUnauthorized *InvalidTokenError        // ValueUnauthorized is a value of 401 Unauthorized response.
}

// decode loads data from *http.Response.
func (result *GetAPISmrtTransactionsResponse) decode(resp *http.Response) error {
	var err error

	dump := bytes.NewBuffer(nil)
	body := io.TeeReader(resp.Body, dump)

	result.StatusCode = resp.StatusCode

	switch resp.StatusCode {
	case http.StatusOK:
		err = json.NewDecoder(body).Decode(&result.ValueOK)
	case http.StatusUnauthorized:
		err = json.NewDecoder(body).Decode(&result.ValueUnauthorized)
	default:
		_, readErr := ioutil.ReadAll(body)
		if readErr != nil {
			err = errors.New("unexpected response status: " + resp.Status +
				", could not read response body: " + readErr.Error())
		} else {
			err = errors.New("unexpected response status: " + resp.Status)
		}
	}

	if err != nil {
		return responseError{
			resp: resp,
			body: dump.Bytes(),
			err:  err,
		}
	}

	return nil
}

// GetAPISmrtTransactions performs REST operation.
func (c *Client) GetAPISmrtTransactions(ctx context.Context, request GetAPISmrtTransactionsRequest) (result GetAPISmrtTransactionsResponse, err error) {
	if c.InstrumentCtxFunc != nil {
		ctx = c.InstrumentCtxFunc(ctx, http.MethodGet, "/api/smrt/transactions", &request)
	}

	if c.Timeout != 0 {
		var cancel func()
		ctx, cancel = context.WithTimeout(ctx, c.Timeout)

		defer cancel()
	}

	req, err := request.encode(ctx, c.BaseURL)
	if err != nil {
		return result, err
	}

	resp, err := c.transport.RoundTrip(req)

	if err != nil {
		return result, err
	}

	defer func() {
		closeErr := resp.Body.Close()
		if closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	err = result.decode(resp)

	return result, err
}
