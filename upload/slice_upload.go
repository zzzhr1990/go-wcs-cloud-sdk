package upload

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

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
	stopFlag     *int32
}

type blockInfo struct {
	// blockIndex, maxBputSize, innerChunkSize, localFilename, lastOffset, uploadToken, key
	blockIndex     int64
	maxBputSize    int64
	innerChunkSize int64
	localFilename  string
	lastOffset     int64
	uploadToken    string
	key            string
	err            error
	ctx            string
}

func createBlockInfo(blockIndex int64, maxBputSize int64, innerChunkSize int64, localFilename string, lastOffset int64, uploadToken string, key string) *blockInfo {
	return &blockInfo{
		blockIndex:     blockIndex,
		maxBputSize:    maxBputSize,
		innerChunkSize: innerChunkSize,
		localFilename:  localFilename,
		lastOffset:     lastOffset,
		uploadToken:    uploadToken,
		key:            key,
		err:            nil,
		ctx:            "",
	}
}

const (
	minBlockBits        = 22
	maxBlockCount       = 16384
	maxBlockSize        = 4 * 1024 * 1024
	maxBputSize   int64 = 2 * 1024 * 1024
)

// CreateNewSliceUpload new module
func CreateNewSliceUpload(uploadPrefix string, stopFlag *int32) *SliceUpload {
	su := &SliceUpload{
		httpManager:  utility.NewHTTPManager(),
		uploadPrefix: uploadPrefix,
		stopFlag:     stopFlag,
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

		btcs := up.blockBits + 1
		var bt int64 = 1
		btcsSize := bt << btcs
		if btcsSize > maxBlockSize {
			count = up.blockCount(fsize)
			break
		}

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
	request, err := utility.CreateCommonRequest("POST", url)
	if nil != err {
		return nil, err
	}
	request.AddData(chunk)

	sec := len(chunk) / 1024 / 100
	if sec < 90 {
		sec = 90
	}
	utility.AddMime(request, "application/octet-stream")
	request.AddHeader("UploadBatch", up.uploadBatch)
	if len(key) > 0 {
		request.AddHeader("Key", utility.URLSafeEncodeString(key))
	}

	request.SetTimeout(time.Duration(sec) * time.Second)

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

	request, err := utility.CreateCommonRequest("POST", up.uploadPrefix+"/bput/"+context+"/"+strconv.FormatInt(offset, 10))
	if nil != err {
		return nil, err
	}
	// s
	request.AddData(chunk)
	sec := len(chunk) / 1024 / 100
	if sec < 90 {
		sec = 90
	}
	request.SetTimeout(time.Duration(sec) * time.Second)
	utility.AddMime(request, "application/octet-stream")
	request.AddHeader("UploadBatch", up.uploadBatch)
	if len(key) > 0 {
		request.AddHeader("Key", utility.URLSafeEncodeString(key))
	}
	result := &MakeBlockBputResult{}
	// startTime := time.Now().UnixNano()
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

	ctx := strings.Join(contexts, ",")
	request, err := utility.CreateCommonRequest("POST", url)
	if nil != err {
		return nil, err
	}
	request.AddStringBody(ctx)
	request.SetTimeout(time.Duration(60) * time.Second)
	utility.AddMime(request, "text/plain;charset=UTF-8")
	request.AddHeader("UploadBatch", up.uploadBatch)
	if len(key) > 0 {
		request.AddHeader("Key", utility.URLSafeEncodeString(key))
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
		log.Errorf("cannot mkfile, filesize %v, ctx %v", size, ctx)
		return nil, err
	}
	return result, nil

}

// UploadFile Upload
func (up *SliceUpload) UploadFile(localFilename string, uploadToken string, key string, maxConcurrent int32) (*MakeFileResult, error) {

	if 128 < maxConcurrent {
		maxConcurrent = 128
	}
	if 0 == len(localFilename) {
		return nil, errors.New("localFilename is empty")
	}

	f, err := os.Open(localFilename)

	if err != nil {
		log.Errorf("cannot open file %v: %v", localFilename, err)
		return nil, err
	}
	defer f.Close()
	fi, err := os.Stat(localFilename)
	if err != nil {
		log.Errorf("cannot stat file %v: %v", localFilename, err)
		return nil, err
	}

	up.resizeFileCount(fi.Size())

	// if fi.Size() ==

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
	if atomic.LoadInt32(up.stopFlag) > 0 {
		return nil, errors.New("upload terminated")
	}
	makeBlockResponse, err := up.MakeBlock(innerblockSize, 0, firstChunk, uploadToken, key)
	if nil != err {
		log.Errorf("cannot make block %v", err)
		return nil, err
	}

	if makeBlockResponse.Crc32 == 0 && firstChunkSize > 0 {
		log.Errorf("cannot make first block, result is empty, fsize %v", fi.Size())
		return nil, errors.New("cannot make block, result is empty")
	}

	innerBlockCount := up.blockCount(fi.Size())
	if innerBlockCount == 0 {
		innerBlockCount = 1
	}
	contexts := make([]string, innerBlockCount)
	contexts[0] = makeBlockResponse.Ctx

	// 上传第 1 个 block 剩下的数据
	if innerblockSize > firstChunkSize {

		blockSizeLeft := innerblockSize - firstChunkSize
		lastOffset := firstChunkSize
		blockInfoList := []*blockInfo{} //  Must have data

		// we create first size

		firstInfo := createBlockInfo(0, maxBputSize, blockSizeLeft, localFilename, lastOffset, uploadToken, key)
		firstInfo.ctx = contexts[0]
		blockInfoList = append(blockInfoList, firstInfo)
		lastOffset += blockSizeLeft

		// UPLOAD OTHER BLOCKS

		//
		for blockIndex := int64(1); blockIndex < innerBlockCount; blockIndex++ {
			pos := innerblockSize * blockIndex
			leftSize := fi.Size() - pos
			var innerChunkSize int64
			if leftSize > innerblockSize {
				innerChunkSize = innerblockSize
			} else {
				innerChunkSize = leftSize
			}

			blockInfoList = append(blockInfoList, createBlockInfo(blockIndex, maxBputSize, innerChunkSize, localFilename, lastOffset, uploadToken, key))

			lastOffset += innerChunkSize
		}
		ql := len(blockInfoList)
		if ql > 0 {
			concurrent := int(maxConcurrent)
			if ql < concurrent {
				concurrent = ql
			}
			queueLen := 0
			infoChan := make(chan *blockInfo, 10)
			for len(blockInfoList) > 0 && queueLen < concurrent {
				// fetch
				task, blockInfoLst := blockInfoList[0], blockInfoList[1:]
				blockInfoList = blockInfoLst
				go func(tsk *blockInfo) {
					up.uploadSingleBlock(tsk)
					infoChan <- task
				}(task)

				queueLen++
			}

			hasError := false
			var e error

			for resBlockResult := range infoChan {
				if resBlockResult.err != nil {
					hasError = true
					e = resBlockResult.err
				}
				queueLen--
				if hasError {
					if queueLen < 1 {
						close(infoChan)
						break
					}
				} else {
					// normal
					contexts[resBlockResult.blockIndex] = resBlockResult.ctx

					rest := len(blockInfoList)
					if rest < 1 && queueLen < 1 {
						close(infoChan)
						break
					} else if rest > 0 {
						task, blockInfoLst := blockInfoList[0], blockInfoList[1:]
						blockInfoList = blockInfoLst
						go func(tsk *blockInfo) {
							up.uploadSingleBlock(tsk)
							infoChan <- task
						}(task)

						queueLen++
					}
				}
			}
			if hasError {
				return nil, e
			}
		}

	}

	if atomic.LoadInt32(up.stopFlag) > 0 {
		return nil, errors.New("upload terminated")
	}

	response, err := up.MakeFile(fi.Size(), key, contexts, uploadToken)
	if nil != err {
		log.Errorf("cannot make block file %v: %v", localFilename, err)
		return nil, err
	}
	return response, nil
}

// UploadFile Upload
func (up *SliceUpload) uploadSingleBlock(info *blockInfo) {
	var lastOffset int64
	f, err := os.Open(info.localFilename)
	makeBlock := false
	lastCtx := ""
	if len(info.ctx) > 0 {
		lastCtx = info.ctx
		info.ctx = ""
		makeBlock = true
		lastOffset = info.lastOffset
	}

	defer f.Close()
	if err != nil {
		log.Errorf("cannot open file %v: %v", info.localFilename, err)
		info.err = err
		return
	}

	blockSizeLeft := info.innerChunkSize

	f.Seek(info.lastOffset, 0)
	for blockSizeLeft > 0 {
		if atomic.LoadInt32(up.stopFlag) > 0 {
			info.err = errors.New("upload terminated")
			return
		}
		bputSize := maxBputSize
		if bputSize > blockSizeLeft {
			bputSize = blockSizeLeft
		}
		leftChunk := make([]byte, bputSize)
		n, err := f.Read(leftChunk)
		if nil != err {
			log.Errorf("cannot read chunk %v: %v", info.localFilename, err)
			info.err = err
			return
		}
		if bputSize != int64(n) {
			err = errors.New("Read size < request size")
			log.Errorf("cannot stat file (size check failed) %v: %v", info.localFilename, err)
			info.err = err
			return
		}
		if !makeBlock {
			makeBlockResponse, err := up.MakeBlock(info.innerChunkSize, info.blockIndex, leftChunk, info.uploadToken, info.key)
			if nil != err {
				log.Errorf("cannot make block file %v [block:%v]: %v", info.localFilename, info.blockIndex, err)
				info.err = err
				return
			}
			lastCtx = makeBlockResponse.Ctx
			makeBlock = true
			lastOffset = makeBlockResponse.Offset
			if len(lastCtx) == 0 {
				log.Errorf("no ctx found from block, terminate")
				info.err = errors.New("no ctx found from make block, terminiate")
				return
			}
		} else {
			if len(lastCtx) == 0 {
				log.Errorf("no ctx found from block")
				info.err = errors.New("no ctx found from make block, cannot bput")
				return
			}
			makeBlockResponse, err := up.Bput(lastCtx, lastOffset, leftChunk, info.uploadToken, info.key)
			if nil != err {
				log.Errorf("cannot make block file %v: %v", info.localFilename, err)
				info.err = err
				return
			}
			lastCtx = makeBlockResponse.Ctx
			lastOffset = makeBlockResponse.Offset
		}
		blockSizeLeft -= bputSize
	}
	info.ctx = lastCtx
}
