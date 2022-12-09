package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
)

// Creates a json file mapping IBC coin hashes to coins
const (
	SUPPLY_URL                         = "https://api.cosmos.network/cosmos/bank/v1beta1/supply?pagination.limit=1000"
	DENOM_TRACES_URL                   = "https://api.cosmos.network/ibc/apps/transfer/v1/denom_traces/"
	CG_COIN_LIST_URL                   = "https://api.coingecko.com/api/v3/coins/list"
	MAX_CONC_HTTP_REQ                  = 50
	COSMOS_HASH2SYMBOL_FILE            = "coins.json"
	CG_COIN2ID_MAP_FILE                = "cgcoins.json"
	COSMOSHASH2COINGECKO_METADATA_FILE = "coinmetadata.json"
)

type cgSymbolEntry struct {
	Id     string `json:"id"`
	Symbol string `json:"symbol"`
	Name   string `json:"name"`
}

type MetadataEntry struct {
	*cgSymbolEntry
	Exp          int    `json:"exp"`
	CosmosSymbol string `json:"cosmos_symbol"`
}

func main() {
	// convertHashToCosmosSymbols()
	// getCGCoinList()
	createCosmosHash2CGMetadata()
}

// gets coin list from CoinGecko
func getCGCoinList() {
	res, err := http.Get(CG_COIN_LIST_URL)
	if err != nil {
		log.Fatalf("Request failed for coin lists! %s could not be generated. Exiting..., %v", CG_COIN2ID_MAP_FILE, err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("Request failed for coin lists! %s could not be generated. Exiting...%v", CG_COIN2ID_MAP_FILE, err)
	}
	js := json.RawMessage(body)
	jsonstr, _ := json.MarshalIndent(js, "", " ")
	ioutil.WriteFile(CG_COIN2ID_MAP_FILE, jsonstr, 0644)
}

func convertHashToCosmosSymbols() {
	supplies := getAllSupplies()

	//for each entry (hash or coin name) map to coin names
	coinMap := map2Coins(supplies)
	//log.Print(coinMap)

	//write to file
	jsonstr, _ := json.MarshalIndent(coinMap, "", " ")

	ioutil.WriteFile(COSMOS_HASH2SYMBOL_FILE, jsonstr, 0644)
}

func createCosmosHash2CGMetadata() {
	var coinIdEntries []cgSymbolEntry
	coin2IdMap := make(map[string]cgSymbolEntry)
	jsonn, err := ioutil.ReadFile(CG_COIN2ID_MAP_FILE)
	if err != nil {
		log.Fatalf("%s could not be found.Uncomment getCGCoinList() in main(), and rerun this app.", CG_COIN2ID_MAP_FILE)
	}

	err = json.Unmarshal(jsonn, &coinIdEntries)
	if err != nil {
		log.Fatalf("%s is corrupted. Uncomment getCGCoinList() in main(), and rerun this app.", CG_COIN2ID_MAP_FILE)
	}

	for _, v := range coinIdEntries {
		coin2IdMap[v.Symbol] = v
	}

	jsonn, err = ioutil.ReadFile(COSMOS_HASH2SYMBOL_FILE)
	var cosmos_hash2symbol_map map[string]string
	if err != nil {
		log.Fatalf("%s could not be found. Uncomment convertHashToCosmosSymbols() in main(), and rerun this app.", COSMOS_HASH2SYMBOL_FILE)
	} else {
		err = json.Unmarshal(jsonn, &cosmos_hash2symbol_map)
		if err != nil {
			log.Fatalf("%s is corrupted. Uncomment convertHashToCosmosSymbols() in main(), and rerun this app.", COSMOS_HASH2SYMBOL_FILE)
		}
	}

	metadata_map := make(map[string]MetadataEntry)
	for hash, val := range cosmos_hash2symbol_map {
		sym, exp := coin2Symbol(val)
		cg, ok := coin2IdMap[sym]
		var entry MetadataEntry
		if !ok {
			log.Printf("%s (converted from %s) does not have a match in CoinGecko symbols. Skipping...", sym, val)
			entry = MetadataEntry{&cg, exp, val}
		} else {
			entry = MetadataEntry{&cg, exp, val}

		}
		metadata_map[hash] = entry
	}

	jsonstr, _ := json.MarshalIndent(metadata_map, "", " ")
	ioutil.WriteFile(COSMOSHASH2COINGECKO_METADATA_FILE, jsonstr, 0644)

}

// Gets all tokens (either with IBC hash or name)
// Exits if error
func getAllSupplies() []string {
	// get supply - returns coin hashes
	//curl -X GET "https://api.cosmos.network/cosmos/bank/v1beta1/supply?pagination.limit=1000" -H "accept: application/json"
	//
	// RESP -> {supply:
	// [
	//   {
	// 	    "denom": "ibc/00255B18FBBC1E36845AAFDCB4CBD460DC45331496A64C2A29CEAFDD3B997B5F",
	// 	    "amount": "1016160059457055549180"
	//   },
	//   {
	//   	"denom": "poolFD005C5AB01714A4B62E87F5213F5D5CDE357773D70712916A93664BCE5A6931",
	// 		"amount": "4821542"
	//   },
	//   {
	// 		"denom": "uatom",
	// 		"amount": "318992623240499"
	//   },
	//  ...
	// ]}

	supplyMap, err := getJsonResponse(SUPPLY_URL, "Supply Request")
	if err != nil {
		log.Fatal(err)
	}
	arr := supplyMap["supply"].([]interface{})
	out := make([]string, len(arr))
	for i, v := range arr {
		out[i] = v.(map[string]interface{})["denom"].(string)
	}
	return out
}

// maps the coin hashes to their denoms if they start with "ibc/" prefix
func map2Coins(supplies []string) map[string]string {
	// curl -X GET "https://api.cosmos.network/ibc/apps/transfer/v1/denom_traces/14F9BC3E44B8A9C1BE1FB08980FAB87034C9905EF17CF2F5008FC085218811CC"

	var lockOut sync.Mutex
	out := make(map[string]string)

	sem := make(chan int, MAX_CONC_HTTP_REQ)
	var wg sync.WaitGroup

	for _, str := range supplies {
		wg.Add(1)
		sem <- 1
		go func(s string) {
			ind := strings.Index(s, "ibc/")
			if ind > -1 {
				desc := fmt.Sprintf("Map2Coins Request for \"%s\"", s)
				m, err := getJsonResponse(DENOM_TRACES_URL+s[4:], desc)
				if err != nil {
					log.Print(err)
				} else {
					val := m["denom_trace"].(map[string]interface{})["base_denom"].(string)
					lockOut.Lock()
					out[s] = val
					lockOut.Unlock()
				}

			} else {
				lockOut.Lock()
				out[s] = s
				lockOut.Unlock()
			}

			<-sem
			wg.Done()
		}(str)

	}
	wg.Wait()
	return out
}

func getJsonResponse(url string, desc string) (map[string]interface{}, error) {
	res, err := http.Get(url)
	if err != nil {
		errMsg := fmt.Sprintf("%s failed with error:%v", desc, err)
		return nil, errors.New(errMsg)
	} else if res.StatusCode != 200 {
		errMsg := fmt.Sprintf("%s returned %s", desc, res.Status)
		return nil, errors.New(errMsg)
	}
	var out map[string]interface{}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		errMsg := fmt.Sprintf("%s response body could not be read %v", desc, err)
		return nil, errors.New(errMsg)
	}
	err = json.Unmarshal(body, &out)
	if err != nil {
		errMsg := fmt.Sprintf("%s response could not be mapped %v", desc, err)
		return nil, errors.New(errMsg)
	}
	return out, nil

}

func coin2Symbol(s string) (string, int) {
	//Best effort to resolve symbol for the
	if strings.Index(s, "usd") == 0 {
		return s, 0
	} else if strings.Index(s, "u") == 0 {
		return s[1:], 6
	} else if strings.Index(s, "nano") == 0 {
		return s[4:], 9
	} else {
		return s, 0
	}
}
