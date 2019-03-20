package main

import (
	"github.com/gojektech/heimdall"
	"github.com/gojektech/heimdall/hystrix"
	"time"
)

func SetClient() *hystrix.Client {

	// First set a backoff mechanism. Constant backoff increases the backoff at a constant rate
	backoffInterval := 5 * time.Millisecond
	// Define a maximum jitter interval. It must be more than 1*time.Millisecond
	maximumJitterInterval := 15 * time.Millisecond

	backoff := heimdall.NewConstantBackoff(backoffInterval, maximumJitterInterval)

	// Create a new retry mechanism with the backoff
	retrier := heimdall.NewRetrier(backoff)

	client := hystrix.NewClient(
		hystrix.WithHTTPTimeout(120*time.Second),
		hystrix.WithHystrixTimeout(time.Second*120),
		hystrix.WithMaxConcurrentRequests(600),
		hystrix.WithErrorPercentThreshold(30),
		hystrix.WithSleepWindow(10),
		hystrix.WithRetrier(retrier),
		hystrix.WithRetryCount(4),
		hystrix.WithCommandName("MyClient"),
	)

	return client

}