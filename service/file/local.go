package file

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/teacat/pathx"
)

type Local struct {
	Directory string
}

func NewLocal(dir string) Handler {
	return &Local{dir}
}

// Upload
func (l *Local) Upload(c context.Context, path string, b []byte) error {
	fullPath := pathx.Join(l.Directory, path)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0644); err != nil {
		return err
	}
	err := os.WriteFile(fullPath, b, 0644)
	if err != nil {
		return err
	}
	return nil
}

// Download
func (l *Local) Download(c context.Context, path string) ([]byte, error) {
	fullPath := pathx.Join(l.Directory, path)
	b, err := os.ReadFile(fullPath)
	if err != nil {
		return []byte{}, err
	}
	return b, nil
}

// Copy
func (l *Local) Head(c context.Context, path string) (os.FileInfo, error) {
	fullPath := pathx.Join(l.Directory, path)
	f, err := os.Open(fullPath)
	if err != nil {
		return nil, err
	}
	if err := f.Close(); err != nil {
		return nil, err
	}
	return f.Stat()
}

// Copy
func (l *Local) Copy(c context.Context, src, dest string) error {
	fullSrc := pathx.Join(l.Directory, src)
	fullDest := pathx.Join(l.Directory, dest)
	//
	b, err := os.ReadFile(fullSrc)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(fullDest), 0644); err != nil {
		return err
	}
	if err := os.WriteFile(fullDest, b, 0644); err != nil {
		return err
	}
	return nil
}

// Delete
func (l *Local) Delete(c context.Context, path string) error {
	fullPath := pathx.Join(l.Directory, path)
	if err := os.Remove(fullPath); err != nil {
		return err
	}
	return nil
}

// StreamUpload
func (l *Local) StreamUpload(c context.Context, path string, reader io.Reader) error {
	fullPath := pathx.Join(l.Directory, path)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0644); err != nil {
		return err
	}
	f, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	if _, err = io.Copy(f, reader); err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return nil
}

// StreamDownload
func (l *Local) StreamDownload(c context.Context, path string, writer io.Writer) error {
	fullPath := pathx.Join(l.Directory, path)
	f, err := os.Open(fullPath)
	if err != nil {
		return err
	}
	if _, err = io.Copy(writer, f); err != nil {
		return err
	}
	return nil
}

// IsNotFound
func (l *Local) IsNotExist(err error) bool {
	if err == nil {
		return false
	}
	return os.IsNotExist(err)
}
