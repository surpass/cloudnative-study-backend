package vo

type LocationVO struct {
	Tid        string  `json:"tid"`
	Stime      int64   `json:"stime"` //采集时间
	Lon        string  `json:"lon"`
	Lat        string  `json:"lat"`
	Speed      float32 `json:"speed"`
	Createtime int64   `json:"createTime"`
}
