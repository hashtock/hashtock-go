package jobs

import (
	"fmt"
	"log"
	"time"

	"github.com/hashtock/hashtock-go/core"
)

type OrderWorker struct {
	storage   core.OrderExecuter
	bank      core.BankStorage
	portfolio core.PortfolioStorage

	interval time.Duration
	ticker   *time.Ticker
}

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

	//TODO(error): Handle errors somehow
	for _, order := range activeOrders {
		if err := o.executeBankOrder(order); err != nil {
			log.Printf("OrderWorker: Could not execute bank order %v. Err: %v", order.UUID, err)
		}
	}

	if len(activeOrders) > 0 {
		log.Printf("OrderWorker: %v bank orders executed", len(activeOrders))
	} else {
		log.Println("OrderWorker: No bank orders to execute")
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
		log.Println("execution of non bank order:", order.Type)
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
			o.storage.OrderCompleted(order.UUID, core.FAILURE, "Not enough founds")
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
			o.storage.OrderCompleted(order.UUID, core.FAILURE, "Not enough shares in users possession")
			msg := fmt.Sprintf("User %v does not have enough shares (%v) to complete %v - %#v", order.UserID, tagShare.Quantity, order.UUID, order.OrderBase)
			return core.NewBadRequestError(msg)
		}
	}

	err = o.storage.OrderCompleted(order.UUID, core.SUCCESS, "")

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
