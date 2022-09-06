package test

import (
	"lru"
    "lru/dll"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBasicPutGet(t *testing.T) {
    obj := lru.Init(2)   // nil

    pair1 := dll.Pair {
        Key: 1,
        Value: 10,
    }

    obj.Put(pair1)

    assert.Equal(t, 10, obj.Get(1), "wrong data returned")
}

func TestLRUCache(t *testing.T) {
    obj := lru.Init(2)   // nil

    pair1 := dll.Pair {
        Key: 1,
        Value: 10,
    }

    pair2 := dll.Pair {
        Key: 2,
        Value: 20,
    }

    pair3 := dll.Pair {
        Key: 3,
        Value: 30,
    }

    pair4 := dll.Pair {
        Key: 4,
        Value: 40,
    }

    obj.Put(pair1)          // nil, linked list: [1:10]

    obj.Put(pair2)          // nil, linked list: [2:20, 1:10]

    assert.Equal(t, 10, obj.Get(1), "wrong data returned") // 10, linked list: [1:10, 2:20]

    obj.Put(pair3)          // nil, linked list: [3:30, 1:10]

    assert.Equal(t, -1, obj.Get(2), "wrong data returned") // -1, linked list: [3:30, 1:10]

    obj.Put(pair4)          // nil, linked list: [4:40, 3:30]

    assert.Equal(t, -1, obj.Get(1), "wrong data returned") // -1, linked list: [4:40, 3:30]

    assert.Equal(t, 30, obj.Get(3), "wrong data returned") // 30, linked list: [3:30, 4:40]
}
