package main

import (
	"net/http"
	"time"

	"github.com/ant0ine/go-json-rest/rest"
)

func SendSSE(w rest.ResponseWriter, r *rest.Request) {
	httpResponseWriter := w.(http.ResponseWriter)
	httpResponseWriter.Header().Set("Content-Type", "text/event-stream")
	httpResponseWriter.Header().Set("Cache-Control", "no-cache")
	httpResponseWriter.Header().Set("Connection", "keep-alive")

	// go interval(httpResponseWriter)
	for {
		if len(pipe["1"]) > 0 {
			pop := pipe["1"][0]
			if len(pipe["1"]) == 1 {
				pipe["1"] = []string{}
			} else {
				pipe["1"] = pipe["1"][1:]
			}
			httpResponseWriter.Write([]byte("data: " + pop + "\n\n"))
		}
		httpResponseWriter.Write([]byte("data: wait...\n\n"))
		time.Sleep(5 * time.Second)
	}

	// httpResponseWriter.Write([]byte("id: " + time.Now().String() + "\n"))
	// httpResponseWriter.Write([]byte("data: " + "test\n\n"))
}

func interval(w http.ResponseWriter) {
	for {
		id := time.Now().String()
		w.Write([]byte("id: " + id + "\n"))
		w.Write([]byte("data: " + "test data" + "\n\n"))
		time.Sleep(5 * time.Second)
	}
}
