package chaining

type ListElement struct {
	key   string
	value string
	next  *ListElement
}

type ChainingHashMap struct {
	nodes int
	count int
	data  []*ListElement
}

func New() *ChainingHashMap {
	c := ChainingHashMap{
		nodes: 32,
		count: 0,
		data:  make([]*ListElement, 32),
	}

	for i := 0; i < c.nodes; i++ {
		c.data[i] = nil
	}

	return &c
}

func (c *ChainingHashMap) Insert(key, value string) {
	idx := int(hash(key)) % c.nodes

	if c.data[idx] == nil {
		c.data[idx] = &ListElement{
			key:   key,
			value: value,
			next:  nil,
		}
		return
	}

	curr := c.data[idx]
	for curr.next != nil {
		if curr.key == key {
			curr.value = value
			return
		}
		curr = curr.next
	}
	curr.next = &ListElement{
		key:   key,
		value: value,
		next:  nil,
	}
}

func (c *ChainingHashMap) Retrieve(key string) (string, bool) {
	idx := int(hash(key)) % c.nodes

	if c.data[idx] == nil {
		return "", false
	}

	curr := c.data[idx]
	for curr != nil {
		if curr.key == key {
			return curr.value, true
		}
		curr = curr.next
	}

	return "", false
}

func hash(s string) uint32 {
	var val uint32 = 31
	for pos, char := range s {
		val = val * (1 + (uint32(pos) * uint32(char)))
	}

	return val
}
