package tests

import (
	"github.com/leanote/leanote/app/db"
	"testing"
	//	. "github.com/leanote/leanote/app/lea"
	"github.com/leanote/leanote/app/service"
	//	"gopkg.in/mgo.v2"
	//	"fmt"
)

func init() {
	db.Init("mongodb://note:123456@192.168.0.254:27017/leanote", "leanote")
	service.InitService()
}

// 测试登录
func TestAuth(t *testing.T) {
	_, err := service.AuthS.Login("admin", "123")
	if err != nil {
		t.Error("Admin User Auth Error")
	}
}
