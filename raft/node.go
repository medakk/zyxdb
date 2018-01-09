package raft

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Node struct {
	Id       int    `yaml: "id"`
	Name     string `yaml: "name"`
	Location string `yaml: "location"`
}

func (node *Node) sendRequestvote(request RequestVoteRequest) (RequestVoteResponse, error) {
	response := RequestVoteResponse{}

	requestBody, err := json.Marshal(request)
	if err != nil {
		return response, err
	}

	url := fmt.Sprintf("http://%s/request-vote/", node.Location)
	httpResponse, err := http.Post(url, "application/json", bytes.NewReader(requestBody))
	if err != nil {
		return response, err
	}

	responseBody, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return response, err
	}

	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return response, err
	}

	return response, nil
}

func (node *Node) sendAppendEntries(request AppendEntriesRequest) (AppendEntriesResponse, error) {
	response := AppendEntriesResponse{}

	requestBody, err := json.Marshal(request)
	if err != nil {
		return response, err
	}

	url := fmt.Sprintf("http://%s/append-entries/", node.Location)
	httpResponse, err := http.Post(url, "application/json", bytes.NewReader(requestBody))
	if err != nil {
		return response, err
	}

	if httpResponse.StatusCode != http.StatusOK {
		err = errors.New(fmt.Sprintf("Bad status code: %d", httpResponse.StatusCode))
		return response, err
	}

	if httpResponse.StatusCode != http.StatusOK {
		err = errors.New(fmt.Sprintf("Bad status code: %d", httpResponse.StatusCode))
		return response, err
	}

	responseBody, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return response, err
	}

	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return response, err
	}

	return response, nil
}
