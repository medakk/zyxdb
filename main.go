package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type RetrievePost struct {
	Key string `json: "key"`
}

type InsertPost struct {
	Key   string `json: "key"`
	Value string `json: "value"`
}

var (
	m = make(map[string]string)
)

func RetrieveHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		msg, _ := json.Marshal(map[string]string{
			"status": "this method is not allowed",
		})
		w.Write(msg)
		return
	}

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		// TODO: Why 500?
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	requestData := RetrievePost{}
	err = json.Unmarshal(b, &requestData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	key := requestData.Key
	value, ok := m[key]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		msg, _ := json.Marshal(map[string]string{
			key: "no value found",
		})
		w.Write(msg)
		return
	}

	w.WriteHeader(http.StatusFound)
	msg, _ := json.Marshal(map[string]string{
		key: value,
	})
	w.Write(msg)
}

func InsertHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		msg, _ := json.Marshal(map[string]string{
			"status": "this method is not allowed",
		})
		w.Write(msg)
		return
	}

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		// TODO: Why 500?
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	requestData := InsertPost{}
	err = json.Unmarshal(b, &requestData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	key, value := requestData.Key, requestData.Value
	m[key] = value

	w.WriteHeader(http.StatusCreated)
	msg, _ := json.Marshal(map[string]string{
		key: value,
	})
	w.Write(msg)
}

func main() {
	host := flag.String("host", "localhost", "host name to bind to")
	port := flag.Int("port", 5002, "port to listen on")
	flag.Parse()

	listenOn := fmt.Sprintf("%s:%d", *host, *port)
	log.Printf("Starting server on %s\n", listenOn)

	http.HandleFunc("/retrieve/", RetrieveHandler)
	http.HandleFunc("/insert/", InsertHandler)
	log.Fatal(http.ListenAndServe(listenOn, nil))
}
