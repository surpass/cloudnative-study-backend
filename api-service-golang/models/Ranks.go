package models

import (
	"fmt"

	"github.com/astaxie/beego/orm"

	_ "github.com/go-sql-driver/mysql" //导入数据库驱动
)

func init() {

}

type Ranks struct {
	SId        string  `json:"id"`
	CId        string  `json:"uid"`
	CcId       string  `json:"dkey"`
	Avatar_url string  `json:"avatar_url"`
	Distance   float64 `json:"distance"`
	RunTime    int64   `json:"runTime"`
	Sname      string  `json:"sname"`
	Cname      string  `json:"cname"`
	Ccname     string  `json:"ccname"`
	RankNo     string  `json:"rankNo"`
}

func GetRanks(types string, timeCycle string, sId string, cId string, ccId string) []*Ranks {
	o := orm.NewOrm()
	ranks := make([]*Ranks, 1, 10)

	ranksql := fmt.Sprintf(`select t.id,t.year,t.month,t.day,t.week,t.sid,t.cid,t.ccid,t.distance,t.run_time,t.rank,t.types,s.name sname,c.name cname from u_ranks t left  join u_students s on s.id = t.sid left join u_classes c on c.id = t.cid`)
	_, err := o.Raw(ranksql).QueryRows(&ranks)
	if err != nil && err == orm.ErrNoRows {
		fmt.Println("ranks is not rows :%s", cId)
	}
	return ranks
}
