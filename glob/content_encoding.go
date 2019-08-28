package glob

import (
	"bytes"
	"compress/flate"
	"compress/zlib"
	"io"
	"io/ioutil"
)

// contentEncoding = response.Header.Get("content-encoding")
func DecodeContent(contentEncoding string, data []byte) (result []byte, err error) {
	switch contentEncoding {
	case "gzip":
		result, err = DecodeGzip(data)
		return
	case "deflate":
		var readCloser io.ReadCloser
		readCloser = flate.NewReader(bytes.NewReader(data))
		defer readCloser.Close()

		result, err = ioutil.ReadAll(readCloser)
		if err != nil {
			return
		}
		return
	case "zlib":
		var readCloser io.ReadCloser
		readCloser, err = zlib.NewReader(bytes.NewReader(data))
		if err != nil {
			return
		}
		defer readCloser.Close()

		result, err = ioutil.ReadAll(readCloser)
		if err != nil {
			return
		}
		return
	}

	result = data
	return
}
