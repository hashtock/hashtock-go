# hashtock-go

## Install
Assuming you have Go installed and set up
```
# Install mongodb
# Clone hashtock-go
git clone git@github.com:hashtock/hashtock-go.git $GOPATH/src/github.com hashtock/hashtock-go

# Install requirements
go get github.com/tools/godep
```

## Serve
```
cd $GOPATH/src/github.com/hashtock/hashtock-go
godep go build
./hashtock-go
```

## Run tests

```
cd $GOPATH/src/github.com/hashtock/hashtock-go
./run_tests.sh
```

## Functionality

**Bank**:
- Buy hash tag
- Sell hash tag
- Current bank value of a given hash tag, and how much does it have for sell
- Currnet bank value of all hash tags it knows about, and how much does it have for sell
- History of bank operations

**Market** (ToDo):
- Accepts buy offers for a given hashtag+price+amount
- Accepts sell offers for a given hashtag+price+amount
- Allows to view any request by ID - owner only
- Allows to cancel any request if not fulfilled
- History of market operations

**Client**:
- Account balance
- Current portfolio of hash tags

## API

URI prefix: `/api/`

| URI               | Method | Description                           |
|-------------------|--------|---------------------------------------|
| /portfolio        | GET    | List of owned tags                    |
| /portfolio/{hash} | GET    | Detailt about owned tag               |
| /balance          | GET    | User's cache balance                  |
| /bank/            | GET    | List of all tags with bank values     |
| /bank/{hash}/     | GET    | Details about the hash tag            |
| /order/           | GET    | List of current orders                |
| /order/           | POST   | Add new order                         |
| /order/{uuid}/    | GET    | Order details                         |
| /order/{uuid}/    | DELETE | Cancel the order                      |
| /order/history/   | GET    | List of all orders                    |
