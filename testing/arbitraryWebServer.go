package main

import (
	"encoding/json"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/_cluster/healthz", jRespond)
	http.HandleFunc("/_cluster/health", jRespondguts)
	http.HandleFunc("/close", closeMe)
	http.ListenAndServe(":8081", nil)

}

func jRespond(w http.ResponseWriter, r *http.Request) {

	jsonResponseBody := []byte(`{"status":"green","number_of_nodes":1,"unassigned_shards":0}`)

	js := json.Valid(jsonResponseBody)

	if js {
		w.Write([]byte("It was a valid json object"))
	} else {
		w.Write([]byte("it was invalid AF"))
	}
}

func jRespondguts(w http.ResponseWriter, r *http.Request) {

	jsonResponseBody := []byte(`{"status":"green","number_of_nodes":1,"unassigned_shards":0}`)

	w.Header().Set("content-type", "application/json; charset=utf-8")
	w.Write(jsonResponseBody)
}

func closeMe(w http.ResponseWriter, r *http.Request) {
	os.Exit(0)
}
