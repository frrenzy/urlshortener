// Package gzip
package gzip

import (
	"compress/gzip"
	"errors"
	"io"
	"net/http"
	"slices"
)

type compressWriter struct {
	http.ResponseWriter
	zw           *gzip.Writer
	contentTypes []string
}

func newCompressWriter(w http.ResponseWriter, contentTypes []string) *compressWriter {
	return &compressWriter{
		ResponseWriter: w,
		zw:             gzip.NewWriter(w),
		contentTypes:   contentTypes,
	}
}

func (c *compressWriter) shouldCompressType() bool {
	contentType := c.Header().Get("Content-Type")
	return slices.Contains(c.contentTypes, contentType)
}

func (c *compressWriter) writer() io.Writer {
	if c.shouldCompressType() {
		return c.zw
	}

	return c.ResponseWriter
}

func (c *compressWriter) WriteHeader(statusCode int) {
	if statusCode < 300 && c.shouldCompressType() {
		c.Header().Set("Content-Encoding", "gzip")
	}

	c.ResponseWriter.WriteHeader(statusCode)
}

func (c *compressWriter) Write(p []byte) (int, error) {
	return c.writer().Write(p)
}

func (c *compressWriter) Close() error {
	if c, ok := c.writer().(io.WriteCloser); ok {
		return c.Close()
	}
	return errors.New("io.WriteCloser is unavailable on the writer")
}

type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &compressReader{
		r:  r,
		zr: zr,
	}, nil
}

func (c compressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}
