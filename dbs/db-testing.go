package dbs

import (
	"net/http"
)

type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func NewHttpClient(responses map[string]*http.Response) *http.Client {

	return &http.Client{
		Transport: RoundTripFunc(func(req *http.Request) *http.Response {

			urlRequest := req.URL.String()

			if _, ok := responses[urlRequest]; ok {
				return responses[urlRequest]
			}

			return &http.Response{
				StatusCode: http.StatusNotFound,
			}
		}),
	}
}
