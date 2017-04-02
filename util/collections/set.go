package collections

// Set is a collection of unique items.
type Set struct {
	items map[interface{}]bool
}

// Items returns set's items
func (s *Set) Items() []interface{} {
	keys := make([]interface{}, len(s.items))

	i := 0
	for k := range s.items {
		keys[i] = k
		i++
	}
	return keys
}

// Add inserts an item.
func (s *Set) Add(item interface{}) {
	if s.items == nil {
		s.items = map[interface{}]bool{}
	}
	if s.items[item] {
		return
	}
	s.items[item] = true
}

// Delete removes an item.
func (s *Set) Delete(key interface{}) {
	if s.items != nil {
		delete(s.items, key)
	}
}
