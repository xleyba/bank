package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
	"flag"
	"os/signal"
	"context"
)

// Configuration
type Configuration struct {
	Port             string
	CalledServiceURL string
}


// Return default message for root routing
func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}

// Handle iterative path and calls iterative calculation service
func echoHandler(calledServiceURL string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)

		url := calledServiceURL + "/echo/" + params["message"]

		tr := &http.Transport{
			//MaxIdleConns:       500,
			//MaxIdleConnsPerHost:  500,
		}

		netClient := &http.Client{Transport: tr}

		req, reqErr := http.NewRequest("GET", url, nil)
		if reqErr != nil {
			log.Fatal("Error en response: ", reqErr)
			fmt.Fprintf(w, "Error en response: %s", reqErr)
			return
		}

		//req.Header.Set("Connection", "close")
		//req.Header.Set("Connection", "Keep-Alive")


		//resp, err := netClient.Get(url)
		resp, err := netClient.Do(req)
		if err != nil {
			log.Fatal("Error en response: ", err)
			fmt.Fprintf(w, "Error en response: %s", err)
			return
		}

		respData, errResp := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		if errResp != nil {
			log.Fatal("Error en RespData", errResp)
			fmt.Fprintf(w, "Error en respData: %s", err)
			return
		}

		fmt.Fprintf(w, "%s", respData)

	}
}

// Handle iterative path and calls iterative calculation service
func factorialIterativeHandler(calledServiceURL string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)

		url := calledServiceURL + "/factorialIterative/" + params["number"]

		tr := &http.Transport{
			//MaxIdleConns:       500,
			//MaxIdleConnsPerHost:  500,
		}

		netClient := &http.Client{Transport: tr}

		req, reqErr := http.NewRequest("GET", url, nil)
		if reqErr != nil {
			log.Fatal("Error en request: ", reqErr)
			fmt.Fprintf(w, "Error en request: %s", reqErr)
			return
		}

		req.Header.Set("Connection", "close")

		//resp, err := http.Get(url)
		resp, errResp := netClient.Do(req)
		if errResp != nil {
			log.Fatal("Error en response: ", errResp)
			fmt.Fprintf(w, "Error en response: %s", errResp)
			return
		}

		respData, errRespData := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		if errRespData != nil {
			log.Fatal(errRespData)
			fmt.Fprintf(w, "Error en respData: %s", errRespData)
			return
		}

		fmt.Fprintf(w, "%s", respData)

	}
}

// Handle recursive path and calls recursive calculation service
func factorialRecursiveHandler(calledServiceURL string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)

		url := calledServiceURL + "/factorialRecursive/" + params["number"]

		tr := &http.Transport{
			//MaxIdleConns:       500,
			//MaxIdleConnsPerHost:  500,
		}

		netClient := &http.Client{Transport: tr}

		req, reqErr := http.NewRequest("GET", url, nil)
		if reqErr != nil {
			log.Fatal("Error en response: ", reqErr)
			fmt.Fprintf(w, "Error en response: %s", reqErr)
			return
		}

		req.Header.Set("Connection", "close")

		//resp, err := http.Get(url)
		resp, err := netClient.Do(req)
		if err != nil {
			log.Fatal(err)
			fmt.Fprintf(w, "%s", err)
			return
		}

		if err != nil {
			log.Fatal(err)
		}

		respData, errResp := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		if errResp != nil {
			log.Fatal(errResp)
			fmt.Fprintf(w, "Error en respData: %s", err)
			return
		}

		fmt.Fprintf(w, "%s", respData)

	}
}

// Main function
func main() {

	// Used for server graceful shutdown
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second * 15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	// Set default values
	port := ":9296"
	calledServiceURL := "http://localhost:9596"

	log.Println("Starting server")

	// Start to read conf file
	log.Println("\n\n")
	log.Println("=============================================")
	log.Println("      Configuration checking - bank v0.7")
	log.Println("=============================================")
	file, err := os.Open("conf.json")

	if err != nil {
		log.Println("-- No conf file, using port: ", port)
		log.Println("-- No conf file, using serive URL: ", calledServiceURL)
	} else {
		defer file.Close()
		decoder := json.NewDecoder(file)
		configuration := Configuration{}
		err := decoder.Decode(&configuration)

		if err != nil {
			fmt.Println("error:", err)
			log.Fatal()
		} else {

			// Check port parameter
			if len(configuration.Port) == 0 {
				log.Println("-- No port inf config file, using: ", port)
			} else {
				port = configuration.Port
				log.Println("-- Using port: ", port)
			}

			// Check service url to call
			if len(configuration.CalledServiceURL) == 0 {
				log.Println("-- No calledServiceURL present in conf, will be set to: ", calledServiceURL)
			} else {
				calledServiceURL := configuration.CalledServiceURL
				log.Println("-- CalledServiceURL set to: ", calledServiceURL)
			}

		}
	}

	log.Println("=============================================")


	router := mux.NewRouter() //.StrictSlash(true)

	router.HandleFunc("/", Index).Methods("GET")
	router.HandleFunc("/echo/{message}", echoHandler(calledServiceURL)).Methods("GET")
	router.HandleFunc("/factorialIterative/{number}", factorialIterativeHandler(calledServiceURL)).Methods("GET")
	router.HandleFunc("/factorialRecursive/{number}", factorialRecursiveHandler(calledServiceURL)).Methods("GET")

	// set timeout
	muxWithMiddlewares := http.TimeoutHandler(router, time.Second*3, "Timeout!")

	srv := &http.Server{
		Addr:         port,
		// Good practice to set timeouts to avoid Slowloris attacks.
		// Using just the read parameter due to this article
		// https://stackoverflow.com/questions/29334407/creating-an-idle-timeout-in-go
		//WriteTimeout: time.Second * 60,
		ReadTimeout:  time.Second * 15,
		//IdleTimeout:  time.Second * 120,
		//Handler: router, // Pass our instance of gorilla/mux in.
		Handler: muxWithMiddlewares,
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
	ctx, cancel := context.WithTimeout(context.Background(), wait)
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
