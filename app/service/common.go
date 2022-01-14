package service

import "strings"

// service 通用方法

// 分页, 排序处理
func parsePageAndSort(pageNumber, pageSize int, sortField string, isAsc bool) (skipNum int, sortFieldList []string) {
	skipNum = (pageNumber - 1) * pageSize
	sortFieldList = parseSort(sortField, isAsc)
	return
}

// 只排序 不分页
func parseSort(sortField string, isAsc bool) (sortFieldList []string) {
	if sortField == "" {
		sortField = "UpdatedTime"
	}
	sortFieldList = strings.Split(sortField, ";")
	for i, s := range sortFieldList {
		if strings.HasPrefix(s, "-") {continue}
		if strings.HasPrefix(s, "+") {
			sortFieldList[i] = s[1:] // 去掉+号
			continue
		}
		if !isAsc { // 其他的根据 isAsc 来确定
			sortFieldList[i] = "-" + s // 降序
		} else {
			sortFieldList[i] = s
		}
	}
	return
}
