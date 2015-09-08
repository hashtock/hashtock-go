package core

type PortfolioStorage interface {
	PortfolioShares(userId string) ([]TagShare, error)
	PortfolioShare(userId string, tag string, strict bool) (*TagShare, error)
	PortfolioBalance(userId string) (Balance, error)
}

type BankStorage interface {
	Tags() ([]HashTag, error)
	Tag(tag string) (*HashTag, error)
	TagSetValue(tag string, value float64) error
}

type OrderStorage interface {
	Orders(filter OrderFilter) ([]Order, error)
	Order(userId string, orderId string) (*Order, error)
	AddOrder(order *Order) error
	DeleteOrder(userId string, orderId string) error
}

type OrderExecuter interface {
	Order(userId string, orderId string) (*Order, error)
	OrdersToExecute() ([]Order, error)
	OrderCompleted(orderId string, status OrderResolution, notes string) error
}
