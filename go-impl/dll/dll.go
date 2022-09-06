package dll

import (
	"fmt"
	"lru/errors"
)

// Representing data stored in DLL
type Pair struct {
    Key int
    Value int
}

// Representing a single node in DLL
type Node struct {
	Data Pair
	Next *Node
	Prev *Node
}

// Doubly Linked List
type DoublyLinkedList struct {
	len  int
	tail *Node
	head *Node
}

func Init() *DoublyLinkedList {
	return &DoublyLinkedList{
		len:  0,
		tail: nil,
		head: nil,
	}
}

func (dll *DoublyLinkedList) PushFront(data Pair) *Node {
	newNode := &Node{
		Data: data,
		Next: nil,
		Prev: nil,
	}

	if dll.head == nil {
		dll.head = newNode
		dll.tail = newNode
	} else {
		dll.head.Prev = newNode
		newNode.Next = dll.head

		dll.head = newNode
	}

	dll.len++
	return newNode
}

func (dll *DoublyLinkedList) PushBack(data Pair) *Node {
	newNode := &Node{
		Data: data,
		Next: nil,
		Prev: nil,
	}

	if dll.tail == nil {
		dll.head = newNode
		dll.tail = newNode
	} else {
		dll.tail.Next = newNode
		newNode.Prev = dll.tail

		dll.tail = newNode
	}

	dll.len++
	return newNode
}

func (dll *DoublyLinkedList) PopFront() {
	if dll.head == nil {
		return
	}

	dll.head = dll.head.Next
	dll.head.Prev = nil

	dll.len--
}

func (dll *DoublyLinkedList) PopBack() (*Pair, bool) {
	if dll.tail == nil {
		return nil, false
	}

    nodeData := dll.tail.Data
	dll.tail = dll.tail.Prev

	// For the case where there is only one element in DLL
	// and pop back is called
	if dll.tail != nil {
		dll.tail.Next = nil
	} else {
		// As DLL is empty make head empty too
		dll.head = nil
	}

	dll.len--
    return &nodeData, true
}

func (dll *DoublyLinkedList) Remove(data Pair) (*Pair, bool) {
	if dll.head == nil {
		return nil, false
	}

	currNode := dll.head
	for currNode != nil {
		if currNode.Data == data {
			prevNode := currNode.Prev
			nextNode := currNode.Next

			// If head node matches
			if currNode.Prev != nil {
				currNode.Prev.Next = nextNode
			}

			// If tail node matches
			if nextNode == nil {
				dll.tail = prevNode
			} else {
				nextNode.Prev = prevNode
			}

			dll.len--

			return &currNode.Data, true
		}

		currNode = currNode.Next
	}

	return nil, false
}

func (d *DoublyLinkedList) Size() int {
	return d.len
}

func (d *DoublyLinkedList) Head() *Node {
	return d.head
}

func (d *DoublyLinkedList) Tail() *Node {
	return d.tail
}

func (d *DoublyLinkedList) MoveNodeToFront(node *Node) error {

	data, present := d.Remove(node.Data)
	if !present {
		return errors.NoNodeWithGivenData
	}

	d.PushFront(*data)

	return nil
}

func (d *DoublyLinkedList) PrintList() {
	if d.head == nil {
		fmt.Println("Empty list")
		return
	}

	currNode := d.head
	fmt.Print(currNode.Data)

	for currNode.Next != nil {
		fmt.Print(" -> ")
		currNode = currNode.Next
		fmt.Print(currNode.Data)
	}

	fmt.Println()
}
