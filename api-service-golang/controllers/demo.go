package controllers

import (
	"encoding/json"
	"fmt"
	"sportApi/models"

	beego "github.com/beego/beego/v2/server/web"
)

// Operations about Users
type DemoController struct {
	beego.Controller
}

type Categories struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Icon string `json:"icon"`
}

type Slides struct {
	Id    string `json:"id"`
	Image string `json:"image"`
	Link  string `json:"link"`
}

var (
	SlidesList []Slides
)

var (
	CategoriesList []Categories
)

// @Title get Slides
// @Description create users
// @Success 200 {int} models.User.Id
// @Failure 403 body is empty
// @router /slides [Get]
func (l *DemoController) Slides() {

	SlidesList = make([]Slides, 2, 2)
	l1 := Slides{"1", "https://www.easyolap.cn/images/banner01.png", ""}
	SlidesList[0] = l1
	l2 := Slides{"2", "https://www.easyolap.cn/images/banner02.png", "/pages/list/list?cat=10"}
	SlidesList[1] = l2

	l.Data["json"] = SlidesList
	l.ServeJSON()
}

// @Title categories
// @Description get all categories
// @Success 200 {object} models.Location
// @router /categories [get]
func (l *DemoController) Categories() {
	CategoriesList = make([]Categories, 9, 9)
	l1 := Categories{"1", "美食", "http://ww1.sinaimg.cn/large/006ThXL5ly1fj8w5i2onyj302u02umwz.jpg"}
	CategoriesList[0] = l1
	l2 := Categories{"2", "快递", "http://ww1.sinaimg.cn/large/006ThXL5ly1fj8w5i2j4dj302u02umwy.jpg"}
	CategoriesList[1] = l2
	l3 := Categories{"3", "Test", "http://ww1.sinaimg.cn/large/006ThXL5ly1fj8w5i56i0j302u02u744.jpg"}
	CategoriesList[2] = l3
	l4 := Categories{"4", "蓝球", "http://ww1.sinaimg.cn/large/006ThXL5ly1fj8w5i2uzvj302u02udfo.jpg"}
	CategoriesList[3] = l4
	l5 := Categories{"5", "找工作", "http://ww1.sinaimg.cn/large/006ThXL5ly1fj8w5i2rnlj302u02umwz.jpg"}
	CategoriesList[4] = l5
	l6 := Categories{"6", "辅导班", "http://ww1.sinaimg.cn/large/006ThXL5ly1fj8w5i2zloj302u02udfn.jpg"}
	CategoriesList[5] = l6
	l7 := Categories{"7", "图书馆", "http://ww1.sinaimg.cn/large/006ThXL5ly1fj8w5i69eij302u02ua9w.jpg"}
	CategoriesList[6] = l7
	l8 := Categories{"8", "租房", "http://ww1.sinaimg.cn/large/006ThXL5ly1fj8w5i6j2lj302u02u0sj.jpg"}
	CategoriesList[7] = l8
	l9 := Categories{"9", "维修", "http://ww1.sinaimg.cn/large/006ThXL5ly1fj8w5i6z1pj302u02ua9u.jpg"}
	CategoriesList[8] = l9

	l.Data["json"] = CategoriesList
	l.ServeJSON()
}

// 定义一个函数
func (c *DemoController) Index() {
	var ob models.Person
	json.Unmarshal(c.Ctx.Input.RequestBody, &ob)
	flag := models.QueryUser(&ob) // 调用models中的QueryUser方法
	fmt.Print(flag)
	c.Data["json"] = flag
	c.ServeJSON()
}
