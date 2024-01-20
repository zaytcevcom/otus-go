package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	item, wasInCache := c.items[key]

	if wasInCache {
		item.Value = value
		c.queue.MoveToFront(item)
	} else {
		item = &ListItem{Value: value}
		c.items[key] = c.queue.PushFront(item.Value)

		if c.queue.Len() > c.capacity {
			back := c.queue.Back()
			c.queue.Remove(back)

			// todo: это сложность O(n)?
			for k, val := range c.items {
				if val == back {
					delete(c.items, k)
					break
				}
			}
		}
	}

	return wasInCache
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	item, ok := c.items[key]

	if ok {
		c.queue.MoveToFront(item)
		return item.Value, true
	}

	return nil, false
}

func (c *lruCache) Clear() {
	c.items = make(map[Key]*ListItem)
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
