package log

import (
	"bytes"
	"github.com/json-iterator/go"
	"io/ioutil"
	"net/http"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

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

func Get_AccessToken(corpid, corpsecret string) string {
	gettoken_url := "https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=" + corpid + "&corpsecret=" + corpsecret
	//print(gettoken_url)
	client := &http.Client{}
	req, _ := client.Get(gettoken_url)
	defer req.Body.Close()
	body, _ := ioutil.ReadAll(req.Body)
	//fmt.Printf("\n%q",string(body))
	var json_str JSON
	json.Unmarshal([]byte(body), &json_str)
	//fmt.Printf("\n%q",json_str.Access_token)
	return json_str.Access_token
}

func Send_Message(access_token, msg string) (string, []byte) {
	send_url := "https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=" + access_token
	//print(send_url)
	client := &http.Client{}
	req, _ := http.NewRequest("POST", send_url, bytes.NewBuffer([]byte(msg)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("charset", "UTF-8")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return resp.Status, body
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
