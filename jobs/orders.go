package jobs

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/hashtock/hashtock-go/core"
	"github.com/hashtock/hashtock-go/validators"
)

type OrderWorker struct {
	storage   core.OrderExecuter
	bank      core.BankStorage
	portfolio core.PortfolioStorage

	interval time.Duration
	ticker   *time.Ticker
}

var (
	ErrNotEnoughFunds      = errors.New("Not enough founds")
	ErrNotEnoughUserShares = errors.New("Not enough shares in users possession")
)

func NewOrderWorker(storage core.OrderExecuter, bank core.BankStorage, portfolio core.PortfolioStorage, interval time.Duration) *OrderWorker {
	return &OrderWorker{
		storage:   storage,
		bank:      bank,
		portfolio: portfolio,
		interval:  interval,
	}
}

func (o *OrderWorker) processOrders() {
	activeOrders, err := o.storage.OrdersToExecute()
	if err != nil {
		log.Println("OrderWorker: Could not fetch active bank orders. Err:", err)
		return
	}

	for _, order := range activeOrders {
		var err error

		switch order.Type {
		case core.TYPE_BANK:
			err = o.executeBankOrder(order)
		case core.TYPE_MARKET_FULFIL:
			err = o.executeMarketOrder(order)
		default:
			err = errors.New("Don't know how to execute this order type")
		}

		if err != nil {
			log.Printf("OrderWorker: Could not execute '%v' order %v. Err: %v", order.Type, order.UUID, err)
		}
	}

	if len(activeOrders) > 0 {
		log.Printf("OrderWorker: %v orders executed", len(activeOrders))
	} else {
		log.Println("OrderWorker: No orders to execute")
	}
}

func (o *OrderWorker) executeBankOrder(order core.Order) (err error) {
	var (
		profileBalance core.Balance
		hashTag        *core.HashTag
		tagShare       *core.TagShare
	)

	// It's time to blow up if asked to execute non bank order here
	if order.Type != core.TYPE_BANK {
		msg := fmt.Sprintf("execution of non bank order: %v", order.Type)
		return core.NewBadRequestError(msg)
	}

	if hashTag, err = o.bank.Tag(order.HashTag); err != nil {
		o.storage.OrderCompleted(order.UUID, core.ERROR, "")
		return
	}

	if profileBalance, err = o.portfolio.PortfolioBalance(order.UserID); err != nil {
		o.storage.OrderCompleted(order.UUID, core.ERROR, "")
		return
	}

	if tagShare, err = o.portfolio.PortfolioShare(order.UserID, order.HashTag, false); err != nil {
		o.storage.OrderCompleted(order.UUID, core.ERROR, "")
		return
	}

	// Buy
	if order.Quantity > 0.0 {
		if profileBalance.Cash < order.Value {
			o.storage.OrderCompleted(order.UUID, core.FAILURE, ErrNotEnoughFunds.Error())
			msg := fmt.Sprintf("User %v does not have enough founds to complete %v", order.UserID, order)
			return core.NewBadRequestError(msg)
		}

		if hashTag.InBank < order.Quantity {
			o.storage.OrderCompleted(order.UUID, core.FAILURE, "Not enough shares in bank")
			msg := fmt.Sprintf("Bank does not have enough shares to complete %v", order)
			return core.NewBadRequestError(msg)
		}
	}

	// Sell
	if order.Quantity < 0.0 {
		if tagShare.Quantity < -order.Quantity {
			o.storage.OrderCompleted(order.UUID, core.FAILURE, ErrNotEnoughUserShares.Error())
			msg := fmt.Sprintf("User %v does not have enough shares (%v) to complete %v - %#v", order.UserID, tagShare.Quantity, order.UUID, order.OrderBase)
			return core.NewBadRequestError(msg)
		}
	}

	err = o.storage.OrderCompleted(order.UUID, core.SUCCESS, "")

	return
}

func (o *OrderWorker) executeMarketOrder(order core.Order) (err error) {
	// It's time to blow up if asked to execute non market order here
	if order.Type != core.TYPE_MARKET_FULFIL {
		msg := fmt.Sprintf("execution of non market fulfil order: %v", order.Type)
		return core.NewBadRequestError(msg)
	}

	orderOwnerId := ""
	orderToFulfil, err := o.storage.Order(orderOwnerId, order.BaseOrderID)
	if err != nil {
		msg := "Order to fulfil does not exist any more"
		o.storage.OrderCompleted(order.UUID, core.FAILURE, msg)
		return core.NewBadRequestError(msg)
	}

	if orderToFulfil.Complete {
		o.storage.OrderCompleted(order.UUID, core.FAILURE, "Order no longer on the market, it's already fulfilled.")
		return
	}

	if err = validators.ValidateMarketOrdersCompatible(order.OrderBase, orderToFulfil.OrderBase); err != nil {
		o.storage.OrderCompleted(order.UUID, core.ERROR, "")
		return
	}

	// Buy (this order, someone on market is selling)
	if order.Quantity > 0.0 {
		ok := false
		ok, err = validators.UserCanAffordToSpend(o.portfolio, order.UserID, order.Value)
		if err != nil {
			o.storage.OrderCompleted(order.UUID, core.ERROR, "")
			return
		} else if !ok {
			o.storage.OrderCompleted(order.UUID, core.FAILURE, ErrNotEnoughFunds.Error())
			return
		}
	}

	// Sell (this order, someone on market is buying)
	if order.Quantity < 0.0 {
		ok := false
		ok, err = validators.UserHaveEnoughToSellTag(o.portfolio, order.UserID, order.HashTag, order.Quantity)
		if err != nil {
			o.storage.OrderCompleted(order.UUID, core.ERROR, "")
			return
		} else if !ok {
			o.storage.OrderCompleted(order.UUID, core.FAILURE, ErrNotEnoughUserShares.Error())
			return
		}
	}

	err = o.storage.OrderCompleted(order.UUID, core.SUCCESS, "")
	err = o.storage.OrderCompleted(orderToFulfil.UUID, core.SUCCESS, "")
	return
}

func (o *OrderWorker) Start() (err error) {
	if o.ticker != nil {
		return
	}

	o.processOrders()
	o.ticker = time.NewTicker(o.interval)

	go func() {
		for range o.ticker.C {
			o.processOrders()
		}
	}()

	return
}

func (o *OrderWorker) Stop() {
	if o.ticker == nil {
		return
	}

	o.ticker.Stop()
	o.ticker = nil
}
