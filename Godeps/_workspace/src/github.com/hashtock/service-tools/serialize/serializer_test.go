package serialize_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hashtock/service-tools/serialize"
)

type brokenMarshalJSON struct{}

func (b brokenMarshalJSON) MarshalJSON() ([]byte, error) {
	return []byte("data"), errors.New("MarshalJSON error")
}

func TestImplementsSerializer(t *testing.T) {
	serializer := serialize.WebAPISerializer{}
	assert.Implements(t, (*serialize.Serializer)(nil), serializer)
}

func TestWebAPISerializerJSONSetsContentType(t *testing.T) {
	serializer := serialize.WebAPISerializer{}
	w := httptest.NewRecorder()

	serializer.JSON(w, http.StatusOK, nil)

	assert.EqualValues(t, []string{"application/json; charset=utf-8"}, w.HeaderMap["Content-Type"])
}

func TestWebAPISerializerJSONSetsStatus(t *testing.T) {
	statuses := []int{
		http.StatusOK,
		http.StatusNotFound,
		http.StatusInternalServerError,
		99999,
		0,
	}
	serializer := serialize.WebAPISerializer{}

	for _, status := range statuses {
		w := httptest.NewRecorder()

		serializer.JSON(w, status, nil)

		assert.EqualValues(t, status, w.Code)
	}
}

func TestWebAPISerializerJSONBody(t *testing.T) {
	kv := struct {
		Key string `json:"map-key"`
	}{
		Key: "value",
	}

	objs := map[interface{}]string{
		"string": `"string"`,
		1:        `1`,
		kv:       `{"map-key":"value"}`,
		nil:      "",
	}
	serializer := serialize.WebAPISerializer{}

	for obj, objStr := range objs {
		w := httptest.NewRecorder()

		serializer.JSON(w, http.StatusOK, obj)

		assert.EqualValues(t, objStr, w.Body.String())
	}
}

func TestWebAPISerializerJSONErrorMessage(t *testing.T) {
	err := errors.New("Error msg")
	serializer := serialize.WebAPISerializer{}
	w := httptest.NewRecorder()

	serializer.JSON(w, 123, err)

	assert.EqualValues(t, 123, w.Code)
	assert.EqualValues(t, `"Error msg"`, w.Body.String())
}

func TestWebAPISerializerJSONNilAndNoContent(t *testing.T) {
	serializer := serialize.WebAPISerializer{}

	type sample struct {
		code int
		obj  interface{}
	}

	matrix := []sample{
		sample{http.StatusOK, nil},
		sample{http.StatusNoContent, nil},
		sample{http.StatusNoContent, 123},
	}

	for _, testSample := range matrix {
		w := httptest.NewRecorder()

		serializer.JSON(w, testSample.code, testSample.obj)

		assert.EqualValues(t, testSample.code, w.Code)
		assert.EqualValues(t, "", w.Body.String())
	}
}

func TestWebAPISerializerJSONMarshalError(t *testing.T) {
	serializer := serialize.WebAPISerializer{}
	w := httptest.NewRecorder()

	obj := brokenMarshalJSON{}
	serializer.JSON(w, http.StatusOK, obj)

	assert.EqualValues(t, http.StatusInternalServerError, w.Code)
	assert.EqualValues(t, "Could not serialize object to JSON\n", w.Body.String())
}
