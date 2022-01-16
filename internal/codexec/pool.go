package codexec

import "time"

type (
	// ContainerNode node that stores the container id, next and prev and metrics information
	ContainerNode struct {
		ID             string
		ExecutionCount int64
		AvgResTime     time.Duration
		MaxResTime     time.Duration
		MinResTime     time.Duration

		Prev *ContainerNode
		Next *ContainerNode
	}

	// ContainerPool stores the container nodes in a structured way
	// Active container nodes stored in a Circular Doubly Linked List and Nodes map
	// Passive, Exited container nodes are removed from CDLL and stored in garbage
	ContainerPool struct {
		Head    *ContainerNode
		Tail    *ContainerNode
		Curr    *ContainerNode
		Nodes   map[string]*ContainerNode
		Garbage map[string]*ContainerNode
	}
)

func NewContainerPool() *ContainerPool {
	return &ContainerPool{
		Nodes:   make(map[string]*ContainerNode),
		Garbage: make(map[string]*ContainerNode),
	}
}

func (p *ContainerPool) Get() *ContainerNode {
	if p.Head == nil {
		return nil
	}

	if p.Curr == nil {
		p.Curr = p.Head
	}

	defer func() { p.Curr = p.Curr.Next }()
	return p.Curr
}

func (p *ContainerPool) Remove(id string) {
	// Remove from nodes and from list
	// Add to garbage to collect later
	nodeToRemove := p.Nodes[id]
	delete(p.Nodes, id)

	if p.Head == p.Tail {
		p.Head = nil
		p.Tail = nil
		p.Curr = nil
		return
	}

	if p.Curr == nodeToRemove {
		p.Curr = nodeToRemove.Next
	}

	if p.Head == nodeToRemove {
		p.Head.Next.Prev = p.Head.Prev
		p.Head = p.Head.Next
		return
	}

	if p.Tail == nodeToRemove {
		p.Tail.Prev.Next = p.Tail.Next
		p.Tail = p.Tail.Prev
		return
	}

	nodeToRemove.Next.Prev = nodeToRemove.Prev
	nodeToRemove.Prev.Next = nodeToRemove.Next

	nodeToRemove.Next = nil
	nodeToRemove.Prev = nil
}

func (p *ContainerPool) Add(id string) {
	node := &ContainerNode{ID: id}
	p.Nodes[id] = node

	if p.Head == nil {
		node.Next = node
		node.Prev = node
		p.Head = node
		p.Tail = node
		p.Curr = node
		return
	}

	if p.Tail == p.Head {
		node.Next = p.Head
		node.Prev = p.Head
		p.Head.Next = node
		p.Head.Prev = node
		p.Tail = node
		return
	}

	node.Next = p.Tail.Next
	node.Prev = p.Tail
	p.Tail.Next.Prev = node
	p.Tail.Next = node
	p.Tail = node
}
