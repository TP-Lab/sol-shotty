package main

import (
	"fmt"
	"net/http"

	"github.com/trustless-engineering/sol-shotty/pkg"
	"github.com/trustless-engineering/sol-shotty/pkg/utils"
)

var endpointList []string

func proxy(w http.ResponseWriter, req *http.Request) {
	// Shotgun the request to all endpoints
	response, err := pkg.Shotgun(endpointList, req)
	if err != nil {
		fmt.Printf("Error shotgunning request: %v\n", err)
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
		return
	}

	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, solana-client")

	// Copy the response headers back to the client
	for key, values := range response.Result.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	// Add CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")

	w.WriteHeader(response.Result.StatusCode)

	// Write the response body to the client's response writer
	_, err = w.Write(response.Body)
	if err != nil {
		fmt.Printf("Error writing response: %v\n", err)
		return
	}

	defer response.Result.Body.Close()

	fmt.Printf("Successful response from %s in %dms\n", response.Endpoint, response.RTT)
}

func main() {
	fmt.Printf("Loading the shotty...\n")
	if endpoints, err := utils.LoadEndpoints(); err != nil {
		fmt.Printf("Error loading endpoints: %v\n", err)
		return
	} else {
		endpointList = endpoints
	}
	http.HandleFunc("/", proxy)
	fmt.Printf("Shotty is ready to fire! Listening at http://127.0.0.1:420\n")
	http.ListenAndServe("172.28.143.6:420", nil)
}
