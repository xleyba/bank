package main

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
	"io/ioutil"
	"strconv"
)

// Return default message for root routing
func Index(ctx *fasthttp.RequestCtx) {
	fmt.Fprint(ctx, "Hello")
}

// Handle echo
func (h *MyHandler) echoHandlerHystrix(ctx *fasthttp.RequestCtx) {

	url := h.calledServiceURL + "/echo/" + fmt.Sprintf("%s", ctx.UserValue("message"))
	log.Debug().Msgf("URL to call: %s", url)


	response, errResp := h.client.Get(url)
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

	ctx.Response.SetStatusCode(fasthttp.StatusOK)
	fmt.Fprintf(ctx, "%s", respBody)
}


// Handle iterative path and calls iterative calculation service
func (h *MyHandler) factorialIterativeHandler(ctx *fasthttp.RequestCtx) {

	url := h.calledServiceURL + "/factorialIterative/" + fmt.Sprintf("%s", ctx.UserValue("number"))
	log.Debug().Msgf("URL to call: %s", url)

	response, errResp := h.client.Get(url)
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

	ctx.Response.SetStatusCode(fasthttp.StatusOK)

	fmt.Fprintf(ctx, "%s", respBody)

}

// Handle recursive path and calls recursive calculation service
func (h *MyHandler) factorialRecursiveHandler(ctx *fasthttp.RequestCtx) {

	url := h.calledServiceURL + "/factorialIterative/" + fmt.Sprintf("%s", ctx.UserValue("number"))

	response, errResp := h.client.Get(url)
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

	ctx.Response.SetStatusCode(fasthttp.StatusOK)
	fmt.Fprintf(ctx, "%s", respBody)

}
