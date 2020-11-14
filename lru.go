package lru

import (
	"container/list"
	"sync"
)

// cache is an LRU cache. It is not safe for concurrent access.
type cache struct {
	maxEntries int
	ll         *list.List
	cache      map[interface{}]*list.Element
}

type Key interface{}
type Value interface{}

type entry struct {
	key   Key
	value Value
}

func (c *cache) set(key Key, value Value) {
	if c.cache == nil {
		c.cache = make(map[interface{}]*list.Element)
		c.ll = list.New()
	}
	if ee, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ee)
		ee.Value.(*entry).value = value
		return
	}
	ele := c.ll.PushFront(&entry{key, value})
	c.cache[key] = ele
	if c.maxEntries != 0 && c.ll.Len() > c.maxEntries {
		c.removeOldest()
	}
}

func (c *cache) get(key Key) (value Value, ok bool) {
	if c.cache == nil {
		return
	}
	if ele, hit := c.cache[key]; hit {
		c.ll.MoveToFront(ele)
		return ele.Value.(*entry).value, true
	}
	return
}

func (c *cache) del(key Key) {
	if c.cache == nil {
		return
	}
	if ele, hit := c.cache[key]; hit {
		c.removeElement(ele)
	}
}

func (c *cache) removeOldest() {
	if c.cache == nil {
		return
	}
	ele := c.ll.Back()
	if ele != nil {
		c.removeElement(ele)
	}
}

func (c *cache) removeElement(e *list.Element) {
	c.ll.Remove(e)
	kv := e.Value.(*entry)
	delete(c.cache, kv.key)
}

func (c *cache) len() int {
	if c.cache == nil {
		return 0
	}
	return c.ll.Len()
}

// Cache is an LRU cache. It is safe for concurrent access.
type Cache struct {
	inner cache
	mutex sync.Mutex
}

// Set adds a value to the cache.
func (c *Cache) Set(key Key, value Value) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.inner.set(key, value)
}

// Get looks up a key's value from the cache.
func (c *Cache) Get(key Key) (Value, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.inner.get(key)
}

// Del removes the provided key from the cache.
func (c *Cache) Del(key Key) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.inner.del(key)
}

// Len returns the number of items in the cache.
func (c *Cache) Len() int {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.inner.len()
}

// New creates a new Cache.
// If maxEntries is zero, the cache has no limit.
func New(maxEntries int) *Cache {
	return &Cache{
		inner: cache{
			maxEntries: maxEntries,
			ll:         list.New(),
			cache:      make(map[interface{}]*list.Element),
		},
		mutex: sync.Mutex{},
	}
}

// Alias for New.
func Lru(maxEntries int) *Cache {
	return New(maxEntries)
}
