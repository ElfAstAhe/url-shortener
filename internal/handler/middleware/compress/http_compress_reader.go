package compress

import (
	"compress/gzip"
	"compress/lzw"
	"compress/zlib"
	"fmt"
	"io"

	"github.com/andybalholm/brotli"
)

type customBrotliReadCloser struct {
	*brotli.Reader
}

func (cbrc *customBrotliReadCloser) Read(p []byte) (n int, err error) {
	return cbrc.Reader.Read(p)
}

func (cbrc *customBrotliReadCloser) Close() error {
	return nil
}

type CustomReader struct {
	SourceReader   io.ReadCloser
	CompressReader io.ReadCloser
}

func NewCustomReader(source io.ReadCloser, encoding string) (*CustomReader, error) {
	compressReader, err := getCompressReader(encoding, source)
	if err != nil {
		return nil, err
	}

	return &CustomReader{
		SourceReader:   source,
		CompressReader: compressReader,
	}, nil
}

func (cr *CustomReader) Read(p []byte) (n int, err error) {
	return cr.CompressReader.Read(p)
}

func (cr *CustomReader) Close() error {
	err := cr.SourceReader.Close()
	if err != nil {
		return err
	}

	return cr.SourceReader.Close()
}

func getCompressReader(encoding string, reader io.Reader) (io.ReadCloser, error) {
	switch encoding {
	case encodingGzip:
		return gzip.NewReader(reader)
	case encodingDeflate:
		return zlib.NewReader(reader)
	case encodingCompress:
		return lzw.NewReader(reader, lzw.LSB, 8), nil
	case encodingBrotli:
		return &customBrotliReadCloser{Reader: brotli.NewReader(reader)}, nil
	}

	return nil, fmt.Errorf("unknown compress encoding: %s", encoding)
}
