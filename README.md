### Overview
Simple currency REST API. Parse https://nationalbank.kz/rss/rates_all.xml?switch=kazakh and save to db. History saves when refreshing too <br/>

### Stacks
1. Go
1. PostgreSQL
1. Cron
1. Item 3a
### Clone
```sh
git clone github.com/gospodinzerkalo/currency_api

```

### Install & Run
```sh
docker-compose build 
docker-compose up
```
### Usage (example)

```dotenv
POSTGRES_HOST: db
POSTGRES_DATABASE: postgres
POSTGRES_PASSWORD: postgres
POSTGRES_USER: postgres
```

 You can change cycle of job scheduler, I have put "@every 10s" for testing. For daily refresh pus "@daily"
```go
_, err = cr.AddFunc("@every 10s", func() {
    if err := currencyService.RefreshCurrencies(); err != nil {
    log.Fatal(err)
}
})
```

### API
### Get all currency
```http request
GET http://localhost:8080/currency
```
response 
```json
[
  {
    "currency_pair": "TRY/KZT",
    "value": "57.38"
  },
  {
    "currency_pair": "UZS/KZT",
    "value": "4.04"
  },
    ...
]
```
### add parameters for cross: base & quoted
```http request
GET http://localhost:8080/currency?base=USD&quoted=RUB
```
response
```json
[
    {
        "currency_pair": "USD/RUB",
        "value": "75.22"
    }
]
```

### Convert currency pair
```http request
POST http://localhost:8080/convert
```

Request body 
```json
{
    "convert_from": "USD",
    "convert_to": "RUB",
    "value" : 1
}
```

Response
```json
{
  "convert_from": "USD",
  "convert_to": "RUB",
  "convert_from_value": 1,
  "convert_to_value": "75.2193"
}
```
### Get history of change 
```http request
GET http://localhost:8080/history
```
Response 
```json
[
    {
        "id": 1,
        "title": "AUD",
        "pub_date": "28.01.2021",
        "change": "+1.44"
    },  ...
]
```
