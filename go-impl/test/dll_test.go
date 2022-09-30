package test

import (
	"lru/dll"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBasicPushFront(t *testing.T) {

	pair2 := dll.Pair{
		Key:   2,
		Value: 20,
	}

	pair3 := dll.Pair{
		Key:   3,
		Value: 30,
	}

	dllist := dll.Init()

	dllist.PushFront(pair3)
	assert.Equal(t, pair3, dllist.Head().Data, "wrong data at wrong place")

	dllist.PushFront(pair2)
	assert.Equal(t, pair2, dllist.Head().Data, "wrong data at wrong place")
}

func TestBasicPushBack(t *testing.T) {
	dllist := dll.Init()

	pair2 := dll.Pair{
		Key:   2,
		Value: 20,
	}

	pair3 := dll.Pair{
		Key:   3,
		Value: 30,
	}

	dllist.PushBack(pair2)
	assert.Equal(t, pair2, dllist.Head().Data, "wrong data at wrong place")

	dllist.PushBack(pair3)
	assert.Equal(t, pair2, dllist.Head().Data, "wrong data at wrong place")
}

func TestPopBack(t *testing.T) {
	dllist := dll.Init()

	pair2 := dll.Pair{
		Key:   2,
		Value: 20,
	}

	pair3 := dll.Pair{
		Key:   3,
		Value: 30,
	}

	dllist.PushBack(pair2)
	dllist.PushBack(pair3)

	dllist.PopBack()

	assert.Equal(t, 1, dllist.Size(), "incorrect size state")

	assert.Equal(t, pair2, dllist.Head().Data, "wrong data at wrong place")
	assert.Equal(t, pair2, dllist.Tail().Data, "head != tail")

	dllist1 := dll.Init()

	dllist1.PushBack(pair2)

	dllist1.PopBack()

	assert.Equal(t, 0, dllist1.Size(), "incorrect size state")

	if dllist1.Head() != nil {
		assert.Fail(t, "head pointer should be nil")
	}

	if dllist1.Tail() != nil {
		assert.Fail(t, "head pointer should be nil")
	}

	// assert.Equal tests are failing with:
	// expected: <nil>(<nil>)
	// actual  : *dll.Node((*dll.Node)(nil))
	// so that's why defined the tests manually above,
	// maybe this is because assert.Equal use
	// runtime reflection and that gives type as
	// *dll.Node((*dll.Node)(nil)), I am happy not
	// tryna replicate this type for now

	// assert.Equal(t, nil, dllist1.Head(), "")
	// assert.Equal(t, nil, dllist1.Tail(), "head != tail")
}

func TestBasicPopFront(t *testing.T) {
	dllist := dll.Init()

	pair2 := dll.Pair{
		Key:   2,
		Value: 20,
	}

	pair3 := dll.Pair{
		Key:   3,
		Value: 30,
	}

	dllist.PushFront(pair3)
	dllist.PushFront(pair2)

	dllist.PopFront()

	assert.Equal(t, 1, dllist.Size(), "incorrect size state")

	assert.Equal(t, pair3, dllist.Head().Data, "wrong data at wrong place")
	assert.Equal(t, pair3, dllist.Tail().Data, "head != tail")
}

func TestBasicRemove(t *testing.T) {
	dllist := dll.Init()

	pair1 := dll.Pair{
		Key:   1,
		Value: 10,
	}

	pair2 := dll.Pair{
		Key:   2,
		Value: 20,
	}

	pair3 := dll.Pair{
		Key:   3,
		Value: 30,
	}

	dllist.PushBack(pair2)
	dllist.PushBack(pair3)
	dllist.PushFront(pair1)

	dllist.Remove(pair2)

	assert.Equal(t, 2, dllist.Size(), "incorrect size state")

	assert.Equal(t, pair1, dllist.Head().Data, "wrong data at wrong place")
	assert.Equal(t, pair3, dllist.Tail().Data, "wrong data at wrong place")
}

func TestRemoveOneElement(t *testing.T) {
	dllist := dll.Init()

	pair1 := dll.Pair{
		Key:   1,
		Value: 10,
	}

	dllist.PushFront(pair1)

	dllist.Remove(pair1)

	assert.Equal(t, 0, dllist.Size(), "incorrect size state")

	assert.Equal(t, (*dll.Node)(nil), dllist.Head(), "wrong data at wrong place")
	assert.Equal(t, (*dll.Node)(nil), dllist.Tail(), "wrong data at wrong place")
}

func TestMovedNodeToFront(t *testing.T) {
	dllist := dll.Init()

	pair1 := dll.Pair{
		Key:   1,
		Value: 10,
	}

	pair2 := dll.Pair{
		Key:   2,
		Value: 20,
	}

	pair3 := dll.Pair{
		Key:   3,
		Value: 30,
	}

	dllist.PushBack(pair1)
	dllist.PushBack(pair2)
	cNode := dllist.PushBack(pair3)

	assert.Equal(t, pair3, cNode.Data, "wrong data stored in node")

	err := dllist.MoveNodeToFront(cNode)
	if err != nil {
		assert.Fail(t, "Failure due to error "+err.Error())
	}

	assert.Equal(t, pair3, dllist.Head().Data, "wrong data at wrong place")

	assert.Equal(t, pair2, dllist.Tail().Data, "wrong data at wrong place")
}

func TestMiscCases1(t *testing.T) {
	dllist := dll.Init()

	pair2 := dll.Pair{
		Key:   2,
		Value: 20,
	}

	pair3 := dll.Pair{
		Key:   3,
		Value: 30,
	}

	pair4 := dll.Pair{
		Key:   4,
		Value: 40,
	}

	dllist.PushFront(pair3)
	assert.Equal(t, pair3, dllist.Head().Data, "wrong data at wrong place")

	dllist.PushFront(pair2)
	assert.Equal(t, pair2, dllist.Head().Data, "wrong data at wrong place")

	dllist.PushBack(pair4)
	assert.Equal(t, 3, dllist.Size(), "incorrect size state")
	assert.Equal(t, pair2, dllist.Head().Data, "wrong data at wrong place")
	assert.Equal(t, pair4, dllist.Tail().Data, "wrong data at wrong place")

	dllist.PopBack()
	assert.Equal(t, 2, dllist.Size(), "incorrect size state")
	assert.Equal(t, pair2, dllist.Head().Data, "wrong data at wrong place")
	assert.Equal(t, pair3, dllist.Tail().Data, "wrong data at wrong place")

	dllist.PopFront()
	assert.Equal(t, 1, dllist.Size(), "incorrect size state")
	assert.Equal(t, pair3, dllist.Head().Data, "wrong data at wrong place")
	assert.Equal(t, pair3, dllist.Tail().Data, "wrong data at wrong place")
}
