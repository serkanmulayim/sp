## SP Cosmos Hub Account Balance Viewer

### Add Dependency

```
go get github.com/shopspring/decimal
```
**Note:** [Goa](https://goa.design/) should be installed.

### Build Server
Prebuilt binary **sp** is built with MacOS. It can be rebuilt with the following command. 
```
go build ./cmd/sp
```

### Coinhelper
Coinhelper behaves as a database keeping coin metadatas. It keeps cosmos to coingecko token mappings. Final file is coinmetadata.json. It is precomputed in the repo. If desired, it can be run to refresh the coin metadata (e.g. if new coins are released in Cosmos, or if they are recognized by CoinGecko). 
```
cd coinhelper
go build coinhelper.go
./coinhelper
cd ..
```

### Run Server
```
./sp
```

### API
```
http://localhost:8080/sp/portfolio/{cosmosAccountID}
```
