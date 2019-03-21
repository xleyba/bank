package main

import (
	"fmt"
	"net/http"
)

func SetClient() *http.Client {

	defaultRoundTripper := http.DefaultTransport
	defaultTransportPointer, ok := defaultRoundTripper.(*http.Transport)
	if !ok {
		panic(fmt.Sprintf("defaultRoundTripper not an *http.Transport"))
	}
	defaultTransport := *defaultTransportPointer // dereference it to get a copy of the struct that the pointer points to
	defaultTransport.MaxIdleConns = 510
	defaultTransport.MaxIdleConnsPerHost = 510

	myClient := &http.Client{Transport: &defaultTransport}

	return myClient

}