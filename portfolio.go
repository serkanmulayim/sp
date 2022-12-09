package sp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	portfolio "sp/gen/portfolio"
	"strings"
	"sync"
	"time"

	"github.com/shopspring/decimal"
)

const (
	ACCOUNT_BALANCE_URL = "https://api.cosmos.network/cosmos/bank/v1beta1/balances/"
	CG_MARKET_URL       = "https://api.coingecko.com/api/v3/coins/markets"
	COIN_METADATA_FILE  = "coinhelper/coinmetadata.json"
)

// portfolio service example implementation.
// The example methods log the requests and return zero values.
type portfoliosrvc struct {
	logger *log.Logger
}

type ReturnType int

const (
	SUCCESS        ReturnType = 0
	NOT_FOUND      ReturnType = 1
	INTERNAL_ERROR ReturnType = 2
)

var (
	coin2IdMap            map[string]MetadataEntry
	internalErrorInstance = portfolio.InternalError{Message: "Internal Error"}
	notFoundErrorInstance = portfolio.NotFound{Message: "Not Found"}
	exponents             = map[int]decimal.Decimal{0: decimal.NewFromInt(1), 3: decimal.New(1, -3), 6: decimal.New(1, -6), 9: decimal.New(1, -9)}
	decimal0              = decimal.NewFromInt(0)
)

// entries for response from 'https://api.coingecko.com/api/v3/coins/list'
type CgSymbolEntry struct {
	Id     string `json:"id"`
	Symbol string `json:"symbol"`
	Name   string `json:"name"`
}

// metadata encapsulating CGSymbol with CosmosSymbol and exponential for the symbol
type MetadataEntry struct {
	*CgSymbolEntry
	Exp          int    `json:"exp"`
	CosmosSymbol string `json:"cosmos_symbol"`
}

// encapsulating Metadata with the amount
type DenomEntry struct {
	*MetadataEntry
	Amount string `json:"amount"`
}

// entries for the response from "https://api.coingecko.com/api/v3/coins/markets"
type MarketResponse struct {
	CurrentPrice float64 `json:"current_price"`
	Symbol       string  `json:"symbol"`
	//Do not care about the rest
}

// NewPortfolio returns the portfolio service implementation.
func NewPortfolio(logger *log.Logger) portfolio.Service {
	var once sync.Once
	once.Do(func() {
		coin2IdMap = make(map[string]MetadataEntry)
		jsonn, err := ioutil.ReadFile(COIN_METADATA_FILE)
		if err != nil {
			log.Fatalf("%s could not be found. It is a required file. To generate it go to coinhelper folder and run coinhelper.go application first...", COIN_METADATA_FILE)
		}

		err = json.Unmarshal(jsonn, &coin2IdMap)
		if err != nil {
			log.Print(err)
			log.Fatalf("%s is corrupted. It is a required file. To generate it go to coinhelper folder and run coinhelper.go application first...", COIN_METADATA_FILE)
		}

	})
	return &portfoliosrvc{logger}
}

// PortfolioSummary implements portfolioSummary.
func (s *portfoliosrvc) PortfolioSummary(ctx context.Context, p *portfolio.PortfolioSummaryPayload) (res *portfolio.PortfolioResult, err error) {

	balanceResp, errType, err := getHttpJsonResponse(ACCOUNT_BALANCE_URL+p.Account, "Cosmos Account Balance")
	if err != nil {
		if errType == NOT_FOUND {
			log.Printf("Not Found for account \"%s\":%v", p.Account, err)
			return nil, &notFoundErrorInstance
		} else {
			log.Printf("Internal Error for account \"%s\":%v", p.Account, err)
			return nil, &internalErrorInstance
		}
	}

	denomEntries, err := parseBalance(balanceResp)
	if err != nil {
		log.Printf("Internal Error for account \"%s\":%v", p.Account, err)
		return nil, &internalErrorInstance
	}

	//Set of missing symbols in the exchange
	missingSet := make(map[string]bool)

	//Set of denom ids to be queried
	ids := make([]string, 0)
	for _, denomEntry := range denomEntries {
		if denomEntry.Id == "" {
			missingSet[denomEntry.CosmosSymbol] = true
		} else {
			ids = append(ids, denomEntry.Id)
		}
	}

	missings := make([]string, len(missingSet))
	ind := 0
	for k, _ := range missingSet {
		missings[ind] = k
		ind++
	}

	var marketResponses []MarketResponse
	if len(ids) > 0 {
		reqUrl := createCoinGeckoURL(ids)
		marketResp, _, err := getHttpJsonResponse(reqUrl, "CoinGecko Market Price")
		if err != nil {
			log.Printf("Internal Error for account \"%s\":%v", p.Account, err)
			return nil, &internalErrorInstance
		}

		err = json.Unmarshal(marketResp, &marketResponses)
		if err != nil {
			log.Printf("Internal Error for account \"%s\":%v", p.Account, err)
			return nil, &internalErrorInstance
		}

	}

	total, port, err := calculatePortfolioCurrencySummary(denomEntries, marketResponses, missingSet)
	if err != nil {
		log.Printf("Internal Error for account \"%s\":%v", p.Account, err)
		return nil, &internalErrorInstance
	}

	res = &portfolio.PortfolioResult{AccountID: p.Account, Portfolio: port, Total: "$" + total.StringFixed(2), Missing: missings}
	s.logger.Print("portfolio.portfolioSummary")
	return
}

func calculatePortfolioCurrencySummary(denomEntries []DenomEntry, marketResponses []MarketResponse, missingSet map[string]bool) (decimal.Decimal, [](*portfolio.Currency), error) {
	type mapentry struct {
		DenomEntry *DenomEntry
		realAmount decimal.Decimal
		price      decimal.Decimal
	}
	if len(marketResponses) == 0 {
		return decimal0, make([]*portfolio.Currency, 0), nil
	}

	//keeps aggregate amount and prices for denomEntries
	denomMap := make(map[string]mapentry)

	//prices coming from the exchange
	prices := make(map[string]decimal.Decimal)

	for _, v := range marketResponses {
		prices[v.Symbol] = decimal.NewFromFloat(v.CurrentPrice)
	}

	for _, denomEntry := range denomEntries {
		denomEntry := denomEntry
		if missingSet[denomEntry.CosmosSymbol] {
			continue
		}
		//"ok" is necessary since there might be multiple ibc denoms pointing to the same currency
		entry, ok := denomMap[denomEntry.Symbol]
		amountDecimal, err := decimal.NewFromString(denomEntry.Amount)
		if err != nil {
			return decimal0, nil, err
		}
		curAmount := amountDecimal.Mul(exponents[denomEntry.Exp])

		if !ok {

			pricee := prices[denomEntry.Symbol]

			denomMap[denomEntry.Symbol] = mapentry{DenomEntry: &denomEntry, realAmount: curAmount, price: pricee}
		} else {
			entry.realAmount = entry.realAmount.Add(curAmount)
		}
	}

	portfolioCurrencies := make([](*portfolio.Currency), 0)
	total := decimal0
	for _, entry := range denomMap {
		entry := entry
		var element portfolio.Currency
		denomTotal := entry.realAmount.Mul(entry.price)
		element.Amount = entry.realAmount.String()
		element.Denom = entry.DenomEntry.Name + " (" + entry.DenomEntry.Symbol + ")"
		element.DenomTotal = "$" + denomTotal.StringFixed(2)
		element.Price = "$" + entry.price.String()
		total = total.Add(denomTotal)

		portfolioCurrencies = append(portfolioCurrencies, &element)
	}

	return total, portfolioCurrencies, nil
}

func createCoinGeckoURL(ids []string) string {
	params := url.Values{}
	reqIds := strings.Join(ids, ",")
	params.Add("vs_currency", "usd")
	params.Add("ids", reqIds)
	return CG_MARKET_URL + "?" + params.Encode()
}

func parseBalance(resp json.RawMessage) ([]DenomEntry, error) {
	type cosmosAccountBalanceEntry struct {
		Denom  string `json:"denom"`
		Amount string `json:"amount"`
	}

	type cosmosBalanceResponse struct {
		Balances []cosmosAccountBalanceEntry `json:"balances"`
		//do not care about pagination
	}
	out := make([]DenomEntry, 0)

	var balances cosmosBalanceResponse

	err := json.Unmarshal(resp, &balances)
	if err != nil {
		return nil, err
	}

	for _, v := range balances.Balances {
		var entry DenomEntry
		entry.Amount = v.Amount
		metadata := coin2IdMap[v.Denom]
		entry.MetadataEntry = &(metadata)

		out = append(out, entry)
	}
	return out, nil
}

func getHttpJsonResponse(url string, desc string) (json.RawMessage, ReturnType, error) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	res, err := client.Get(url)
	if err != nil {
		errMsg := fmt.Sprintf("%s failed with error:%v", desc, err)
		return nil, INTERNAL_ERROR, errors.New(errMsg)
	} else if res.StatusCode >= 500 {
		errMsg := fmt.Sprintf("%s returned %s", desc, res.Status)
		return nil, INTERNAL_ERROR, errors.New(errMsg)
	} else if res.StatusCode != 200 {
		errMsg := fmt.Sprintf("%s returned %s", desc, res.Status)
		return nil, NOT_FOUND, errors.New(errMsg)
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		errMsg := fmt.Sprintf("%s response body could not be read %v", desc, err)
		return nil, INTERNAL_ERROR, errors.New(errMsg)
	}

	return body, SUCCESS, nil
}
