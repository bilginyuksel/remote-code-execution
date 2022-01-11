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

	// Pool stores the container nodes in a structured way
	// Active container nodes stored in a Circular Doubly Linked List and Nodes map
	// Passive, Exited container nodes are removed from CDLL and stored in garbage
	Pool struct {
		Head    *ContainerNode
		Tail    *ContainerNode
		Curr    *ContainerNode
		Nodes   map[string]*ContainerNode
		Garbage map[string]*ContainerNode
	}
)

func NewPool() *Pool {
	return &Pool{
		Nodes:   make(map[string]*ContainerNode),
		Garbage: make(map[string]*ContainerNode),
	}
}

func (p *Pool) Get() *ContainerNode {
	if p.Head == nil {
		return nil
	}

	if p.Curr == nil {
		p.Curr = p.Head
	}

	defer func() { p.Curr = p.Curr.Next }()
	return p.Curr
}

func (p *Pool) Add(id string) {
	node := &ContainerNode{ID: id}
	p.Nodes[id] = node

	if p.Head == nil {
		node.Next = node
		node.Prev = node
		p.Head = node
		p.Tail = node
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
