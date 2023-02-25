package tests

import (
	"bytes"
	"net/http"
	"net/url"

	urlpkg "net/url"
)

type NewRequestOptions struct {
	Body       string
	Headers    map[string]string
	GetValues  url.Values
	PostValues *url.Values
}

// NewRequest creates a new Request for testing, but adds RequestURI
// as the default imlemented in GoLang does not add the RequestURI
// and leaves it to the end user to implement
func NewRequest(method string, url string, opts NewRequestOptions) (*http.Request, error) {
	if method == "POST" && opts.PostValues != nil {
		if opts.Headers == nil {
			opts.Headers = map[string]string{}
		}
		opts.Body = opts.PostValues.Encode()
		opts.Headers["Content-Type"] = "application/x-www-form-urlencoded"
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer([]byte(opts.Body)))
	if err != nil {
		return nil, err
	}

	u, err := urlpkg.Parse(url)
	if err != nil {
		return nil, err
	}

	req.RequestURI = u.Path

	u.RawQuery = opts.GetValues.Encode()
	req.URL.RawQuery = u.RawQuery

	for key, value := range opts.Headers {
		req.Header.Set(key, value)
	}

	return req, nil
}
