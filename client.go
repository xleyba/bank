package main

import (
	"fmt"
	"github.com/spf13/viper"
	"net/http"
)

func SetClient() *http.Client {

	defaultRoundTripper := http.DefaultTransport
	defaultTransportPointer, ok := defaultRoundTripper.(*http.Transport)
	if !ok {
		panic(fmt.Sprintf("defaultRoundTripper not an *http.Transport"))
	}
	defaultTransport := *defaultTransportPointer // dereference it to get a copy of the struct that the pointer points to
	defaultTransport.MaxIdleConns = viper.GetInt("MaxIdleConns")
	defaultTransport.MaxIdleConnsPerHost = viper.GetInt("MaxIdleConnsPerHost")

	myClient := &http.Client{Transport: &defaultTransport}

	return myClient

}
