package webapp

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/hashtock/service-tools/serialize"
	"github.com/pborman/uuid"

	"github.com/hashtock/hashtock-go/core"
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
	o.listOrders(false, rw, req)
}

func (o *orderService) CompletedOrders(rw http.ResponseWriter, req *http.Request) {
	o.listOrders(true, rw, req)
}

func (o *orderService) listOrders(completed bool, rw http.ResponseWriter, req *http.Request) {
	id, err := userId(req)
	if err != nil {
		o.serializer.JSON(rw, http.StatusInternalServerError, err.Error())
		return
	}

	// Additional filters
	queryValues := req.URL.Query()
	tag := queryValues.Get("tag")
	resolution := queryValues.Get("resolution")

	orders, err := o.storage.Orders(id, completed, tag, resolution)
	if err != nil {
		o.serializer.JSON(rw, http.StatusInternalServerError, err.Error())
		return
	}

	o.serializer.JSON(rw, http.StatusOK, orders)
}

func (o *orderService) NewOrder(rw http.ResponseWriter, req *http.Request) {
	id, err := userId(req)
	if err != nil {
		o.serializer.JSON(rw, http.StatusInternalServerError, err.Error())
		return
	}

	baseOrder := core.OrderBase{}
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&baseOrder); err != nil {
		status, err := core.ErrToErrorer(err)
		o.serializer.JSON(rw, status, err)
		return
	}

	if err := o.validateOrder(baseOrder); err != nil {
		status, err := core.ErrToErrorer(err)
		o.serializer.JSON(rw, status, err)
		return
	}

	systemOrder, err := o.newOrderSystem(id, baseOrder.HashTag, baseOrder.Quantity)
	if err != nil {
		status, err := core.ErrToErrorer(err)
		o.serializer.JSON(rw, status, err)
		return
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

func (o *orderService) newOrderSystem(userId string, tag string, quantity float64) (order core.OrderSystem, err error) {
	var hashTag *core.HashTag

	if hashTag, err = o.bank.Tag(tag); err != nil {
		return
	}

	order = core.OrderSystem{
		UUID:       uuid.New(),
		UserID:     userId,
		Complete:   false,
		CreatedAt:  time.Now(),
		Resolution: core.PENDING,
		Value:      quantity * hashTag.Value * -1.0,
	}

	return
}

func (o *orderService) validateOrder(order core.OrderBase) (err error) {
	if order.Type != core.TYPE_BANK {
		return core.NewBadRequestError("Order type not supported")
	}

	return
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
