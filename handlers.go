package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
	"io/ioutil"
	"net/http"
	"strconv"
)

// Return default message for root routing
func Index(ctx *fasthttp.RequestCtx) {
	fmt.Fprint(ctx, "Hello")
}

// Handle echo with Hystrix
func (h *MyHandler) echoHandlerHystrix(ctx *fasthttp.RequestCtx) {

	url := h.calledServiceURL + "/echo/" + fmt.Sprintf("%s", ctx.UserValue("message"))

	headers := http.Header{}
	//headers.Set("Content-Type", "application/json")
	headers.Set("Connection", "close")

	response, errResp := h.client.Get(url, headers)
	if errResp != nil {
		log.Error().Msgf("Error en response: %s", errResp.Error())
		ctx.Error(errResp.Error(), fasthttp.StatusInternalServerError)
		ctx.Response.Header.Set("Status", strconv.Itoa(fasthttp.StatusInternalServerError))
		fmt.Fprintf(ctx, "Error en response: %s", errResp.Error())
		return
	}

	defer response.Body.Close()

	respBody, errData := ioutil.ReadAll(response.Body)
	if errData != nil {
		log.Error().Msgf("failed to read response body %s", errData.Error())
	}

	fmt.Fprintf(ctx, "%s", respBody)
}


// Handle iterative path and calls iterative calculation service
func (h *MyHandler) factorialIterativeHandler(ctx *fasthttp.RequestCtx) {

	params := mux.Vars(r)

	url := h.calledServiceURL + "/factorialIterative/" + params["number"]

	req, reqErr := http.NewRequest("GET", url, nil)
	if reqErr != nil {
		log.Error().Msgf("Error en request: %s", reqErr)
		fmt.Fprintf(w, "Error en request: %s", reqErr)
		return
	}

	req.Header.Set("Connection", "close")
	//req.Header.Set("Connection", "Keep-Alive")

	resp, errResp := h.client.Do(req)
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

// Handle recursive path and calls recursive calculation service
func (h *MyHandler) factorialRecursiveHandler(ctx *fasthttp.RequestCtx) {

	params := mux.Vars(r)

	url := h.calledServiceURL + "/factorialRecursive/" + params["number"]

	req, reqErr := http.NewRequest("GET", url, nil)
	if reqErr != nil {
		log.Fatal().Msgf("Error en response: %s", reqErr)
		fmt.Fprintf(w, "Error en response: %s", reqErr)
		return
	}

	req.Header.Set("Connection", "close")
	//req.Header.Set("Connection", "Keep-Alive")

	resp, errResp := h.client.Do(req)
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
