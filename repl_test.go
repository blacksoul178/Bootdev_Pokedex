package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "basic",
			input:    "hello world",
			expected: []string{"hello", "world"},
		},
		{
			name:     "medium",
			input:    "this is a test",
			expected: []string{"this", "is", "a", "test"},
		},
		{name: "caps",
			input:    "Lets Make things Funky",
			expected: []string{"lets", "make", "things", "funky"},
		},
		{
			name:     "extra spaces",
			input:    " this  is gonna     be   weird   ",
			expected: []string{"this", "is", "gonna", "be", "weird"},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := cleanInput(c.input)

			if len(actual) != len(c.expected) {
				t.Errorf("Actual length does not match Expected length; expecting %d, got %v", len(c.expected), len(actual))
				return
			}

			for i := range actual {
				word := actual[i]
				expectedWord := c.expected[i]
				if word != expectedWord {
					t.Errorf("Expecting %s, got %s", expectedWord, word)
				}
			}
			t.Logf("passed %q -> %v", c.input, actual)
		})
	}

}
