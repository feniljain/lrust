package test

import (
	"lru"
	"lru/dll"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBasicPutGet(t *testing.T) {
	obj := lru.Init(1)

	pair1 := dll.Pair{
		Key:   1,
		Value: 10,
	}

	deletedPair := obj.Put(pair1)

	var emptyDeletedPair *dll.Pair
	assert.Equal(t, emptyDeletedPair, deletedPair, "wrong data returned")

	assert.Equal(t, 10, obj.Get(1), "wrong data returned")
}

func TestBasicPutGetWithTwoElementsWithSameData(t *testing.T) {
	obj := lru.Init(1)

	pair1 := dll.Pair{
		Key:   1,
		Value: 10,
	}

	pair2 := dll.Pair{
		Key:   1,
		Value: 10,
	}

	deletedPair := obj.Put(pair1)

	var emptyDeletedPair *dll.Pair
	assert.Equal(t, emptyDeletedPair, deletedPair, "wrong data returned")

	assert.Equal(t, 10, obj.Get(1), "wrong data returned")

	deletedPair = obj.Put(pair2)

	assert.Equal(t, emptyDeletedPair, deletedPair, "wrong data returned")

	assert.Equal(t, 10, obj.Get(1), "wrong data returned")
}

func TestBasicPutGetWithTwoElementsWithDifferentData(t *testing.T) {
	obj := lru.Init(1)

	pair1 := dll.Pair{
		Key:   1,
		Value: 10,
	}

	pair2 := dll.Pair{
		Key:   2,
		Value: 20,
	}

	deletedPair := obj.Put(pair1)

	var emptyDeletedPair *dll.Pair
	assert.Equal(t, emptyDeletedPair, deletedPair, "wrong data returned")

	assert.Equal(t, 10, obj.Get(1), "wrong data returned")

	deletedPair = obj.Put(pair2)

	assert.Equal(t, &pair1, deletedPair, "wrong data returned")

	assert.Equal(t, 20, obj.Get(2), "wrong data returned")
}

func TestLRUCache(t *testing.T) {
	obj := lru.Init(2) // nil

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

	pair4 := dll.Pair{
		Key:   4,
		Value: 40,
	}

	deletedPair := obj.Put(pair1) // nil, linked list: [1:10]

	var emptyDeletedPair *dll.Pair
	assert.Equal(t, emptyDeletedPair, deletedPair, "wrong data returned")

	deletedPair = obj.Put(pair2) // nil, linked list: [2:20, 1:10]

	assert.Equal(t, emptyDeletedPair, deletedPair, "wrong data returned")
	assert.Equal(t, 10, obj.Get(1), "wrong data returned") // 10, linked list: [1:10, 2:20]

	deletedPair = obj.Put(pair3) // nil, linked list: [3:30, 1:10]

	assert.Equal(t, &dll.Pair{Key: 2, Value: 20}, deletedPair, "wrong data returned")
	assert.Equal(t, -1, obj.Get(2), "wrong data returned") // -1, linked list: [3:30, 1:10]

	deletedPair = obj.Put(pair4) // nil, linked list: [4:40, 3:30]

	assert.Equal(t, &dll.Pair{Key: 1, Value: 10}, deletedPair, "wrong data returned")
	assert.Equal(t, -1, obj.Get(1), "wrong data returned") // -1, linked list: [4:40, 3:30]

	assert.Equal(t, 30, obj.Get(3), "wrong data returned") // 30, linked list: [3:30, 4:40]
}
