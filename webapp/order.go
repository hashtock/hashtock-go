package webapp

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/hashtock/service-tools/serialize"
	"github.com/pborman/uuid"

	"github.com/hashtock/hashtock-go/core"
	"github.com/hashtock/hashtock-go/validators"
)

type orderService struct {
	storage    core.OrderStorage
	bank       core.BankStorage
	serializer serialize.Serializer
}

func (o *orderService) OrderDetails(rw http.ResponseWriter, req *http.Request) {
	id, err := userId(req)
	if err != nil {
		o.serializer.JSON(rw, http.StatusInternalServerError, err.Error())
		return
	}

	orderId := req.URL.Query().Get(":uuid")

	order, err := o.storage.Order(id, orderId)
	if err != nil {
		o.serializer.JSON(rw, http.StatusInternalServerError, err.Error())
		return
	}

	o.serializer.JSON(rw, http.StatusOK, order)
}

func (o *orderService) ActiveOrders(rw http.ResponseWriter, req *http.Request) {
	filters, err := orderFilterFromRequest(req)
	if err != nil {
		o.serializer.JSON(rw, http.StatusInternalServerError, err.Error())
		return
	}
	filters.Complete = false

	// To list market orders of everyone
	if filters.Type == core.TYPE_MARKET {
		filters.UserID = ""
	}

	o.listOrders(filters, rw, req)
}

func (o *orderService) CompletedOrders(rw http.ResponseWriter, req *http.Request) {
	filters, err := orderFilterFromRequest(req)
	if err != nil {
		o.serializer.JSON(rw, http.StatusInternalServerError, err.Error())
		return
	}
	filters.Complete = true

	o.listOrders(filters, rw, req)
}

func (o *orderService) listOrders(filters core.OrderFilter, rw http.ResponseWriter, req *http.Request) {
	orders, err := o.storage.Orders(filters)
	if err != nil {
		o.serializer.JSON(rw, http.StatusInternalServerError, err.Error())
		return
	}

	o.serializer.JSON(rw, http.StatusOK, orders)
}

func (o *orderService) NewOrder(rw http.ResponseWriter, req *http.Request) {
	baseOrder := core.OrderBase{}
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&baseOrder); err != nil {
		status, err := core.ErrToErrorer(err)
		o.serializer.JSON(rw, status, err)
		return
	}

	if baseOrder.Type != core.TYPE_BANK && baseOrder.Type != core.TYPE_MARKET {
		err := core.NewBadRequestError("Only bank and initial market orders are allowed where")
		o.serializer.JSON(rw, err.ErrCode(), err.Error())
		return
	}

	o.createNewOrder(baseOrder, rw, req)
}

func (o *orderService) FulfilOrder(rw http.ResponseWriter, req *http.Request) {
	id, err := userId(req)
	if err != nil {
		status, err := core.ErrToErrorer(err)
		o.serializer.JSON(rw, status, err)
		return
	}

	orderId := req.URL.Query().Get(":uuid")

	orderOwnerId := ""
	orderToFulfil, err := o.storage.Order(orderOwnerId, orderId)
	if err != nil {
		status, err := core.ErrToErrorer(err)
		o.serializer.JSON(rw, status, err)
		return
	}

	if orderToFulfil.Type != core.TYPE_MARKET {
		err = core.NewBadRequestError("Only market orders can be fulfilled")
		status, err := core.ErrToErrorer(err)
		o.serializer.JSON(rw, status, err)
		return
	}

	if id == orderToFulfil.UserID {
		err = core.NewBadRequestError("It does not make sense to fulfil your own order")
		status, err := core.ErrToErrorer(err)
		o.serializer.JSON(rw, status, err)
		return
	}

	// Full copy and then modify to make it a valid fulfil order
	baseOrder := orderToFulfil.OrderBase
	baseOrder.Type = core.TYPE_MARKET_FULFIL
	baseOrder.BaseOrderID = orderToFulfil.UUID
	baseOrder.Quantity *= -1.0

	o.createNewOrder(baseOrder, rw, req)
}

func (o *orderService) createNewOrder(baseOrder core.OrderBase, rw http.ResponseWriter, req *http.Request) {
	if err := validators.ValidateIncommingOrderSchema(baseOrder); err != nil {
		status, err := core.ErrToErrorer(err)
		o.serializer.JSON(rw, status, err)
		return
	}

	if err := validators.ValidateIncommingOrderSanity(o.bank, o.storage, baseOrder); err != nil {
		status, err := core.ErrToErrorer(err)
		o.serializer.JSON(rw, status, err)
		return
	}

	uid, err := userId(req)
	if err != nil {
		o.serializer.JSON(rw, http.StatusInternalServerError, err.Error())
		return
	}

	unitPrice := 0.0
	switch baseOrder.Type {
	case core.TYPE_BANK:
		hashTag, err := o.bank.Tag(baseOrder.HashTag)
		if err != nil {
			status, err := core.ErrToErrorer(err)
			o.serializer.JSON(rw, status, err)
			return
		}
		unitPrice = hashTag.Value
	case core.TYPE_MARKET:
		unitPrice = baseOrder.UnitPrice
	case core.TYPE_MARKET_FULFIL:
		unitPrice = baseOrder.UnitPrice
	}

	systemOrder := core.OrderSystem{
		UUID:       uuid.New(),
		UserID:     uid,
		Complete:   false,
		CreatedAt:  time.Now(),
		Resolution: core.PENDING,
		Value:      baseOrder.Quantity * unitPrice * -1.0,
	}

	order := &core.Order{
		OrderBase:   baseOrder,
		OrderSystem: systemOrder,
	}

	if err := o.storage.AddOrder(order); err != nil {
		status, err := core.ErrToErrorer(err)
		o.serializer.JSON(rw, status, err)
		return
	}

	o.serializer.JSON(rw, http.StatusCreated, order)
}

func (o *orderService) CancelOrder(rw http.ResponseWriter, req *http.Request) {
	id, err := userId(req)
	if err != nil {
		status, err := core.ErrToErrorer(err)
		o.serializer.JSON(rw, status, err)
		return
	}

	orderId := req.URL.Query().Get(":uuid")

	order, err := o.storage.Order(id, orderId)
	if err != nil {
		status, err := core.ErrToErrorer(err)
		o.serializer.JSON(rw, status, err)
		return
	}

	if order.Complete {
		err = core.NewBadRequestError("Order can not be cancelled any more")
		status, err := core.ErrToErrorer(err)
		o.serializer.JSON(rw, status, err)
		return
	}

	if err = o.storage.DeleteOrder(id, orderId); err != nil {
		status, err := core.ErrToErrorer(err)
		o.serializer.JSON(rw, status, err)
	}

	rw.WriteHeader(http.StatusNoContent)
}
