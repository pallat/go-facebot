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
	go waiting(httpResponseWriter)

	httpResponseWriter.Write([]byte("id: " + time.Now().String() + "\n"))
	httpResponseWriter.Write([]byte("data: wait...\n\n"))

	// httpResponseWriter.Write([]byte("id: " + time.Now().String() + "\n"))
	// httpResponseWriter.Write([]byte("data: " + "test\n\n"))
}

func waiting(w http.ResponseWriter) {
	for {
		if len(pipe["1"]) > 0 {
			pop := pipe["1"][0]
			if len(pipe["1"]) == 1 {
				pipe["1"] = []string{}
			} else {
				pipe["1"] = pipe["1"][1:]
			}
			w.Write([]byte("id: " + time.Now().String() + "\n"))
			w.Write([]byte("data: " + pop + "\n\n"))
		}
		w.Write([]byte("id: " + time.Now().String() + "\n"))
		w.Write([]byte("data: wait...\n\n"))
		time.Sleep(5 * time.Second)
	}
}

func interval(w http.ResponseWriter) {
	for {
		id := time.Now().String()
		w.Write([]byte("id: " + id + "\n"))
		w.Write([]byte("data: " + "test data" + "\n\n"))
		time.Sleep(5 * time.Second)
	}
}
