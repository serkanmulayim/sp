// Code generated by goa v3.10.2, DO NOT EDIT.
//
// portfolio endpoints
//
// Command:
// $ goa gen sp/design

package portfolio

import (
	"context"

	goa "goa.design/goa/v3/pkg"
)

// Endpoints wraps the "portfolio" service endpoints.
type Endpoints struct {
	PortfolioSummary goa.Endpoint
}

// NewEndpoints wraps the methods of the "portfolio" service with endpoints.
func NewEndpoints(s Service) *Endpoints {
	return &Endpoints{
		PortfolioSummary: NewPortfolioSummaryEndpoint(s),
	}
}

// Use applies the given middleware to all the "portfolio" service endpoints.
func (e *Endpoints) Use(m func(goa.Endpoint) goa.Endpoint) {
	e.PortfolioSummary = m(e.PortfolioSummary)
}

// NewPortfolioSummaryEndpoint returns an endpoint function that calls the
// method "portfolioSummary" of service "portfolio".
func NewPortfolioSummaryEndpoint(s Service) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*PortfolioSummaryPayload)
		res, err := s.PortfolioSummary(ctx, p)
		if err != nil {
			return nil, err
		}
		vres := NewViewedPortfolioResult(res, "default")
		return vres, nil
	}
}
