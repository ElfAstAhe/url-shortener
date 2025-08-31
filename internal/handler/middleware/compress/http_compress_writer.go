package compress

import (
	"compress/gzip"
	"compress/lzw"
	"compress/zlib"
	"fmt"
	"io"
	"net/http"

	"github.com/andybalholm/brotli"
)

type customBrotliWriteCloser struct {
	*brotli.Writer
}

func (cbwc *customBrotliWriteCloser) Close() error {
	return cbwc.Writer.Close()
}

type CustomResponseWriter struct {
	SourceWriter    http.ResponseWriter
	CompressWriter  io.WriteCloser
	contentEncoding string
}

func NewCustomResponseWriter(rw http.ResponseWriter, encoding string) (*CustomResponseWriter, error) {
	compressWriter, err := getCompressWriter(encoding, rw)
	if err != nil {
		return nil, err
	}

	return &CustomResponseWriter{
		SourceWriter:    rw,
		CompressWriter:  compressWriter,
		contentEncoding: encoding,
	}, nil
}

func (crw *CustomResponseWriter) Header() http.Header {
	return crw.SourceWriter.Header()
}

func (crw *CustomResponseWriter) WriteHeader(statusCode int) {
	if statusCode < http.StatusMultipleChoices {
		crw.Header().Set("Content-Encoding", crw.contentEncoding)
	}

	crw.SourceWriter.WriteHeader(statusCode)
}

func (crw *CustomResponseWriter) Write(p []byte) (n int, err error) {
	// check for Content-Type
	contentType := crw.Header().Get("Content-Type")
	if contentTypeApplicationJSON == contentType || contentTypeTextHTML == contentType {
		return crw.CompressWriter.Write(p)
	}

	return crw.SourceWriter.Write(p)
}

func (crw *CustomResponseWriter) Close() error {
	return crw.CompressWriter.Close()
}

func getCompressWriter(encoding string, responseWriter http.ResponseWriter) (io.WriteCloser, error) {
	switch encoding {
	case encodingBrotli:
		return &customBrotliWriteCloser{brotli.NewWriterLevel(responseWriter, brotli.BestCompression)}, nil
	case encodingGzip:
		return gzip.NewWriterLevel(responseWriter, gzip.BestCompression)
	case encodingCompress:
		return lzw.NewWriter(responseWriter, lzw.LSB, 8), nil
	case encodingDeflate:
		return zlib.NewWriterLevel(responseWriter, zlib.BestCompression)
	}

	return nil, fmt.Errorf("unknown compress encoding: %s", encoding)
}
