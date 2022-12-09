// Code generated by goa v3.10.2, DO NOT EDIT.
//
// portfolio views
//
// Command:
// $ goa gen sp/design

package views

import (
	goa "goa.design/goa/v3/pkg"
)

// PortfolioResult is the viewed result type that is projected based on a view.
type PortfolioResult struct {
	// Type to project
	Projected *PortfolioResultView
	// View to render
	View string
}

// PortfolioResultView is a type that runs validations on a projected type.
type PortfolioResultView struct {
	// Account ID
	AccountID *string
	// List of currency information
	Portfolio []*CurrencyView
	// Total value of the portfolio in dollars
	Total *string
	// List of tokens which could not be resolved and did not contribute to the
	// total
	Missing []string
}

// CurrencyView is a type that runs validations on a projected type.
type CurrencyView struct {
	// Name of the currency
	Denom *string
	// Price of the currency
	Price *string
	// Amount of the currency in the account
	Amount *string
	// Total value of the denomination in the portfolio
	DenomTotal *string
}

var (
	// PortfolioResultMap is a map indexing the attribute names of PortfolioResult
	// by view name.
	PortfolioResultMap = map[string][]string{
		"default": {
			"AccountID",
			"Portfolio",
			"Total",
			"Missing",
		},
	}
)

// ValidatePortfolioResult runs the validations defined on the viewed result
// type PortfolioResult.
func ValidatePortfolioResult(result *PortfolioResult) (err error) {
	switch result.View {
	case "default", "":
		err = ValidatePortfolioResultView(result.Projected)
	default:
		err = goa.InvalidEnumValueError("view", result.View, []interface{}{"default"})
	}
	return
}

// ValidatePortfolioResultView runs the validations defined on
// PortfolioResultView using the "default" view.
func ValidatePortfolioResultView(result *PortfolioResultView) (err error) {
	if result.AccountID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("AccountID", "result"))
	}
	if result.Portfolio == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("Portfolio", "result"))
	}
	if result.Total == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("Total", "result"))
	}
	if result.Missing == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("Missing", "result"))
	}
	for _, e := range result.Portfolio {
		if e != nil {
			if err2 := ValidateCurrencyView(e); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
	}
	return
}

// ValidateCurrencyView runs the validations defined on CurrencyView.
func ValidateCurrencyView(result *CurrencyView) (err error) {
	if result.Denom == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("Denom", "result"))
	}
	if result.Price == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("Price", "result"))
	}
	if result.Amount == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("Amount", "result"))
	}
	if result.DenomTotal == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("DenomTotal", "result"))
	}
	return
}
