/*=================================== LICENSE =======================================

                                   Apache License
                             Version 2.0, January 2004
                          http://www.apache.org/licenses/

============================== BRIEF FILE DESCRIPTION ===============================

This file contains the definition of the `diskHandler` interface, which provides a
layer of abstraction over disk IO. Every disk-resource (files) have a key which
uniquely identifies it. The `diskHandler` interface allows to retrieve either a
reader or a writer for a given key. The write operation is temporary and must be
confiremd by calling the finalize function, which is a callback provided by
the `getWriter` method as the second return value.
==================================================================================*/

package persistence

import (
	"io"
)

type diskHandler interface {
	getWriter(key string) (io.Writer, func(), error)
	getReader(key string) (io.ByteReader, bool)
}
