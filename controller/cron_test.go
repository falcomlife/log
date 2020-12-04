package controller

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"k8s.io/log-controller/log"
	"strings"
	"testing"
)

func TestCron(t *testing.T) {
	crontab := cron.New()
	task := func() {
		totag := 1
		agentid := log.AgentId   //企业号中的应用id。
		corpid := log.CorpId     //企业号的标识
		corpsecret := log.Secret ///企业号中的应用的Secret
		accessToken := log.GetAccessToken(corpid, corpsecret)
		content := "1234"
		//  fmt.Println(content)
		// 序列化成json之后，\n会被转义，也就是变成了\\n，使用str替换，替换掉转义
		msg := strings.Replace(log.Messages("", "", totag, agentid, content), "\\\\", "\\", -1)
		fmt.Println(strings.Replace(msg, "\\\\", "\\", -1))
		log.SendMessage(accessToken, msg)
	}
	crontab.AddFunc("*/5 * * * * *", task)
	crontab.Start()
	defer crontab.Stop()
	select {}
}
