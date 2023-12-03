package tests

import (
	"github.com/leanote/leanote/app/db"
	"testing"
	//	. "github.com/leanote/leanote/app/lea"
	//	"github.com/leanote/leanote/app/service"
	//	"gopkg.in/mgo.v2"
	//	"fmt"
)

func TestDBConnect(t *testing.T) {
	db.Init("mongodb://note:123456@192.168.0.254:27017/leanote", "leanote")
}
