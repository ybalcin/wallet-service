# Wallet Service

A REST API to provide wallet use cases

## Getting Started

If this is your first time encountering Go, please follow [the instructions](https://go.dev/doc/install) to
install Go on your computer.

## Usage
```shell
git clone https://github.com/ybalcin/wallet-service.git

cd wallet-service

# Runs with test db
make run

# Run tests
make run-test
```

After `make run` API server is running on `http://127.0.0.1:8080`

### Create Wallet

#### Request

`POST /api/wallets/`

    curl -i -H 'Accept: application/json' -d 'username=ybalcin' http://127.0.0.1:8080/api/wallets/

#### Response
    
    HTTP/1.1 200 OK
    Date: Thu, 27 Jul 2023 16:12:24 GMT
    Content-Type: application/json
    Content-Length: 45

    {"id":"7eadc3e1-c0d6-4653-b5eb-6b25d76d3446"}

### Deposit Money

#### Request

`PUT /api/wallets/7eadc3e1-c0d6-4653-b5eb-6b25d76d3446/deposit`

    curl -i -H 'Accept: application/json' -d 'amount=10' -X PUT http://127.0.0.1:8080/api/wallets/7eadc3e1-c0d6-4653-b5eb-6b25d76d3446/deposit

#### Response
    
    HTTP/1.1 200 OK
    Date: Thu, 27 Jul 2023 16:18:16 GMT
    Content-Type: application/json
    Content-Length: 90

    {"id":"7eadc3e1-c0d6-4653-b5eb-6b25d76d3446","username":"ybalcin","balance":{"amount":10}}

### Withdraw Money

#### Request

`PUT /api/wallets/7eadc3e1-c0d6-4653-b5eb-6b25d76d3446/withdraw`

    curl -i -H 'Accept: application/json' -d 'amount=10' -X PUT http://127.0.0.1:8080/api/wallets/7eadc3e1-c0d6-4653-b5eb-6b25d76d3446/withdraw

#### Response

    HTTP/1.1 200 OK
    Date: Thu, 27 Jul 2023 16:20:16 GMT
    Content-Type: application/json
    Content-Length: 89

    {"id":"7eadc3e1-c0d6-4653-b5eb-6b25d76d3446","username":"ybalcin","balance":{"amount":0}}

### Get Wallet

#### Request

`GET /api/wallets/7eadc3e1-c0d6-4653-b5eb-6b25d76d3446`

    curl -i -H 'Accept: application/json' http://127.0.0.1:8080/api/wallets/7eadc3e1-c0d6-4653-b5eb-6b25d76d3446

#### Response

    HTTP/1.1 200 OK
    Date: Thu, 27 Jul 2023 16:24:10 GMT
    Content-Type: application/json
    Content-Length: 89

    {"id":"7eadc3e1-c0d6-4653-b5eb-6b25d76d3446","username":"ybalcin","balance":{"amount":0}}