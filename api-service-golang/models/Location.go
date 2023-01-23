package models

import (
	"fmt"
	"log"
	"math/rand"
	"sportApi/utils"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/gocql/gocql"
	uuid "github.com/satori/go.uuid"

	_ "github.com/go-sql-driver/mysql" //导入数据库驱动
)

var (
	LocationList map[string]*Location
)

var session *gocql.Session
var cluster *gocql.ClusterConfig

func init() {
	LocationList = make(map[string]*Location)

	// cosmosCassandraContactPoint := os.Getenv("COSMOSDB_CASSANDRA_CONTACT_POINT")
	// cosmosCassandraPort := os.Getenv("COSMOSDB_CASSANDRA_PORT")
	// cosmosCassandraUser := os.Getenv("COSMOSDB_CASSANDRA_USER")
	// cosmosCassandraPassword := os.Getenv("COSMOSDB_CASSANDRA_PASSWORD")
	// if cosmosCassandraPort == "" {
	// 	cosmosCassandraPort = "9042"
	// }
	// if cosmosCassandraContactPoint == "" || cosmosCassandraUser == "" || cosmosCassandraPassword == "" {
	// 	log.Fatal("missing mandatory environment variables")
	// }
	cluster = gocql.NewCluster("nosql.easyolap.cn")
	/*cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: "cassandra",
		Password: "111111",
	}*/

	cluster.Keyspace = "sport"

	cluster.Timeout = 5 * time.Second

	cluster.ProtoVersion = 4

	lsession, err := cluster.CreateSession()
	if err != nil {
		log.Fatalf("Could not connect to cassandra cluster: %v", err)
	}
	session = lsession

	KeySpaceMeta, _ := session.KeyspaceMetadata("sport")
	if _, exists := KeySpaceMeta.Tables["u_rundata"]; exists != true {
		session.Query("CREATE TABLE u_rundata ( " +
			"	tid text," +
			"	stime bigint," + //采集时间
			"	lon text," +
			"	lat text," +
			"	speed float," +
			"	ctime timestamp," +
			"	PRIMARY KEY(tid,stime)" +
			")").Exec()
		// 创建索引
		//session.Query("CREATE INDEX ON u_rundata(stime)").Exec()
	}

	// if _, exists := KeySpaceMeta.Tables["s_userinfo"]; exists != true {
	// 	session.Query("CREATE TABLE s_userinfo ( " +
	// 		"	id text," +
	// 		"	code text," +
	// 		"	nick_name text," +
	// 		"	language text," +
	// 		"	avatarUrl text," +
	// 		"	country text," +
	// 		"	province text," +
	// 		"	city text," +
	// 		"	gender text," +
	// 		"	phone text," +
	// 		"	ctime timestamp," +
	// 		"	utime timestamp," +
	// 		"	PRIMARY KEY(id)" +
	// 		")").Exec()
	// }

}

type Trip struct {
	Id       string `json:"id"`
	Uid      string `json:"uid"`
	Dkey     string `json:"dkey"`
	Stime    int64  `json:"stime"`
	Slon     string `json:"slon"`
	Slat     string `json:"slat"`
	Etime    int64  `json:"etime"`
	Elon     string `json:"elon"`
	Elat     string `json:"elat"`
	Distance int64  `json:"distance"`
	Ctime    int64  `json:"ctime"`
}

type Location struct {
	Tid        string  `json:"tid"`
	Stime      int64   `json:"stime"` //采集时间
	Lon        string  `json:"lon"`
	Lat        string  `json:"lat"`
	Speed      float32 `json:"speed"`
	Createtime int64   `json:"createTime"`
}

func UploadLocation(l Location) Location {
	currentTime := time.Now()
	session.Query("INSERT INTO u_rundata (tid,stime,lon,lat,speed,ctime) values(?,?,?,?,?,?)",
		l.Tid, l.Stime, l.Lon, l.Lat, l.Speed, currentTime).Exec()

	return l
}

func GetAllLocations() map[string]*Location {

	var tid string
	var stime int64
	var lon string
	var lat string
	var speed float32
	var ctime int64
	iter := session.Query("select tid,stime,lon,lat,speed,ctime from u_rundata").Iter()
	defer func() {
		if iter != nil {
			iter.Close()
		}
	}()
	for iter.Scan(&tid, &stime, &lon, &lat, &speed, &ctime) {
		log.Printf("lon:%v,lat:%v,speed:%v", lon, lat, speed)
		l := Location{tid, stime, lon, lat, speed, ctime}
		var sid = tid + "_" + strconv.FormatInt(int64(stime), 10)
		LocationList[sid] = &l
	}
	return LocationList
}

func GetLocationsByTid(tid string, startTime time.Time, endTime time.Time) []*Location {
	locations := make([]*Location, 0)
	var dbtid string
	var stime int64
	var lon string
	var lat string
	var speed float32
	var ctime int64
	//sql := fmt.Sprintf(`select * from sport.u_rundata where tid = '%s' and  stime >= %d and stime <= %d`, tid, startTime.Unix(), endTime.Unix())
	sql := fmt.Sprintf(`select tid,stime,lon,lat,speed,ctime from sport.u_rundata where tid = '%s' and  stime >= %d;`, tid, startTime.Unix())

	fmt.Println("GetLocationsByTid tid:%s,%d,%d", tid, startTime, endTime)
	fmt.Println("sql:%s", sql)

	iter := session.Query(sql).Iter()
	defer func() {
		if iter != nil {
			iter.Close()
		}
	}()
	for iter.Scan(&dbtid, &stime, &lon, &lat, &speed, &ctime) {
		fmt.Println("lon:%v,lat:%v,speed:%v", strings.Replace(lon, "'", "", -1), strings.Replace(lat, "'", "", -1), speed)
		l := Location{dbtid, stime, strings.Replace(lon, "'", "", -1), strings.Replace(lat, "'", "", -1), speed, ctime}
		locations = append(locations, &l)
	}
	return locations
}

func CreateTrip(t SqlTrip) SqlTrip {

	o := orm.NewOrm()
	id := uuid.NewV4()
	t.Id = id.String()
	currentTime := time.Now()
	rand.Seed(time.Now().UnixNano())
	num := rand.Intn(30001)
	t.Dkey = strconv.FormatInt(int64(num), 10) + "a"
	t.CreateDate = currentTime
	qs := o.QueryTable(TRIP_TABLE_NAME)
	i, _ := qs.PrepareInsert()

	loginId, lerr := i.Insert(&t)
	if lerr != nil {
		fmt.Println("wx_login error:%v", lerr)
	} else {
		fmt.Println("wx_login success:%v", loginId)
	}

	return t
}

func EndTrip(t SqlTrip) SqlTrip {
	currentTime := time.Now()
	o := orm.NewOrm()
	t.UpdateDate = currentTime
	num, err := o.Update(&t)
	if err != nil {
		fmt.Println("SqlTrip error:%v", err)
	} else {
		fmt.Println("SqlTrip success:%v,%d", t, num)
	}
	return t
}

func GetTrip(tId string) SqlTrip {
	o := orm.NewOrm()
	var trip SqlTrip
	sql := fmt.Sprintf(`select * from %s where id=? `, TRIP_TABLE_NAME)
	err := o.Raw(sql, tId).QueryRow(&trip)
	if err != nil && err == orm.ErrNoRows {
		fmt.Println("sql not rows :%s", tId)
	}
	return trip
}

func Select_Trip(sId string, pageNum int16, pageSize int16) ([]*SqlTrip, int64) {
	o := orm.NewOrm()
	var trips []*SqlTrip

	var count int64
	//countsql := fmt.Sprintf(`select count(*) from %s where s_id=? `, TRIP_TABLE_NAME)
	//_, cerr := o.Raw(countsql, sId).QueryRows(&count)
	qs := o.QueryTable(TRIP_TABLE_NAME)
	qs.Filter("s_id", sId)

	count, cerr := qs.Count()
	if cerr != nil && cerr == orm.ErrNoRows {
		fmt.Println("Total_Trip sql not rows :%s", sId)
	}

	//	总页数
	pageNumTotal := utils.Page(pageSize, pageNum, count)
	fmt.Println("total page is :%d", pageNumTotal)
	startRow := pageSize*(pageNum-1) + 0

	//获取分页数据
	_, err := qs.Limit(pageSize, startRow).All(&trips)
	if err != nil && err == orm.ErrNoRows {
		fmt.Println("sql not rows :%s", sId)
	}

	return trips, count
}

func Total_Trip(sId string) int64 {
	o := orm.NewOrm()
	var count int64
	sql := fmt.Sprintf(`select count(*) from %s where s_id=? `, TRIP_TABLE_NAME)
	_, err := o.Raw(sql, sId).QueryRows(&count)
	if err != nil && err == orm.ErrNoRows {
		fmt.Println("Total_Trip sql not rows :%s", sId)
	}

	return count
}
