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

func ZipIterators[T, U any](firstIterator iter.Seq[T], secondIterator iter.Seq[U]) iter.Seq2[T, U] {
	nextFirst, stopFirst := iter.Pull(firstIterator)
	nextSecond, stopSecond := iter.Pull(secondIterator)
	return func(yield func(T, U) bool) {
		firstElem, firstExists := nextFirst()
		secondElem, secondExists := nextSecond()
		for firstExists && secondExists {
			if !yield(firstElem, secondElem) {
				break
			}
			firstElem, firstExists = nextFirst()
			secondElem, secondExists = nextSecond()
		}
		stopFirst()
		stopSecond()
	}
}

func CountIterations[T any](seq iter.Seq[T]) int {
	count := 0
	for range seq {
		count++
	}
	return count
}

func ZipSlices[T, U any](firstSlice []T, secondSlice []U) iter.Seq2[T, U] {
	return ZipIterators(
		NewSliceIterator(firstSlice),
		NewSliceIterator(secondSlice),
	)
}

func CollectAsSlice[T any](seq iter.Seq[T]) []T {
	var result []T
	for value := range seq {
		result = append(result, value)
	}
	return result
}

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
