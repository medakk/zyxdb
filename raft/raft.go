package raft

import (
	"log"
	"math/rand"
	"net/http"
	"sync/atomic"
	"time"
)

const (
	// base time in ms, actual timeout could +/- the sway amount
	electionTimeout     = 1000
	electionTimeoutSway = 800

	// Leader sends empty AppendEntries request every time the interval of time passes
	heartbeatInterval = 200 * time.Millisecond
)

type NodeState int

const (
	StateLeader NodeState = iota
	StateCandidate
	StateFollower
)

var (
	// tries reading configs in this order, with the first one that is present
	// being used
	configFilepaths = []string{"./zyxdb.yml", "/etc/zyxdb/zyxdb.yml"}
)

type Log struct {
	term    int
	command string
}

type RaftCtx struct {
	// Node details
	name     string
	selfNode *Node

	// State of the node
	state NodeState

	// persistent state
	currentTerm int
	votedFor    *Node
	log         []Log

	// volatile state on all servers
	commitIndex int
	lastApplied int

	// volatile state on leader
	nextIndex  []int
	matchIndex []int

	// keep track of who the current leader is
	currentLeader *Node

	// configuration where the list of nodes is stored
	config ZyxdbConfig

	// integer to register number of entries since the last election timeout
	// always access this through sync/atomic functions
	heartbeat uint64
}

type AppendEntriesRequest struct {
	Term         int      `json:"term"`
	LeaderName   string   `json:"leader"`
	PrevLogIndex int      `json:"prev_log_index"`
	PrevLogTerm  int      `json:"prev_log_term"`
	Entries      []string `json:"entries"`
	LeaderCommit int      `json:"leader_commit"`
}

type AppendEntriesResponse struct {
}

type RequestVoteRequest struct {
	Term          int    `json:"term"`
	CandidateName string `json:"candidate"`
	LastLogIndex  int    `json:"last_log_index"`
	LastLogTerm   int    `json:"last_log_term"`
}

type RequestVoteResponse struct {
	Term        int  `json:"term"`
	VoteGranted bool `json:"vote_granted"`
}

func New(name string) *RaftCtx {
	zyxdbConfig, err := loadConfig()
	if err != nil {
		panic(err.Error())
	}

	c := RaftCtx{
		currentTerm: 1,
		votedFor:    nil,
		name:        name,
		selfNode:    zyxdbConfig.getNodeByName(name),
		state:       StateFollower,
		config:      zyxdbConfig,
	}

	go c.runTickerEvents()

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

	//TODO: Check for LastLogIndex
	if request.Term < c.currentTerm {
		log.Printf("False leader: %v", request.LeaderName)
	} else {
		log.Printf("%v is leader", request.LeaderName)
		c.currentTerm = request.Term
		c.votedFor = nil

		// Turns out there is a new leader, and its not me.
		if c.state != StateFollower {
			c.state = StateFollower
		}
		c.currentLeader = c.config.getNodeByName(request.LeaderName)
	}

	return response
}

func (c *RaftCtx) RequestVote(request RequestVoteRequest) RequestVoteResponse {
	log.Printf("someone wants a vote: %v\n", request)

	response := RequestVoteResponse{
		Term:        c.currentTerm,
		VoteGranted: false,
	}

	if c.state == StateCandidate {
		log.Printf("Won't vote for %v\n", request.CandidateName)
		return response
	}

	if request.Term > c.currentTerm {
		log.Printf("OK, I vote for %v\n", request.CandidateName)
		response.VoteGranted = true
		response.Term = request.Term
		c.votedFor = c.config.getNodeByName(request.CandidateName)
	}

	return response
}

func (c *RaftCtx) State() NodeState {
	return c.state
}

func (c *RaftCtx) GetLeader() *Node {
	return c.currentLeader
}

// runTickerEvents starts running time based events: the election timeout and
// the heartbeat
func (c *RaftCtx) runTickerEvents() {
	tickerElectionTimeout := time.NewTicker(electionTimeoutDuration())
	tickerHeartbeat := time.NewTicker(heartbeatInterval)

	for {
		select {
		case <-tickerElectionTimeout.C:
			//TODO: Do something better here
			// Skip if not follower
			if c.state != StateFollower {
				continue
			}

			heartbeatCount := atomic.LoadUint64(&c.heartbeat)
			if heartbeatCount == 0 {
				c.attemptLeadership()
			}

			// Reset to 0
			atomic.StoreUint64(&c.heartbeat, 0)

		case <-tickerHeartbeat.C:
			// TODO: Avoid this somehow
			if c.state != StateLeader {
				continue
			}

			request := AppendEntriesRequest{
				LeaderName: c.name,
				Term:       c.currentTerm,
			}

			log.Printf("Sending heartbeats")
			for name, node := range c.config.Nodes {
				if name != c.name {
					_, err := node.sendAppendEntries(request)
					if err != nil {
						log.Printf("Error while sending append-entries to %v: %v\n", node, err)
						continue
					}
				}
			}

		}
	}
}

func (c *RaftCtx) attemptLeadership() {
	c.state = StateCandidate

	c.currentTerm++                           // Increment term
	c.votedFor = c.selfNode                   // Vote for self
	totalVotes := 1                           // Keep tally of total votes
	requiredVotes := c.config.nodeCount() / 2 // Majority required

	request := RequestVoteRequest{
		Term:          c.currentTerm,
		CandidateName: c.name,
		// TODO: Add these:
		LastLogIndex: 0,
		LastLogTerm:  0,
	}

	for name, node := range c.config.Nodes {
		// Don't send to self
		if name == c.name {
			continue
		}

		response, err := node.sendRequestvote(request)
		if err != nil {
			log.Printf("Error while sending request-vote to %v: %v\n", node, err)
			continue
		}

		if response.VoteGranted {
			totalVotes++
		}
	}

	if totalVotes >= requiredVotes {
		c.state = StateLeader
		log.Printf("I am leader now.")
	} else {
		c.state = StateFollower
	}
}

func electionTimeoutDuration() time.Duration {
	sway := rand.Intn(electionTimeoutSway*2) - electionTimeoutSway
	return time.Duration(electionTimeout+sway) * time.Millisecond
}
