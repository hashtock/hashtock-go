# hashtock-go

## Install
Assuming you have Go installed and set up
```
# Install mongodb
# Clone hashtock-go
git clone git@github.com:hashtock/hashtock-go.git $GOPATH/src/github.com hashtock/hashtock-go

# Install requirements
go get .
```

## Serve
```
cd $GOPATH/src/github.com/hashtock/hashtock-go
go build
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

**Market**:
- Accepts buy offers for a given hashtag+price+amount
- Accepts sell offers for a given hashtag+price+amount
- Allows to view any request by ID - owner only
- Allows to cancel any request if not fulfilled
- History of market operations

**Client**:
- Account balance
- Current portfolio of hash tags

**Admin**:
- Add new tag to bank

## API

URI prefix: `/api/`

| URI             | Method | Description                           | Done? |
|-----------------|--------|---------------------------------------|-------|
| /               | GET    | Main entry points to resouces         |  [x]  |
| /order/         | GET    | List of current orders                |  [x]  |
| /order/         | POST   | Add new order                         |  [x]  |
| /order/{uuid}/  | GET    | Order details                         |  [x]  |
| /order/{uuid}/  | DELETE | Cancel the order                      |  [x]  |
| /order/history/ | GET    | List of all orders                    |  [x]  |
| /tag/           | GET    | List of all tags with bank values     |  [x]  |
| /tag/           | POST   | Add new tag (admin)                   |  [x]  |
| /tag/{hash}/    | GET    | Details about the hash tag            |  [x]  |
| /tag/{hash}/    | PUT    | Update tag Value (admin)              |  [x]  |
| /user/          | GET    | High level user details               |  [x]  |
| /user/tags/     | GET    | List of users shares of tags          |  [x]  |
