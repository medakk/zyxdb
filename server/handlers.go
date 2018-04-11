package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"gitlab.com/medakk/zyxdb/raft"
)

func AppendEntriesHandler(w http.ResponseWriter, r *http.Request, raftCtx *raft.RaftCtx) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		// TODO: Why 500?
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	requestData := raft.AppendEntriesRequest{}
	err = json.Unmarshal(b, &requestData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := raftCtx.AppendEntries(requestData)
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		panic(fmt.Sprintf("Failed to marshall raft response: %v", resp))
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
}

func RequestVoteHandler(w http.ResponseWriter, r *http.Request, raftCtx *raft.RaftCtx) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		// TODO: Why 500?
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	requestData := raft.RequestVoteRequest{}
	err = json.Unmarshal(b, &requestData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := raftCtx.RequestVote(requestData)
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		panic(fmt.Sprintf("Failed to marshall raft response: %v", resp))
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)

}

func PingHandler(w http.ResponseWriter, r *http.Request) {
	msg, _ := json.Marshal(map[string]string{
		"ping": "pong",
	})

	w.WriteHeader(http.StatusOK)
	w.Write(msg)
}

func InsertToLogHandler(w http.ResponseWriter, r *http.Request, raftCtx *raft.RaftCtx) {
	if raftCtx.State() == raft.StateFollower {
		leaderLocation := raftCtx.GetLeader().Location
		http.Redirect(w, r, leaderLocation+"/insert/", http.StatusPermanentRedirect)
		return
	}

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		// TODO: Why 500?
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	requestData := raft.InsertToLogRequest{}
	err = json.Unmarshal(b, &requestData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := raftCtx.InsertToLog(requestData)
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		panic(fmt.Sprintf("Failed to marshall raft response: %v", resp))
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
}

func DebugGetEntriesHandler(w http.ResponseWriter, r *http.Request, raftCtx *raft.RaftCtx) {
	b, err := json.Marshal(raftCtx.DebugGetEntries())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(b)
}
