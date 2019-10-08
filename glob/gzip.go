package glob

import (
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
)

func EncodeGzip(in []byte) (out []byte, err error) {
	var buffer bytes.Buffer
	writer := gzip.NewWriter(&buffer)
	_, err = writer.Write(in)
	if err != nil {
		return
	}
	err = writer.Close()
	if err != nil {
		return
	}

	out = buffer.Bytes()
	return
}

func DecodeGzip(in []byte) (out []byte, err error) {
	var readCloser io.ReadCloser
	readCloser, err = gzip.NewReader(bytes.NewReader(in))
	if err != nil {
		return
	}
	defer func() {
		closeErr := readCloser.Close()
		if err == nil {
			err = closeErr
		}
	}()

	out, err = ioutil.ReadAll(readCloser)
	if err != nil {
		return
	}
	return
}
