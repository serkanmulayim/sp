package design

import (
	. "goa.design/goa/v3/dsl"
)

// API declaration
var _ = API("sp", func() {
	Title("SP Portfolio Calculator")
	Description("This api returns the summary of the portfolio of an account")
	Server("sp", func() {
		Host("localhost", func() {
			URI("http://localhost:8080/sp")
		})
	})
})

var _ = Service("portfolio", func() {
	Description("The portfolio service allows for accessing the portfolio summary")
	Method("portfolioSummary", func() {
		Payload(func() {
			Field(1, "account", String, "Account ID")
			Required("account")
		})
		Result(PortfolioResult)
		Error("not_found", NotFound, "Account not found")
		Error("internal_error", InternalError, "Internal error")
		HTTP(func() {
			GET("/sp/portfolio/{account}")
			Response(StatusOK)
			Response("not_found", StatusNotFound)
			Response("internal_error", StatusInternalServerError)
		})
	})
	Files("/openapi.json", "./gen/http/openapi.json")
})

var PortfolioResult = ResultType("application/vnd.sp", func() {
	Description("Portfolio Summary of a given account")
	Reference(Portfolio)
	TypeName("PortfolioResult")

	Attributes(func() {
		Attribute("AccountID", String, "Account ID")
		Attribute("Portfolio", ArrayOf(Currency), "List of currency information")
		Attribute("Total", String, "Total value of the portfolio in dollars")
		Attribute("Missing", ArrayOf(String), "List of tokens which could not be resolved and did not contribute to the total")
	})

	View("default", func() {
		Attribute("AccountID")
		Attribute("Portfolio")
		Attribute("Total")
		Attribute("Missing")
	})

	Required("AccountID", "Portfolio", "Total", "Missing")
})

var Portfolio = Type("Portfolio", func() {
	Description("Portfolio of the account")
	Attribute("AccountID", String, "Account ID")
	Attribute("Portfolio", ArrayOf(Currency), "List of currency information")
	Attribute("Total", String, "Total value of the portfolio in dollars")
	Attribute("Missing", ArrayOf(String), "List of tokens which could not be resolved and did not contribute to the total")
	Required("AccountID", "Portfolio", "Total", "Missing")
})

var Currency = Type("Currency", func() {
	Description("Details of the cryptocurrency in the account")
	Attribute("Denom", String, "Name of the currency")
	Attribute("Price", String, "Price of the currency")
	Attribute("Amount", String, "Amount of the currency in the account")
	Attribute("DenomTotal", String, "Total value of the denomination in the portfolio")
	Required("Denom", "Price", "Amount", "DenomTotal")
})

var NotFound = Type("NotFound", func() {
	Description("NotFound is the error returned when the requested account that does not exist.")
	Attribute("message", String, "Message of error", func() {
		Example("Account cosmos1e0w5t53nrq7p66fye6c8p0ynyhf6y24l4yuxd7 not found")
	})

	Required("message")
})

var InternalError = Type("InternalError", func() {
	Description("InternalError is the error returned when there is an internal error in the application.")
	Attribute("message", String, "Message of error", func() {
		Example("Internal Error")
	})

	Required("message")
})
