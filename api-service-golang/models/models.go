package models

import (
	"fmt"
	"time"

	beego "github.com/beego/beego/v2/server/web"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql" //导入数据库驱动
)

var (
	mysqlhost     string = ""
	mysqlport     int    = 3306
	mysqluser     string = ""
	mysqlpassword string = ""
	mysqldbname   string = ""

	max_conn  int = 40
	idle_conn int = 10
)

func init() {
	mysqlhost, _ := beego.AppConfig.String("mysql_host")
	mysqlport, _ = beego.AppConfig.Int("mysql_port")
	mysqluser, _ = beego.AppConfig.String("mysql_user")
	mysqlpassword, _ = beego.AppConfig.String("mysql_password")
	mysqldbname, _ = beego.AppConfig.String("mysql_dbname")

	connectString := fmt.Sprintf(`%s:%s@tcp(%s:%d)/%s?charset=utf8`,
		mysqluser, mysqlpassword, mysqlhost, mysqlport, mysqldbname)

	fmt.Println("db con info:", connectString)

	orm.RegisterDriver("mysql", orm.DRMySQL)
	//注册数据库，一共有五个参数，后面连个用于连接池操作
	_ = orm.RegisterDataBase("default", "mysql", connectString)
	//_ = orm.RegisterDataBase("default", "mysql",
	//	"root:123456@tcp(sql.easyolap.cn:3306)/sport?charset=utf8")
	//进行注册模型结构，可以有多个，用逗号分隔
	orm.RegisterModel(new(Person), new(WxLogin), new(WxUser), new(University), new(Colleges), new(Classes), new(Students), new(SqlTrip), new(RankEngity))
	//创建表，默认为default，只建立一次，后面再执行这个会忽略
	_ = orm.RunSyncdb("default", false, true)

	orm.SetMaxIdleConns("default", 10)
	orm.SetMaxOpenConns("default", 100)

	orm.Debug = true

	// 设置为 UTC 时间
	//orm.DefaultTimeLoc = time.UTC
	// 设置为 东8区 时间
	timelocal := time.FixedZone("CST", 3600*8)
	time.Local = timelocal
	orm.DefaultTimeLoc = timelocal

	//使用表名前缀
	//orm.RegisterModelWithPrefix("prefix_", new(User))

}

// func getMysqlDB() (*sql.DB, error) {
// 	connectString := fmt.Sprintf(`%s:%s@tcp(%s:%s)/%s?charset=utf8`,
// 		user, password, host, port, dbname)

// 	db, err := sql.Open("mysql", connectString)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return db, nil
// }
