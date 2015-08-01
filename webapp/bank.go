package webapp

import (
	"net/http"

	"github.com/hashtock/service-tools/serialize"

	"github.com/hashtock/hashtock-go/core"
)

type bankService struct {
	storage    core.BankStorage
	serializer serialize.Serializer
}

func (b *bankService) TagInfo(rw http.ResponseWriter, req *http.Request) {
	hashTagName := req.URL.Query().Get(":tag")

	tag, err := b.storage.Tag(hashTagName)
	if err != nil {
		b.serializer.JSON(rw, http.StatusInternalServerError, err.Error())
		return
	}

	b.serializer.JSON(rw, http.StatusOK, tag)
}

func (b *bankService) ListOfAllHashTags(rw http.ResponseWriter, req *http.Request) {
	tags, err := b.storage.Tags()
	if err != nil {
		b.serializer.JSON(rw, http.StatusInternalServerError, err.Error())
		return
	}

	b.serializer.JSON(rw, http.StatusOK, tags)
}
