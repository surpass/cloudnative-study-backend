package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sportApi/models"
	"sportApi/utils"
	"strings"

	beego "github.com/beego/beego/v2/server/web"
)

// Operations about Users
type UserController struct {
	beego.Controller
}

// @Title CreateUser
// @Description create users
// @Param	body		body 	models.User	true		"body for user content"
// @Success 200 {int} models.User.Id
// @Failure 403 body is empty
// @router / [post]
func (u *UserController) Post() {
	var user models.User
	json.Unmarshal(u.Ctx.Input.RequestBody, &user)
	uid := models.AddUser(user)
	u.Data["json"] = map[string]string{"uid": uid}
	u.ServeJSON()
}

// @Title GetAll
// @Description get all Users
// @Success 200 {object} models.User
// @router / [get]
func (u *UserController) GetAll() {
	users := models.GetAllUsers()
	u.Data["json"] = users
	u.ServeJSON()
}

// @Title Get
// @Description get user by uid
// @Param	uid		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.User
// @Failure 403 :uid is empty
// @router /:uid [get]
func (u *UserController) Get() {
	uid := u.GetString(":uid")
	if uid != "" {
		user, err := models.GetUser(uid)
		if err != nil {
			u.Data["json"] = err.Error()
		} else {
			u.Data["json"] = user
		}
	}
	u.ServeJSON()
}

// @Title Update
// @Description update the user
// @Param	uid		path 	string	true		"The uid you want to update"
// @Param	body		body 	models.User	true		"body for user content"
// @Success 200 {object} models.User
// @Failure 403 :uid is not int
// @router /:uid [put]
func (u *UserController) Put() {
	uid := u.GetString(":uid")
	if uid != "" {
		var user models.User
		json.Unmarshal(u.Ctx.Input.RequestBody, &user)
		uu, err := models.UpdateUser(uid, &user)
		if err != nil {
			u.Data["json"] = err.Error()
		} else {
			u.Data["json"] = uu
		}
	}
	u.ServeJSON()
}

// @Title Delete
// @Description delete the user
// @Param	uid		path 	string	true		"The uid you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 uid is empty
// @router /:uid [delete]
func (u *UserController) Delete() {
	uid := u.GetString(":uid")
	models.DeleteUser(uid)
	u.Data["json"] = "delete success!"
	u.ServeJSON()
}

// @Title Login
// @Description Logs user into the system
// @Param	username		query 	string	true		"The username for login"
// @Param	password		query 	string	true		"The password for login"
// @Success 200 {string} login success
// @Failure 403 user not exist
// @router /login [get]
func (u *UserController) Login() {
	username := u.GetString("username")
	password := u.GetString("password")
	if models.Login(username, password) {
		u.Data["json"] = "login success"
	} else {
		u.Data["json"] = "user not exist"
	}
	u.ServeJSON()
}

// @Title logout
// @Description Logs out current logged in user session
// @Success 200 {string} logout success
// @router /logout [get]
func (u *UserController) Logout() {
	u.Data["json"] = "logout success"
	u.ServeJSON()
}

// Operations about login
type LoginController struct {
	beego.Controller
}

func (l *LoginController) Loginwx() {
	var login models.WxLogin
	json.Unmarshal(l.Ctx.Input.RequestBody, &login)
	code := login.Code
	fmt.Println("wei xin login form info:%v", login)
	appID, _ := beego.AppConfig.String("appID")
	secret, _ := beego.AppConfig.String("secret")
	codeToSessURL, _ := beego.AppConfig.String("codeToSessURL")

	fmt.Println("wei xin login info:%v,%v", appID, codeToSessURL)

	codeToSessURL = strings.Replace(codeToSessURL, "{appid}", appID, -1)
	codeToSessURL = strings.Replace(codeToSessURL, "{secret}", secret, -1)
	codeToSessURL = strings.Replace(codeToSessURL, "{code}", code, -1)

	resp, err := http.Get(codeToSessURL)
	if err != nil {
		fmt.Println(err.Error())
		l.Data["json"] = err.Error()
		l.ServeJSON()
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		l.Data["json"] = resp.StatusCode
		l.ServeJSON()
		return
	}
	fmt.Println("wei xin login sucess body:%v", resp.Body)

	var data map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		fmt.Println(err.Error())
		l.Data["json"] = err.Error()
		l.ServeJSON()
		return
	}

	if _, ok := data["session_key"]; !ok {
		fmt.Println("session_key 不存在")
		fmt.Println(data)
		l.Data["json"] = "session_key 不存在"
		l.ServeJSON()
		return
	}

	var openID string
	var sessionKey string
	openID = data["openid"].(string)
	sessionKey = data["session_key"].(string)

	encryptedData := l.GetString("encryptedData")
	iv := l.GetString("iv")
	fmt.Println("get phone number sucess:%v,%v", encryptedData, iv)

	fmt.Println("wei xin login sucess:%s,%s", openID, sessionKey)

	loginVO := models.GetOrCreateUser(openID, login)

	defaultPwd, _ := beego.AppConfig.String("defaultPwd")
	loginToken, err := models.GetToken(loginVO.Phone, defaultPwd)
	if err != nil {
		log.Fatalf("failed to KeyCloak Token: %v", err)
	}
	loginVO.JwtToken = loginToken

	l.SetSession("token", loginToken)
	l.SetSession("openID", openID)
	l.SetSession("udata", loginVO)
	l.SetSession("sessionKey", sessionKey)

	var res models.ResWxLogin
	res.Code = 200
	res.OpenId = openID
	res.SessionKey = sessionKey
	res.Udata = loginVO

	l.Data["json"] = res
	l.ServeJSON()
}

func (l *LoginController) GetPhonenumber() {
	sessionID := l.GetString("sessionID")
	encryptedData := l.GetString("encryptedData")
	iv := l.GetString("iv")

	src, err := utils.Dncrypt(sessionID, encryptedData, iv)
	fmt.Println(err)
	var s = map[string]interface{}{}
	json.Unmarshal([]byte(src), &s)
	fmt.Printf("== %+v", src)
	fmt.Printf("cc== %+v", s)

	//var res map[string][]*string
	//res = make(map[string][]*string, 3)
	//res["phoneNumber"] = s[phoneNumber]

	l.Data["json"] = s
	l.ServeJSON()
}

func (l *LoginController) AutoLoginwxByOpenid() {
	openID := l.GetString("openId")
	loginVO := models.GetUserByOpenId(openID)

	defaultPwd, _ := beego.AppConfig.String("defaultPwd")
	loginToken, err := models.GetToken(loginVO.Phone, defaultPwd)
	if err != nil {
		log.Fatalf("failed to KeyCloak Token: %v", err)
	}
	loginVO.JwtToken = loginToken

	l.SetSession("token", loginToken)
	l.SetSession("openID", openID)
	l.SetSession("udata", loginVO)

	var res models.ResWxLogin
	res.Code = 200
	res.OpenId = openID
	res.Udata = loginVO

	l.Data["json"] = res
	l.ServeJSON()
}
