# SP Cosmos Hub Account Balance Viewer

## add dependency

```
go get github.com/shopspring/decimal
```
**Note:** [Goa](https://goa.design/) should be installed.

### build server
Prebuilt binary **sp** is built with MacOS. It can be rebuilt with the following. 
```
go build ./cmd/sp
```

### coinhelper
Coinhelper behaves as a database keeping coin metadata. It keeps cosmos to coingecko token mappings. Final file is cgcoins.json. It is precomputed in the repo. If desired, it can be run to refresh the coin metadata (e.g. if new coins are released). 
```
cd coinhelper
go build coinhelper.go
./coinhelper
cd ..
```

### run server
```
./sp
```

### API
```
http://localhost:8080/sp/portfolio/{cosmosAccountID}
```
