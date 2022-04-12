package service

import (
	// "github.com/leanote/leanote/app/db"
	// "fmt"
	// "io/ioutil"
	"time"

	"github.com/leanote/leanote/app/info"
	. "github.com/leanote/leanote/app/lea"

	// "gopkg.in/mgo.v2/bson"

	"html/template"
	"os"
)

type StaticBlogService struct{}

func (this *StaticBlogService) SetBlogUrl(userInfo info.User, userBlog info.UserBlog, ViewArgs map[string]interface{}) {
	siteUrl := configService.GetSiteUrl()
	blogUrls := blogService.GetBlogUrls(&userBlog, &userInfo)

	ViewArgs["siteUrl"] = siteUrl
	ViewArgs["indexUrl"] = blogUrls.IndexUrl
	ViewArgs["postUrl"] = blogUrls.PostUrl
	ViewArgs["searchUrl"] = blogUrls.SearchUrl
	ViewArgs["singleUrl"] = blogUrls.SingleUrl // 单页
	ViewArgs["archiveUrl"] = blogUrls.ArchiveUrl
	ViewArgs["archivesUrl"] = blogUrls.ArchiveUrl // 别名
	ViewArgs["tagsUrl"] = blogUrls.TagsUrl
	ViewArgs["tagPostsUrl"] = blogUrls.TagPostsUrl
	ViewArgs["tagUrl"] = blogUrls.TagPostsUrl // 别名
	ViewArgs["cateUrl"] = blogUrls.CateUrl
	ViewArgs["catesUrl"] = blogUrls.CatesUrl
	ViewArgs["rssUrl"] = blogUrls.RSSUrl

	ViewArgs["themeBaseUrl"] = "/" + userBlog.ThemePath

	ViewArgs["jQueryUrl"] = "/public/libs/jquery/jquery-1.9.0.min.js"
	ViewArgs["prettifyJsUrl"] = "/public/libs/google-code-prettify/prettify.js"
	ViewArgs["prettifyCssUrl"] = "/public/libs/google-code-prettify/prettify.css"
	ViewArgs["blogCommonJsUrl"] = "/public/blog/js/common.js"
	ViewArgs["shareCommentCssUrl"] = "/public/blog/css/share_comment.css"
	ViewArgs["shareCommentJsUrl"] = "/public/blog/js/share_comment.js"
	ViewArgs["fontAwesomeUrl"] = "/public/libs/font-awesome-4.2.0/css/font-awesome.css"
	ViewArgs["bootstrapCssUrl"] = "/public/libs/bootstrap/bootstrap.css"
	ViewArgs["bootstrapJsUrl"] = "/public/libs/bootstrap/bootstrap-min.js"

	return
}

func (this StaticBlogService) GetBlogInfo(userBlog info.UserBlog, userInfo info.User) (blogInfo map[string]interface{}) {
	blogInfo = map[string]interface{}{
		"UserId":      userBlog.UserId.Hex(),
		"Username":    userInfo.Username,
		"UserLogo":    userInfo.Logo,
		"Title":       userBlog.Title,
		"SubTitle":    userBlog.SubTitle,
		"Logo":        userBlog.Logo,
		"OpenComment": userBlog.CanComment,
		"CommentType": userBlog.CommentType, // leanote, or disqus
		"DisqusId":    userBlog.DisqusId,
		"ThemeId":     userBlog.ThemeId,
		"SubDomain":   userBlog.SubDomain,
		"Domain":      userBlog.Domain,
	}
	return
}

func (this *StaticBlogService) BlogCommon(userInfo info.User) (userBlog info.UserBlog, ViewArgs map[string]interface{}) {
	userId := userInfo.UserId.Hex()
	_, recentBlogs := blogService.ListBlogs(userId, "", 1, 5, "UpdatedTime", false)
	userBlog = blogService.GetUserBlog(userId)

	this.SetBlogUrl(userInfo, userBlog, ViewArgs)

	ViewArgs["blogInfo"] = this.GetBlogInfo(userBlog, userInfo)

	ViewArgs["recentPosts"] = blogService.FixBlogs(recentBlogs)
	ViewArgs["tags"] = blogService.GetBlogTags(userId)

	ViewArgs["singles"] = blogService.GetSingles(userId)

	themeInfo := themeService.GetThemeInfo(userBlog.ThemeId.Hex(), userBlog.Style)
	ViewArgs["themeInfo"] = themeInfo

	return
}

// -- 产生静态Blog
func getTemplate(themePath, tName string) (t *template.Template) { //themePath是完整路径
	var err error
	// themePath := themeService.getUserThemePath2(userId, themeId)
	t, err = template.ParseFiles(themePath + "/" + tName)
	if err != nil {
		LogE("获取模版", tName, "失败")
	}
	return
}

func mkdir(dir string) bool {
	if err := os.MkdirAll(dir, 0655); err != nil {
		LogE("创建 ", dir, " 失败")
		return false
	}
	return true
}

func (this *StaticBlogService) GenerateRSS(userInfo info.User, userBlog info.UserBlog) string {
	userId := userInfo.UserId.Hex()

	// staticPath := "public/blog/" + userId
	// if !mkdir(staticPath) {
	// 	return
	// }

	postURL := configService.GetSiteUrl() + "/blog/post/" + userInfo.Username
	rss := `<?xml version="1.0" ?>
<rss version="2.0">
<channel>
  <title>` + userBlog.Title + `</title>
  <link>` + configService.GetSiteUrl() + `</link>
  <description>` + userBlog.SubTitle + `</description>`

	_, blogs := blogService.ListBlogs(userId, "", 1, userBlog.PerPageSize, userBlog.SortField, false) // page 从1开始计数，最新的PerPageSize篇， 降序
	for _, b := range blogs {
		var pubDate time.Time
		if userBlog.SortField == "CreatedTime" {
			pubDate = b.CreatedTime
		} else if userBlog.SortField == "UpdatedTime" {
			pubDate = b.UpdatedTime
		} else {
			pubDate = b.PublicTime
		}

		rss += `
  <item>
    <title>` + b.Title + `</title>
    <link>` + postURL + "/" + b.UrlTitle + `</link>
    <description>` + b.Desc + `</description>
    <pubDate>` + pubDate.Format(time.RFC822) + `</pubDate>
    <guid>` + postURL + "/" + b.UrlTitle + `</guid>
  </item>`
	}

	rss += `
</channel>
</rss>`

	// if e := ioutil.WriteFile(staticPath+"/rss.xml", []byte(rss), 0644); e != nil {
	// 	LogE("写 rss.xml 文件失败：" + e.Error())
	// }
	return rss
}

// 未写完
func (this *StaticBlogService) GenerateStaticHtmls(username string) {
	userInfo := userService.GetUserInfoByUsername(username)
	userId := userInfo.UserId.Hex()
	// userBlog := this.GetUserBlog(userId)
	// db.Users.Find(bson.M{"Username": username}).One(&userInfo)
	// db.UserBlogs.FindId(userInfo.UserId).One(&userBlog)
	// this.fixUserBlog(&userBlog)

	staticPath := "public/blog/" + userId
	if !mkdir(staticPath) {
		return
	}

	// 准备数据
	userBlog, ViewArgs := this.BlogCommon(userInfo)

	// 生成首页
	f, err := os.OpenFile(staticPath+"/index.html", os.O_CREATE|os.O_WRONLY, 0655)
	if err != nil {
		LogE("创建 ", staticPath+"/index.html 失败")
		return
	}
	t := getTemplate(userBlog.ThemePath, "index.html")
	ViewArgs["curIsIndex"] = true
	t.Execute(f, ViewArgs)
	f.Close()

	ViewArgs["curIsIndex"] = false
}
