package array

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExceptItems(t *testing.T) {
	type input struct {
		parent []interface{}
		child  []interface{}
	}

	type expected struct {
		result []interface{}
	}

	testCases := []struct {
		name     string
		input    input
		expected expected
	}{
		{
			name: "Succeed",
			input: input{
				parent: []interface{}{1, 2, 3, 5, 6},
				child:  []interface{}{2, 3, 4, 5},
			},
			expected: expected{
				result: []interface{}{1, 6},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			result := ExceptItems(tc.input.parent, tc.input.child)

			require.Equal(t, tc.expected.result, result)
		})
	}
}

func TestRemoveDuplicatesString(t *testing.T) {
	testCases := []struct {
		name     string
		array    []string
		expected []string
	}{
		{
			name:     "with_empty_array",
			array:    []string{},
			expected: []string{},
		},
		{
			name:     "with_array_has_not_dup_value",
			array:    []string{"a", "b", "c", "d"},
			expected: []string{"a", "b", "c", "d"},
		},
		{
			name:     "with_array_has_dup_value",
			array:    []string{"a", "b", "c", "d", "a"},
			expected: []string{"a", "b", "c", "d"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := RemoveDuplicatesString(tc.array)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestGetItemsByLen(t *testing.T) {
	testCases := []struct {
		name     string
		array    []string
		len      int
		expected map[int]interface{}
	}{
		{
			name:     "with_empty_array",
			array:    []string{},
			expected: map[int]interface{}{},
		},
		{
			name:     "get_slice_shorter_than_array",
			array:    []string{"a", "b", "c", "d"},
			len:      3,
			expected: map[int]interface{}{0: "a", 1: "b", 2: "c"},
		},
		{
			name:     "get_slice_longer_than_array",
			array:    []string{"a", "b", "c", "d", "a"},
			len:      7,
			expected: map[int]interface{}{0: "a", 1: "b", 2: "c", 3: "d", 4: "a"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := GetItemsByLen(tc.array, tc.len)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestItemExists(t *testing.T) {
	testCases := []struct {
		name     string
		array    interface{}
		item     interface{}
		expected bool
	}{
		{
			name:     "with_array_string_and_item_exist",
			array:    []string{"a", "b", "c"},
			item:     "a",
			expected: true,
		},
		{
			name:     "with_array_string_and_item_not_exist",
			array:    []string{"a", "b", "c"},
			item:     "d",
			expected: false,
		},
		{
			name:     "with_array_int_and_item_exist",
			array:    []int{1, 2, 3, 4},
			item:     1,
			expected: true,
		},
		{
			name:     "with_array_int_and_item_not_exist",
			array:    []int{1, 2, 3, 4},
			item:     5,
			expected: false,
		},
		{
			name: "with_array_struct_and_item_exist",
			array: []struct{ test int }{
				{test: 1},
				{test: 2},
				{test: 3},
			},
			item: struct{ test int }{
				test: 1,
			},
			expected: true,
		},
		{
			name: "with_array_struct_and_item_not_exist",
			array: []struct{ test int }{
				{test: 1},
				{test: 2},
				{test: 3},
			},
			item: struct{ test int }{
				test: 5,
			},
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := ItemExists(tc.array, tc.item)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestIntersection(t *testing.T) {
	type input struct {
		f []interface{}
		s []interface{}
	}

	type expected struct {
		result bool
		vals   []interface{}
	}

	testCases := []struct {
		name     string
		input    input
		expected expected
	}{
		{
			name: "succeed",
			input: input{
				f: []interface{}{1, 2, 3},
				s: []interface{}{3, 4, 5},
			},
			expected: expected{
				result: true,
				vals:   []interface{}{3},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			result, vals := Intersection(tc.input.f, tc.input.s)

			require.Equal(t, tc.expected.result, result)
			require.Equal(t, tc.expected.vals, vals)
		})
	}
}
