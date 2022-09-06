package lru

import (
	"lru/dll"
)

type LRUCache struct {
	cap int
	m   map[int]*dll.Node
	l   *dll.DoublyLinkedList
}

func Init(capacity int) LRUCache {
	return LRUCache{
		cap: capacity,
		m:   make(map[int]*dll.Node, capacity),
		l:   dll.Init(),
	}
}

func (c *LRUCache) Get(key int) int {
	node, present := c.m[key]
	if !present {
		return -1
	}

	c.l.MoveNodeToFront(node)
	return node.Data.Value
}

func (c *LRUCache) Put(pair dll.Pair) {

	// If the key is already present in map
	if node, present := c.m[pair.Key]; present {
		// Push it to the the front
		c.l.MoveNodeToFront(node)

		// Update the value in dll and map
		c.m[pair.Key].Data = pair
	}

	// If list is full delete the last node
	// of dll (and also from map)
	if c.l.Size() == c.cap {
		pair, _ := c.l.PopBack()
		delete(c.m, pair.Key)
	}

	// Create a new node and push it in the front
	// Set the value in map

	newNode := c.l.PushFront(pair)

	c.m[pair.Key] = newNode
}
