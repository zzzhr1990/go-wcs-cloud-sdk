package upload

import (
	"bytes"
	"errors"
	"net/http"
	"os"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/zzzhr1990/go-wcs-cloud-sdk/utility"
)

// SliceUpload used for udp
type SliceUpload struct {
	// config       *core.Config
	httpManager  *utility.HTTPManager
	uploadBatch  string
	uploadPrefix string
	blockBits    uint64
	blockSize    int64
	blockMask    int64
}

const (
	minBlockBits  = 22
	maxBlockCount = 100
)

// CreateNewSliceUpload new module
func CreateNewSliceUpload(uploadPrefix string) *SliceUpload {
	su := &SliceUpload{
		httpManager:  utility.NewHTTPManager(),
		uploadPrefix: uploadPrefix,
	}
	su.resizeFileCount(1)
	return su
}

func (up *SliceUpload) blockCount(fsize int64) int64 {
	return (fsize + up.blockMask) >> up.blockBits
}

func (up *SliceUpload) resizeFileCount(fsize int64) {
	up.blockBits = minBlockBits
	up.blockSize = 1 << up.blockBits
	up.blockMask = 1<<up.blockBits - 1

	count := up.blockCount(fsize)
	for count > maxBlockCount {
		up.blockBits++
		up.blockSize = 1 << up.blockBits
		up.blockMask = 1<<up.blockBits - 1
		count = up.blockCount(fsize)
	}
}

// MakeBlock make a block to upload
func (up *SliceUpload) MakeBlock(blockSize int64, blockOrder int64, chunk []byte, uploadToken string, key string) (*MakeBlockBputResult, error) {
	if blockSize < 0 || up.blockSize < blockSize {
		return nil, errors.New("innerblockSize is invalid")
	}
	if 0 == len(uploadToken) {
		return nil, errors.New("upload_token is empty")
	}

	url := up.uploadPrefix + "/mkblk/" + strconv.FormatInt(blockSize, 10) + "/" + strconv.FormatInt(blockOrder, 10)
	request, err := http.NewRequest("POST", url, bytes.NewReader(chunk))
	if nil != err {
		return nil, err
	}

	utility.AddMime(request, "application/octet-stream")
	request.Header.Set("UploadBatch", up.uploadBatch)
	if len(key) > 0 {
		request.Header.Set("Key", utility.URLSafeEncodeString(key))
	}

	result := &MakeBlockBputResult{}
	err = up.httpManager.DoWithTokenAndRetry(request, uploadToken, result, 10)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Bput put ctx
func (up *SliceUpload) Bput(context string, offset int64, chunk []byte, uploadToken string, key string) (*MakeBlockBputResult, error) {
	if 0 == len(context) {
		return nil, errors.New("context is empty")
	}
	if 0 == len(uploadToken) {
		return nil, errors.New("upload_token is empty")
	}

	request, err := http.NewRequest("POST", up.uploadPrefix+"/bput/"+context+"/"+strconv.FormatInt(offset, 10),
		bytes.NewReader(chunk))
	if nil != err {
		return nil, err
	}

	utility.AddMime(request, "application/octet-stream")
	request.Header.Set("UploadBatch", up.uploadBatch)
	if len(key) > 0 {
		request.Header.Set("Key", utility.URLSafeEncodeString(key))
	}
	result := &MakeBlockBputResult{}
	err = up.httpManager.DoWithTokenAndRetry(request, uploadToken, result, 10)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// MakeFile make a file, no extra
func (up *SliceUpload) MakeFile(size int64, key string, contexts []string, uploadToken string) (*MakeFileResult, error) {
	if size < 0 {
		return nil, errors.New("size is invalid")
	}
	if nil == contexts {
		return nil, errors.New("contexts is empty")
	}
	if 0 == len(uploadToken) {
		return nil, errors.New("upload_token is empty")
	}

	url := up.uploadPrefix + "/mkfile/" + strconv.FormatInt(size, 10)
	/*
		if nil != put_extra && nil != put_extra.Params {
			for k, v := range put_extra.Params {
				if strings.HasPrefix(k, "x:") && len(v) > 0 {
					url += "/" + k + "/" + utility.UrlSafeEncodeString(v)
				}
			}
		}
	*/

	ctx := ""
	for _, c := range contexts {
		ctx += "," + c
	}
	request, err := http.NewRequest("POST", url, strings.NewReader(ctx[1:]))
	if nil != err {
		return nil, err
	}

	utility.AddMime(request, "text/plain;charset=UTF-8")
	request.Header.Set("UploadBatch", up.uploadBatch)
	if len(key) > 0 {
		request.Header.Set("Key", utility.URLSafeEncodeString(key))
	}
	/*
		if nil != put_extra {
			if len(put_extra.MimeType) > 0 {
				request.Header.Set("MimeType", put_extra.MimeType)
			}
			if -1 != put_extra.Deadline {
				request.Header.Set("Deadline", strconv.Itoa(put_extra.Deadline))
			}
		}
	*/
	result := &MakeFileResult{}
	err = up.httpManager.DoWithTokenAndRetry(request, uploadToken, result, 10)
	if err != nil {
		return nil, err
	}
	return result, nil

}

// UploadFile Upload
func (up *SliceUpload) UploadFile(localFilename string, uploadToken string, key string) (*MakeFileResult, error) {
	if 0 == len(localFilename) {
		return nil, errors.New("localFilename is empty")
	}

	f, err := os.Open(localFilename)
	if err != nil {
		log.Errorf("cannot open file %v: %v", localFilename, err)
		return nil, err
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		log.Errorf("cannot stat file %v: %v", localFilename, err)
		return nil, err
	}

	up.resizeFileCount(fi.Size())

	var innerblockSize int64
	// 第一个分片不宜太大，因为可能遇到错误，上传太大是白费流量和时间！
	var firstChunkSize int64

	if fi.Size() < 1024 {
		innerblockSize = fi.Size()
		firstChunkSize = fi.Size()
	} else {
		if fi.Size() < up.blockSize {
			innerblockSize = fi.Size()
		} else {
			innerblockSize = up.blockSize
		}
		firstChunkSize = 1024
	}

	firstChunk := make([]byte, firstChunkSize)
	n, err := f.Read(firstChunk)
	if nil != err {
		log.Errorf("cannot read chunk %v", err)
		return nil, err
	}
	if firstChunkSize != int64(n) {
		err = errors.New("Read size < request size")
		return nil, err
	}

	makeBlockResponse, err := up.MakeBlock(innerblockSize, 0, firstChunk, uploadToken, key)
	if nil != err {
		log.Errorf("cannot make block %v", err)
		return nil, err
	}

	if len(makeBlockResponse.Crc32) == 0 {
		log.Errorf("cannot make block, result is empty")
		return nil, errors.New("cannot make block, result is empty")
	}

	innerBlockCount := up.blockCount(fi.Size())
	contexts := make([]string, innerBlockCount)
	contexts[0] = makeBlockResponse.Ctx

	// 上传第 1 个 block 剩下的数据
	if innerblockSize > firstChunkSize {
		firstBlockLeftSize := innerblockSize - firstChunkSize
		leftChunk := make([]byte, firstBlockLeftSize)
		n, err = f.Read(leftChunk)
		if nil != err {
			log.Errorf("cannot read chunk %v: %v", localFilename, err)
			return nil, err
		}
		if firstBlockLeftSize != int64(n) {
			err = errors.New("Read size < request size")
			log.Errorf("cannot stat file (size check failed) %v: %v", localFilename, err)
			return nil, err
		}
		makeBlockResponse, err = up.Bput(contexts[0], firstChunkSize, leftChunk, uploadToken, key)
		if nil != err {
			log.Errorf("cannot make block file %v: %v", localFilename, err)
			return nil, err
		}

		contexts[0] = makeBlockResponse.Ctx
		// 上传后续 block，每次都是一整块上传
		for blockIndex := int64(1); blockIndex < innerBlockCount; blockIndex++ {
			pos := innerblockSize * blockIndex
			leftSize := fi.Size() - pos
			var innerChunkSize int64
			if leftSize > innerblockSize {
				innerChunkSize = innerblockSize
			} else {
				innerChunkSize = leftSize
			}
			block := make([]byte, innerChunkSize)
			n, err = f.Read(block)
			if nil != err {
				log.Errorf("cannot read block file %v: %v", localFilename, err)
				return nil, err
			}
			if innerChunkSize != int64(n) {
				err = errors.New("Read size < request size")
				log.Errorf("cannot read chunk file %v: %v", localFilename, err)
				return nil, err
			}
			makeBlockResponse, err = up.MakeBlock(innerChunkSize, blockIndex, block, uploadToken, key)
			if nil != err {
				log.Errorf("cannot make block file %v [block:%v]: %v", localFilename, block, err)
				return nil, err
			}
			contexts[blockIndex] = makeBlockResponse.Ctx
		}
	}

	response, err := up.MakeFile(fi.Size(), key, contexts, uploadToken)
	if nil != err {
		log.Errorf("cannot make block file %v: %v", localFilename, err)
		return nil, err
	}
	return response, nil
}
