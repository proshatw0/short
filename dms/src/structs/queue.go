package structs

import "errors"

type Node struct {
	Data string
	Next *Node
}

type Queue struct {
	Head *Node
	Tail *Node
}

func (queue *Queue) Qpush(val string) error {
	if val == "" {
		return errors.New("non value")
	}
	node := &Node{Data: val}
	if queue.Head == nil {
		queue.Head = node
		queue.Tail = node
	} else {
		queue.Tail.Next = node
		queue.Tail = node
	}
	return nil
}

func (queue *Queue) Qpop() (string, error) {
	if queue.Head == nil {
		return "", errors.New("is empty")
	} else {
		val := queue.Head.Data
		queue.Head = queue.Head.Next
		return val, nil
	}
}
