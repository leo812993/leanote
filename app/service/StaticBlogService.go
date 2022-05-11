package service

import (
	// "github.com/leanote/leanote/app/db"
	// "fmt"
	// "io/ioutil"
	"encoding/xml"
	"time"

	"github.com/leanote/leanote/app/info"
	. "github.com/leanote/leanote/app/lea"

	// "gopkg.in/mgo.v2/bson"

	"html/template"
	"os"
)

type StaticBlogService struct{}

type SiteMap struct {
	XMLName xml.Name      `xml:"urlset"`
	Xmlns   string        `xml:"xmlns,attr"`
	Urls    []SiteMapItem `xml:"url"`
}

type SiteMapItem struct {
	XMLName xml.Name `xml:"url"`
	Loc     string   `xml:"loc"`
	Lastod  string   `xml:"lastmod"`
}

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

	ViewArgs["jQueryUrl"] = "/public/libs/jquery/jquery.min.js"
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

	_, blogs := blogService.ListBlogs(userId, "", 1, userBlog.PerPageSize, "PublicTime", false) // page 从1开始计数，最新的PerPageSize篇， 降序, userBlog.SortField
	for _, b := range blogs {
		pubDate := b.PublicTime
		// var pubDate time.Time
		// if userBlog.SortField == "CreatedTime" {
		// 	pubDate = b.CreatedTime
		// } else if userBlog.SortField == "UpdatedTime" {
		// 	pubDate = b.UpdatedTime
		// } else {
		// 	pubDate = b.PublicTime
		// }

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

func (this *StaticBlogService) GenerateSiteMap(userInfo info.User, userBlog info.UserBlog) SiteMap {
	userId := userInfo.UserId.Hex()

	postURL := configService.GetSiteUrl() + "/blog/post/" + userInfo.Username
	_, blogs := blogService.ListBlogs(userId, "", 1, 50000, "PublicTime", false) // 最多只能支持 50,000 个链接
	SiteMapItems := make([]SiteMapItem, 0, len(blogs))
	for _, b := range blogs {
		SiteMapItems = append(SiteMapItems, SiteMapItem{Loc: postURL + "/" + b.UrlTitle + "/", Lastod: b.UpdatedTime.Format("2006-01-02")})
	}
	SiteMap := SiteMap{Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9", Urls: SiteMapItems}
	return SiteMap
}

// 都可以识别，谷歌要在sitemap.xml后加 '/'
func (this *StaticBlogService) GenerateSiteMapTXT(userInfo info.User, userBlog info.UserBlog) string {
	userId := userInfo.UserId.Hex()

	postURL := configService.GetSiteUrl() + "/blog/post/" + userInfo.Username
	sitemap := `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`

	_, blogs := blogService.ListBlogs(userId, "", 1, 50000, "PublicTime", false) // 最多只能支持 50,000 个链接
	for _, b := range blogs {
		if len(postURL+"/"+b.UrlTitle) >= 2048 { // + https:// 整体长度需要小于2048个字符
			continue
		}
		sitemap += `
  <url>
    <loc>` + postURL + "/" + b.UrlTitle + "/" + `</loc>
    <lastmod>` + b.UpdatedTime.Format("2006-01-02") + `</lastmod>` // <changefreq>monthly</changefreq>
		// 	if b.IsTop {
		// 		sitemap += `
		// <priority>1</priority>`
		// 	}
		sitemap += `
  </url>`
	}

	sitemap += `
</urlset>`

	// staticPath := "public/blog/" + userId
	// if !mkdir(staticPath) {
	// 	return ""
	// }
	// if e := ioutil.WriteFile(staticPath+"/sitemap.xml", []byte(sitemap), 0644); e != nil {
	// 	LogE("写 sitemap.xml 文件失败：" + e.Error())
	// }
	return sitemap

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
