package main

import (
	"fmt"
	"github.com/gojektech/heimdall/hystrix"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"html"
	"io/ioutil"
	"net/http"
)

// Return default message for root routing
func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}

// Handle echo with Hystrix
func echoHandlerHystrix(calledServiceURL string, client *hystrix.Client) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)

		url := calledServiceURL + "/echo/" + params["message"]

		headers := http.Header{}
		//headers.Set("Content-Type", "application/json")
		headers.Set("Connection", "close")

		response, errResp := client.Get(url, headers)
		if errResp != nil {
			log.Error().Msgf("Error en response: %s")
			fmt.Fprintf(w, "Error en response: %s", errResp)
			return
		}

		defer response.Body.Close()

		respBody, errData := ioutil.ReadAll(response.Body)
		if errData != nil {
			log.Error().Msgf("failed to read response body %s", errData)
		}

		fmt.Fprintf(w, "%s", respBody)

	}
}

// Handle iterative path and calls iterative calculation service
func echoHandler(calledServiceURL string, client *hystrix.Client) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)

		url := calledServiceURL + "/echo/" + params["message"]

		req, reqErr := http.NewRequest("GET", url, nil)
		if reqErr != nil {
			log.Error().Msgf("Error en response: %s", reqErr)
			w.Header().Set("Status", "504")
			fmt.Fprint(w, "504 - Gateway timeout")
			return
		}

		req.Header.Set("Connection", "close")
		//req.Header.Set("Connection", "Keep-Alive")

		//resp, err := netClient.Do(req)
		resp, errResp := client.Do(req)
		if errResp != nil {
			log.Error().Msgf("Error en response: %s", errResp)
			fmt.Fprintf(w, "Error en response: %s", errResp)
			return
		}

		respData, errResp := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if errResp != nil {
			log.Error().Msgf("Error en RespData: %s", errResp)
			fmt.Fprintf(w, "Error en respData: %s", errResp)
			return
		}

		log.Info().Msgf("Status %s\n", resp.Header.Get("Status"))
		log.Info().Msgf("Ok - Resp: %v\n", respData)
		fmt.Fprintf(w, "%s", respData)

	}
}

// Handle iterative path and calls iterative calculation service
func factorialIterativeHandler(calledServiceURL string, netClient *http.Client) func(w http.ResponseWriter,
	r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)

		url := calledServiceURL + "/factorialIterative/" + params["number"]

		req, reqErr := http.NewRequest("GET", url, nil)
		if reqErr != nil {
			log.Error().Msgf("Error en request: %s", reqErr)
			fmt.Fprintf(w, "Error en request: %s", reqErr)
			return
		}

		req.Header.Set("Connection", "close")
		//req.Header.Set("Connection", "Keep-Alive")

		resp, errResp := netClient.Do(req)
		if errResp != nil {
			log.Error().Msgf("Error en response: %s", errResp)
			fmt.Fprintf(w, "Error en response: %s", errResp)
			return
		}

		respData, errRespData := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if errRespData != nil {
			log.Error().Msgf("Error en respData: %s", errRespData)
			fmt.Fprintf(w, "Error en respData: %s", errRespData)
			return
		}

		fmt.Fprintf(w, "%s", respData)

	}
}

// Handle recursive path and calls recursive calculation service
func factorialRecursiveHandler(calledServiceURL string, netClient *http.Client) func(w http.ResponseWriter,
	r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)

		url := calledServiceURL + "/factorialRecursive/" + params["number"]

		req, reqErr := http.NewRequest("GET", url, nil)
		if reqErr != nil {
			log.Fatal().Msgf("Error en response: %s", reqErr)
			fmt.Fprintf(w, "Error en response: %s", reqErr)
			return
		}

		req.Header.Set("Connection", "close")
		//req.Header.Set("Connection", "Keep-Alive")

		resp, errResp := netClient.Do(req)
		if errResp != nil {
			log.Error().Msgf("%s", errResp)
			fmt.Fprintf(w, "%s", errResp)
			return
		}

		respData, errData := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if errData != nil {
			log.Error().Msgf("%s", errData)
			fmt.Fprintf(w, "Error en respData: %s", errData)
			return
		}

		fmt.Fprintf(w, "%s", respData)

	}
}
