package vo

type LocationVO struct {
	Tid        string  `json:"tid"`
	Stime      int64   `json:"stime"` //ιιζΆι΄
	Lon        string  `json:"lon"`
	Lat        string  `json:"lat"`
	Speed      float32 `json:"speed"`
	Createtime int64   `json:"createTime"`
}
