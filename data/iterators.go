/*=================================== LICENSE =======================================

                                   Apache License
                             Version 2.0, January 2004
                          http://www.apache.org/licenses/

============================== BRIEF FILE DESCRIPTION ===============================

This files contains some slice-related utilities that are often needed in the codebase.
Such utilities are not specific to a particular package, but are used in multiple
packages. Most of them are simple wrappers/adapters around the "iter" package.
==================================================================================*/

package data

import (
	"bufio"
	"iter"
	"os"
	"strings"
)

func NewSliceIterator[T any](slice []T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, value := range slice {
			if !yield(value) {
				break
			}
		}
	}
}

func NewStringIterator(inlineText string) iter.Seq[string] {
	return func(yield func(string) bool) {
		for _, value := range strings.Fields(inlineText) {
			if !yield(value) {
				break
			}
		}
	}
}

func NewFileReaderIterator(file *os.File) iter.Seq[string] {
	return func(yield func(string) bool) {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			if !yield(scanner.Text()) {
				break
			}
		}
	}
}
