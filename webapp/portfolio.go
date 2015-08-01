package webapp

import (
	"log"
	"net/http"

	"github.com/hashtock/service-tools/serialize"

	"github.com/hashtock/hashtock-go/core"
)

type portfolioService struct {
	storage    core.PortfolioStorage
	serializer serialize.Serializer
}

func (p *portfolioService) Portfolio(rw http.ResponseWriter, req *http.Request) {
	id, err := userId(req)
	if err != nil {
		status, err := core.ErrToErrorer(err)
		p.serializer.JSON(rw, status, err)
		return
	}

	tags, err := p.storage.PortfolioShares(id)
	if err != nil {
		status, err := core.ErrToErrorer(err)
		p.serializer.JSON(rw, status, err)
		return
	}

	p.serializer.JSON(rw, http.StatusOK, tags)
}

func (p *portfolioService) PortfolioTagInfo(rw http.ResponseWriter, req *http.Request) {
	id, err := userId(req)
	if err != nil {
		status, err := core.ErrToErrorer(err)
		p.serializer.JSON(rw, status, err)
		return
	}

	hashTagName := req.URL.Query().Get(":tag")

	tag, err := p.storage.PortfolioShare(id, hashTagName, true)
	if err != nil {
		status, err := core.ErrToErrorer(err)
		p.serializer.JSON(rw, status, err)
		return
	}

	p.serializer.JSON(rw, http.StatusOK, tag)
}

func (p *portfolioService) PortfolioBalance(rw http.ResponseWriter, req *http.Request) {
	id, err := userId(req)
	if err != nil {
		status, err := core.ErrToErrorer(err)
		p.serializer.JSON(rw, status, err)
		return
	}

	balance, err := p.storage.PortfolioBalance(id)
	if err != nil {
		status, err := core.ErrToErrorer(err)
		log.Panicln("Err:", err.Error())
		p.serializer.JSON(rw, status, err)
		return
	}

	p.serializer.JSON(rw, http.StatusOK, balance)
}
