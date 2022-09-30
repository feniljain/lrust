package test

import (
	"lru/dll"
	mglru "lru/mglru"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBasicPutGetMGLRU(t *testing.T) {
	obj := mglru.Init(2, 1) // nil

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

	assert.Equal(t, emptyDeletedPair, deletedPair, "wrong data returned")

	assert.Equal(t, 20, obj.Get(2), "wrong data returned")
}

func TestBasicPutGetMGLRUWithOneLenAndCapacity(t *testing.T) {
	obj := mglru.Init(1, 1) // nil

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

	assert.Equal(t, pair1.Key, deletedPair.Key, "wrong data returned")
	assert.Equal(t, pair1.Value, deletedPair.Value, "wrong data returned")

	assert.Equal(t, 20, obj.Get(2), "wrong data returned")
}

func TestMLRUCache(t *testing.T) {
	obj := mglru.Init(2, 2)

	nPairs := 5

	pairs := make([]dll.Pair, nPairs)

	i := 0
	for i < nPairs {
		pairs[i] = dll.Pair{
			Key:   i + 1,
			Value: (i + 1) * 10,
		}
		i++
	}

	deletedPair := obj.Put(pairs[0])

	var emptyDeletedPair *dll.Pair
	assert.Equal(t, emptyDeletedPair, deletedPair, "wrong data returned")

	deletedPair = obj.Put(pairs[1])

	assert.Equal(t, emptyDeletedPair, deletedPair, "wrong data returned")
	assert.Equal(t, 10, obj.Get(1), "wrong data returned")
	assert.Equal(t, 1, obj.GetFirstElement().Key, "wrong data returned")

	deletedPair = obj.Put(pairs[2])

	assert.Equal(t, emptyDeletedPair, deletedPair, "wrong data returned")
	assert.Equal(t, 3, obj.GetFirstElement().Key, "wrong data returned")

	deletedPair = obj.Put(pairs[3])

	assert.Equal(t, emptyDeletedPair, deletedPair, "wrong data returned")
	assert.Equal(t, 4, obj.GetFirstElement().Key, "wrong data returned")

	deletedPair = obj.Put(pairs[4])

	assert.Equal(t, &dll.Pair{Key: 2, Value: 20}, deletedPair, "wrong data returned")
	assert.Equal(t, -1, obj.Get(2), "wrong data returned")

	assert.Equal(t, 30, obj.Get(3), "wrong data returned")
	assert.Equal(t, 3, obj.GetFirstElement().Key, "wrong data returned")
}
