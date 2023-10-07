package blob

import "io"

type StatBlobResult struct {
	Size int64
}

type BlobRange struct {
	Start int64
	End   int64
}

type ReadBlobOptions struct {
	Range *BlobRange
}

type BlobProvider interface {
	ReadBlob(name string, options ReadBlobOptions) (io.ReadCloser, error)
	WriteBlob(name string, reader io.Reader) error
	StatBlob(name string) (StatBlobResult, error)
}
