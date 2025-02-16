package models

type User struct {
	Id           int64  `json:"id" gorm:"primaryKey"`
	Username     string `json:"username" gorm:"unique"`
	Password     string `json:"password"`
	LoginType    int    `json:"loginType"`
	FullName     string `json:"fullName"`
	Role         int    `json:"role"`
	Status       int    `json:"status"`
	CreatedAt    int64  `json:"createdAt"`
	UpdatedAt    int64  `json:"updatedAt"`
	LastLogin    int64  `json:"lastLogin"`
	CredsArrJson string `json:"credsArrJson"`
}
