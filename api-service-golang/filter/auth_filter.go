package filter

import (
	"log"
	"net/http"
	"sportApi/models"
	"strings"

	context "github.com/beego/beego/v2/server/web/context"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/filter/cors"
)

// 初始化，添加路由鉴权
func init() {
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		AllowCredentials: true,
	}))
	log.Println("beego.InsertFilter before")

	//ctx := context.Background()
	beego.InsertFilter("/v1/location/*", beego.BeforeRouter, authApi)
	log.Println("beego.InsertFilter after")

	//....
}

// 校验函数
func authApi(ctx *context.Context) {
	log.Println("beego authApi before===========")
	AccessToken := ctx.Input.Header("Authorization")
	if AccessToken == "" {
		ctx.Output.JSON(map[string]interface{}{
			"status": http.StatusUnauthorized, "description": http.StatusText(http.StatusUnauthorized)},
			false, false)
		return
	}

	if strings.Index(AccessToken, "Bearer ") == 0 {
		AccessToken = AccessToken[7:]
	}
	//log.Println("====AccessToken:", AccessToken)

	if models.VerifyAccessToken(AccessToken) {
		//log.Println("====AccessToken 有效")
	} else {
		log.Println("====AccessToken 无效")
	}

	UserInfo, err := models.GetUserInfo(AccessToken)
	if err != nil || UserInfo == nil {
		log.Println(err)
		ctx.Output.JSON(map[string]interface{}{
			"status": http.StatusUnauthorized, "description": http.StatusText(http.StatusUnauthorized)},
			false, false)
		return
	}
	log.Println("====UserInfo:", UserInfo)
	/*
			userType, err := models.GetUserRole(*UserInfo.Sub)
			if err != nil || userType == models.UnAuthorizedUserType {
				log.Println(err)
				ctx.Output.JSON(map[string]interface{}{
					"status": http.StatusUnauthorized, "description": http.StatusText(http.StatusUnauthorized) + " 未查询到使用授权信息。"},
					false, false)
				return
			}
		ctx.Input.SetData("UserType", userType)
	*/
	ctx.Input.SetData("UserId", *UserInfo.Sub)
	ctx.Input.SetData("UserName", *UserInfo.PreferredUsername)

	// 具体业务代码，如获取用户名、根据用户角色进行不同的鉴权处理。
	// ....
}
