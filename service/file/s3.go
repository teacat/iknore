package file

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type S3 struct {
	Session *session.Session
}

// NewS3
func NewS3() Handler {
	sess := session.Must(session.NewSession(&aws.Config{}))
	return &S3{sess}
}

// Upload
func (s *S3) Upload(c context.Context, path string, b []byte) error {
	_, err := s3manager.NewUploader(s.Session).Upload(&s3manager.UploadInput{
		Bucket: aws.String("teacat"),
		Key:    aws.String(path),
		Body:   bytes.NewReader(b),
	})
	if err != nil {
		return err
	}
	return nil
}

// Download
func (s *S3) Download(c context.Context, path string) ([]byte, error) {
	output, err := s3.New(s.Session).GetObject(&s3.GetObjectInput{
		Bucket: aws.String("teacat"),
		Key:    aws.String(path),
	})
	if err != nil {
		return nil, err
	}
	b, err := io.ReadAll(output.Body)
	if err != nil {
		return nil, err
	}
	return b, nil
}

type fileInfo struct {
	path   string
	output *s3.HeadObjectOutput
}

func (f *fileInfo) Name() string {
	return filepath.Base(f.path)
}

func (f *fileInfo) Size() int64 {
	return *f.output.ContentLength
}

// TODO: https://docs.aws.amazon.com/fsx/latest/LustreGuide/attach-s3-posix-permissions.html
// "Metadata": { "file-permissions": "0100664" }
func (f *fileInfo) Mode() os.FileMode {
	return 0644
}

func (f *fileInfo) ModTime() time.Time {
	return *f.output.LastModified
}

func (f *fileInfo) IsDir() bool {
	return strings.ToLower(*f.output.ContentType) == "application/x-directory"
}

func (f *fileInfo) Sys() any {
	return nil
}

// Copy
func (s *S3) Head(c context.Context, path string) (os.FileInfo, error) {
	output, err := s3.New(s.Session).HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String("teacat"),
		Key:    aws.String(path),
	})
	if err != nil {
		return nil, err
	}
	return &fileInfo{path, output}, nil
}

// Copy
func (s *S3) Copy(c context.Context, src, dest string) error {
	_, err := s3.New(s.Session).CopyObject(&s3.CopyObjectInput{
		Bucket:     aws.String("teacat"),
		CopySource: aws.String(src),
		Key:        aws.String(dest),
	})
	if err != nil {
		return err
	}
	return nil
}

// Delete
func (s *S3) Delete(c context.Context, path string) error {
	_, err := s3.New(s.Session).DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String("teacat"),
		Key:    aws.String(path),
	})
	if err != nil {
		return err
	}
	return nil
}

// StreamUpload
func (s *S3) StreamUpload(c context.Context, path string, reader io.Reader) error {
	uploader := s3manager.NewUploader(s.Session, func(u *s3manager.Uploader) {
		u.PartSize = 5 * 1024 * 1024
		u.Concurrency = 2
	})
	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String("teacat"),
		Key:    aws.String(path),
		Body:   reader,
	})
	if err != nil {
		return err
	}
	return nil
}

type fakeWriterAt struct {
	w io.Writer
}

func (fw fakeWriterAt) WriteAt(p []byte, offset int64) (n int, err error) {
	return fw.w.Write(p)
}

// StreamDownload
func (s *S3) StreamDownload(c context.Context, path string, writer io.Writer) error {
	downloader := s3manager.NewDownloader(s.Session, func(options *s3manager.Downloader) {
		options.Concurrency = 1
	})
	_, err := downloader.Download(fakeWriterAt{writer}, &s3.GetObjectInput{
		Bucket: aws.String("teacat"),
		Key:    aws.String(path),
	})
	if err != nil {
		return err
	}
	return nil
}

// IsNotFound
func (s *S3) IsNotExist(err error) bool {
	if err == nil {
		return false
	}
	if aerr, ok := err.(awserr.Error); ok {
		return aerr.Code() == s3.ErrCodeNoSuchKey
	}
	return false
}
