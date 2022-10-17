package api

import (
	"github.com/leanote/leanote/app/info"

	"github.com/revel/revel"
)

type ApiAdmin struct {
	ApiBaseContrller
}

// 获取博客数量，笔记数量，用户数量等
func (c ApiAdmin) GetStatistic() revel.Result {
	return c.RenderJSON(info.Statistic{Ok: true, UserCount: userService.CountUser(), BlogCount: noteService.CountBlog(""), NoteCount: noteService.CountNote("")})
}
