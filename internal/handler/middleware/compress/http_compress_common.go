package compress

// Content-Type
const (
	contentTypeApplicationJSON = "application/json"
	contentTypeTextHTML        = "text/html"
)

// Accept-Encoding, Content-Encoding
const (
	encodingGzip     = "gzip"     // gzip
	encodingCompress = "compress" // lzw
	encodingDeflate  = "deflate"  // zlib
	encodingBrotli   = "br"       // brotli
)
