package lru

import "sync"

type Entry struct {
	Key   string
	Value any
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

func (c *Cache) Put(key string, v any) {
	c.mu.Lock()
	defer c.mu.Unlock()
	// put this key exist, update index hot
	if e, has := c.cache[key]; has {
		e.Key = key
		e.Value = v
		c.moveToHead(e)
		return
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

func (c *Cache) Get(key string) (v any) {
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
	// 将 e 自身从链表中当前位置摘除
	// e.pre link e.next(entry or nil)
	e.pre.next = e.next
	if e == c.tail {
		// if e is tail, c.tail = e.pre
		c.tail = e.pre
	}

	// 将 e 自身放入链表头部
	c.head.pre = e  // 原头部节点的前指针指向新节点
	c.head = e      // 原头部节点替换成新节点
	e.next = c.head // 新节点的后指针指向头部节点
	e.pre = nil     // 新节点的前指针置空
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

func (c *Cache) removeTail() {
	removeEntry := c.tail
	c.tail = c.tail.pre
	removeEntry.pre = nil
	c.tail.next = nil
	delete(c.cache, removeEntry.Key)
}
