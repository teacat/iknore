package file

import (
	"context"
	"io"
	"os"
)

// Handler
type Handler interface {
	Upload(c context.Context, path string, b []byte) error
	Download(c context.Context, path string) ([]byte, error)
	Head(c context.Context, path string) (os.FileInfo, error)
	Copy(c context.Context, src, dest string) error
	Delete(c context.Context, path string) error
	StreamUpload(c context.Context, path string, reader io.Reader) error
	StreamDownload(c context.Context, path string, writer io.Writer) error
	IsNotExist(err error) bool
}
