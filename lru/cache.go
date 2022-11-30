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

func (c *Cache) PutHead(key string, v interface{}) {
	c.Put(key, v)
}

// todo
func (c *Cache) PutTail(key string, v interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	// put this key exist, update index
	if e, has := c.cache[key]; has {
		c.moveToTail(e)
	}
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
	c.putNewEntryHead(entry)
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
	// e.pre link e.next(entry or nil)
	e.pre.next = e.next
	if e == c.tail {
		// if e is tail, c.tail = e.pre
		c.tail = e.pre
	}
	e.next = c.head
	c.head.pre = e
	c.head = e
	e.pre = nil
}

// todo
func (c *Cache) moveToTail(e *Entry) {
	//if e == c.tail {
	//	// already in tail, return
	//	return
	//}
	//// e.pre link e.next(entry or nil)
	//e.pre.next = e.next
	//if e == c.head {
	//	// if e is head, c.head = e.next
	//	c.head = e.next
	//}
	//c.tail.next = e
	//e.pre = c.tail
	//c.tail = e
	//e.next = nil
}

func (c *Cache) putNewEntryHead(e *Entry) {
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

// todo
func (c *Cache) putNewEntryTail(e *Entry) {

}

func (c *Cache) removeTail() {
	removeEntry := c.tail
	c.tail = c.tail.pre
	removeEntry.pre = nil
	c.tail.next = nil
	delete(c.cache, removeEntry.Key)
}

// todo
func (c *Cache) removeHead() {

}
