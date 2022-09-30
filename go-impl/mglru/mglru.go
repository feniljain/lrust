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

func Init(length, capacity int) MGLRUCache {
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
		keys:     make(map[int]int, length*capacity),
	}
}

func (l *MGLRUCache) Get(key int) int {

	l.MoveNodeToFront(key)

	if lruIdx, present := l.keys[key]; present {
		lru := l.lrus[lruIdx]

		// TODO: Promote received pair from the below GET
		// to highest LRU, with of course the same details
		// of operation as PUT, but add an assert/check
		// that at last we do not get any deleted element
		// becuase here we are re-shuffling stuff
		// - Also don't forget to remove the same from current
		// LRU

		// returning Value field of pair
		return lru.Get(key)
	}
	return -1
}

func (l *MGLRUCache) Put(pair dll.Pair) *dll.Pair {
	if _, present := l.keys[pair.Key]; present {
		// TODO: Instead of doing this, promote the key
		// to the topmost position in youngest LRU
		//
		// This will require deleting key from the given LRU,
		// then readjusting all the entries after that according
		// to capacity of current LRU and LRUs to come after
		//
		// Next step is taking this deleted key and inserting
		// it into top most ( youngest ) LRU at top most position

		log.Println("Key already exists in LRU")

		pair := l.MoveNodeToFront(pair.Key)

		return pair
	}

	return l.InsertElement(pair)
}

func (l *MGLRUCache) InsertElement(pair dll.Pair) *dll.Pair {

	// Make a loop over LRUs and check if current LRU
	// is full, if it is not, pop the last item,
	// return it and then insert it OR just insert it
	// in LRU
	i := 0

	var emptyElement *dll.Pair

	element := &pair

	for i < l.length {
		element = l.lrus[i].Put(*element)
		if element == emptyElement {
			// Making this empty because if we find an emptyElement
			// which signifies that Put operation was successful
			// without evicting any other elements in the current
			// LRU, then we need to clear any previous values stored
			// in this variable
			element = emptyElement

			break
		}

		l.keys[element.Key] = i+1

		i++
	}

	// Insert the key into set of keys for quick
	// search in Get function
	l.keys[pair.Key] = 0

	if element != emptyElement {
		// log.Println("Evicted element: ", element)
		delete(l.keys, element.Key)
		return element
	}

	return nil
}

func (l *MGLRUCache) MoveNodeToFront(key int) *dll.Pair {

	// If key already exists in MGLRU, remove it from
	// the holding cache and push it to front using
	// `InsertElement`. Also update map to indicate
	// it has been promoted to youngest generation
	// i.e. LRU at postition 0

	// If not present, return nil

	if lruIdx, present := l.keys[key]; present {
		if deletedPair, present := l.lrus[lruIdx].Remove(key); present {

			l.InsertElement(*deletedPair)

			l.keys[key] = 0

			return deletedPair
		}
	}

	return nil
}

// I know I should have written unit tests instead of integration tests,
// so that I don't have to expose all these details, but it's fine

// For testing puposes
func (l *MGLRUCache) GetFirstElement() dll.Pair {
    return l.lrus[0].GetFirstElement().Data
}

// For debugging puposes
func (l *MGLRUCache) PrintAllLRUs() {
	i := 0

	for i < l.length {
		l.lrus[i].PrintList()
		i++
	}
}

// For debugging puposes
func (l *MGLRUCache) PrintState() {
    log.Println("Map: ", l.keys)
    l.PrintAllLRUs()
}

