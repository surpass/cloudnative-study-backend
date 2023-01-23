package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/orm"
	beego "github.com/beego/beego/v2/server/web"
	"sportApi/models"
	"time"
)

const (
	YYYYMMDDHHMISS = "2006-01-02 15:04:05" //常规类型
)

type JSONTime struct {
	time.Time
}

// MarshalJSON on JSONTime format Time field with %Y-%m-%d %H:%M:%S
func (t JSONTime) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf("\"%s\"", t.Format(YYYYMMDDHHMISS))
	return []byte(formatted), nil
}

// Value insert timestamp into mysql need this function.
// func (t JSONTime) Value() (driver.Value, error) {
// 	var zeroTime time.Time
// 	if t.Time.UnixNano() == zeroTime.UnixNano() {
// 		return nil, nil
// 	}
// 	return t.Time, nil
// }

// Scan valueof time.Time
func (t *JSONTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = JSONTime{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// The time is expected to be a quoted string in RFC 3339 format.
func (t *JSONTime) UnmarshalJSON(data []byte) error {
	// Ignore null, like in the main JSON package.
	if string(data) == "null" {
		return nil
	}
	// Fractional seconds are handled implicitly by Parse.
	var err error
	(*t).Time, err = time.ParseInLocation(`"`+YYYYMMDDHHMISS+`"`, string(data), time.Local)
	return err
}

type PersonVo struct {
	UserId     int       `json:"userId"`
	Username   string    `json:"userName"`
	Gender     int       `json:"gender"`
	Email      string    `json:"email"`
	CreateDate time.Time `json:"createDate"`
}

type SelectAll struct {
	beego.Controller
}

//全局变量
var stu models.Person

//查询所有信息
func (c *SelectAll) Get() {
	var maps []orm.Params
	maps = models.SelectAll()
	c.Data["json"] = maps
	c.ServeJSON()
}

//查询所有信息
func (c *SelectAll) QueryUser() {
	var ob models.Person
	json.Unmarshal(c.Ctx.Input.RequestBody, &ob)
	persons := models.QueryUser(&ob)
	c.Data["json"] = persons
	c.ServeJSON()
}

type SelectbyId struct {
	beego.Controller
}

//根据id查询用户
func (c *SelectbyId) GetById() {
	id, _ := c.GetInt(":id")

	idp := c.Ctx.Input.Param(":id")

	fmt.Println("select by id :%v,%v", id, idp)
	stu = models.SelectById(id)
	c.Data["json"] = stu
	c.ServeJSON()
}

//根据id删除用户
type DeletebyId struct {
	beego.Controller
}

func (c *DeletebyId) Get() {
	id, _ := c.GetInt(":id")
	models.Deletebyid(id)
	str := models.SelectAll()
	c.Data["json"] = str
	c.ServeJSON()
}

//根据id更新信息
type Updatebyid struct {
	beego.Controller
}

func (c *Updatebyid) UpdatePerson() {

	var person models.Person
	var personVo PersonVo

	json.Unmarshal(c.Ctx.Input.RequestBody, &personVo)
	person.UserId = personVo.UserId
	person.Username = personVo.Username
	fmt.Println("update person:%v,%v", person.UserId, person.Username)
	fmt.Println("update person:%v", person)

	models.Updatebyid(person.UserId, person)
	stu = models.SelectById(person.UserId)
	c.Data["json"] = stu
	c.ServeJSON()
}

//插入数据
type InsertData struct {
	beego.Controller
}

// @Title Create person
// @Description create person
// @Param	body		body 	models.Person	true		"body for person content"
// @Success 200 {person} models.Person
// @Failure 403 body is empty
// @router /add [post]
func (c *InsertData) Post() {
	fmt.Println("=====InsertData=======")
	//fmt.Println("request body:%v", c.Ctx.Input.RequestBody)
	var person models.Person
	var personVo PersonVo

	json.Unmarshal(c.Ctx.Input.RequestBody, &personVo)
	person.UserId = personVo.UserId
	person.Username = personVo.Username
	person.Gender = personVo.Gender
	person.Email = personVo.Email
	person.CreateDate = personVo.CreateDate
	fmt.Println("add person:%v,%v", person.UserId, person.Username)
	fmt.Println("add person:%v", person)

	userId := models.InsertPerson(person)
	c.Data["json"] = userId
	c.ServeJSON()
}

//模糊查询
type SelectUnClear struct {
	beego.Controller
}

func (c *SelectUnClear) Get() {
	name := c.GetString(":name")
	var maps []orm.Params
	maps = models.SelectUnClear(name)
	c.Data["json"] = maps
	c.ServeJSON()
}
