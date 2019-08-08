package main

import (
	"log"

	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"strings"

	"strconv"
	"time"

	"github.com/zzzhr1990/go-wcs-cloud-sdk/bucket"
	"github.com/zzzhr1990/go-wcs-cloud-sdk/core"
	"github.com/zzzhr1990/go-wcs-cloud-sdk/entity"
	"github.com/zzzhr1990/go-wcs-cloud-sdk/upload"
	"github.com/zzzhr1990/go-wcs-cloud-sdk/utility"
)

func main() {
	if false {
		s := utility.ComputeEtag([]byte{})
		log.Print(s)
		r, _ := utility.ComputeFileEtag("test.h")
		log.Print(r)
		return
	}
	ak := ""
	sk := ""
	file := "/"
	policy := &entity.UploadPolicy{}
	current := time.Now()
	policy.Deadline = strconv.FormatInt((time.Now().UnixNano()/int64(time.Millisecond))+1000*60*60*6, 10)
	key := "test/" + current.Format("2006-01-02") + "/" + strconv.FormatInt(1992, 10) + "/" + "Hellboy.2019.1080p.AMZN.WEBRip.DD5.1.x264-NTG.mkv"
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
	res, err := upl.UploadFile(file, token, "")
	if err != nil {
		log.Printf("error up, %v", err)
	}
	log.Printf("remote hash: %v", res.Hash)
	// log.Println(res.Hash)
	etag, _ := utility.ComputeFileEtag(file)
	log.Printf("ori hash: %v", etag)
	// prepare to post...
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
	auth := utility.NewAuth("8758804f90558e3a9222174725ee5d36ab9c7208", "5dbd1a40a8334d285261d78562640f4667ff90a9")

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
