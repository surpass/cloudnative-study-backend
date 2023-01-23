package models

import (
	"fmt"
	"sportApi/vo"
	"time"

	"github.com/Nerzal/gocloak/v12"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql" //导入数据库驱动
	uuid "github.com/satori/go.uuid"

	"context"
	"log"
)

const LOGIN_TABLE_NAME string = "wx_login"
const WXUSER_TABLE_NAME string = "wx_user"

type WxLogin struct {
	LoginId     string    `orm:"size(100);column(login_id);pk" json:"login_id"`
	Code        string    `orm:"column(code)" json:"code"`
	Mobile      string    `orm:"column(mobile)" json:"mobile"`
	RealName    string    `orm:"column(realName)" json:"realName"`
	StudentCode string    `orm:"column(student_code)" json:"studentCode"` //学号
	UserId      string    `orm:"column(user_id)" json:"userId"`
	OpenId      string    `orm:"column(open_id)" json:"openId"`
	NickName    string    `orm:"size(100);column(nick_name)" json:"nickName"`
	AvatarUrl   string    `orm:"size(500);column(avatar_url)" json:"avatar_url"`
	Language    string    `orm:"size(100);column(language)" json:"language"`
	Province    string    `orm:"size(100);column(province)" json:"province"`
	City        string    `orm:"size(100);column(city)" json:"city"`
	Country     string    `orm:"size(100);column(country)" json:"country"`
	Gender      string    `orm:"size(100);column(gender)" json:"gender"`
	CreateDate  time.Time `orm:"type(timestamp);column(create_date)"  json:"createDate"`
	UpdateDate  time.Time `orm:"type(timestamp);column(update_date)"  json:"updateDate"`
}

func (u *WxLogin) TableName() string {
	return LOGIN_TABLE_NAME
}

type WxUser struct {
	OpenId     string    `orm:"size(100);column(open_id);pk" json:"openId"`
	StudentId  string    `orm:"column(student_id)" json:"studentId"`
	CreateDate time.Time `orm:"type(timestamp);column(create_date)"  json:"createDate"`
	UpdateDate time.Time `orm:"type(timestamp);column(update_date)"  json:"updateDate"`
}

func (u *WxUser) TableName() string {
	return WXUSER_TABLE_NAME
}

type ResWxLogin struct {
	Code       int32      `json:"code"`
	Token      string     `json:"token"`
	OpenId     string     `json:"openId"`
	SessionKey string     `json:"sessionKey"`
	Udata      vo.LoginVO `json:"udata"`
}

func GetOrCreateUser(openID string, login WxLogin) vo.LoginVO {
	var loginVO vo.LoginVO
	//1.用openId查询数据库,确认用户是否存在
	//2.如果存在返回uid
	//3.如果不存在添加一条记录，返回uid
	id := uuid.NewV4()
	timeNow := time.Now()
	currentTime := timeNow.Format("2006-01-02 05:04:03")
	fmt.Println("GetOrCreateUser:", openID, currentTime)
	orm.Debug = true
	o := orm.NewOrm()

	login.LoginId = id.String()
	login.OpenId = openID
	login.CreateDate = timeNow
	login.UpdateDate = timeNow
	//loginsql := fmt.Sprintf(`insert into wx_login(login_id,code,open_id,create_date,update_date) values(?,?,?,?,?)`)
	//lerr := o.Raw(loginsql, login.LoginId, login.Code, login.OpenId, login.CreateDate, login.UpdateDate)
	qs := o.QueryTable(LOGIN_TABLE_NAME)
	i, _ := qs.PrepareInsert()

	loginId, lerr := i.Insert(&login)
	if lerr != nil {
		fmt.Println("wx_login error:%v", lerr)
	} else {
		fmt.Println("wx_login success:%s", loginId)
	}

	// if lerr != nil {
	// 	fmt.Println("wx_login error:%v", lerr)
	// }
	var wxUser WxUser

	//err := o.Raw("select * from wx_user where open_id=?", openID).QueryRow(&wxUser)
	qt := orm.NewOrm().QueryTable(WXUSER_TABLE_NAME)
	// 根据参数拼接条件
	if openID != "" {
		qt = qt.Filter("open_id", openID)
	}
	var student Students
	err := qt.One(&wxUser)
	if err != nil && err == orm.ErrNoRows {
		fmt.Println("get -----one %v error:%v", WXUSER_TABLE_NAME, err, wxUser)

		//通过微信或手机号查询学生信息

		studentSql := fmt.Sprintf(`select * from %s where wx_code=? or phone=? or Code = ?`, STUDENTS_TABLE_NAME)
		serr := o.Raw(studentSql, login.Code, login.Mobile, login.StudentCode).QueryRow(&student)
		if serr != nil && serr == orm.ErrNoRows {
			//loginVO.ErrorMsg = "userNOtBing"
			//return loginVO
			//用户不存在，添加用户
			id2 := uuid.NewV4()
			student.Id = id2.String()
			student.Cid = "default"
			student.Ccid = "default"
			student.Code = login.Mobile
			student.Name = login.Mobile
			student.WxCode = login.Mobile
			student.Phone = login.Mobile
			o.Insert(&student)
			syncCreateCloakUser(login.Mobile)
		}

		wxUser.OpenId = openID
		wxUser.StudentId = student.Id
		wxUser.CreateDate = timeNow
		wxUser.UpdateDate = timeNow
		//err := o.Raw("insert into "+WXUSER_TABLE_NAME+"(open_id,user_id,create_date) values(?,?,?)",
		//	wxUser.OpenId, wxUser.Uid, currentTime) //.QueryRow(&wxUser)
		userid, err := o.Insert(&wxUser)
		if err != nil {
			fmt.Println("add %v error:%v", WXUSER_TABLE_NAME, err)
		} else {
			fmt.Println("add success %v ,%v", WXUSER_TABLE_NAME, userid)
		}
	} else {
		fmt.Println("get one %s no error:%v", WXUSER_TABLE_NAME, wxUser)
		studentSql := fmt.Sprintf(`select * from %s where id=?`, STUDENTS_TABLE_NAME)
		serr := o.Raw(studentSql, wxUser.StudentId).QueryRow(&student)
		if serr != nil && serr == orm.ErrNoRows {
			//loginVO.ErrorMsg = "userNOtBing"
			//return loginVO
			//用户不存在，添加用户
			id2 := uuid.NewV4()
			student.Id = id2.String()
			student.Cid = "default"
			student.Ccid = "default"
			student.Code = login.Mobile
			student.Name = login.Mobile
			student.WxCode = login.Mobile
			student.Phone = login.Mobile
			o.Insert(&student)
			syncCreateCloakUser(login.Mobile)
		}
	}
	loginVO.StudentId = wxUser.StudentId
	loginVO.Cid = student.Cid
	loginVO.Ccid = student.Ccid
	loginVO.Phone = student.Phone
	return loginVO
}

func GetUserByOpenId(openID string) vo.LoginVO {
	var loginVO vo.LoginVO
	//1.用openId查询数据库,确认用户是否存在
	//2.如果存在返回uid
	//3.如果不存在添加一条记录，返回uid

	timeNow := time.Now()
	currentTime := timeNow.Format("2006-01-02 05:04:03")
	fmt.Println("GetUserByOpenId:", openID, currentTime)

	o := orm.NewOrm()

	var wxUser WxUser

	//err := o.Raw("select * from wx_user where open_id=?", openID).QueryRow(&wxUser)
	qt := orm.NewOrm().QueryTable(WXUSER_TABLE_NAME)
	// 根据参数拼接条件
	if openID != "" {
		qt = qt.Filter("open_id", openID)
	}
	var student Students
	err := qt.One(&wxUser)
	if err != nil && err == orm.ErrNoRows {
		fmt.Println("get -----one %v error:%v", WXUSER_TABLE_NAME, err, wxUser)

	} else {
		fmt.Println("get one %v no error:%v", WXUSER_TABLE_NAME, err, wxUser)
		studentSql := fmt.Sprintf(`select * from %s where id=?`, STUDENTS_TABLE_NAME)
		serr := o.Raw(studentSql, wxUser.StudentId).QueryRow(&student)
		if serr != nil && serr == orm.ErrNoRows {
			loginVO.ErrorMsg = "userNOtBing"
			return loginVO
		}
	}
	loginVO.StudentId = student.Id
	loginVO.Cid = student.Cid
	loginVO.Ccid = student.Ccid
	loginVO.Phone = student.Phone
	return loginVO
}

func syncCreateCloakUser(phone string) {
	client := gocloak.NewClient("https://www.easyolap.cn/auth/")
	ctx := context.Background()

	realmName := "master"
	//登陆，输入用户名、密码、领域，返回toekn
	token, err := client.LoginAdmin(ctx, "admin", "123456", realmName)
	if err != nil {
		log.Println("LoginAdmin error: ", err.Error())
		return
	}
	demoRealm := "golang"
	enable := true
	name := phone
	//创建用户
	str, err := client.CreateUser(ctx, token.AccessToken, demoRealm, gocloak.User{
		Username: &name,
		Enabled:  &enable,
	})
	if err != nil {
		log.Println("CreateUser error: ", err.Error())
		return
	}
	log.Println("CreateUser str: ", str)
	//设置密码
	err = client.SetPassword(ctx, token.AccessToken, str, demoRealm, "123456", false)
	if err != nil {
		log.Println("SetPassword error: ", err.Error())
		return
	}
	log.Println("SetPassword success")

}
