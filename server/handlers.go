package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"gitlab.com/medakk/zyxdb/storage"
)

type RetrievePost struct {
	Key string `json: "key"`
}

type InsertPost struct {
	Key   string `json: "key"`
	Value string `json: "value"`
}

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
	value, ok := storage.Retrieve(key)
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
	storage.Insert(key, value)

	w.WriteHeader(http.StatusCreated)
	msg, _ := json.Marshal(map[string]string{
		key: value,
	})
	w.Write(msg)
}

func PingHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	msg, _ := json.Marshal(map[string]string{
		"status": "pong",
	})
	w.Write(msg)
}
