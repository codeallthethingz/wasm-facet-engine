package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInitialize(t *testing.T) {
	facetGroups, err := initializeObjects("{}", "[]")
	if err != nil {
		panic(err)
	}
	require.Equal(t, "{}", facetGroups)
}
func TestQuery(t *testing.T) {
	initializeObjects("{}", "[]")
	ids, facetGroups, err := query()
	require.Nil(t, err)
	require.Equal(t, "[]", ids)
	require.Equal(t, "{}", facetGroups)
}
func TestFilter(t *testing.T) {
	initializeObjects("{}", "[]")
	err := addFilter("facetGroupName", "facetName", true, 0, true, 10)
	require.Nil(t, err)
}
func TestFilterError(t *testing.T) {
	initializeObjects("{}", "[]")
	err := addFilter(" ", "facetName", true, 0, true, 10)
	require.Error(t, err)
	err = addFilter("group", " ", true, 0, true, 10)
	require.Error(t, err)
}
func TestClearFilter(t *testing.T) {
	initializeObjects("{}", "[]")
	JSClearFilters(nil)
}
