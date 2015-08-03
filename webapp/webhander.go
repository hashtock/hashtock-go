package webapp

import (
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/pat"
	authClient "github.com/hashtock/auth/client"
	authCore "github.com/hashtock/auth/core"
	"github.com/hashtock/service-tools/serialize"

	"github.com/hashtock/hashtock-go/core"
)

const UserContextKey = "user"

type Options struct {
	PortfolioStorage core.PortfolioStorage
	BankStorage      core.BankStorage
	OrderStorage     core.OrderStorage
	Serializer       serialize.Serializer
	WhoClient        authCore.Who
}

func Handlers(options Options) http.Handler {
	n := negroni.New(
		negroni.NewRecovery(),
		negroni.NewLogger(),
		authClient.NewAuthMiddleware(options.WhoClient),
	)

	ps := portfolioService{options.PortfolioStorage, options.Serializer}
	bs := bankService{options.BankStorage, options.Serializer}
	os := orderService{
		storage:    options.OrderStorage,
		bank:       options.BankStorage,
		serializer: options.Serializer,
	}

	m := pat.New()
	m.Get("/portfolio/{tag}/", ps.PortfolioTagInfo).Name("Portfolio:TagInfo")
	m.Get("/portfolio/", ps.Portfolio).Name("Portfolio:All")
	m.Get("/balance/", ps.PortfolioBalance).Name("Portfolio:Balance")

	m.Get("/bank/{tag}/", bs.TagInfo).Name("Bank:TagInfo")
	m.Get("/bank/", bs.ListOfAllHashTags).Name("Bank:Tags")

	m.Get("/order/history/", os.CompletedOrders).Name("Order:CompletedOrders")
	m.Get("/order/{uuid}/", os.OrderDetails).Name("Order:OrderDetails")
	m.Get("/order/", os.ActiveOrders).Name("Order:Orders")
	m.Delete("/order/{uuid}/", os.CancelOrder).Name("Order:CancelOrder")
	m.Post("/order/", os.NewOrder).Name("Order:NewOrder")

	n.UseHandler(m)

	return n
}
