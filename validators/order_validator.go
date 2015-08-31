package validators

import (
	"github.com/hashtock/hashtock-go/core"
)

func ValidateIncommingOrderSchema(order core.OrderBase) (err error) {
	hasPrice := order.UnitPrice != 0
	hasOrderRef := order.BaseOrderID != ""
	hasQuantity := order.Quantity != 0

	if order.HashTag == "" {
		return core.NewBadRequestError("Order must specify hash tag")
	}

	if hasQuantity == false {
		return core.NewBadRequestError("Order quantity has to be specified")
	}

	if order.Type == core.TYPE_BANK {
		if hasPrice {
			return core.NewBadRequestError("Bank order does not allow to specify unit price")
		}
		if hasOrderRef {
			return core.NewBadRequestError("Bank order does not allow to reference order id")
		}
	} else if order.Type == core.TYPE_MARKET {
		hasPriceAndOrderRef := hasOrderRef && hasPrice
		hasNoPriceAndNoOrderRef := (hasOrderRef || hasPrice) == false

		if hasPriceAndOrderRef || hasNoPriceAndNoOrderRef { // You can either fulfil the order or place new one
			return core.NewBadRequestError("Market order must have either unit price or reference order id")
		}

		if hasPrice && order.UnitPrice < 0 {
			return core.NewBadRequestError("Unit price has to be greater then 0")
		}
	} else {
		return core.NewBadRequestError("Order type not supported")
	}

	return nil
}

// ValidateIncommingOrderSanity has to check is order actually executable:
// - Tag exists
// - Reference order exists and yields the same quantity
func ValidateIncommingOrderSanity(bank core.BankStorage, orders core.OrderStorage, order core.OrderBase) (err error) {
	if _, err := bank.Tag(order.HashTag); err != nil {
		return core.NewBadRequestError("Can't place order for unrecognised tag")
	}

	if order.Type == core.TYPE_MARKET && order.BaseOrderID != "" {
		refOrder, err := orders.Order("", order.BaseOrderID)
		if err != nil {
			return core.NewBadRequestError("Base order can't be found")
		}
		if err := ValidateMarketOrdersCompatible(order, refOrder.OrderBase); err != nil {
			return err
		}
	}

	return nil
}

func ValidateMarketOrdersCompatible(order core.OrderBase, refOrder core.OrderBase) (err error) {
	if order.HashTag == refOrder.HashTag {
		return core.NewBadRequestError("Orders have different tags to trade")
	}
	if order.Quantity == -refOrder.Quantity {
		return core.NewBadRequestError("Orders have different quantities to trade")
	}
	return nil
}
