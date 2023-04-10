package page_loader

import (
	"io"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	addr        = "127.0.0.1:8000"
	baseUrl     = "http://" + addr
	headerKey   = "X-Custom-Header"
	headerValue = "42"
)

func TestHttpPageLoader(t *testing.T) {
	stopServer := make(chan any)
	go runDemoServer(stopServer)
	defer func() {
		stopServer <- nil
	}()

	// wait for server startup
	time.Sleep(10 * time.Microsecond)
	t.Run("Loads page content with correct request header", func(t *testing.T) {
		header := http.Header{}
		header.Add(headerKey, headerValue)
		loader := NewHttpPageLoader(header)

		u, _ := url.JoinPath(baseUrl, "/ok")
		res, err := loader.LoadPage(u)

		assert.Nil(t, err)
		assert.Equal(t, "hey", res)
	})

	t.Run("Returns error if status code is not OK", func(t *testing.T) {
		loader := NewHttpPageLoader(nil)

		u, _ := url.JoinPath(baseUrl, "/not-ok")
		_, err := loader.LoadPage(u)

		assert.Error(t, err)
	})

	t.Run("Returns error if URL is invalid", func(t *testing.T) {
		loader := NewHttpPageLoader(nil)

		_, err := loader.LoadPage("invalid url")

		assert.ErrorContains(t, err, "unsupported protocol scheme")
	})
}

func runDemoServer(stopServer chan any) {
	mux := http.NewServeMux()

	mux.HandleFunc("/ok", func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get(headerKey)
		if header != headerValue {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "hey")
	})

	mux.HandleFunc("/mnot-ok", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	server := http.Server{
		Addr:    addr,
		Handler: mux,
	}

	stopSignalReceived := false
	go func() {
		<-stopServer
		stopSignalReceived = true
		server.Close()
	}()

	if err := server.ListenAndServe(); err != nil && !stopSignalReceived {
		panic(err)
	}
}
