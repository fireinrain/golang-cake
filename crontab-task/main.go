package main

import (
	"fmt"
	"github.com/robfig/cron"
)

func main() {
	c := cron.New()
	//这个库有点问题 没法支持秒
	// 标准的cron实现不包含秒
	c.AddFunc("0/1 * * * * ? *", func() {
		fmt.Println("1 * * * * * *")
	})
	//c.AddFunc("/5 * * * * *", crontabTask.DoCloudflareIPScanner)
	c.Start()
	select {}
}
