package follow_along

type Set struct {
	items map[interface{}]bool
}

func (s *Set) Items() []interface{} {
	keys := make([]interface{}, len(s.items))

	i := 0
	for k := range s.items {
		keys[i] = k
		i++
	}
	return keys
}

func (s *Set) add(item interface{}) {
	if s.items == nil {
		s.items = map[interface{}]bool{}
	}
	if s.items[item] {
		return
	}
	s.items[item] = true
}
