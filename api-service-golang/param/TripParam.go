package param

type TripParam struct {
	Id            string  `json:"id"`
	Ccid          string  `json:"ccid"`
	Cid           string  `json:"cid"`
	Sid           string  `json:"sid"`
	Dkey          string  `json:"dkey"`
	FromDateStamp int64   `json:"fromDate"`
	ThruDateStamp int64   `json:"thruDate"`
	Flon          string  `json:"flon"`
	Flat          string  `json:"flat"`
	Tlon          string  `json:"tlon"`
	Tlat          string  `json:"tlat"`
	Distance      float64 `json:"distance"`
	Runtime       int64   `json:"runTime"`
	AvgeSpeed     float64 `json:"avgeSpeed"`
	MaxSpeed      float64 `json:"maxSpeed"`
}
