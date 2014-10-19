# hashtock-go

## Functionality

**Bank**:
- Buy hash tag
- Sell hash tag
- Current bank value of a given hash tag, and how much does it have for sell
- Currnet bank value of all hash tags it knows about, and how much does it have for sell
- History of bank operations (**ToDo**)

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
| /order/         | GET    | List of current orders                |  [ ]  |
| /order/         | POST   | Add new order                         |  [ ]  |
| /order/{uuid}/  | GET    | Order details                         |  [ ]  |
| /order/{uuid}/  | DELETE | Cancel the order                      |  [ ]  |
| /order/history/ | GET    | List of all orders                    |  [ ]  |
| /tag/           | GET    | List of all tags with bank values     |  [x]  |
| /tag/{hash}/    | GET    | Details about the hash tag            |  [x]  |
| /tag/{hash}/    | POST   | Update info about the tag (admin)     |  [ ]  |
| /user/          | GET    | High level user details               |  [x]  |
| /user/tags/     | GET    | List of users shares of tags          |  [x]  |
