package controllers

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sportApi/models"
	"sportApi/param"
	"sportApi/vo"
	"time"

	beego "github.com/beego/beego/v2/server/web"
)

// Operations about Users
type LocationController struct {
	beego.Controller
}

// @Title Create Cur Location
// @Description create users
// @Param	body		body 	models.Location	true		"body for user content"
// @Success 200 {int} models.User.Id
// @Failure 403 body is empty
// @router / [post]
func (l *LocationController) Post() {
	var runLocation models.Location
	//fmt.Println("request body:%v", l.Ctx.Input.RequestBody)
	json.Unmarshal(l.Ctx.Input.RequestBody, &runLocation)
	fmt.Println("request body:%v", &runLocation)

	location := models.UploadLocation(runLocation)
	var locationVO vo.LocationVO
	locationVO.Tid = location.Tid
	locationVO.Stime = location.Stime
	l.Data["json"] = locationVO
	l.ServeJSON()
}

// @Title GetAll
// @Description get all Locations
// @Success 200 {object} models.Location
// @router / [get]
func (l *LocationController) GetAll() {
	locations := models.GetAllLocations()
	l.Data["json"] = locations
	l.ServeJSON()
}

// @Title GetLocation
// @Description get all Locations
// @Success 200 {object} models.Location
// @router /getLocationByTid [get]
func (l *LocationController) GetLocationByTid(tid string) {
	tId := l.GetString("tid")
	fmt.Println("get trip detail by tId:%s", tId)

	trip := models.GetTrip(tId)

	locations := models.GetLocationsByTid(trip.Id, trip.FromDate, trip.ThruDate)
	var res GetLocationRes
	res.R = 200
	res.Trip = &trip
	res.Locations = locations

	l.Data["json"] = res
	l.ServeJSON()
}

type GetLocationRes struct {
	R         int                `json:"r"`
	Trip      *models.SqlTrip    `json:"trip"`
	Locations []*models.Location `json:"locations"`
}

type TripController struct {
	beego.Controller
}

// @Title Create Cur Trip
// @Description create Trip
// @Param	body		body 	models.Trip	true		"body for Trip content"
// @Success 200 {int} models.Trip
// @Failure 403 body is empty
// @router /trip [post]
func (l *TripController) CreateTrip() {

	var tripParam param.TripParam
	json.Unmarshal(l.Ctx.Input.RequestBody, &tripParam)
	fmt.Println("tripParam request body:%v", &tripParam)
	var trip models.SqlTrip
	CopyStruct(&trip, &tripParam)
	fmt.Println("trip request body:%v", &trip)
	trip.FromDate = time.Unix(tripParam.FromDateStamp, 0)
	fmt.Println("trip request body:%v", &trip)

	tripDb := models.CreateTrip(trip)
	l.Data["json"] = tripDb
	l.ServeJSON()
}

// @Title Close Cur Trip
// @Description Close Trip
// @Param	body		body 	models.Trip	true		"body for Trip content"
// @Success 200 {int} models.Trip
// @Failure 403 body is empty
// @router /trip [put]
func (l *TripController) CloseTrip() {
	var tripParam param.TripParam
	json.Unmarshal(l.Ctx.Input.RequestBody, &tripParam)

	json.Unmarshal(l.Ctx.Input.RequestBody, &tripParam)
	fmt.Println("upate tripParam request body:%v", tripParam)

	tripObj := models.GetTrip(tripParam.Id)
	tripObj.ThruDate = time.Unix(tripParam.ThruDateStamp, 0)
	tripObj.Tlon = tripParam.Tlon
	tripObj.Tlat = tripParam.Tlat
	tripObj.Runtime = tripParam.Runtime
	tripObj.MaxSpeed = tripParam.MaxSpeed
	tripObj.AvgeSpeed = tripParam.AvgeSpeed
	tripObj.Distance = tripParam.Distance

	tripId := models.EndTrip(tripObj)

	fmt.Println("upated trip request body:%v", tripId)
	//tripObj := models.GetTrip(tripId.Id)
	l.Data["json"] = tripObj
	l.ServeJSON()
}

type GetRes struct {
	R     int               `json:"r"`
	Total int64             `json:"total"`
	Datas []*models.SqlTrip `json:"datas"`
}
type GetTrip struct {
	beego.Controller
}

func (l *GetTrip) GetTrip() {
	id := l.GetString("sid")
	fmt.Println("get my trip by id:%s", id)
	sid := l.Ctx.Input.Param("sid")
	fmt.Println("get my trip by sid:%s", sid)

	page, _ := l.GetInt16("_page")
	limit, _ := l.GetInt16("_limit")
	fmt.Println("get my trip page:%d limit:%d", page, limit)

	trips, total := models.Select_Trip(id, page, limit)

	//total := models.Total_Trip(id)

	var res GetRes
	res.R = 200
	res.Total = total
	res.Datas = trips

	l.Data["json"] = res
	l.ServeJSON()
}

type Rank struct {
	beego.Controller
}

type RankRes struct {
	R     int             `json:"r"`
	Datas []*models.Ranks `json:"ranks"`
}

func (l *Rank) RankList() {
	sid := l.GetString("sid")
	cId := l.GetString("cId")
	ccId := l.GetString("ccId")

	types := l.GetString("type")
	timeCycle := l.GetString("timeCycle")

	fmt.Println("get RankList by sid:%s,%s,%s", sid, cId, ccId, types, timeCycle)

	rankLists := models.GetRanks(types, timeCycle, sid, cId, ccId)

	var res RankRes
	res.R = 200
	res.Datas = rankLists

	l.Data["json"] = res
	l.ServeJSON()
}

// CopyStruct
// dst 目标结构体，src 源结构体
// 必须传入指针，且不能为nil
// 它会把src与dst的相同字段名的值，复制到dst中
func CopyStruct(dst, src interface{}) {
	dstValue := reflect.ValueOf(dst).Elem()
	srcValue := reflect.ValueOf(src).Elem()

	for i := 0; i < srcValue.NumField(); i++ {
		srcField := srcValue.Field(i)
		srcName := srcValue.Type().Field(i).Name
		dstFieldByName := dstValue.FieldByName(srcName)

		if dstFieldByName.IsValid() {
			switch dstFieldByName.Kind() {
			case reflect.Ptr:
				switch srcField.Kind() {
				case reflect.Ptr:
					if srcField.IsNil() {
						dstFieldByName.Set(reflect.New(dstFieldByName.Type().Elem()))
					} else {
						dstFieldByName.Set(srcField)
					}
				default:
					dstFieldByName.Set(srcField.Addr())
				}
			default:
				switch srcField.Kind() {
				case reflect.Ptr:
					if srcField.IsNil() {
						dstFieldByName.Set(reflect.Zero(dstFieldByName.Type()))
					} else {
						dstFieldByName.Set(srcField.Elem())
					}
				default:
					dstFieldByName.Set(srcField)
				}
			}
		}
	}
}
