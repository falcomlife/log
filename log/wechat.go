package log

import (
	"github.com/json-iterator/go"
	"k8s.io/klog/v2"
	"k8s.io/log-controller/common"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

const (
	WechatDomain    = "qyapi.weixin.qq.com"
	AccessTokenPath = "/cgi-bin/gettoken?corpid="
	SendMessagePath = "/cgi-bin/message/send?access_token="
)

type JSON struct {
	Access_token string `json:"access_token"`
}

type MESSAGES struct {
	Touser  string `json:"touser"`
	Toparty string `json:"toparty"`
	Totag   int    `json:"totag"`
	Msgtype string `json:"msgtype"`
	Agentid int    `json:"agentid"`
	Text    struct {
		//Subject string `json:"subject"`
		Content string `json:"content"`
	} `json:"text"`
	Safe int `json:"safe"`
}

func GetAccessToken(corpid, corpsecret string) string {
	client := common.HttpClient{
		Protocol: "https",
		Host:     WechatDomain,
	}
	var json_str JSON = JSON{}
	err := client.Get(AccessTokenPath+corpid+"&corpsecret="+corpsecret, &json_str)
	if err != nil {
		return err.Error()
	}
	return json_str.Access_token
}

func SendMessage(access_token, msg string) []byte {
	client := common.HttpClient{
		Protocol: "https",
		Host:     WechatDomain,
	}
	b, err := client.Post(SendMessagePath+access_token, msg)
	if err != nil {
		klog.Error(err)
	}
	return b
}

func Messages(touser string, toparty string, totag int, agentid int, content string) string {
	msg := MESSAGES{
		Totag:   totag,
		Msgtype: "text",
		Agentid: agentid,
		Safe:    0,
		Text: struct {
			//Subject string `json:"subject"`
			Content string `json:"content"`
		}{Content: content},
	}
	sed_msg, _ := json.Marshal(msg)
	return string(sed_msg)
}
