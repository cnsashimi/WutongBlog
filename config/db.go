package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
)

type Config struct {
	Mysql Mysql `json:"mysql"`
	Redis Redis `json:"redis"`
	Base  Base  `json:"base"`
}

type Mysql struct {
	Url      string
	Port     string
	User     string
	Password string
	Dbname   string
}

type Redis struct {
	Host     string
	Port     string
	Passowrd string
	Db       int
	Poolsize int
}

type Base struct {
	Name        string
	Beian       string
	Logourl     string
	Description string
	Keywords    string
	Baidutongji string
}

var Connect *gorm.DB

func GetYml() Config {
	dataBytes, err := os.ReadFile("config.yml")
	if err != nil {
		fmt.Println("读取文件失败：", err)
	}
	configfile := Config{}
	err = yaml.Unmarshal(dataBytes, &configfile)
	if err != nil {
		fmt.Println("解析 yaml 文件失败：", err)
	}
	mp := make(map[string]any, 2)
	err = yaml.Unmarshal(dataBytes, mp)
	if err != nil {
		fmt.Println("解析 yaml 文件失败：", err)
	}
	return configfile
}
func Sqlinit() {
	configfile := GetYml()
	var err error
	dns := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		configfile.Mysql.User,
		configfile.Mysql.Password,
		configfile.Mysql.Url,
		configfile.Mysql.Port,
		configfile.Mysql.Dbname,
	)

	Connect, err = gorm.Open(mysql.Open(dns), &gorm.Config{})
	if err != nil {
		fmt.Println("mysql Connect faild, err:", err)
		return
	}

	return
}
