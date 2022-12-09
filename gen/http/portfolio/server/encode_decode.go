// Code generated by goa v3.10.2, DO NOT EDIT.
//
// portfolio HTTP server encoders and decoders
//
// Command:
// $ goa gen sp/design

package server

import (
	"context"
	"errors"
	"net/http"
	portfolio "sp/gen/portfolio"
	portfolioviews "sp/gen/portfolio/views"

	goahttp "goa.design/goa/v3/http"
	goa "goa.design/goa/v3/pkg"
)

// EncodePortfolioSummaryResponse returns an encoder for responses returned by
// the portfolio portfolioSummary endpoint.
func EncodePortfolioSummaryResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, interface{}) error {
	return func(ctx context.Context, w http.ResponseWriter, v interface{}) error {
		res := v.(*portfolioviews.PortfolioResult)
		enc := encoder(ctx, w)
		body := NewPortfolioSummaryResponseBody(res.Projected)
		w.WriteHeader(http.StatusOK)
		return enc.Encode(body)
	}
}

// DecodePortfolioSummaryRequest returns a decoder for requests sent to the
// portfolio portfolioSummary endpoint.
func DecodePortfolioSummaryRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	return func(r *http.Request) (interface{}, error) {
		var (
			account string

			params = mux.Vars(r)
		)
		account = params["account"]
		payload := NewPortfolioSummaryPayload(account)

		return payload, nil
	}
}

// EncodePortfolioSummaryError returns an encoder for errors returned by the
// portfolioSummary portfolio endpoint.
func EncodePortfolioSummaryError(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder, formatter func(ctx context.Context, err error) goahttp.Statuser) func(context.Context, http.ResponseWriter, error) error {
	encodeError := goahttp.ErrorEncoder(encoder, formatter)
	return func(ctx context.Context, w http.ResponseWriter, v error) error {
		var en goa.GoaErrorNamer
		if !errors.As(v, &en) {
			return encodeError(ctx, w, v)
		}
		switch en.GoaErrorName() {
		case "internal_error":
			var res *portfolio.InternalError
			errors.As(v, &res)
			enc := encoder(ctx, w)
			var body interface{}
			if formatter != nil {
				body = formatter(ctx, res)
			} else {
				body = NewPortfolioSummaryInternalErrorResponseBody(res)
			}
			w.Header().Set("goa-error", res.GoaErrorName())
			w.WriteHeader(http.StatusInternalServerError)
			return enc.Encode(body)
		case "not_found":
			var res *portfolio.NotFound
			errors.As(v, &res)
			enc := encoder(ctx, w)
			var body interface{}
			if formatter != nil {
				body = formatter(ctx, res)
			} else {
				body = NewPortfolioSummaryNotFoundResponseBody(res)
			}
			w.Header().Set("goa-error", res.GoaErrorName())
			w.WriteHeader(http.StatusNotFound)
			return enc.Encode(body)
		default:
			return encodeError(ctx, w, v)
		}
	}
}

// marshalPortfolioviewsCurrencyViewToCurrencyResponseBody builds a value of
// type *CurrencyResponseBody from a value of type *portfolioviews.CurrencyView.
func marshalPortfolioviewsCurrencyViewToCurrencyResponseBody(v *portfolioviews.CurrencyView) *CurrencyResponseBody {
	res := &CurrencyResponseBody{
		Denom:      *v.Denom,
		Price:      *v.Price,
		Amount:     *v.Amount,
		DenomTotal: *v.DenomTotal,
	}

	return res
}