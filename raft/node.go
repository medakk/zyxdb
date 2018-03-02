package raft

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	timeout = 50 * time.Millisecond
)

type Node struct {
	Location string `yaml:"location"`
}

func (node *Node) sendRequestvote(request RequestVoteRequest) (RequestVoteResponse, error) {
	response := RequestVoteResponse{}

	requestBody, err := json.Marshal(request)
	if err != nil {
		return response, err
	}

	url := fmt.Sprintf("http://%s/request-vote/", node.Location)
	httpClient := http.Client{
		Timeout: timeout,
	}
	httpResponse, err := httpClient.Post(url, "application/json", bytes.NewReader(requestBody))
	if err != nil {
		return response, err
	}

	if httpResponse.StatusCode != http.StatusOK {
		err = fmt.Errorf("Bad status code: %d", httpResponse.StatusCode)
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
	httpClient := http.Client{
		Timeout: timeout,
	}
	httpResponse, err := httpClient.Post(url, "application/json", bytes.NewReader(requestBody))
	if err != nil {
		return response, err
	}

	if httpResponse.StatusCode != http.StatusOK {
		err = fmt.Errorf("Bad status code: %d", httpResponse.StatusCode)
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
