package mglru

import (
	"log"
	"lru"
	"lru/dll"
)

type MGLRUCache struct {
	lrus []lru.LRUCache
	// Capacity of each LRU
	capacity int
	// Number of LRUs in MGLRU
	length int
	// This if for quick access of all the keys, first
	// field is the key and second field is the index
	// of LRU where it is present
	keys map[int]int
}

func NewMGLRU(length, capacity int) MGLRUCache {
	var lrus []lru.LRUCache

	i := 0
	for i < length {
		lrus = append(lrus, lru.Init(capacity))
		i++
	}

	return MGLRUCache{
		lrus:     lrus,
		capacity: capacity,
		length:   length,
	}
}

func (l *MGLRUCache) Get(key int) int {
	if lruIdx, present := l.keys[key]; present {
		lru := l.lrus[lruIdx]

        // TODO: Promote received pair from the below GET
        // to highest LRU, with of course the same details
        // of operation as PUT, but add an assert/check
        // that at last we do not get any deleted element
        // becuase here we are re-shuffling stuff

        // returning Value field of pair
        return lru.Get(key)
	}
	return -1
}

func (l *MGLRUCache) Put(pair dll.Pair) int {
    if _, present := l.keys[pair.Key]; present {
        log.Println("Key already exists in LRU")
        return -1
    }

	// Make a loop over LRUs and check if current LRU
	// is full, if it is not, pop the last item,
	// return it and then insert it OR just insert it
	// in LRU
	i := 0

	var emptyElement *dll.Pair

	element := pair

	for i < l.length {
		element = *l.lrus[i].Put(element)
		if &element == emptyElement {
			// Making this empty because if we find an emptyElement
			// which signifies that Put operation was successful
			// without evicting any other elements in the current
			// LRU, then we need to clear any previous values stored
			// in this variable
			element = *emptyElement

			break
		}

		i++
	}

	// Insert the key into set of keys for quick
	// search in Get function
	l.keys[pair.Key] = i

	log.Println("Inserted element: ", pair)

	if &element != emptyElement {
		log.Println("Evicted element: ", element)
        delete(l.keys, element.Key)
	}

	return pair.Key
}
