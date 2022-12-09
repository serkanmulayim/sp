// Code generated by goa v3.10.2, DO NOT EDIT.
//
// HTTP request path constructors for the portfolio service.
//
// Command:
// $ goa gen sp/design

package server

import (
	"fmt"
)

// PortfolioSummaryPortfolioPath returns the URL path to the portfolio service portfolioSummary HTTP endpoint.
func PortfolioSummaryPortfolioPath(account string) string {
	return fmt.Sprintf("/sp/portfolio/%v", account)
}
