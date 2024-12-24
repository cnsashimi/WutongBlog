package model

import (
	"errors"
	"gorm.io/gorm"
	"imzixuan/config"
)

type User struct {
	Id       int64  `gorm:"column:id;primary_key;AUTO_INCREMENT"`
	Username string `gorm:"column:username"`
	Level    int    `gorm:"column:level;default:0"`
	Password string `gorm:"column:password"`
	Nickname string `gorm:"column:nickname"`
	Avatar   string `gorm:"column:avatar"`
	Aboutme  string `gorm:"column:aboutme"`
}

func (m *User) TableName() string {
	return "blog.user"
}
func GetUser(username string) (User, error) {
	var user User
	result := config.Connect.First(&user, "username = ?", username)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return User{}, errors.New("user not found")
		}
		return User{}, result.Error
	}
	return user, nil
}
func GetUserbyid(uid string) (User, error) {
	var user User
	result := config.Connect.First(&user, "id = ?", uid)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return User{}, errors.New("user not found")
		}
		return User{}, result.Error
	}
	return user, nil
}
func GetUsername(uid int64) string {
	var user User
	var nickname string

	err := config.Connect.Model(&user).Select("nickname").Where("id = ?", uid).Scan(&nickname)
	if err.Error != nil {
		return "无"
	}
	return nickname
}
func UpdateUser(user User) bool {
	if user.Nickname == "" {
		return false
	}
	config.Connect.Select("Nickname", "Avatar", "Aboutme").Save(&user)

	return true
}
func UpdateResetPassword(user User) bool {
	config.Connect.Select("Password").Save(&user)
	return true
}
