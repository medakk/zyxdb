package raft

import (
	"math/rand"
	"net/http"
	"sync/atomic"
	"time"
)

const (
	// base time in ms, actual timeout could +/- the sway amount
	electionTimeout     = 100
	electionTimeoutSway = 20

	// Leader sends empty AppendEntries request every time the interval of time passes
	heartbeatInterval = 10 * time.Millisecond
)

type RaftNodeState int

const (
	STATE_LEADER RaftNodeState = iota
	STATE_CANDIDATE
	STATE_FOLLOWER
)

var (
	// tries reading configs in this order, with the first one that is present
	// being used
	configFilepaths = []string{"./zyxdb.yml", "/etc/zyxdb/zyxdb.yml"}
)

type RaftCtx struct {
	// Node details
	selfNode Node

	// State of the node
	state RaftNodeState

	// persistent state
	currentTerm int
	votedFor    int
	log         []string

	// configuration where the list of nodes is stored
	config ZyxdbConfig

	// integer to register number of entries since the last election timeout
	// always access this through sync/atomic functions
	heartbeat uint64
}

type AppendEntriesRequest struct {
	Term         int      `json: "term"`
	LeaderId     int      `json: "leader_id"`
	PrevLogIndex int      `json: "prev_log_index"`
	PrevLogTerm  int      `json: "prev_log_term"`
	Entries      []string `json: "entries"`
	LeaderCommit int      `json: "leader_commit"`
}

type AppendEntriesResponse struct {
}

type RequestVoteRequest struct {
	Term         int `json: "term"`
	CandidateId  int `json: "candidate_id"`
	LastLogIndex int `json: "last_log_index"`
	LastLogTerm  int `json: "last_log_term"`
}

type RequestVoteResponse struct {
	Term        int  `json: "term"`
	VoteGranted bool `json: "vote_granted"`
}

func New(name string) *RaftCtx {
	zyxdbConfig, err := loadConfig()
	if err != nil {
		panic(err.Error())
	}

	c := RaftCtx{
		currentTerm: 1,
		selfNode:    zyxdbConfig.getNodeByName(name),
		state:       STATE_FOLLOWER,
		config:      zyxdbConfig,
	}

	c.startElectionTimeoutCheck()

	return &c
}

func (c *RaftCtx) Middleware(handler func(http.ResponseWriter, *http.Request, *RaftCtx)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, c)
	}
}

func (c *RaftCtx) AppendEntries(request AppendEntriesRequest) AppendEntriesResponse {
	atomic.AddUint64(&c.heartbeat, 1)

	response := AppendEntriesResponse{}
	return response
}

func (c *RaftCtx) RequestVote(request RequestVoteRequest) RequestVoteResponse {
	response := RequestVoteResponse{}
	return response
}

func (c *RaftCtx) startElectionTimeoutCheck() {
	go func() {
		for {
			// Wait for timeout
			time.Sleep(electionTimeoutDuration())

			//TODO: Do something better here
			// Skip if leader
			if c.state == STATE_LEADER {
				continue
			}

			heartbeatCount := atomic.LoadUint64(&c.heartbeat)
			if heartbeatCount == 0 {
				c.attemptLeadership()
			}

			// Reset to 0
			atomic.StoreUint64(&c.heartbeat, 0)
		}
	}()
}

func (c *RaftCtx) attemptLeadership() {
	c.state = STATE_CANDIDATE

	c.currentTerm++                             // Increment term
	c.votedFor = c.selfNode.Id                  // Vote for self
	totalVotes := 1                             // Keep tally of total votes
	requiredVotes := c.config.nodeCount()/2 + 1 // Majority required

	request := RequestVoteRequest{
		Term:        c.currentTerm,
		CandidateId: c.selfNode.Id,
		// TODO: Add these:
		LastLogIndex: 0,
		LastLogTerm:  0,
	}

	for _, node := range c.config.Nodes {
		// Don't send to self
		if node.Id == c.selfNode.Id {
			continue
		}

		response, err := node.sendRequestvote(request)
		if err != nil {
			continue
		}

		if response.VoteGranted {
			totalVotes++
		}
	}

	if totalVotes >= requiredVotes {
		c.state = STATE_LEADER
	} else {
		c.state = STATE_FOLLOWER
	}
}

func electionTimeoutDuration() time.Duration {
	sway := rand.Intn(electionTimeoutSway*2) - electionTimeoutSway
	return time.Duration(electionTimeout+sway) * time.Millisecond
}
