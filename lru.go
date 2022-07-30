// Package lru implements an LRU cache.
package lru

import (
	"container/list"
	"sync"
)

// Nut is an LRU cache. It is not safe for concurrent access.
type Nut struct {
	// MaxEntries is the maximum number of cache entries before
	// an item is evicted. Zero means no limit.
	MaxEntries int
	ll         *list.List
	cache      map[any]*list.Element
}

// NewNut creates a new Cache. If maxEntries is zero, the cache has no limit.
func NewNut(maxEntries int) *Nut {
	return &Nut{
		MaxEntries: maxEntries,
		ll:         list.New(),
		cache:      make(map[any]*list.Element),
	}
}

// A Key may be any value that is comparable. See http://golang.org/ref/spec#Comparison_operators
type Key any

type entry struct {
	key   Key
	value any
}

// Set adds a value to the cache.
func (c *Nut) Set(key Key, value any) {
	if c.cache == nil {
		c.cache = make(map[any]*list.Element)
		c.ll = list.New()
	}
	if ee, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ee)
		ee.Value.(*entry).value = value
		return
	}
	ele := c.ll.PushFront(&entry{key, value})
	c.cache[key] = ele
	if c.MaxEntries != 0 && c.ll.Len() > c.MaxEntries {
		c.removeOldest()
	}
}

// Get looks up a key's value from the cache.
func (c *Nut) Get(key Key) (value any, ok bool) {
	if c.cache == nil {
		return
	}
	if ele, hit := c.cache[key]; hit {
		c.ll.MoveToFront(ele)
		return ele.Value.(*entry).value, true
	}
	return
}

// Del removes the provided key from the cache.
func (c *Nut) Del(key Key) {
	if c.cache == nil {
		return
	}
	if ele, hit := c.cache[key]; hit {
		c.removeElement(ele)
	}
}

func (c *Nut) removeOldest() {
	if c.cache == nil {
		return
	}
	ele := c.ll.Back()
	if ele != nil {
		c.removeElement(ele)
	}
}

func (c *Nut) removeElement(e *list.Element) {
	c.ll.Remove(e)
	kv := e.Value.(*entry)
	delete(c.cache, kv.key)
}

// Len returns the number of items in the cache.
func (c *Nut) Len() int {
	if c.cache == nil {
		return 0
	}
	return c.ll.Len()
}

// Clr all.
func (c *Nut) Clr() {
	c.cache = nil
}

// Lru is an LRU cache.
type Lru struct {
	// MaxEntries is the maximum number of cache entries before
	// an item is evicted. Zero means no limit.
	MaxEntries int
	inner      *Nut
	mutex      sync.Mutex
}

// Set adds a value to the cache.
func (c *Lru) Set(key Key, value any) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.inner.Set(key, value)
}

// Get looks up a key's value from the cache.
func (c *Lru) Get(key Key) (any, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.inner.Get(key)
}

// Del removes the provided key from the cache.
func (c *Lru) Del(key Key) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.inner.Del(key)
}

// Len returns the number of items in the cache.
func (c *Lru) Len() int {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.inner.Len()
}

// Clr all
func (c *Lru) Clr() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.inner.Clr()
}

// NewLru creates a new Cache. If maxEntries is zero, the cache has no limit.
func NewLru(maxEntries int) *Lru {
	return &Lru{
		MaxEntries: maxEntries,
		inner:      NewNut(maxEntries),
		mutex:      sync.Mutex{},
	}
}
