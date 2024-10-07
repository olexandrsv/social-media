package transport

import (
	"bytes"
	"io"
	"net/http"
)

type request struct {
	method string
	url    string
	header string
	form   map[string]string
	cookie map[string]string
	body   io.Reader
}

func NewRequest(method, url string) *request {
	return &request{
		method: method,
		url:    url,
		body:   new(bytes.Buffer),
	}
}

func (r *request) WithHeader(header string) *request {
	r.header = header
	return r
}

func (r *request) WithForm(form map[string]string) *request {
	r.form = form
	return r
}

func (r *request) WithCookie(cookie map[string]string) *request {
	r.cookie = cookie
	return r
}

func (r *request) Execute() (response string, code int, err error) {
	if r.form != nil {
		r.body, r.header, err = createForm(r.form)
		if err != nil {
			return
		}
	}
	req, err := http.NewRequest(r.method, r.url, r.body)
	if err != nil {
		return
	}
	if r.header != "" {
		req.Header.Add("Content-Type", r.header)
	}

	for k, v := range r.cookie {
		req.AddCookie(&http.Cookie{
			Name:  k,
			Value: v,
		})
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	return string(bytes), resp.StatusCode, nil
}
