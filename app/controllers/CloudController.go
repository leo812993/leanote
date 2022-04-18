package controllers

import (
	"github.com/revel/revel"
	//	"encoding/json"
	// "fmt"
	"github.com/leanote/leanote/app/info"
	"github.com/leanote/leanote/app/service"
	// . "github.com/leanote/leanote/app/lea"
	// "github.com/leanote/leanote/app/lea/netutil"
	// "gopkg.in/mgo.v2/bson"
	//	"strconv"
	// "strings"
)

type Cloud struct {
	BaseController
}

// 获取所有图片
func (c Cloud) GetImages() revel.Result {
	return c.RenderJSON(service.CloudS.GetImages())
}

// 获取Token
func (c Cloud) GenUpToken(albumName string) revel.Result {
	return c.RenderText(service.CloudS.GenUpToken(albumName))
}

// 重命名相册
func (c Cloud) RenameAlbum(curName, newName string) revel.Result {
	re := info.NewRe()
	if len(curName) == 0 || len(newName) == 0 {
		return c.RenderJSON(re)
	}
	re.Ok, re.Msg = service.CloudS.RenameAlbum(curName, newName)
	return c.RenderJSON(re)
}

// 删除相册
func (c Cloud) DeleteAlbum(albumId string) revel.Result {
	re := info.NewRe()
	re.Ok, re.Msg = service.CloudS.DeleteAlbum(albumId)
	return c.RenderJSON(re)
}

// 删除图片
func (c Cloud) DeleteImage(key string) revel.Result {
	re := info.NewRe()
	re.Ok, re.Msg = service.CloudS.DeleteImage(key)
	return c.RenderJSON(re)
}

// 重命名图片
func (c Cloud) UpdateImageTitle(old, new string) revel.Result {
	re := info.NewRe()
	re.Ok, re.Msg = service.CloudS.UpdateImageTitle(old, new)
	return c.RenderJSON(re)
}
