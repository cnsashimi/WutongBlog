package main

import (
	"fmt"
	"imzixuan/config"
	"imzixuan/controller"
	"imzixuan/task"
)

func main() {
	//sql初始
	config.Sqlinit()
	//redis初始
	config.InitRedisClient()
	//计划任务
	go task.CronExec()
	fmt.Println("Hello, WuTongBlog!")
	//控制器载入
	controller.InitController()

}
