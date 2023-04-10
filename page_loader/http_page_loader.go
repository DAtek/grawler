package page_loader

import (
	"bytes"
	"errors"
	"io"
	"net/http"

	"github.com/DAtek/gotils"
)

type httpPageLoader struct {
	header http.Header
	client *http.Client
}

func NewHttpPageLoader(header http.Header) IPageLoader {
	return &httpPageLoader{
		header: header,
		client: &http.Client{},
	}
}

func (loader *httpPageLoader) LoadPage(url string) (string, error) {
	req := gotils.ResultOrPanic(http.NewRequest("GET", url, &bytes.Buffer{}))

	if loader.header != nil {
		req.Header = loader.header
	}

	resp, respErr := loader.client.Do(req)
	if respErr != nil {
		return "", respErr
	}

	if resp.StatusCode != http.StatusOK {
		return "", errors.New(resp.Status)
	}

	buf := &bytes.Buffer{}
	io.Copy(buf, resp.Body)
	return buf.String(), nil
}
