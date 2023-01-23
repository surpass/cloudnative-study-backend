package controllers

import (
	"fmt"
	"log"
	"sportApi/models"

	gocloak "github.com/Nerzal/gocloak/v12"
	beego "github.com/beego/beego/v2/server/web"
)

// Operations about Users
type OauthController struct {
	beego.Controller
}

// login by username and passwd
func (l *OauthController) GetToken() {
	username := l.GetString("username")
	passwd := l.GetString("passwd")

	loginToken, err := models.GetToken(username, passwd)
	if err != nil {
		log.Fatalf("failed to KeyCloak Token: %v", err)
	}
	log.Println("token:", loginToken.AccessToken)
	l.Data["json"] = loginToken
	l.ServeJSON()
}

type TokenRes struct {
	Code  int          `json:"code"`
	Datas *gocloak.JWT `json:"datas"`
}

// login by username and passwd
func (l *OauthController) RefreshToken() {
	refreshToken := l.GetString("refreshToken")

	fmt.Println("refreshToken:%s", refreshToken)

	loginToken, err := models.RefreshToken(refreshToken)
	if err != nil {
		log.Fatalf("failed to refresh KeyCloak Token: %v", err)
	}
	log.Println("token:", loginToken.AccessToken)
	var res TokenRes
	res.Code = 200
	res.Datas = loginToken
	l.Data["json"] = res
	l.ServeJSON()
}

// login by username and passwd
func (l *OauthController) VerifyAccessToken() {
	accessToken := l.GetString("accessToken")

	booleanResult := models.VerifyAccessToken(accessToken)

	log.Println("token:%s,result:%v", accessToken, booleanResult)
	l.Data["json"] = booleanResult
	l.ServeJSON()
}
