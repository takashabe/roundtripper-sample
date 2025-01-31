package main

import (
	"log"
	"net/http"
)

var _ http.RoundTripper = &SimpleRoundTripper{}

type SimpleRoundTripper struct{}

func (rt *SimpleRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	log.Println("Sending request")
	return http.DefaultTransport.RoundTrip(req)
}
