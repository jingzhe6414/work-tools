package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
	"strconv"
	"webhook/model"
)

var (
	db, _ = strconv.Atoi(os.Getenv("REDIS_DB"))
	r     = model.RedisConfig{
		IP:  os.Getenv("REDIS_IP"),
		PWD: os.Getenv("REDIS_PWD"),
		DB:  db,
	}
	//db = 0
	//r  = model.RedisConfig{
	//	//IP: "192.168.2.103:31379",
	//	IP: "192.168.2.116:6379",
	//	DB: db,
	//}
	c       = r.InitRedis()
	msgList = make(map[string]string)
)

func Alerts(c *gin.Context) {
	var json model.Alertmanager
	if err := c.ShouldBind(&json); err != nil {
		panic(err.Error())
	}
	for _, alert := range json.Alerts {
		msg := fmt.Sprintf("{alertname: %s,namespace: %s,alert %s,msg: %s} \n", alert.Labels.Alertname, alert.Labels.Namespace, alert.Labels.Severity, alert.Annotations.Message)
		if alert.Labels.Namespace == "" {
			return
		}
		if json.Status == "firing" {
			fmt.Println("开始添加")
			firing(msg, alert.Labels.Alertname, alert.Labels.Namespace)
		}
		if json.Status == "resolved" {
			fmt.Println("开始删除")
			resolved(alert.Labels.Namespace, msg)
		}
	}

}
func firing(msg, name, ns string) {
	var ls []string

	val, _ := c.LRange(ns, 0, -1).Result()
	fmt.Sprintf("redis中key %s的值为： %s", ns, val)
	if len(val) != 0 {
		for _, v := range val {
			if v == msg {
				return
			}
		}
	}

	ls = append(ls, msg)
	fmt.Println("redis存数据：", ls)
	err := c.LPush(ns, ls).Err()

	if err != nil {
		panic(err)
	}

}

// 恢复 json.Status=resolved
func resolved(ns, msg string) {
	val, _ := c.LRange(ns, 0, -1).Result()
	fmt.Sprintf("redis中key %s的值为： %s", ns, val)
	if len(val) != 0 {
		for _, v := range val {
			if v == msg {
				err := c.LRem(ns, 0, msg).Err()
				if err != nil {
					panic(err)
				}
				fmt.Println("成功删除：", msg)
			}
		}
	}
}
