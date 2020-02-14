package util

import (
	"bytes"
	"compress/flate"
	"github.com/pkg/errors"
	"io/ioutil"
)

// 压缩
func GzDeflate(input []byte) ([]byte, error) {
	buf := new(bytes.Buffer)
	w, err := flate.NewWriter(buf, 4)
	if err != nil {
		return nil, err
	}
	if _, err = w.Write(input); err != nil {
		return nil, err
	}
	if err = w.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), err
}

// 解压
func GzInflate(in []byte) (out []byte, err error) {
	if len(in) == 0 {
		err = errors.New("input error")
		return
	}
	out, err = ioutil.ReadAll(flate.NewReader(bytes.NewReader(in)))
	if err != nil {
		return
	}
	return
}
