package jobs

import (
	"log"

	"github.com/nats-io/nats"

	"github.com/hashtock/hashtock-go/core"
)

type TagValueWorker struct {
	storage core.BankStorage
	subject string
	natsUrl string
	conn    *nats.EncodedConn
	sub     *nats.Subscription
}

func NewTagValueWorker(setter core.BankStorage, natsUrl string, msgSubject string) *TagValueWorker {
	return &TagValueWorker{
		storage: setter,
		subject: msgSubject,
		natsUrl: natsUrl,
	}
}

func (t TagValueWorker) processMsg(counts map[string]int) {
	for name, count := range counts {
		if err := t.storage.TagSetValue(name, float64(count)); err != nil {
			log.Printf("Error while updating value for tag %v to %v. %v", name, count, err.Error())
		}
	}
}

func (t *TagValueWorker) Start() (err error) {
	if t.conn != nil || t.sub != nil {
		return
	}

	if err = t.connect(); err != nil {
		return
	}

	t.sub, err = t.conn.Subscribe(t.subject, t.processMsg)
	return
}

func (t *TagValueWorker) Stop() {
	if t.conn != nil {
		t.conn.Close()
		t.conn = nil
	}

	if t.sub != nil && t.sub.IsValid() {
		t.sub.Unsubscribe()
		t.sub = nil
	}
}

func (t *TagValueWorker) connect() error {
	if t.conn != nil {
		return nil
	}

	nc, err := nats.Connect(t.natsUrl)
	if err != nil {
		return err
	}
	conn, err := nats.NewEncodedConn(nc, "json")
	if err != nil {
		return err
	}

	t.conn = conn
	return nil
}
