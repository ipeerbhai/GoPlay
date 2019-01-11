package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// SimpleHandler -- literally just a simple forwarder as an http handler.
func SimpleHandler(res http.ResponseWriter, req *http.Request) {
	defer req.Body.Close() // setup to close in case things go wrong.

	// Let's start with accepting a connect from them and parsing out the URI I need to connect to.
	fromClientURI := req.RequestURI
	forwardToBaseAddress := "http://www.youtube.com"
	fullURI := forwardToBaseAddress + fromClientURI

	// let's go make the same get
	fromForwardedServer, serverErr := http.Get(fullURI)
	if serverErr != nil {
		panic(serverErr)
	}
	defer fromForwardedServer.Body.Close() // make sure we close what we got whenever we return/whatever

	// Gotta copy the headers...
	for headerKey, headerValueArray := range fromForwardedServer.Header {
		for _, thisHeaderValueSlice := range headerValueArray {
			res.Header().Add(headerKey, thisHeaderValueSlice)
		}
	}

	// the status code
	res.WriteHeader(fromForwardedServer.StatusCode)

	// Now, send the body, byte for byte.
	body, bodyErr := ioutil.ReadAll(fromForwardedServer.Body)
	if bodyErr != nil {
		panic(bodyErr)
	}
	res.Write(body)
}

// main -- the main entry point.
func main() {
	// starup a server and have it pass everything to a simple handler
	myHTTPServer := &http.Server{
		Addr:           ":8000",
		Handler:        http.HandlerFunc(SimpleHandler),
		ReadTimeout:    1 * time.Second,
		WriteTimeout:   1 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	err := myHTTPServer.ListenAndServe()
	if err != nil {
		fmt.Println("Server failed: ", err.Error())
	}
}
