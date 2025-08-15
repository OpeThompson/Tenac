package raft

import (
    "log"
    "math/rand"
    "sync"
    "time"
)

// State represents the Raft node state.
type State int

const (
    Follower State = iota
    Candidate
    Leader
)

func (s State) String() string {
    switch s {
    case Follower:
        return "Follower"
    case Candidate:
        return "Candidate"
    case Leader:
        return "Leader"
    default:
        return "Unknown"
    }
}

// Config holds basic Raft timing config.
type Config struct {
    MinElectionTimeoutMs int // e.g., 150
    MaxElectionTimeoutMs int // e.g., 300
    HeartbeatMs          int // e.g., 50-100
}

// Node is a minimal scaffold for Raft behaviors (Week 2 starter).
// This DOES NOT implement RPCs yet; it only demonstrates timers and state transitions for learning.
type Node struct {
    mu       sync.Mutex
    id       string
    state    State
    cfg      Config
    stopCh   chan struct{}
    stopped  bool
}

// NewNode creates a Raft node in Follower state.
func NewNode(id string, cfg Config) *Node {
    return &Node{
        id:     id,
        state:  Follower,
        cfg:    cfg,
        stopCh: make(chan struct{}),
    }
}

// Start launches the election timer loop in a goroutine.
func (n *Node) Start() {
    go n.run()
}

func (n *Node) run() {
    rand.Seed(time.Now().UnixNano())
    for {
        timeout := time.Duration(rand.Intn(n.cfg.MaxElectionTimeoutMs-n.cfg.MinElectionTimeoutMs)+n.cfg.MinElectionTimeoutMs) * time.Millisecond
        select {
        case <-time.After(timeout):
            n.mu.Lock()
            if n.state == Follower {
                n.state = Candidate
                log.Printf("[raft:%s] election timeout -> become %s", n.id, n.state.String())
                // In real Raft: increment term, vote for self, send RequestVote RPCs
                // For Week 2: immediately simulate winning and become Leader
                n.state = Leader
                log.Printf("[raft:%s] simulated election win -> become %s", n.id, n.state.String())
            } else if n.state == Candidate {
                // retry election
                log.Printf("[raft:%s] still Candidate, retrying election", n.id)
            } else if n.state == Leader {
                // Leaders would send heartbeats; we just log for now
                log.Printf("[raft:%s] sending heartbeats (simulated)", n.id)
            }
            n.mu.Unlock()
        case <-n.stopCh:
            n.mu.Lock()
            n.stopped = true
            n.mu.Unlock()
            return
        }
    }
}

// Stop stops the node loops.
func (n *Node) Stop() {
    close(n.stopCh)
}
