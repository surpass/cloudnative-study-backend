package models

import (
	"context"
	"time"

	"log"

	gocloak "github.com/Nerzal/gocloak/v12"
	beego "github.com/beego/beego/v2/server/web"
)

type UserType int

const (
	ApplicationAdminType UserType = iota
	SuperAdminType                //超级管理员
	UnAuthorizedUserType          //未授权用户
)

var (
	keycloakUrl, _  = beego.AppConfig.String("keycloakUrl")
	clientID, _     = beego.AppConfig.String("clientID")
	clientSecret, _ = beego.AppConfig.String("clientSecret")
	realmName, _    = beego.AppConfig.String("realmName")
	defaultPwd, _   = beego.AppConfig.String("defaultPwd")

	adminRealm, _ = beego.AppConfig.String("adminRealm")
	adminUser, _  = beego.AppConfig.String("adminUser")
	adminPwd, _   = beego.AppConfig.String("adminPwd")

	userId    string
	client    = gocloak.NewClient(keycloakUrl)
	clientJWT *gocloak.JWT
	// retrospecTokenResult存放了clientJWT的过期时间(Exp)等
	retrospecTokenResult *gocloak.IntroSpectTokenResult
)

func init() {
	updateClient()
}

/**
 * 登录到client，并获取clientJWT与retrospecTokenResult。
 */
func updateClient() {
	var err error

	clientJWT, err = client.LoginClient(context.Background(), clientID, clientSecret, realmName)
	if err != nil || clientJWT == nil {
		log.Fatalf("failed to Login KeyCloak Client: %v", err)
		panic("failed to Login KeyCloak Client ")
	}
	log.Println("clientJWT is :", clientJWT)

	retrospecTokenResult, err = client.RetrospectToken(context.Background(), clientJWT.AccessToken, clientID, clientSecret, realmName)
	if err != nil || retrospecTokenResult == nil {
		log.Fatalf("failed to Retrospect KeyCloak Token: %v", err)
		panic("failed to Retrospect KeyCloak Token ")
	}
	log.Println("retrospecTokenResult is:", retrospecTokenResult)

}

// 获取用户基础信息
func GetUserInfo(accessToken string) (user *gocloak.UserInfo, err error) {
	user, err = client.GetUserInfo(context.Background(), accessToken, realmName)
	return
}

// 获取用户角色信息
func GetUserRole(userId string) (userType UserType, err error) {
	userType = UnAuthorizedUserType
	//判断是否过期
	if int64(*retrospecTokenResult.Exp) < time.Now().Unix() {
		updateClient()
	}
	mappingsRepresentation, err := client.GetRoleMappingByUserID(context.Background(), clientJWT.AccessToken, realmName, userId)
	for _, v := range *mappingsRepresentation.RealmMappings {
		if *v.Name == "demo_admin_role" {
			userType = SuperAdminType
			break
		}
		if *v.Name == "demo_user_role" {
			userType = ApplicationAdminType
		}
	}
	return userType, err
}

// get JWT token by username and passwd
func GetToken(username string, passwd string) (*gocloak.JWT, error) {
	loginToken, err := client.GetToken(
		context.Background(),
		realmName,
		gocloak.TokenOptions{
			ClientID:      &clientID,
			ClientSecret:  &clientSecret,
			Username:      &username,
			Password:      &passwd,
			GrantType:     gocloak.StringP("password"),
			ResponseTypes: &[]string{"token", "id_token"},
			Scopes:        &[]string{"openid", "offline_access"},
		},
	)
	if err != nil {
		panic("Inspection failed:" + err.Error())
	}
	log.Println("token:", loginToken.AccessToken)
	log.Println("RefreshToken:", loginToken.RefreshToken)
	return loginToken, err
}

func RefreshToken(refreshToken string) (*gocloak.JWT, error) {
	token, err := client.RefreshToken(
		context.Background(),
		refreshToken,
		clientID,
		clientSecret,
		realmName)
	if err != nil {
		log.Println("Inspection failed:" + err.Error())
	}

	log.Println("refresh result  Token : ", token)
	log.Println("refresh token result token is :", token.AccessToken)

	return token, err
}

// login admin
func GetAdminToken() (*gocloak.JWT, error) {
	//登陆，输入用户名、密码、领域，返回toekn
	token, err := client.LoginAdmin(context.Background(), adminUser, adminPwd, adminRealm)
	if err != nil {
		log.Println("LoginAdmin error: ", err.Error())
	}
	return token, err
}

// create keycloak user
func CreateUser(username string) (*gocloak.User, error) {
	token, _ := GetAdminToken()

	user := gocloak.User{
		Username: &username,
		Enabled:  gocloak.BoolP(true),
		Attributes: &map[string][]string{
			"foo": {"bar", "alice", "bob", "roflcopter"},
			"bar": {"baz"},
		},
	}

	userID, err := client.CreateUser(
		context.Background(),
		token.AccessToken,
		realmName,
		user)
	user.ID = &userID

	//设置密码
	err = client.SetPassword(context.Background(), token.AccessToken, userID, realmName, defaultPwd, false)
	if err != nil {
		log.Println("SetPassword error: ", err.Error())
	}
	return &user, err
}

func VerifyAccessToken(accessToken string) bool {
	rptResult, err := client.RetrospectToken(context.Background(), accessToken, clientID, clientSecret, realmName)
	if err != nil {
		panic("Inspection failed:" + err.Error())
	}
	log.Println("========RptResult:", rptResult)
	if !*rptResult.Active {
		return false
	}
	return true
}
