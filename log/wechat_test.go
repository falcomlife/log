package log

import (
	"fmt"
	"strings"
	"testing"
)

func Test_wechat(t *testing.T) {
	totag := 1
	agentid := AgentId   //企业号中的应用id。
	corpid := CorpId     //企业号的标识
	corpsecret := Secret ///企业号中的应用的Secret
	accessToken := Get_AccessToken(corpid, corpsecret)
	content := "1234"
	//  fmt.Println(content)
	// 序列化成json之后，\n会被转义，也就是变成了\\n，使用str替换，替换掉转义
	msg := strings.Replace(Messages("", "", totag, agentid, content), "\\\\", "\\", -1)
	fmt.Println(strings.Replace(msg, "\\\\", "\\", -1))
	Send_Message(accessToken, msg)
}
