package sqlite

import (
	//"github.com/go-gorp/gorp"

	"time"

	"github.com/leanote/leanote/app/info"
	. "github.com/leanote/leanote/app/lea"
	rgorp "github.com/revel/modules/orm/gorp/app"
)

var Notebooks string
var Notes string
var NoteContents string
var NoteContentHistories string

var ShareNotes string
var ShareNotebooks string
var HasShareNotes string
var Blogs string
var Users string
var Groups string
var GroupUsers string

var Tags string
var NoteTags string
var TagCounts string

var UserBlogs string

var Tokens string

var Suggestions string

// Album & file(image)
var Albums string
var Files string
var Attachs string

var NoteImages string
var Configs string
var EmailLogs string

// blog
var BlogLikes string
var BlogComments string
var Reports string
var BlogSingles string
var Themes string

// session
var Sessions string

func InitSqlite() {
	Dbm := rgorp.Db.Map

	// notebook
	Notebooks = "notebooks"

	// notes
	Notes = "notes"

	// noteContents
	NoteContents = "note_contents"
	NoteContentHistories = "note_content_histories"

	// share
	ShareNotes = "share_notes"
	ShareNotebooks = "share_notebooks"
	HasShareNotes = "has_share_notes"

	// user
	Users = "users"
	// group
	Groups = "groups"
	GroupUsers = "group_users"

	// blog
	Blogs = "blogs"

	// tag
	Tags = "tags"
	NoteTags = "note_tags"
	TagCounts = "tag_count"

	// blog
	UserBlogs = "user_blogs"
	BlogSingles = "blog_singles"
	Themes = "themes"

	// find password
	Tokens = "tokens"

	// Suggestion
	Suggestions = "suggestions"

	// Album & file
	Albums = "albums"
	Files = "files"
	Attachs = "attachs"

	NoteImages = "note_images"

	Configs = "configs"
	EmailLogs = "email_logs"

	// 社交
	BlogLikes = "blog_likes"
	BlogComments = "blog_comments"

	// 举报
	Reports = "reports"

	// session
	Sessions = "sessions"
	Logf("sssssss:%s %v", "InitSqlite", Dbm)
}

func AddUser() error {
	rgorp.Db.Map.AddTable(info.Users{}).SetKeys(false, "Id")
	users := &info.Users{}
	users.Username = "123456"
	users.CreatedTime = time.Now().Unix()
	users.Email = "abc@qq.com"
	//if err := rgorp.Db.Map.Insert(users); err != nil {
	//panic(err)
	//	Logf("err:%v", err)

	//}

	users1 := &info.Users{}
	selectSQL1 := "select * from " + "Users" + " where Id = :Id"
	params := map[string]interface{}{"Id": "5fdf2c0f1711a70026000012"}
	//rgorp.Db.SqlStatementBuilder.Select("*").From("User").Where("Username=?", "dd")
	err := rgorp.Db.Map.SelectOne(users1, selectSQL1, params)
	Logf("users1:%s %v", users1.Email, err)
	return nil
}
