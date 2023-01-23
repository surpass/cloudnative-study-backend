package models

import (
	"fmt"
	"time"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql" //导入数据库驱动
)

type Person struct {
	UserId     int       `orm:"column(user_id);pk" json:"userId"`
	Username   string    `orm:"size(100)"column(user_name) json:"userName"`
	Gender     int       `orm:"size(1);column(gender)" json:"gender"`
	Email      string    `json:"email"`
	CreateDate time.Time `orm:"type(date);column(create_date)"  json:"createDate"`
}

var person Person

// func init() {
// 	//注册数据库，一共有五个参数，后面连个用于连接池操作
// 	_ = orm.RegisterDataBase("default", "mysql",
// 		"root:123456@tcp(sql.easyolap.cn:3306)/sport?charset=utf8")
// 	//进行注册模型结构，可以有多个，用逗号分隔
// 	orm.RegisterModel(new(Person))
// 	//创建表，默认为default，只建立一次，后面再执行这个会忽略
// 	_ = orm.RunSyncdb("default", false, true)
// }

func QueryUser(person *Person) []Person {
	// 构造查询
	var personData []Person
	// 获取 QueryBuilder 对象. 需要指定数据库驱动参数。
	// 第二个返回值是错误对象，在这里略过
	qb, _ := orm.NewQueryBuilder("mysql")
	// 构建查询对象
	qb.Select("*").
		From("person").
		Where("user_id > ?").
		Limit(10).Offset(0)
	// 导出 SQL 语句
	sql := qb.String()
	// 执行 SQL 语句
	o := orm.NewOrm()
	o.Raw(sql, 0).QueryRows(&personData)

	fmt.Println(personData)
	return personData
}

// 查询所有学生信息
func SelectAll() []orm.Params {
	o := orm.NewOrm()
	var maps []orm.Params
	_, _ = o.Raw("select * from person").Values(&maps)
	return maps
}

// 根据id查询学生信息
func SelectById(id int) Person {
	o := orm.NewOrm()
	err := o.Raw("select * from person where user_id=?", id).QueryRow(&person)
	if err != nil {
		fmt.Println("get one person error:%v", err)
	}
	return person
}

// 根据id删除用户
func Deletebyid(id int) {
	o := orm.NewOrm()
	err := o.Raw("delete from person where user_id = ?", id).QueryRow(&person)
	if err != nil {
		fmt.Println("del person error:%v", err)
	}
}

// 根据id更新信息
func Updatebyid(id int, person Person) {
	o := orm.NewOrm()
	err := o.Raw("update person set username=? where user_id=?", person.Username, id).QueryRow(&person)
	if err != nil {
		fmt.Println("update person error:%v", err)
	}
}

// 插入数据
func InsertPerson(person Person) int {
	fmt.Println("insert person:%v", person)
	o := orm.NewOrm()
	err := o.Raw("insert into person(user_id,username,gender,email,create_date) values(?,?,?,?,?)",
		person.UserId, person.Username, person.Gender, person.Email, person.CreateDate).QueryRow(&person)
	if err != nil {
		fmt.Println("add person error:%v", err)
	}
	return person.UserId
}

// 根据name进行模糊查询
func SelectUnClear(name string) []orm.Params {
	o := orm.NewOrm()
	var maps []orm.Params
	_, _ = o.Raw("select * from person where user_name like ?", "%"+name+"%").Values(&maps)
	return maps
}
