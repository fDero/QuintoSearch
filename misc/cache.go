package misc

type LRUCache[T any] struct {
	storage  map[string]*cacheBucket[T]
	mru      *cacheBucket[T]
	lru      *cacheBucket[T]
	maxItems int
}

type cacheBucket[T any] struct {
	item T
	key  string
	next *cacheBucket[T]
	prev *cacheBucket[T]
}

func NewLRUCache[T any](maxItems int) *LRUCache[T] {
	return &LRUCache[T]{
		storage:  make(map[string]*cacheBucket[T]),
		mru:      nil,
		lru:      nil,
		maxItems: maxItems,
	}
}

func newCacheBucket[T any](key string, item T) *cacheBucket[T] {
	return &cacheBucket[T]{
		item: item,
		key:  key,
		next: nil,
		prev: nil,
	}
}

func (node *cacheBucket[T]) detach() {
	if node.prev != nil {
		node.prev.next = node.next
	}
	if node.next != nil {
		node.next.prev = node.prev
	}
	node.next = nil
	node.prev = nil
}

func (cache *LRUCache[T]) lookup(key string) *cacheBucket[T] {
	bucket, exists := cache.storage[key]
	if !exists {
		return nil
	}
	if cache.lru == bucket {
		cache.lru = cache.lru.next
	}
	bucket.detach()
	if len(cache.storage) > 1 {
		bucket.prev = cache.mru
		cache.mru.next = bucket
	}
	cache.mru = bucket
	return cache.mru
}

func (cache *LRUCache[T]) Retrieve(key string) *T {
	tmp := cache.lookup(key)
	if tmp == nil {
		return nil
	}
	return &tmp.item
}

func (cache *LRUCache[T]) evict() {
	for len(cache.storage) >= cache.maxItems {
		cache.lru = cache.lru.next
		delete(cache.storage, cache.lru.prev.key)
		cache.lru.prev = nil
	}
}

func (cache *LRUCache[T]) Store(key string, item T) {
	if node := cache.lookup(key); node != nil {
		(*node).item = item
		return
	}
	newNode := newCacheBucket(key, item)
	cache.evict()
	cache.storage[key] = newNode
	if len(cache.storage) == 1 {
		cache.mru = newNode
		cache.lru = newNode
		return
	}
	newNode.prev = cache.mru
	cache.mru.next = newNode
	cache.mru = newNode
}
