package lru

import (
	"container/list"
	"fmt"
	"time"
)

// Lru :an LRU cache
type Lru struct {
	capacity     int64
	currentCount int64
	Ll           *list.List
	cache        map[string]*list.Element
	ttl          int64
	// [optional] executed when an kvMap is purged
	OnRemove func(key string, value Value)
}

type KvMap struct {
	Key   string
	Value Value
	exp   int64
}

// Value use Len to count how many bytes it takes
// 		course Value's type is interface{} which can't use len()
type Value interface {
	Len() int
}

// ------------------------------------------------------------------------------------------------

// LlLen the number of cache entries
func (lru *Lru) LlLen() int {
	return lru.Ll.Len()
}

// InitLru Init InitLru is the Constructor of Cache
func InitLru(c int64, ttl int64, remover func(string, Value)) *Lru {
	return &Lru{
		capacity: c,
		Ll:       list.New(),
		ttl:      ttl,
		cache:    make(map[string]*list.Element),
		OnRemove: remover,
	}
}

func (lru *Lru) Set(key string, val Value) {
	newExp := lru.getExpTime()

	if pa, ok := lru.cache[key]; ok {
		lru.Ll.MoveToFront(pa)
		kv := pa.Value.(*KvMap)

		kv.Value = val
		kv.exp = newExp
	} else {
		pa := lru.Ll.PushFront(&KvMap{key, val, newExp})

		lru.cache[key] = pa
		lru.currentCount += 1
	}

	lru.checkSize()
}

func (lru *Lru) getExpTime() int64 {
	if lru.ttl == 0 {
		return 0
	}

	return lru.ttl + time.Now().Unix()
}

func (lru *Lru) checkSize() {
	for lru.capacity > 0 && lru.currentCount > lru.capacity {
		lru.lruRemove()
	}
}

func (lru *Lru) lruRemove() {
	fmt.Println("Trigger lruRemove")

	if lru.LlLen() == 0 {
		return
	}

	// Back returns the last element of list l or nil if the list is empty.
	pa := lru.Ll.Back()

	if pa == nil {
		return
	}

	lru.Ll.Remove(pa)
	kv := pa.Value.(*KvMap)
	delete(lru.cache, kv.Key)
	lru.currentCount -= 1

	if lru.OnRemove != nil {
		lru.OnRemove(kv.Key, kv.Value)
	}
}

// Get look ups a key's value
func (lru *Lru) Get(key string) (value Value, ok bool) {
	if pa, ok := lru.cache[key]; ok {
		lru.Ll.MoveToFront(pa)
		kv := pa.Value.(*KvMap)

		// check expire
		if kv.exp > 0 && kv.exp < time.Now().Unix() {
			lru.onExpire(pa)

			return nil, false
		}

		return kv.Value, true
	}

	return
}

func (lru *Lru) onExpire(pa *list.Element) {
	fmt.Println("Trigger onExpire")
	lru.Ll.Remove(pa)

	kv := pa.Value.(*KvMap)

	delete(lru.cache, kv.Key)

	lru.currentCount -= 1
}
