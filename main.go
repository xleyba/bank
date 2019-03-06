package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"html"
	"log"
	"net/http"
	"os"
	"io/ioutil"
)

// Configuration
type Configuration struct {
	Port             string
	CalledServiceURL string
}

var port string
var calledServiceURL string

// Return default message for root routing
func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}


// Return default message for root routing
func echoHandler(calledServiceURL string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)

		url := calledServiceURL + "/echo/" + params["message"]

		resp, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
			fmt.Fprintf(w, "%s", err)
			return
		}

		if err != nil {
			log.Fatal(err)
		}

		robots, err := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		if err != nil {
			log.Fatal(err)
		}


		fmt.Fprintf(w, "%s", robots)
	}
}

// Handle iterative path and calls iterative calculation service
func factorialIterativeHandler(calledServiceURL string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)

		url := calledServiceURL + "/factorialIterative/" + params["number"]

		resp, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
			fmt.Fprintf(w, "%s", err)
			return
		}

		if err != nil {
			log.Fatal(err)
		}

		robots, err := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		if err != nil {
			log.Fatal(err)
		}


		fmt.Fprintf(w, "%s", robots)

	}
}

// Handle recursive path and calls recursive calculation service
func factorialRecursiveHandler(calledServiceURL string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)

		url := calledServiceURL + "/factorialRecursive/" + params["number"]

		resp, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
			fmt.Fprintf(w, "%s", err)
			return
		}

		if err != nil {
			log.Fatal(err)
		}

		robots, err := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		if err != nil {
			log.Fatal(err)
		}


		fmt.Fprintf(w, "%s", robots)

	}
}

// Main function
func main() {

	// Set default values
	port := ":9296"
	calledServiceURL := "http://localhost:9596"

	log.Println("Starting server")

	// Start to read conf file
	log.Println("bank v0.5")
	log.Println("=============================================")
	log.Println("         Configuration checking")
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

	log.Println("Running server....")

	log.Fatal(http.ListenAndServe(port, router))
}
