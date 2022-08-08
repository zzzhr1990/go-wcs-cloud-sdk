package utility

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"io"
	"os"
)

// HTTP 协议：ETag == URL 的 Entity Tag，用于标示 URL 对象是否改变，区分不同语言和 Session 等等。

// ComputeEtag cp
func ComputeEtag(data []byte) (etag string) {
	tag := make([]byte, 0, 1+sha1.Size)
	h := sha1.New()
	if len(data) < csBlockSize {
		tag = append(tag, 0x16)
		h.Write(data)
		tag = h.Sum(tag)
	} else {
		tag = append(tag, 0x96)
		innerBlockCount := blockCount(int64(len(data)))
		allBlocksSha1 := make([]byte, 0, innerBlockCount*sha1.Size)
		for i := 0; i < innerBlockCount; i++ {
			var readBytes int
			if i < innerBlockCount-1 {
				readBytes = csBlockSize
			} else {
				readBytes = len(data) - csBlockSize*i
			}
			h.Write(data[csBlockSize*i : csBlockSize*i+readBytes])
			allBlocksSha1 = h.Sum(allBlocksSha1)
			h.Reset()
		}
		h.Write(allBlocksSha1)
		tag = h.Sum(tag)
	}
	return base64.URLEncoding.EncodeToString(tag)
}

// ComputeFileEtag atir
func ComputeFileEtag(filename string) (etag string, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return
	}

	fsize := fi.Size()
	innerBlockCount := blockCount(fsize)
	var tag []byte

	if innerBlockCount <= 1 { // file size <= 4M
		tag, err = computeSha1([]byte{0x16}, f)
		if err != nil {
			return
		}
	} else { // file size > 4M
		allBlocksSha1 := []byte{}

		for i := 0; i < innerBlockCount; i++ {
			body := io.LimitReader(f, csBlockSize)
			allBlocksSha1, err = computeSha1(allBlocksSha1, body)
			if err != nil {
				return
			}
		}

		tag, _ = computeSha1([]byte{0x96}, bytes.NewReader(allBlocksSha1))
	}

	etag = base64.URLEncoding.EncodeToString(tag)
	return
}

// private:
const (
	csBlockBits = 22               // 2 ^ 22 = 4M
	csBlockSize = 1 << csBlockBits // 4M
)

func blockCount(size int64) int {
	return int((size + (csBlockSize - 1)) >> csBlockBits)
}

func computeSha1(b []byte, r io.Reader) ([]byte, error) {
	h := sha1.New()
	_, err := io.Copy(h, r)
	if err != nil {
		return nil, err
	}
	return h.Sum(b), nil
}
