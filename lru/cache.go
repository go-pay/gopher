package lru

import "sync"

type Entry struct {
	Key   string
	Value interface{}
	pre   *Entry
	next  *Entry
}

type Cache struct {
	mu    sync.RWMutex
	cache map[string]*Entry
	cap   int
	head  *Entry
	tail  *Entry
}

func NewCache(cap int) *Cache {
	return &Cache{cache: make(map[string]*Entry), cap: cap}
}

func (c *Cache) Put(key string, v interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	// put this key exist, update index
	if e, has := c.cache[key]; has {
		c.moveToHead(e)
	}
	// generate Entry
	entry := &Entry{Key: key, Value: v, pre: nil, next: c.head}
	// put new Entry
	c.putNewEntry(entry)
	// not full
	if len(c.cache) <= c.cap {
		return
	}
	// full, remove tail entry
	c.removeTail()
}

func (c *Cache) Get(key string) (v interface{}) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if e, has := c.cache[key]; has {
		// link to head
		c.moveToHead(e)
		return e.Value
	}
	return nil
}

func (c *Cache) moveToHead(e *Entry) {
	if e == c.head {
		// already in front, return
		return
	}
	// cv pre link pre next(entry or nil)
	e.pre.next = e.next
	if e == c.tail {
		// cv is tail, tail = cv.pre
		c.tail = e.pre
	}
	e.next = c.head
	c.head.pre = e
	c.head = e
	e.pre = nil
}

func (c *Cache) putNewEntry(e *Entry) {
	// link to head
	if c.head != nil {
		c.head.pre = e
	}
	// reset head entry
	c.head = e
	if c.tail == nil {
		// tail == nil, stack is empty
		c.tail = e
	}
	c.cache[e.Key] = e
}

func (c *Cache) removeTail() {
	removeEntry := c.tail
	c.tail = c.tail.pre
	removeEntry.pre = nil
	c.tail.next = nil
	delete(c.cache, removeEntry.Key)
}
