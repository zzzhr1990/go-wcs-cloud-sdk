package main

import (
	"log"

	"encoding/json"

	"github.com/zzzhr1990/go-wcs-cloud-sdk/bucket"
	"github.com/zzzhr1990/go-wcs-cloud-sdk/core"
	"github.com/zzzhr1990/go-wcs-cloud-sdk/utility"
)

func main() {
	auth := utility.NewAuth("", "")
	testMap := map[string]string{
		"习近平1.png": "https://config.6pan.cn/test/%E4%B9%A0%E8%BF%91%E5%B9%B31.png",
		"习近平2.png": "https://config.6pan.cn/test/%E4%B9%A0%E8%BF%91%E5%B9%B32.png",
		"习近平3.png": "https://config.6pan.cn/test/%E4%B9%A0%E8%BF%91%E5%B9%B33.png",
		"习近平4.png": "https://config.6pan.cn/test/%E4%B9%A0%E8%BF%91%E5%B9%B34.png",
		"习近平5.png": "https://config.6pan.cn/test/%E4%B9%A0%E8%BF%91%E5%B9%B35.png",
		"江泽民8.png": "https://config.6pan.cn/test/%E6%B1%9F%E6%B3%BD%E6%B0%918.png",
		"江泽民7.png": "https://config.6pan.cn/test/%E6%B1%9F%E6%B3%BD%E6%B0%917.png",
		"江泽民6.png": "https://config.6pan.cn/test/%E6%B1%9F%E6%B3%BD%E6%B0%916.png",
		"江泽民5.png": "https://config.6pan.cn/test/%E6%B1%9F%E6%B3%BD%E6%B0%915.png",
		"江泽民4.png": "https://config.6pan.cn/test/%E6%B1%9F%E6%B3%BD%E6%B0%914.png",
		"江泽民3.png": "https://config.6pan.cn/test/%E6%B1%9F%E6%B3%BD%E6%B0%913.png",
		"江泽民2.png": "https://config.6pan.cn/test/%E6%B1%9F%E6%B3%BD%E6%B0%912.png",
		"江泽民1.png": "https://config.6pan.cn/test/%E6%B1%9F%E6%B3%BD%E6%B0%911.png",
	}
	mng := bucket.NewBucketManager(auth, core.NewConfig(false, "upl", ""))
	for k, v := range testMap {
		detectType := "political"
		res, err := mng.ImageDetect(v, "other-storage", detectType)
		if err != nil {
			log.Printf("error detect: %v", err)
		} else {
			if detectType == "porn" {
				log.Printf("porn detect, result: %v, rate: %v", res.Results[0].PornDetect.Label, res.Results[0].PornDetect.Rate)
			} else if detectType == "political" {
				b, _ := json.Marshal(res)
				log.Printf("detect political: %v: %v", k, string(b))
			}
		}
	}

}
