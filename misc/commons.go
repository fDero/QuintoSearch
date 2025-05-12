/*=================================== LICENSE =======================================

                                   Apache License
                             Version 2.0, January 2004
                          http://www.apache.org/licenses/

============================== BRIEF FILE DESCRIPTION ===============================

This files contains some utilities that are often needed in the codebase. Such
utilities are not specific to a particular package, but are used in multiple packages.
==================================================================================*/

package misc

import (
	"iter"
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
