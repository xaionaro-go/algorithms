package solution

type Node struct {
	Index int
	Key int
	Value int
	UseCount int
}

type Nodes []*Node

func (s Nodes) Get(idx int) **Node {
	return &s[idx]
}

func (s Nodes) LeftOf(idx int) **Node {
	return &s[idx*2]
}

func (s Nodes) RightOf(idx int) **Node {
	return &s[idx*2+1]
}

func (s Nodes) MoveDownIfRequired(idx int) {
	parent := s.Get(idx)
	left := s.LeftOf(idx)
	right := s.RightOf(idx)

	child := left
	if (*child).UseCount > (*right).UseCount {
		child = right
	}

	if (*parent).UseCount > (*child).UseCount {
		*parent, *child = *child, *parent
		(*parent).Index, (*child).Index = (*child).Index, (*parent).Index
	}
}

func (s Nodes) MoveToTop(idx int) {

}

func (s Nodes) GetIndex(node *Node) int {
	return node.Index
}

type LRUCache struct {
	Count int
	Storage []Node
	Heap Nodes
	Map map[int]*Node
}

func New() *LRUCache {
	return &LRUCache{
		Map: map[int]*Node{},
	}
}

func (cache *LRUCache) SetCapacity(newCapacity int) {
	cache.Storage = make([]Node, newCapacity)
	cache.Heap = make(Nodes, newCapacity)
	for idx := range cache.Storage {
		cache.Heap[idx] = &cache.Storage[idx]
		cache.Heap[idx].Index = idx
	}
}

func (cache *LRUCache) Get(key int) int {
	v := cache.Map[key]
	if v == nil {
		return -1
	}
	v.UseCount++
	idx := cache.Heap.GetIndex(v)
	cache.Heap.MoveDownIfRequired(idx)
	return v.Value
}


func (cache *LRUCache) Put(key, value int) {
	v := cache.Map[key]
	if v != nil {
		v.UseCount = 0
		cache.Heap.MoveToTop(cache.Heap.GetIndex(v))
		v.Value = value
		return
	}
	if cache.Count < len(cache.Storage) {

	}
}