package ui

import (
	"fmt"
	"github.com/ktr0731/go-fuzzyfinder"
)

// SelectString pops up a fuzzy finder for a simple string slice
func SelectString(label string, items []string) (string, error) {
	idx, err := fuzzyfinder.Find(
		items,
		func(i int) string {
			return items[i]
		},
		fuzzyfinder.WithPromptString(label + " > "),
	)
	if err != nil {
		return "", err
	}
	return items[idx], nil
}
