package service

import (
	// "fmt"
	"bytes"
	"context"
	"strings"

	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"

	"github.com/leanote/leanote/app/info"
	. "github.com/leanote/leanote/app/lea"
)

type QiniuService struct {
	bucketManager *storage.BucketManager
	bucket        string
}

type CloudService interface {
	InitCloud()
	GetImages() info.Page
	RenameAlbum(cur, name string) (ok bool, msg string)
	DeleteAlbum(albumId string) (ok bool, msg string)
	UpdateImageTitle(old, new string) (ok bool, msg string)
	DeleteImage(key string) (ok bool, msg string)
	GenUpToken(albumName string) (token string)
}

func InitCloudService() {
	ChangeCloud()
}

func (q *QiniuService) InitCloud() {
	config := configService.GetGlobalMapConfig("Cloud_Qiniu")
	q.bucket, _ = config["Bucket"]

	mac, cfg := q.generateConfig()
	q.bucketManager = storage.NewBucketManager(mac, &cfg)
}

// 初始marker为 "" query 1次
func (q QiniuService) GetImagesByMark(marker string) []storage.ListItem {
	entries, _, _, _ := q._getImagesByMark("")
	return entries
}

func (q QiniuService) _getImagesByMark(marker string) ([]storage.ListItem, string, bool, error) {
	entries, _, nextMarker, hasNext, e := q.bucketManager.ListFiles(q.bucket, "", "", marker, 1000)
	if e != nil {
		LogE(e.Error())
	}

	return entries, nextMarker, hasNext, e
}

// 多次query以获取所有图片
func (q QiniuService) GetImages() info.Page {
	images := []storage.ListItem{}
	marker := ""
	for {
		entries, nextMarker, hasNext, e := q._getImagesByMark(marker)
		if e != nil {
			LogE(e.Error())
			break
		}
		images = append(images, entries...)

		if hasNext {
			marker = nextMarker
		} else {
			break //list end
		}
	}

	return info.Page{Count: len(images), List: images}
}

// 未写完 由前端上传图片
func (q QiniuService) UploadImage() {
	putPolicy := storage.PutPolicy{
		Scope:   q.bucket,
		Expires: 12 * 3600, // token 过期时间
		SaveKey: configService.GetGlobalMapConfig("Cloud_Qiniu")["SaveKey"],
	}
	mac, cfg := q.generateConfig()
	upToken := putPolicy.UploadToken(mac)

	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	data := []byte("hello, this is qiniu cloud")
	err := formUploader.Put(context.Background(), &ret, upToken, "key", bytes.NewReader(data), int64(len(data)), nil)
	if err != nil {
		LogE(err.Error())
		return
	}
}

func (q QiniuService) RenameAlbum(cur, new string) (bool, string) {
	force := false

	marker := ""
	for {
		entries, nextMarker, hasNext, e := q._getImagesByMark(marker) // 取出数据
		if e != nil {
			return false, e.Error()
		}

		moveKeys := map[string]string{}
		for _, e := range entries {
			if strings.HasPrefix(e.Key, cur+"/") {
				moveKeys[e.Key] = new + "/" + strings.Join(strings.Split(e.Key, "/")[1:], "/")
			}
		}

		moveOps := make([]string, 0, len(moveKeys))
		for srcKey, destKey := range moveKeys {
			moveOps = append(moveOps, storage.URIMove(q.bucket, srcKey, q.bucket, destKey, force))
		}

		_, err := q.bucketManager.Batch(moveOps) //每个batch的操作数量不可以超过1000个，如果总数量超过1000，需要分批发送
		if err != nil {
			return false, err.Error()
		}

		if hasNext {
			marker = nextMarker
		} else {
			break //list end
		}
	}
	return true, ""
}

func (q QiniuService) DeleteAlbum(albumId string) (bool, string) {
	marker := ""
	for {
		entries, nextMarker, hasNext, e := q._getImagesByMark(marker) // 取出数据
		if e != nil {
			return false, e.Error()
		}

		deleteKeys := []string{}
		for _, e := range entries {
			if strings.HasPrefix(e.Key, albumId+"/") {
				deleteKeys = append(deleteKeys, e.Key)
			}
		}

		deleteOps := make([]string, 0, len(deleteKeys))
		for _, key := range deleteKeys {
			deleteOps = append(deleteOps, storage.URIDelete(q.bucket, key))
		}

		_, err := q.bucketManager.Batch(deleteOps) //每个batch的操作数量不可以超过1000个，如果总数量超过1000，需要分批发送
		if err != nil {
			return false, err.Error()
		}

		if hasNext {
			marker = nextMarker
		} else {
			break //list end
		}
	}
	return true, ""
}

func (q QiniuService) DeleteImage(key string) (bool, string) {
	e := q.bucketManager.Delete(q.bucket, key)
	if e == nil {
		return true, ""
	}
	return false, e.Error()
}

func (q QiniuService) UpdateImageTitle(old, new string) (bool, string) {
	err := q.bucketManager.Move(q.bucket, old, q.bucket, new, false)
	if err == nil {
		return true, ""
	}
	return false, err.Error()
}

// 七牛云生成上传Token
func (q QiniuService) GenUpToken(albumName string) string {
	saveKey := configService.GetGlobalMapConfig("Cloud_Qiniu")["SaveKey"]
	if len(albumName) > 0 {
		saveKey = albumName + "/" + saveKey
	}

	putPolicy := storage.PutPolicy{
		Scope:      q.bucket,
		Expires:    30 * 24 * 3600, // token 过期时间 一个月，单位秒
		SaveKey:    saveKey,
		ReturnBody: `{"fsize":$(fsize),"key":"$(key)","hash":"$(etag)","mimeType":"$(mimeType)","endUser":"$(endUser)"}`,
	}
	mac, _ := q.generateConfig()
	return putPolicy.UploadToken(mac)
}

func (q QiniuService) generateConfig() (*auth.Credentials, storage.Config) {
	config := configService.GetGlobalMapConfig("Cloud_Qiniu")
	mac := qbox.NewMac(config["AccessKey"], config["SecretKey"])
	cfg := storage.Config{
		UseHTTPS:      false, // 是否使用https域名进行资源管理
		UseCdnDomains: false, // 是否使用CDN上传加速
	}
	// 指定空间所在的区域，如果不指定将自动探测 如果没有特殊需求，默认不需要指定
	zone, _ := config["Area"]
	if zone == "华南" {
		cfg.Zone = &storage.ZoneHuanan
	} else if zone == "华东" {
		cfg.Zone = &storage.ZoneHuadong
	} else if zone == "华北" {
		cfg.Zone = &storage.ZoneHuabei
	} else if zone == "北美" {
		cfg.Zone = &storage.ZoneBeimei
	} else if zone == "新加坡" {
		cfg.Zone = &storage.ZoneXinjiapo
	} else if zone == "东南亚" {
		cfg.Zone = &storage.ZoneFogCnEast1
	} else if zone == "华东-浙江2" {
		if z, ok := storage.GetRegionByID(storage.RIDHuadongZheJiang); ok {
			cfg.Zone = &z
		} else {
			cfg.Zone = &storage.ZoneHuadong
		}
	}

	return mac, cfg
}

func ChangeCloud() {
	switch configService.GetGlobalStringConfig("ImageBedType") {
	case "", "Local":
		return
	case "Qiniu":
		CloudS = &QiniuService{}
	default:
		return
	}

	cloudService = CloudS
	CloudS.InitCloud()
}

func isQiniu() bool {
	return configService.GetGlobalStringConfig("ImageBedType") == "Qiniu"
}

func isCloud() bool {
	t := configService.GetGlobalStringConfig("ImageBedType")
	return ((t != "") && (t != "Local"))
}
