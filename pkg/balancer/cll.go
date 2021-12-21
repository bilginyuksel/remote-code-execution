package balancer

type CircularLinkedListNode struct {
	value string
	next  *CircularLinkedListNode
	prev  *CircularLinkedListNode
}

type CircularLinkedList struct {
	head *CircularLinkedListNode
	curr *CircularLinkedListNode
}

func NewCLL(headValue string) *CircularLinkedList {
	headNode := &CircularLinkedListNode{value: headValue}
	headNode.next = headNode
	headNode.prev = headNode

	return &CircularLinkedList{
		head: headNode,
		curr: headNode,
	}
}

func (c *CircularLinkedList) Insert(node *CircularLinkedListNode, value string) {
	if node == nil {
		panic("can't insert after nil node")
	}

	node.next = &CircularLinkedListNode{value: value, prev: node, next: node.next}
}

func (c *CircularLinkedList) Remove(node *CircularLinkedListNode) {
	prevNode := node.prev
	nextNode := node.next

	prevNode.next = nextNode

	// cleanup
	node.next = nil
	node.prev = nil
}

func (c *CircularLinkedList) PushFront(value string) {
	c.Insert(c.head, value)
}

func (c *CircularLinkedList) GetCurr() *CircularLinkedListNode {
	return c.curr
}

func (c *CircularLinkedList) GetCurrAndNext() *CircularLinkedListNode {
	curr := c.curr
	c.curr = c.curr.next
	return curr
}
