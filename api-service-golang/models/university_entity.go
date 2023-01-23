package models

import (
	"fmt"
	"time"

	"github.com/astaxie/beego/orm"
	uuid "github.com/satori/go.uuid"
)

func init() {

	intDefault()

	var u University
	id := uuid.NewV4()
	o := orm.NewOrm()
	u.Name = "沈阳工程学院"
	u.Code = "EA"
	uSql := fmt.Sprintf(`select * from %s where name=? or code = ? `, UNIVERSITY_TABLE_NAME)
	err := o.Raw(uSql, u.Name, u.Code).QueryRow(&u)
	if err != nil && err == orm.ErrNoRows {
		fmt.Println("查询不到，进行初始化")
		u.Id = id.String()

		u.Contacts = "Liu Mr."

		uId, err := o.Insert(&u)

		if err != nil {
			fmt.Println("add %s error:%s", UNIVERSITY_TABLE_NAME, err)
		} else {
			fmt.Println("add success %s ,%s", UNIVERSITY_TABLE_NAME, uId)
		}
	} else if err == nil {
		fmt.Println("已初始化，跳过")
		return
	} else {
		fmt.Println("error:%s", err)
	}

	var c Colleges
	c.Id = id.String()
	c.Uid = u.Id
	c.Code = "Info"
	c.Name = "信息工程"
	cId, cerr := o.Insert(&c)

	if cerr != nil {
		fmt.Println("add %s error:%s", COLLEGES_TABLE_NAME, cerr)
	} else {
		fmt.Println("add success ", COLLEGES_TABLE_NAME, cId)
	}

	var cc Classes
	cc.Id = id.String()
	cc.Cid = c.Id
	cc.Code = "992"
	cc.Name = "九九二"
	ccId, ccerr := o.Insert(&cc)

	if ccerr != nil {
		fmt.Println("add %s error:%s", CLASSES_TABLE_NAME, ccerr)
	} else {
		fmt.Println("add success ", CLASSES_TABLE_NAME, ccId)
	}

	id2 := uuid.NewV4()
	var cc2 Classes
	cc2.Id = id2.String()
	cc2.Cid = c.Id
	cc2.Code = "002"
	cc2.Name = "零零二"
	o.Insert(&cc2)

	var s Students
	s.Id = id.String()
	s.Cid = c.Id
	s.Ccid = cc.Id
	s.Code = "N01"
	s.Name = "老六"
	s.WxCode = "LJC7339901"
	s.Phone = "LJC7339901"
	o.Insert(&s)

	CreateUser(s.Phone)

	var s2 Students
	s2.Id = id2.String()
	s2.Cid = c.Id
	s2.Ccid = cc.Id
	s2.Code = "N02"
	s2.Name = "老8"
	s2.WxCode = "s46488820"
	s2.Phone = "13478875695"
	o.Insert(&s2)
	CreateUser(s2.Phone)
}
func intDefault() {

	var udefault University
	iddefault := uuid.NewV4()
	o := orm.NewOrm()
	udefault.Name = "中国"
	udefault.Code = "CN"
	uSql := fmt.Sprintf(`select * from %s where name=? or code = ? `, UNIVERSITY_TABLE_NAME)
	defaulterr := o.Raw(uSql, udefault.Name, udefault.Code).QueryRow(&udefault)
	if defaulterr != nil && defaulterr == orm.ErrNoRows {
		fmt.Println("查询不到，进行初始化")
		udefault.Id = iddefault.String()

		udefault.Contacts = "Li Mr."

		o.Insert(&udefault)

	} else if defaulterr == nil {
		fmt.Println("default已初始化，跳过")
		return
	} else {
		fmt.Println("error:%s", defaulterr)
	}

	var cdefault Colleges
	cdefault.Id = "default"
	cdefault.Uid = udefault.Id
	cdefault.Code = "Info"
	cdefault.Name = "社会大学"
	cIddefault, cdefaulterr := o.Insert(&cdefault)

	if cdefaulterr != nil {
		fmt.Println("add %s error:%s", COLLEGES_TABLE_NAME, cdefaulterr)
	} else {
		fmt.Println("add success ", COLLEGES_TABLE_NAME, cIddefault)
	}

	var ccdefault Classes
	ccdefault.Id = "default"
	ccdefault.Cid = cdefault.Id
	ccdefault.Code = "china"
	ccdefault.Name = "运动达人"
	ccIddefault, ccdefaulterr := o.Insert(&ccdefault)

	if ccdefaulterr != nil {
		fmt.Println("add %s error:%s", CLASSES_TABLE_NAME, ccdefaulterr)
	} else {
		fmt.Println("add success ", CLASSES_TABLE_NAME, ccIddefault)
	}

}

const UNIVERSITY_TABLE_NAME string = "u_university"
const COLLEGES_TABLE_NAME string = "u_colleges"
const CLASSES_TABLE_NAME string = "u_classes"
const STUDENTS_TABLE_NAME string = "u_students"
const TRIP_TABLE_NAME = "u_trip"
const RANK_TABLE_NAME = "u_ranks"

type University struct {
	Id         string    `orm:"size(100);column(id);pk" json:"id"`
	Code       string    `orm:"null;size(100);column(code)" description:(codeornum) json:"code"`
	Name       string    `orm:"null;size(100);column(name)" json:"name"`
	Contacts   string    `orm:"null;size(50);column(contacts)" json:"contacts"`
	Add        string    `orm:"null;size(500);column(add)" json:"add"`
	Phone      string    `orm:"null;size(200);column(phone)" json:"phone"`
	Email      string    `orm:"null;size(100);column(email)" json:"email"`
	CreateUser string    `orm:"null;size(100);column(create_user)" null json:"createUser"`
	CreateDate time.Time `orm:"null;auto_now_add;type(timestamp);column(create_date)"  json:"createDate"`
	UpdateUser string    `orm:"null;size(50);column(update_user)" json:"updateUser"`
	UpdateDate time.Time `orm:"null;auto_now;type(timestamp);column(update_date)"  json:"updateDate"`
	//Colleges   []*Colleges `orm:"reverse(many)"` // 设置一对多的反向关系
}

func (u *University) TableName() string {
	return UNIVERSITY_TABLE_NAME
}

type Colleges struct {
	Id         string    `orm:"size(100);column(id);pk" json:"id"`
	Uid        string    `orm:"column(uid)" json:"uid"`
	Code       string    `orm:"null;size(100)"column(code) json:"code"`
	Name       string    `orm:"null;size(100)"column(name) json:"name"`
	Contacts   string    `orm:"null;size(50)"column(contacts) json:"contacts"`
	Add        string    `orm:"null;size(500)"column(add) json:"add"`
	Phone      string    `orm:"null;size(200)"column(phone) json:"phone"`
	Email      string    `orm:"null;size(100)"column(email) json:"email"`
	CreateUser string    `orm:"null;size(100);column(create_user)" null json:"createUser"`
	CreateDate time.Time `orm:"null;auto_now_add;type(timestamp);column(create_date)"  json:"createDate"`
	UpdateUser string    `orm:"null;size(50);column(update_user)" json:"updateUser"`
	UpdateDate time.Time `orm:"null;auto_now;type(timestamp);column(update_date)"  json:"updateDate"`
	//Classes    []*Classes `orm:"reverse(many)"` // 设置一对多的反向关系
}

func (u *Colleges) TableName() string {
	return COLLEGES_TABLE_NAME
}

type Classes struct {
	Id         string    `orm:"size(100);column(id);pk" json:"id"`
	Cid        string    `orm:"column(cid);" json:"cid"`
	Code       string    `orm:"null;size(100)"column(code) json:"code"`
	Name       string    `orm:"null;size(100)"column(name) json:"name"`
	Contacts   string    `orm:"null;size(50)"column(contacts) json:"contacts"`
	Add        string    `orm:"null;size(500)"column(add) json:"add"`
	Phone      string    `orm:"null;size(200)"column(phone) json:"phone"`
	Email      string    `orm:"null;size(100)"column(email) json:"email"`
	FromDate   time.Time `orm:"null;type(timestamp);column(from_date)"  json:"fromDate"`
	ThruDate   time.Time `orm:"null;type(timestamp);column(thru_date)"  json:"thruDate"`
	CreateUser string    `orm:"null;size(100);column(create_user)" null json:"createUser"`
	CreateDate time.Time `orm:"null;auto_now_add;type(timestamp);column(create_date)"  json:"createDate"`
	UpdateUser string    `orm:"null;size(50);column(update_user)" json:"updateUser"`
	UpdateDate time.Time `orm:"null;auto_now;type(timestamp);column(update_date)"  json:"updateDate"`
	//Student    []*Students `orm:"reverse(many)"` // 设置一对多的反向关系
}

func (u *Classes) TableName() string {
	return CLASSES_TABLE_NAME
}

type Students struct {
	Id         string    `orm:"size(100);column(id);pk" json:"id"`
	Cid        string    `orm:"column(cid);" json:"cid"`
	Ccid       string    `orm:"column(ccid);" json:"ccid"`
	Code       string    `orm:"null;size(100)"column(code) json:"code"`
	Name       string    `orm:"null;size(100)"column(name) json:"name"`
	Phone      string    `orm:"null;size(200)"column(phone) json:"phone"`
	WxCode     string    `orm:"null;size(200)"column(wx_code) json:"wxCode"`
	Email      string    `orm:"null;size(100)"column(email) json:"email"`
	FromDate   time.Time `orm:"null;type(timestamp);column(from_date)"  json:"fromDate"`
	ThruDate   time.Time `orm:"null;type(timestamp);column(thru_date)"  json:"thruDate"`
	CreateUser string    `orm:"null;size(100);column(create_user)" null json:"createUser"`
	CreateDate time.Time `orm:"null;auto_now_add;type(timestamp);column(create_date)"  json:"createDate"`
	UpdateUser string    `orm:"null;size(50);column(update_user)" json:"updateUser"`
	UpdateDate time.Time `orm:"null;auto_now;type(timestamp);column(update_date)"  json:"updateDate"`
}

func (u *Students) TableName() string {
	return STUDENTS_TABLE_NAME
}

type SqlTrip struct {
	Id         string    `orm:"size(100);column(id);pk" json:"id"`
	Ccid       string    `orm:"column(ccid);" json:"ccid"` // 班级id
	Cid        string    `orm:"null;size(100);column(cid);description:(Colleges id);" description:(Colleges id) json:"cid"`
	Sid        string    `orm:"null;size(100);column(sid);description:(student id);" description:(student id) json:"sid"`
	Dkey       string    `orm:"null;size(50);column(d_key);;description:(data encode key);" json:"dkey"`
	FromDate   time.Time `orm:"null;type(timestamp);column(from_date)"  json:"fromDate"`
	ThruDate   time.Time `orm:"null;type(timestamp);column(thru_date)"  json:"thruDate"`
	Flon       string    `orm:"null;size(100);column(f_lon);"json:"flon"`
	Flat       string    `orm:"null;size(100);column(f_lat);" json:"flat"`
	Tlon       string    `orm:"null;size(100);column(t_lon);" json:"tlon"`
	Tlat       string    `orm:"null;size(100);column(t_lat);" json:"tlat"`
	Distance   float64   `orm:"null;digits(15);decimals(4);column(distance);" json:"distance"`
	Runtime    int64     `orm:"null;size(10);column(run_time);" json:"runTime"`
	AvgeSpeed  float64   `orm:"null;digits(12);decimals(4);column(avge_speed);" json:"avgeSpeed"`
	MaxSpeed   float64   `orm:"null;digits(12);decimals(4);column(max_speed);" json:"maxSpeed"`
	CreateUser string    `orm:"null;size(100);column(create_user)" null json:"createUser"`
	CreateDate time.Time `orm:"null;auto_now_add;type(timestamp);column(create_date)"  json:"createDate"`
	UpdateUser string    `orm:"null;size(50);column(update_user)" json:"updateUser"`
	UpdateDate time.Time `orm:"null;auto_now;type(timestamp);column(update_date)"  json:"updateDate"`
}

func (u *SqlTrip) TableName() string {
	return TRIP_TABLE_NAME
}

type RankEngity struct {
	Id         string    `orm:"size(100);column(id);pk" json:"id"`
	Year       string    `orm:"size(100);column(year);" json:"year"`
	Month      string    `orm:"size(100);column(month);" json:"month"`
	Day        string    `orm:"null;size(100);column(day);" json:"day"`
	Week       string    `orm:"null;size(50);column(week);;description:(week field);" json:"week"`
	SId        string    `orm:"size(100);column(sid);" json:"sid"`
	Cid        string    `orm:"null;column(cid);" json:"cid"`
	Ccid       string    `orm:"null;column(ccid);" json:"ccid"`
	Distance   float64   `orm:"null;digits(15);decimals(4);column(distance);" json:"distance"`
	Runtime    int64     `orm:"null;size(10);column(run_time);" json:"runTime"`
	Rank       int       `orm:"column(rank);" json:"rank"`
	Types      int       `orm:"column(types);" json:"types"`
	CreateDate time.Time `orm:"auto_now_add;type(timestamp);column(create_date)"  json:"createDate"`
	UpdateDate time.Time `orm:"null;auto_now;type(timestamp);column(update_date)"  json:"updateDate"`
}

func (u *RankEngity) TableName() string {
	return RANK_TABLE_NAME
}
