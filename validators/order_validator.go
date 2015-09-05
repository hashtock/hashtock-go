package validators

import (
	"github.com/hashtock/hashtock-go/core"
)

func ValidateIncommingOrderSchema(order core.OrderBase) (err error) {
	hasPrice := order.UnitPrice != 0
	priceCorrect := order.UnitPrice > 0
	hasOrderRef := order.BaseOrderID != ""
	hasQuantity := order.Quantity != 0

	if order.HashTag == "" {
		return core.NewBadRequestError("Order must specify hash tag")
	}

	if hasQuantity == false {
		return core.NewBadRequestError("Order quantity has to be specified")
	}

	switch order.Type {
	case core.TYPE_BANK:
		if hasPrice {
			return core.NewBadRequestError("Bank order does not allow to specify unit price")
		}
		if hasOrderRef {
			return core.NewBadRequestError("Bank order does not allow to reference market order")
		}
	case core.TYPE_MARKET:
		if priceCorrect == false {
			return core.NewBadRequestError("Unit price has to be greater then 0")
		}

		if hasOrderRef {
			return core.NewBadRequestError("Marker order can't specify other order")
		}
	case core.TYPE_MARKET_FULFIL:
		if priceCorrect == false {
			return core.NewBadRequestError("Unit price has to be greater then 0")
		}

		if hasOrderRef == false {
			return core.NewBadRequestError("Order to fulfil has to be specified")
		}
	default:
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

	if order.Type == core.TYPE_MARKET_FULFIL {
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
	if order.Type != core.TYPE_MARKET_FULFIL {
		return core.NewBadRequestError("Only fulfil order is allowed to fulfil other market order")
	}
	if refOrder.Type != core.TYPE_MARKET {
		return core.NewBadRequestError("Only market order can be fulfilled")
	}
	if order.HashTag != refOrder.HashTag {
		return core.NewBadRequestError("Orders have different tags to trade")
	}
	if order.Quantity == -refOrder.Quantity {
		return core.NewBadRequestError("Orders have different quantities to trade")
	}
	if order.UnitPrice != refOrder.UnitPrice {
		return core.NewBadRequestError("Unit price when fulfilling order have to match")
	}
	return nil
}
