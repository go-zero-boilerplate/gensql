package main

import "strings"

type additionalImportsAppender struct {
	EntityStruct   []string
	EntityIterator []string
}

func (a *additionalImportsAppender) appendToSlice(slice []string, importPath string) []string {
	if strings.TrimSpace(importPath) == "" {
		return slice
	}

	for _, p := range slice {
		if p == importPath {
			return slice
		}
	}
	return append(slice, importPath)
}

func (a *additionalImportsAppender) AddEntityIteratorPath(importPath string) {
	a.EntityIterator = a.appendToSlice(a.EntityIterator, importPath)
}

func (a *additionalImportsAppender) AddEntityStructPath(importPath string) {
	a.EntityStruct = a.appendToSlice(a.EntityStruct, importPath)
}

func (a *additionalImportsAppender) MergeIntoSelf(other *additionalImportsAppender) {
	for _, otherPath := range other.EntityStruct {
		a.AddEntityStructPath(otherPath)
	}
	for _, otherPath := range other.EntityIterator {
		a.AddEntityIteratorPath(otherPath)
	}
}
