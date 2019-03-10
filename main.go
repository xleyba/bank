package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
	"github.com/spf13/viper"
)


// Return default message for root routing
func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}

// Handle iterative path and calls iterative calculation service
func echoHandler(calledServiceURL string, netClient *http.Client) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)

		url := calledServiceURL + "/echo/" + params["message"]

		req, reqErr := http.NewRequest("GET", url, nil)
		if reqErr != nil {
			log.Fatal("Error en response: ", reqErr)
			fmt.Fprintf(w, "Error en response: %s", reqErr)
			return
		}

		req.Header.Set("Connection", "close")
		//req.Header.Set("Connection", "Keep-Alive")

		resp, err := netClient.Do(req)
		if err != nil {
			log.Fatal("Error en response: ", err)
			fmt.Fprintf(w, "Error en response: %s", err)
			return
		}

		respData, errResp := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if errResp != nil {
			log.Fatal("Error en RespData", errResp)
			fmt.Fprintf(w, "Error en respData: %s", err)
			return
		}

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
			log.Fatal("Error en request: ", reqErr)
			fmt.Fprintf(w, "Error en request: %s", reqErr)
			return
		}

		req.Header.Set("Connection", "close")
		//req.Header.Set("Connection", "Keep-Alive")

		resp, errResp := netClient.Do(req)
		if errResp != nil {
			log.Fatal("Error en response: ", errResp)
			fmt.Fprintf(w, "Error en response: %s", errResp)
			return
		}

		respData, errRespData := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if errRespData != nil {
			log.Fatal(errRespData)
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
			log.Fatal("Error en response: ", reqErr)
			fmt.Fprintf(w, "Error en response: %s", reqErr)
			return
		}

		req.Header.Set("Connection", "close")
		//req.Header.Set("Connection", "Keep-Alive")

		resp, errResp := netClient.Do(req)
		if errResp != nil {
			log.Fatal(errResp)
			fmt.Fprintf(w, "%s", errResp)
			return
		}

		respData, errData := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if errData != nil {
			log.Fatal(errData)
			fmt.Fprintf(w, "Error en respData: %s", errData)
			return
		}

		fmt.Fprintf(w, "%s", respData)

	}
}

// Main function
func main() {

	log.Println("Starting server")

	// Start to read conf file
	log.Println("\n\n")
	log.Println("=============================================")
	log.Println("      Configuration checking - bank v0.8")
	log.Println("=============================================")

	// loading configuration
	viper.SetConfigName("conf")                                   // name of config file (without ext)
	viper.AddConfigPath(".")                                      // default path for conf file
	viper.SetDefault("port", ":9296")                             // default port value
	viper.SetDefault("calledServiceURL", "http://localhost:9596") // default calledServiceURL
	err := viper.ReadInConfig()                                   // Find and read the config file
	if err != nil {                                               // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	log.Println("-- Using port: ", viper.GetString("port"))
	log.Println("-- CalledServiceURL set to: ", viper.GetString("calledServiceURL"))

	log.Println("=============================================")

	// Define just one transport for all calls
	tr := &http.Transport{
		MaxIdleConns:        10000, // iddle conn pool size
		MaxIdleConnsPerHost: 10000, // iddle conn pool size per host
		DisableKeepAlives:   true,  // disable keep alive
	}

	netClient := &http.Client{Transport: tr}

	router := mux.NewRouter()

	router.HandleFunc("/", Index).Methods("GET")
	router.HandleFunc("/echo/{message}", echoHandler(viper.GetString("calledServiceURL"),
		netClient)).Methods("GET")
	router.HandleFunc("/factorialIterative/{number}",
		factorialIterativeHandler(viper.GetString("calledServiceURL"), netClient)).Methods("GET")
	router.HandleFunc("/factorialRecursive/{number}",
		factorialRecursiveHandler(viper.GetString("calledServiceURL"), netClient)).Methods("GET")

	// set timeout
	//muxWithMiddlewares := http.TimeoutHandler(router, time.Second*3, "Timeout!")

	srv := &http.Server{
		Addr: viper.GetString("port"),
		// https://stackoverflow.com/questions/29334407/creating-an-idle-timeout-in-go
		//WriteTimeout: time.Second * 60,
		ReadTimeout: time.Second * 15,
		//IdleTimeout:  time.Second * 120,
		//Handler: router, // Pass our instance of gorilla/mux in.
		Handler: router,
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		log.Println("Running server....")

		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	fmt.Println("\n\n")
	log.Println("shutting down")
	log.Println("Goddbye!....")
	os.Exit(0)

}
