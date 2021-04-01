package wxwork

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/imroc/req"
	"math/rand"
	"time"
)

func GetToken(rds *redis.Client, corpid, corpsecret string) (string, error) {
	ctx := context.Background()
	n, err := rds.Exists(ctx, "wxwork_token").Result()
	if err != nil {
		return "", err
	}
	if n > 0 {
		result, _ := rds.Get(ctx, "wxwork_token").Result()
		return result, nil
	}

	url := "https://qyapi.weixin.qq.com/cgi-bin/gettoken"
	param := req.Param{
		"corpid":     corpid,
		"corpsecret": corpsecret,
	}
	res, err := req.Get(url, param)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	foo := map[string]string{}
	res.ToJSON(&foo)
	rds.Set(ctx, "wxwork_token", foo["access_token"], time.Second*290)
	return foo["access_token"], nil
}

func SendMsg(user, msg, token string, agentid int) error {
	url := "https://qyapi.weixin.qq.com/cgi-bin/message/send"
	message := fmt.Sprintf("### 您的VPN账号已开通\n"+
		">用户名：<font color=\"info\">%s</font>\n"+
		">密  码：<font color=\"info\">%s</font>\n"+
		">[下载：VPN使用手册（必看）](http://www.baidu.com)\n"+
		">[下载：谷歌身份验证器-安卓（必备）](http://tools.peogoo.com/download/google.authenticator2_5.10.apk)\n"+
		">[下载：VPN系统App-安卓系统](http://tools.peogoo.com/download/openvpn-1597289328.apk)\n", user, msg)
	content := map[string]interface{}{
		"content": message,
	}

	param := req.QueryParam{
		"access_token": token,
	}

	body := map[string]interface{}{
		"touser":                   user,
		"msgtype":                  "markdown",
		"agentid":                  agentid,
		"markdown":                 content,
		"enable_duplicate_check":   0,
		"duplicate_check_interval": 1800,
	}

	_, err := req.Post(url, param, req.BodyJSON(&body))
	return err
}

func GetPasswd(length int) string {
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_@+"
	count := len(letters)
	chars := make([]byte, length)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < length; i++ {
		chars[i] = letters[rand.Intn(count)]
	}
	return string(chars)
}
