// Code generated by goa v3.10.2, DO NOT EDIT.
//
// portfolio service
//
// Command:
// $ goa gen sp/design

package portfolio

import (
	"context"
	portfolioviews "sp/gen/portfolio/views"
)

// The portfolio service allows for accessing the portfolio summary
type Service interface {
	// PortfolioSummary implements portfolioSummary.
	PortfolioSummary(context.Context, *PortfolioSummaryPayload) (res *PortfolioResult, err error)
}

// ServiceName is the name of the service as defined in the design. This is the
// same value that is set in the endpoint request contexts under the ServiceKey
// key.
const ServiceName = "portfolio"

// MethodNames lists the service method names as defined in the design. These
// are the same values that are set in the endpoint request contexts under the
// MethodKey key.
var MethodNames = [1]string{"portfolioSummary"}

// Details of the cryptocurrency in the account
type Currency struct {
	// Name of the currency
	Denom string
	// Price of the currency
	Price string
	// Amount of the currency in the account
	Amount string
	// Total value of the denomination in the portfolio
	DenomTotal string
}

// InternalError is the error returned when there is an internal error in the
// application.
type InternalError struct {
	// Message of error
	Message string
}

// NotFound is the error returned when the requested account that does not
// exist.
type NotFound struct {
	// Message of error
	Message string
}

// PortfolioResult is the result type of the portfolio service portfolioSummary
// method.
type PortfolioResult struct {
	// Account ID
	AccountID string
	// List of currency information
	Portfolio []*Currency
	// Total value of the portfolio in dollars
	Total string
	// List of tokens which could not be resolved and did not contribute to the
	// total
	Missing []string
}

// PortfolioSummaryPayload is the payload type of the portfolio service
// portfolioSummary method.
type PortfolioSummaryPayload struct {
	// Account ID
	Account string
}

// Error returns an error description.
func (e *InternalError) Error() string {
	return "InternalError is the error returned when there is an internal error in the application."
}

// ErrorName returns "InternalError".
//
// Deprecated: Use GoaErrorName - https://github.com/goadesign/goa/issues/3105
func (e *InternalError) ErrorName() string {
	return e.GoaErrorName()
}

// GoaErrorName returns "InternalError".
func (e *InternalError) GoaErrorName() string {
	return "internal_error"
}

// Error returns an error description.
func (e *NotFound) Error() string {
	return "NotFound is the error returned when the requested account that does not exist."
}

// ErrorName returns "NotFound".
//
// Deprecated: Use GoaErrorName - https://github.com/goadesign/goa/issues/3105
func (e *NotFound) ErrorName() string {
	return e.GoaErrorName()
}

// GoaErrorName returns "NotFound".
func (e *NotFound) GoaErrorName() string {
	return "not_found"
}

// NewPortfolioResult initializes result type PortfolioResult from viewed
// result type PortfolioResult.
func NewPortfolioResult(vres *portfolioviews.PortfolioResult) *PortfolioResult {
	return newPortfolioResult(vres.Projected)
}

// NewViewedPortfolioResult initializes viewed result type PortfolioResult from
// result type PortfolioResult using the given view.
func NewViewedPortfolioResult(res *PortfolioResult, view string) *portfolioviews.PortfolioResult {
	p := newPortfolioResultView(res)
	return &portfolioviews.PortfolioResult{Projected: p, View: "default"}
}

// newPortfolioResult converts projected type PortfolioResult to service type
// PortfolioResult.
func newPortfolioResult(vres *portfolioviews.PortfolioResultView) *PortfolioResult {
	res := &PortfolioResult{}
	if vres.AccountID != nil {
		res.AccountID = *vres.AccountID
	}
	if vres.Total != nil {
		res.Total = *vres.Total
	}
	if vres.Portfolio != nil {
		res.Portfolio = make([]*Currency, len(vres.Portfolio))
		for i, val := range vres.Portfolio {
			res.Portfolio[i] = transformPortfolioviewsCurrencyViewToCurrency(val)
		}
	}
	if vres.Missing != nil {
		res.Missing = make([]string, len(vres.Missing))
		for i, val := range vres.Missing {
			res.Missing[i] = val
		}
	}
	return res
}

// newPortfolioResultView projects result type PortfolioResult to projected
// type PortfolioResultView using the "default" view.
func newPortfolioResultView(res *PortfolioResult) *portfolioviews.PortfolioResultView {
	vres := &portfolioviews.PortfolioResultView{
		AccountID: &res.AccountID,
		Total:     &res.Total,
	}
	if res.Portfolio != nil {
		vres.Portfolio = make([]*portfolioviews.CurrencyView, len(res.Portfolio))
		for i, val := range res.Portfolio {
			vres.Portfolio[i] = transformCurrencyToPortfolioviewsCurrencyView(val)
		}
	}
	if res.Missing != nil {
		vres.Missing = make([]string, len(res.Missing))
		for i, val := range res.Missing {
			vres.Missing[i] = val
		}
	}
	return vres
}

// transformPortfolioviewsCurrencyViewToCurrency builds a value of type
// *Currency from a value of type *portfolioviews.CurrencyView.
func transformPortfolioviewsCurrencyViewToCurrency(v *portfolioviews.CurrencyView) *Currency {
	if v == nil {
		return nil
	}
	res := &Currency{
		Denom:      *v.Denom,
		Price:      *v.Price,
		Amount:     *v.Amount,
		DenomTotal: *v.DenomTotal,
	}

	return res
}

// transformCurrencyToPortfolioviewsCurrencyView builds a value of type
// *portfolioviews.CurrencyView from a value of type *Currency.
func transformCurrencyToPortfolioviewsCurrencyView(v *Currency) *portfolioviews.CurrencyView {
	res := &portfolioviews.CurrencyView{
		Denom:      &v.Denom,
		Price:      &v.Price,
		Amount:     &v.Amount,
		DenomTotal: &v.DenomTotal,
	}

	return res
}
