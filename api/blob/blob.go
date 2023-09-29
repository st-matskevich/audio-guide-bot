package blob

import "io"

type BlobProvider interface {
	ReadBlob(name string, writer io.Writer) error
	WriteBlob(name string, reader io.Reader) error
}
