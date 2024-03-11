package ttlqueue

import (
	"sync"
	"time"
)

type Node struct {
	Container_id string
	TTL          time.Time
	LHS          *Node
	RHS          *Node
}

func NewNode(container_id string, ttl time.Time) *Node {
	return &Node{
		Container_id: container_id,
		TTL:          ttl,
	}
}

type Client struct {
	len      int
	m        map[string]*Node
	head     *Node
	tail     *Node
	baseTime time.Duration
	mu       sync.RWMutex
}

func New(baseTime time.Duration) *Client {
	head := NewNode("head", time.Now())
	tail := NewNode("tail", time.Now())
	head.RHS = tail
	tail.LHS = head
	return &Client{
		len:      0,
		baseTime: baseTime,
		m:        make(map[string]*Node),
		head:     head,
		tail:     tail,
	}
}

func (c *Client) GetLength() int {
	c.mu.RLock()
	len := c.len
	c.mu.RUnlock()
	return len
}

func (c *Client) IsNewContainer(container_id string) bool {
	c.mu.RLock()
	_, exist := c.m[container_id]
	c.mu.RUnlock()
	return !exist
}

// Return True if len has changed (container_id is new)
// Assume ttl is newer than head
func (c *Client) UpdateContainer(container_id string) bool {
	c.mu.Lock()
	node, exist := c.m[container_id]
	if !exist {
		ttl := time.Now().Add(c.baseTime)
		node := NewNode(container_id, ttl)
		c.m[container_id] = node
		c.appendLeft(node)
		c.len++
	} else {
		node.TTL = time.Now().Add(c.baseTime)
		c.popNode(node)
		c.appendLeft(node)
	}
	c.mu.Unlock()
	return !exist
}

// Return True if len has changed
func (c *Client) RemoveExpire() bool {
	c.mu.Lock()
	removed := false
	now := time.Now()
	for c.tail.LHS.Container_id != "head" && now.After(c.tail.LHS.TTL) {
		removed = true
		node := c.tail.LHS
		c.popNode(node)
		delete(c.m, node.Container_id)
		c.len--
	}
	c.mu.Unlock()
	return removed
}

// Return True if len has changed
// Do not use this because it introduce a bug where self cannot clean up all containers when alone
func (c *Client) CleanupOnExpire() bool {
	c.mu.RLock()
	node := c.tail.LHS
	if node.Container_id == "head" {
		return false
	}
	nap := time.Until(node.TTL)
	c.mu.RUnlock()
	if nap.Seconds() <= 0 {
		return c.RemoveExpire()
	}
	time.Sleep(nap)
	return c.RemoveExpire()
}

// Return True if len has changed
func (c *Client) CleanupOnContainerExpiration(container_id string) bool {
	c.mu.RLock()
	node, exist := c.m[container_id]
	if !exist {
		c.mu.RUnlock()
		return false
	} else {
		nap := time.Until(node.TTL)
		c.mu.RUnlock()
		if nap.Seconds() <= 0 {
			return c.RemoveExpire()
		}
		time.Sleep(nap)
		return c.RemoveExpire()
	}
}

func (c *Client) appendLeft(node *Node) {
	c.head.RHS.LHS = node
	node.RHS = c.head.RHS
	node.LHS = c.head
	c.head.RHS = node
}

func (c *Client) popNode(node *Node) {
	left := node.LHS
	right := node.RHS
	left.RHS = right
	right.LHS = left
}
