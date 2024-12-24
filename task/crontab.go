package task

import (
	"encoding/json"
	"fmt"
	"github.com/robfig/cron"
	"imzixuan/config"
	"imzixuan/model"
	"time"
)

func CronExec() {
	c := cron.New()
	spec := "0 */3 * * * *" //每3分钟执行一次  最左位 为秒
	//TODO::注意 执行周期时间不要大于  mysql的wait_timeout  时间

	c.AddFunc(spec, func() {
		Cachecategories()
		CacheTags()
	})

	spec2 := "0 */10 * * * *" //每10分钟执行一次  最左位 为秒
	c.AddFunc(spec2, func() {

		CacheTopblog()
	})

	c.Start()

}
func CacheTopblog() string {
	top3blog := model.GetBlogTop3()

	jsonStu, err := json.Marshal(top3blog)
	if err != nil {
		fmt.Println("生成json字符串错误")
	}
	config.RedisSet("top3blog", string(jsonStu), time.Minute*150)
	return string(jsonStu)

}
func Cachecategories() string {
	clists := model.GetAllCates()
	jsonStu, err := json.Marshal(clists)
	if err != nil {
		fmt.Println("生成json字符串错误")
	}
	config.RedisSet("categories", string(jsonStu), time.Minute*150)
	return string(jsonStu)
}
func CacheTags() string {
	tlists := model.GetAllTags()
	jsonStu, err := json.Marshal(tlists)
	if err != nil {
		fmt.Println("生成json字符串错误")
	}
	config.RedisSet("tags", string(jsonStu), time.Minute*150)
	return string(jsonStu)
}
