// Code generated by goa v3.10.2, DO NOT EDIT.
//
// portfolio client HTTP transport
//
// Command:
// $ goa gen sp/design

package client

import (
	"context"
	"net/http"

	goahttp "goa.design/goa/v3/http"
	goa "goa.design/goa/v3/pkg"
)

// Client lists the portfolio service endpoint HTTP clients.
type Client struct {
	// PortfolioSummary Doer is the HTTP client used to make requests to the
	// portfolioSummary endpoint.
	PortfolioSummaryDoer goahttp.Doer

	// RestoreResponseBody controls whether the response bodies are reset after
	// decoding so they can be read again.
	RestoreResponseBody bool

	scheme  string
	host    string
	encoder func(*http.Request) goahttp.Encoder
	decoder func(*http.Response) goahttp.Decoder
}

// NewClient instantiates HTTP clients for all the portfolio service servers.
func NewClient(
	scheme string,
	host string,
	doer goahttp.Doer,
	enc func(*http.Request) goahttp.Encoder,
	dec func(*http.Response) goahttp.Decoder,
	restoreBody bool,
) *Client {
	return &Client{
		PortfolioSummaryDoer: doer,
		RestoreResponseBody:  restoreBody,
		scheme:               scheme,
		host:                 host,
		decoder:              dec,
		encoder:              enc,
	}
}

// PortfolioSummary returns an endpoint that makes HTTP requests to the
// portfolio service portfolioSummary server.
func (c *Client) PortfolioSummary() goa.Endpoint {
	var (
		decodeResponse = DecodePortfolioSummaryResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildPortfolioSummaryRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.PortfolioSummaryDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("portfolio", "portfolioSummary", err)
		}
		return decodeResponse(resp)
	}
}
