package verb

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Transaction struct {
	Request  http.Request
	Response http.Response
}

func (t *Transaction) UseBody(body []byte) {
	if body == nil {
		return
	}
	t.Request.Body = ioutil.NopCloser(bytes.NewReader(body))
	t.Request.GetBody = func() (io.ReadCloser, error) {
		return ioutil.NopCloser(bytes.NewReader(body)), nil
	}
	t.Request.ContentLength = int64(len(body))
}

func (t *Transaction) UseURL(u *url.URL) {
	t.Request.Host = u.Host
	t.Request.URL = u
}

func (t *Transaction) Do(c Client) error {
	response, err := c.Do(&t.Request)
	if err != nil {
		return fmt.Errorf("verb: could not send request; %w", err)
	}

	t.Response = *response

	return nil
}

type Op func(*Transaction) error

func (t *Transaction) Try(ops ...Op) error {
	for _, fn := range ops {
		if err := fn(t); err != nil {
			return err
		}
	}
	return nil
}

func Tx(method string) *Transaction {
	var t Transaction
	t.Request = http.Request{
		Header:     make(http.Header),
		Method:     method,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
	}
	return &t
}
