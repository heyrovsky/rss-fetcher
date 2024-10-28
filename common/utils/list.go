// Package utils provides utility functions for string slice operations.
package utils

import (
	"fmt"
	"strings"
)

var ErrNilInput = fmt.Errorf("nil input slice provided")

type StringSet map[string]struct{}

func newStringSet(capacity int) StringSet {
	return make(StringSet, capacity)
}

func (s StringSet) add(str string) {
	if trimmed := strings.TrimSpace(str); trimmed != "" {
		s[trimmed] = struct{}{}
	}
}

func (s StringSet) toSlice() []string {
	if len(s) == 0 {
		return []string{}
	}

	result := make([]string, 0, len(s))
	for key := range s {
		result = append(result, key)
	}
	return result
}

func ListUnique(list1, list2 []string) ([]string, error) {
	if list1 == nil && list2 == nil {
		return nil, ErrNilInput
	}

	uniqueSet := newStringSet(len(list1) + len(list2))

	lists := [][]string{list1, list2}
	for _, list := range lists {
		for _, item := range list {
			uniqueSet.add(item)
		}
	}

	return uniqueSet.toSlice(), nil
}

func ListUniqueList2(list1, list2 []string) ([]string, error) {
	if list1 == nil && list2 == nil {
		return nil, ErrNilInput
	}

	existingItems := newStringSet(len(list1))
	for _, item := range list1 {
		existingItems.add(item)
	}

	uniqueItems := newStringSet(len(list2))
	for _, item := range list2 {
		if trimmed := strings.TrimSpace(item); trimmed != "" {
			if _, exists := existingItems[trimmed]; !exists {
				uniqueItems[trimmed] = struct{}{}
			}
		}
	}

	return uniqueItems.toSlice(), nil
}
