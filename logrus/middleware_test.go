package logrus_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Akagi201/light"
	mlogrus "github.com/Akagi201/middleware/logrus"
	"github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func newRequest(method, url string) *http.Request {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		panic(err)
	}
	return req
}

func TestHandler(t *testing.T) {
	var buf bytes.Buffer

	app := light.New()

	logger := logrus.New()
	logger.Level = logrus.InfoLevel
	logger.Formatter = &logrus.JSONFormatter{}
	logger.Out = &buf
	app.Use(mlogrus.NewLogger(logger, "web"))

	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(200)
		fmt.Fprint(w, "Hello World\n")
	})

	app.Get("/", handler)

	app.ServeHTTP(httptest.NewRecorder(), newRequest("GET", "/"))

	assert.True(t, buf.Len() > 0, "buffer should not be empty")
	assert.True(t, strings.Contains(buf.String(), `"method":"GET"`), "method wrong")
	assert.True(t, strings.Contains(buf.String(), `"status":200`), "status not 200")
}
