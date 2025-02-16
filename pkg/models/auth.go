package models

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

type AuthClaims struct {
	jwt.StandardClaims
	Id        int64  `json:"id"`
	Username  string `json:"username"`
	FullName  string `json:"fullName"`
	LoginType int    `json:"loginType"`
	Expire    int64  `json:"expire"`
	Role      int    `json:"role"`
	CreatedAt int64  `json:"createdAt"`
	LastLogin int64  `json:"lastLogin"`
}

type UserInfo struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
	Role     int    `json:"role"`
}

func (c AuthClaims) Valid() error {
	timestamp := time.Now().Unix()
	if timestamp >= c.Expire {
		return fmt.Errorf("%s", "The credential is expired")
	}
	return nil
}
