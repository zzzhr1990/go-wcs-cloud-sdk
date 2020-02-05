package main

import (
	log "github.com/sirupsen/logrus"

	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"strconv"
	"time"

	"os"
	"path/filepath"

	"flag"

	"runtime"

	"github.com/zzzhr1990/go-wcs-cloud-sdk/bucket"
	"github.com/zzzhr1990/go-wcs-cloud-sdk/core"
	"github.com/zzzhr1990/go-wcs-cloud-sdk/entity"
	"github.com/zzzhr1990/go-wcs-cloud-sdk/upload"
	"github.com/zzzhr1990/go-wcs-cloud-sdk/utility"
)

func walkfunc(path string, info os.FileInfo, err error) error {
	//

	if !info.IsDir() {
		// log.Printf("check file %v", path)
		err := main00(path, info)
		if err != nil {
			log.Fatalf("cannot upload %v", err)
			panic("upload failed!!!")
			// return err
		}
	}
	return nil
}

func main() {
	if false {

		startPath := "/"
		flag.StringVar(&startPath, "path", "/", "scan path")
		flag.Parse()
		log.Info(startPath)
		ifo, err := os.Stat(startPath)
		if err != nil {
			log.Errorf("open file err %v", err)
			return
		}

		if ifo.IsDir() {
			log.Infof("walk file: %v", ifo.Name())
			filepath.Walk(startPath, walkfunc)
		} else {
			main00(startPath, ifo)
		}

	} else {
		runtime.GOMAXPROCS(48)
		startPath := "test.bin"
		ifo, err := os.Stat(startPath)
		if err != nil {
			log.Fatalf("Fail: %v", err)
		} else {
			main00(startPath, ifo)
		}

	}
}

func main00(file string, info os.FileInfo) error {
	runtime.GOMAXPROCS(48)
	ak := ""
	sk := ""
	policy := &entity.UploadPolicy{}
	// current := time.Now()
	policy.Deadline = strconv.FormatInt((time.Now().UnixNano()/int64(time.Millisecond))+1000*60*60*6, 10)
	// key := "test/" + current.Format("2006-01-02") + "/" + strconv.FormatInt(1992, 10) + "/" + "Hellboy.2019.1080p.AMZN.WEBRip.DD5.1.x264-NTG.mkv"
	key := "test/upload_test/" + "test.bin"
	// key := "test/" + current.Format("2006-01-02") + "/" + strconv.FormatInt(1992, 10) + "/" + "Hellboy.2019.1080p.AMZN.WEBRip.DD5.1.x264-NTG.mkv"
	// set scope
	policy.Scope = "other-storage" + ":" + key
	// Set overwrite
	policy.Overwrite = 1
	policy.Separate = "0"
	// policy.CallbackURL = s.config.Wcs.CallbackURL
	// Calc token
	data, _ := json.Marshal(policy)
	encodedData := base64.URLEncoding.EncodeToString(data)
	token := ak + ":" + encodeSign([]byte(encodedData), sk) + ":" + encodedData

	url := "https://upload-vod-v1.qiecdn.com"

	// prepare upload..

	var xxp int32

	upl := upload.CreateNewSliceUpload(url, &xxp)
	start := time.Now()
	res, err := upl.UploadFile(file, token, "", 48)
	// read file size

	if err != nil {
		log.Printf("error up, %v", err)
		return err
	}

	elapsed := time.Since(start)
	speed := info.Size()
	sec := int64(elapsed.Seconds())
	if sec > 0 {
		speed = speed / sec
	}
	log.Printf("%v took %v, speed: %v/sec", info.Name(), sec, byteCountDecimal(speed))

	// log.Printf("remote hash: %v", res.Hash)
	// log.Println(res.Hash)
	etag, _ := utility.ComputeFileEtag(file)
	if res.Hash != etag {
		log.Fatalf("ori hash: %v, remote hash: %v", etag, res.Hash)
	}

	// prepare to post...
	return nil
}

func byteCountDecimal(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "kMGTPE"[exp])
}

func encodeSign(data []byte, sk string) (sign string) {
	hm := hmac.New(sha1.New, []byte(sk))
	hm.Write(data)
	sum := hm.Sum(nil)
	hexString := make([]byte, hex.EncodedLen(len(sum)))
	hex.Encode(hexString, sum)
	return base64.URLEncoding.EncodeToString(hexString)
}

// https://other.qiecdn.com/oriStore/ltGeGqA87FyZ9yYFDV83cYcwaMju?op=imageView2&mode=2&height=256&format=jpg

func main6() {
	auth := utility.NewAuth("", "")

	mng := bucket.NewBucketManager(auth, core.NewConfig(false, "upl", "qietv.mgr33.v1.wcsapi.com"))
	// we‘d test 　m
	detectType := "porn"
	flag := 0
	add := "https://other.qiecdn.com/oriStore/ltGeGqA87FyZ9yYFDV83cYcwaMju?op=imageView2&mode=2&height=256&format=jpg"
	response, err := mng.ImageDetect(add, "other-storage", detectType)
	if err != nil {
		log.Printf("cannot detect img %v: %v", add, err)
		return
	}

	for _, det := range response.Results {
		if det.PornDetect.Label == 0 {
			flag = 1 | flag
			log.Printf("porn detect...%v, rate: %v, need review: %v", add, det.PornDetect.Rate, det.PornDetect.Review)
			break
		}
	}
	log.Printf("detect: %v", flag)
}

func main5() {
	auth := utility.NewAuth("", "")

	mng := bucket.NewBucketManager(auth, core.NewConfig(false, "upl", ""))
	// we‘d test 　m
	detectType := "porn"
	flag := 0
	checkArr := []string{"other-storage:pre/image-video/preview/lhrkTG0MBEkRDDVqQa0_qp_i-Kvg/1562650936305/snapshot_00001.jpg", "other-storage:pre/image-video/preview/lhrkTG0MBEkRDDVqQa0_qp_i-Kvg/1562650936305/snapshot_00002.jpg", "other-storage:pre/image-video/preview/lhrkTG0MBEkRDDVqQa0_qp_i-Kvg/1562650936305/snapshot_00003.jpg", "other-storage:pre/image-video/preview/lhrkTG0MBEkRDDVqQa0_qp_i-Kvg/1562650936305/snapshot_00004.jpg", "other-storage:pre/image-video/preview/lhrkTG0MBEkRDDVqQa0_qp_i-Kvg/1562650936305/snapshot_00005.jpg", "other-storage:pre/image-video/preview/lhrkTG0MBEkRDDVqQa0_qp_i-Kvg/1562650936305/snapshot_00006.jpg", "other-storage:pre/image-video/preview/lhrkTG0MBEkRDDVqQa0_qp_i-Kvg/1562650936305/snapshot_00007.jpg", "other-storage:pre/image-video/preview/lhrkTG0MBEkRDDVqQa0_qp_i-Kvg/1562650936305/snapshot_00008.jpg", "other-storage:pre/image-video/preview/lhrkTG0MBEkRDDVqQa0_qp_i-Kvg/1562650936305/snapshot_00009.jpg", "other-storage:pre/image-video/preview/lhrkTG0MBEkRDDVqQa0_qp_i-Kvg/1562650936305/snapshot_00010.jpg", "other-storage:pre/image-video/preview/lhrkTG0MBEkRDDVqQa0_qp_i-Kvg/1562650936305/snapshot_00011.jpg", "other-storage:pre/image-video/preview/lhrkTG0MBEkRDDVqQa0_qp_i-Kvg/1562650936305/snapshot_00012.jpg", "other-storage:pre/image-video/preview/lhrkTG0MBEkRDDVqQa0_qp_i-Kvg/1562650936305/snapshot_00013.jpg", "other-storage:pre/image-video/preview/lhrkTG0MBEkRDDVqQa0_qp_i-Kvg/1562650936305/snapshot_00014.jpg", "other-storage:pre/image-video/preview/lhrkTG0MBEkRDDVqQa0_qp_i-Kvg/1562650936305/snapshot_00015.jpg", "other-storage:pre/image-video/preview/lhrkTG0MBEkRDDVqQa0_qp_i-Kvg/1562650936305/snapshot_00016.jpg", "other-storage:pre/image-video/preview/lhrkTG0MBEkRDDVqQa0_qp_i-Kvg/1562650936305/snapshot_00017.jpg", "other-storage:pre/image-video/preview/lhrkTG0MBEkRDDVqQa0_qp_i-Kvg/1562650936305/snapshot_00018.jpg", "other-storage:pre/image-video/preview/lhrkTG0MBEkRDDVqQa0_qp_i-Kvg/1562650936305/snapshot_00019.jpg", "other-storage:pre/image-video/preview/lhrkTG0MBEkRDDVqQa0_qp_i-Kvg/1562650936305/snapshot_00020.jpg", "other-storage:pre/image-video/preview/lhrkTG0MBEkRDDVqQa0_qp_i-Kvg/1562650936305/snapshot_00021.jpg", "other-storage:pre/image-video/preview/lhrkTG0MBEkRDDVqQa0_qp_i-Kvg/1562650936305/snapshot_00022.jpg", "other-storage:pre/image-video/preview/lhrkTG0MBEkRDDVqQa0_qp_i-Kvg/1562650936305/snapshot_00023.jpg", "other-storage:pre/image-video/preview/lhrkTG0MBEkRDDVqQa0_qp_i-Kvg/1562650936305/snapshot_00024.jpg", "other-storage:pre/image-video/preview/lhrkTG0MBEkRDDVqQa0_qp_i-Kvg/1562650936305/snapshot_00025.jpg", "other-storage:pre/image-video/preview/lhrkTG0MBEkRDDVqQa0_qp_i-Kvg/1562650936305/snapshot_00026.jpg", "other-storage:pre/image-video/preview/lhrkTG0MBEkRDDVqQa0_qp_i-Kvg/1562650936305/snapshot_00027.jpg", "other-storage:pre/image-video/preview/lhrkTG0MBEkRDDVqQa0_qp_i-Kvg/1562650936305/snapshot_00028.jpg", "other-storage:pre/image-video/preview/lhrkTG0MBEkRDDVqQa0_qp_i-Kvg/1562650936305/snapshot_00029.jpg", "other-storage:pre/image-video/preview/lhrkTG0MBEkRDDVqQa0_qp_i-Kvg/1562650936305/snapshot_00030.jpg", "other-storage:pre/image-video/preview/lhrkTG0MBEkRDDVqQa0_qp_i-Kvg/1562650936305/snapshot_00031.jpg", "other-storage:pre/image-video/preview/lhrkTG0MBEkRDDVqQa0_qp_i-Kvg/1562650936305/snapshot_00032.jpg"}
	for _, add := range checkArr {
		if flag > 0 {
			break
		}
		add = strings.ReplaceAll(add, "other-storage:", "")
		response, err := mng.ImageDetect(add, "other-storage", detectType)
		if err != nil {
			log.Printf("cannot detect img %v: %v", add, err)
			return
		}

		for _, det := range response.Results {
			if det.PornDetect.Label == 0 {
				flag = 1 | flag
				log.Printf("porn detect...%v, rate: %v, need review: %v", add, det.PornDetect.Rate, det.PornDetect.Review)
				break
			}
		}
	}

	if flag == 0 {
		log.Println("not porn")
	} else {
		log.Println("porn")
	}
}
