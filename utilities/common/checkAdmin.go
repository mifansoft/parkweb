package common

const (
	id      = "0"
	account = "admin"
)

// 检查是否为超级管理员
func IsAdmin(str string) bool {
	return str == id || str == account
}
