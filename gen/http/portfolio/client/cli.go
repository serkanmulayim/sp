// Code generated by goa v3.10.2, DO NOT EDIT.
//
// portfolio HTTP client CLI support package
//
// Command:
// $ goa gen sp/design

package client

import (
	portfolio "sp/gen/portfolio"
)

// BuildPortfolioSummaryPayload builds the payload for the portfolio
// portfolioSummary endpoint from CLI flags.
func BuildPortfolioSummaryPayload(portfolioPortfolioSummaryAccount string) (*portfolio.PortfolioSummaryPayload, error) {
	var account string
	{
		account = portfolioPortfolioSummaryAccount
	}
	v := &portfolio.PortfolioSummaryPayload{}
	v.Account = account

	return v, nil
}