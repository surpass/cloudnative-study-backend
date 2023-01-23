// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"sportApi/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/object",
			beego.NSInclude(
				&controllers.ObjectController{},
			),
		),
		beego.NSNamespace("/user",
			beego.NSInclude(
				&controllers.UserController{},
			),
		),
		beego.NSNamespace("/location",
			beego.NSInclude(
				&controllers.LocationController{},
			),
			beego.NSRouter("/trip", &controllers.TripController{}, "Post:CreateTrip"),
			beego.NSRouter("/closetrip", &controllers.TripController{}, "Post:CloseTrip"),
			beego.NSRouter("/getTrip", &controllers.GetTrip{}, "*:GetTrip"),
			beego.NSRouter("/rank", &controllers.Rank{}, "*:RankList"),
		),
		beego.NSNamespace("/data",
			beego.NSRouter("/slides", &controllers.DemoController{}, "Get:Slides"),
			beego.NSRouter("/categories", &controllers.DemoController{}, "Get:Categories"),
			beego.NSRouter("/index", &controllers.DemoController{}, "*:Index"),
		),
		beego.NSNamespace("/person",
			beego.NSRouter("/",
				&controllers.SelectAll{},
			),
			beego.NSRouter("/id/?:id", &controllers.SelectbyId{}, "Get:GetById"),
			beego.NSRouter("/did?:id",
				&controllers.DeletebyId{},
			),
			beego.NSRouter("/uid", &controllers.Updatebyid{}, "Put:UpdatePerson"),
			beego.NSRouter("/add", &controllers.InsertData{}, "POST:Post"),

			beego.NSRouter("/getByName:name",
				&controllers.SelectUnClear{},
			),
		),
		beego.NSNamespace("/wx",
			beego.NSRouter("/login", &controllers.LoginController{}, "*:Loginwx"),
			beego.NSRouter("/getPhonenumber", &controllers.LoginController{}, "*:GetPhonenumber"),
			beego.NSRouter("/autoLoginwxByOpenid", &controllers.LoginController{}, "*:AutoLoginwxByOpenid"),
		),
		beego.NSNamespace("/oauth",
			beego.NSRouter("/getToken", &controllers.OauthController{}, "*:GetToken"),
			beego.NSRouter("/refreshToken", &controllers.OauthController{}, "*:RefreshToken"),
			beego.NSRouter("/verifyAccessToken", &controllers.OauthController{}, "*:VerifyAccessToken"),
		),
	)
	beego.AddNamespace(ns)
}
