package utils

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

type UserRole int
type UserStatus int
type LoginType int

const (
	UserListSessionKey = "userList"
	AliveSessionHours  = 24
)

const (
	RoleSuperAdmin UserRole = iota
	RoleRegular
)

const (
	StatusDeactive UserStatus = iota
	StatusActive
)

const (
	LoginWithPasskey LoginType = iota
	LoginWithPassword
)

type ResponseData struct {
	IsError   bool        `json:"error"`
	ErrorCode string      `json:"errorCode"`
	Msg       string      `json:"msg"`
	Data      interface{} `json:"data"`
}
