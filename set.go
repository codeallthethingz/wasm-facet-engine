package main

import (
	"encoding/json"
	"sort"
)

// Set simple set type
type Set struct {
	list map[string]struct{}
}

// Len return the length of the set
func (s *Set) Len() int {
	return len(s.list)
}

// Contains does this value exist in the set.
func (s *Set) Contains(v string) bool {
	_, ok := s.list[v]
	return ok
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

// MarshalJSON Set is optimized for memory usage and lookup but is really a list of unique elements.
func (s *Set) MarshalJSON() ([]byte, error) {
	sorted := s.ToArray()
	sort.Strings(sorted)
	return json.Marshal(sorted)
}

// UnmarshalJSON Set is optimized for memory usage and lookup but is really a list of unique elements.
func (s *Set) UnmarshalJSON(j []byte) error {
	var list []string
	err := json.Unmarshal(j, &list)
	if err != nil {
		return err
	}
	s.list = map[string]struct{}{}
	for _, item := range list {
		s.Add(item)
	}
	return nil
}
