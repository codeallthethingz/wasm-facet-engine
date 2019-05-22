package engine

// Set simple set type
type Set struct {
	list map[string]struct{} //empty structs occupy 0 memory
}

// Has a value
func (s *Set) Has(v string) bool {
	_, ok := s.list[v]
	return ok
}

// Add a value
func (s *Set) Add(v string) {
	s.list[v] = struct{}{}
}

// Remove a value
func (s *Set) Remove(v string) {
	delete(s.list, v)
}

// Clear all values
func (s *Set) Clear() {
	s.list = make(map[string]struct{})
}

// Size of set
func (s *Set) Size() int {
	return len(s.list)
}

// ToArray return the values
func (s *Set) ToArray() []string {
	keys := make([]string, len(s.list))

	i := 0
	for k := range s.list {
		keys[i] = k
		i++
	}
	return keys
}

// NewSet new set.
func NewSet() *Set {
	s := &Set{}
	s.list = make(map[string]struct{})
	return s
}
