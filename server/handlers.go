package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.com/medakk/zyxdb/raft"
	"gitlab.com/medakk/zyxdb/storage"
)

type AppendEntriesPost struct {
	Term         int      `json: "term"`
	LeaderId     int      `json: "leader_id"`
	PrevLogIndex int      `json: "prev_log_index"`
	PrevLogTerm  int      `json: "prev_log_term"`
	Entries      []string `json: "entries"`
	LeaderCommit int      `json: "leader_commit"`
}

type RequestVotePost struct {
	Term         int `json: "term"`
	CandidateId  int `json: "candidate_id"`
	LastLogIndex int `json: "last_log_index"`
	LastLogTerm  int `json: "last_log_term"`
}

type RetrievePost struct {
}

type InsertPost struct {
	Value string `json: "value"`
}

func AppendEntriesHandler(w http.ResponseWriter, r *http.Request, c *raft.RaftState) {

}

func RequestVoteHandler(w http.ResponseWriter, r *http.Request, c *raft.RaftState) {

}

func RetrieveHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

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

	key := vars["key"]
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
	vars := mux.Vars(r)

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

	key := vars["key"]
	value := requestData.Value
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
