package structs

import (
	"errors"
)

type Hash_Table struct {
	Table []*Doubly_Connected_Table
	Size  int
	Cout  int
}

type Pair struct {
	Key   string
	Value string
}

type Node_Table struct {
	Data     Pair
	Next     *Node_Table
	Previous *Node_Table
}

type Doubly_Connected_Table struct {
	Lenght int
	Head   *Node_Table
	Tail   *Node_Table
}

func NewHashTable(size int) Hash_Table {
	table := make([]*Doubly_Connected_Table, size)
	for i := range table {
		table[i] = &Doubly_Connected_Table{}
	}
	return Hash_Table{
		Table: table,
		Size:  size,
	}
}

func (ht *Hash_Table) Hash(key string) int {
	keyInt := 0
	prime := 31
	for _, symbol := range key {
		keyInt = (keyInt*prime + int(symbol)) % ht.Size
	}
	return keyInt % ht.Size
}

func (ht *Hash_Table) Hset(key string, value string) error {
	val := &Pair{Key: key, Value: value}
	hash := ht.Hash(val.Key)
	if ht.Table[hash].Lenght < 500 {
		ht.Cout++
		return ht.Table[hash].dpush(*val)
	} else {
		oldSize := ht.Size
		newHT := NewHashTable(oldSize * 2)
		for i := 0; i < oldSize; i++ {
			currentNode := ht.Table[i].Head
			for currentNode != nil {
				new_hash := newHT.Hash(currentNode.Data.Key)
				ht.Cout++
				newHT.Table[new_hash].dpush(currentNode.Data)
				currentNode = currentNode.Next
			}
		}
		*ht = newHT
		new_hash := ht.Hash(val.Key)
		ht.Cout++
		return ht.Table[new_hash].dpush(*val)
	}
}

func (ht *Hash_Table) Hdel(key string) (string, error) {
	hash := ht.Hash(key)
	pair, err := ht.Table[hash].ddel(key)
	return pair.Value, err
}

func (ht *Hash_Table) Hget(key string) (string, error) {
	hash := ht.Hash(key)
	currentNode := ht.Table[hash].Head
	for currentNode != nil {
		if currentNode.Data.Key == key {
			return currentNode.Data.Value, nil
		}
		currentNode = currentNode.Next
	}
	return "", errors.New("element not found")
}

func (doubly_connected *Doubly_Connected_Table) dpush(val Pair) error {
	node_hesh := &Node_Table{Data: val}
	if doubly_connected.Head == nil {
		doubly_connected.Head = node_hesh
		doubly_connected.Tail = node_hesh
	} else {
		currentNode := doubly_connected.Head
		for currentNode != nil {
			if currentNode.Data.Key == val.Key {
				return errors.New("key already exists")
			}
			currentNode = currentNode.Next
		}
		doubly_connected.Tail.Next = node_hesh
		node_hesh.Previous = doubly_connected.Tail
		doubly_connected.Tail = node_hesh
	}
	doubly_connected.Lenght++
	return nil
}

func (doubly_connected *Doubly_Connected_Table) ddel(val string) (Pair, error) {
	currentNode := doubly_connected.Head
	if currentNode == nil {
		return Pair{}, errors.New("list is clear")
	}
	for currentNode != nil {
		if currentNode.Data.Key == val {
			if currentNode == doubly_connected.Head {
				doubly_connected.Head = currentNode.Next
				if doubly_connected.Head != nil {
					doubly_connected.Head.Previous = nil
				}
			} else if currentNode == doubly_connected.Tail {
				doubly_connected.Tail = currentNode.Previous
				if doubly_connected.Tail != nil {
					doubly_connected.Tail.Next = nil
				}
			} else {
				currentNode.Previous.Next = currentNode.Next
				currentNode.Next.Previous = currentNode.Previous
			}
			doubly_connected.Lenght--
			return currentNode.Data, nil
		}

		currentNode = currentNode.Next
	}
	return Pair{}, errors.New("key not founde")
}
