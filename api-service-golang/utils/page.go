package utils

import (
	"fmt"
	"math"
)

type PageParam struct {
	PageSize int //每页显示多少条
	PageNum  int //第几页
}

func Page(pageSize int16, pageNu int16, count int64) (pageNum int16) {
	pageCount := math.Ceil((float64(count) / float64(pageSize)))
	fmt.Println("总页数pageCount", pageCount)
	//	获取当前页码
	pageNum = pageNu
	fmt.Println("获取当前页码pageNum", pageNum)
	if pageNum == 0 || pageNum == -1 {
		pageNum = 1
	}
	return pageNum
}
