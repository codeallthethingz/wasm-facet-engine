package engine

// Set simple set type
type Set struct {
	list map[string]struct{} //empty structs occupy 0 memory
}

// Add a value
func (s *Set) Add(v string) {
	s.list[v] = struct{}{}
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
