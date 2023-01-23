package vo

import (
	gocloak "github.com/Nerzal/gocloak/v12"
)

type LoginVO struct {
	StudentId string `json:"studentId"`
	Cid       string `json:"cId"`
	Ccid      string `json:"ccId"`
	Phone     string `json:"phone"`
	JwtToken  *gocloak.JWT
	ErrorMsg  string `json:"errorMsg"`
}
