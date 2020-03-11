package parser

import "strings"

type byKey []string

// Len length of array
func (a byKey) Len() int {
	return len(a)
}

// Swap 2 array positions
func (a byKey) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

// Less checks which position is less then another
func (a byKey) Less(i, j int) bool {
	p1 := strings.Index(a[i], "=") - 1
	if p1 < 0 {
		p1 = len(a[i])
	}

	p2 := strings.Index(a[j], "=") - 1
	if p2 < 0 {
		p2 = len(a[j])
	}

	return strings.Compare(a[i][:p1+1], a[j][:p2+1]) < 0
}
