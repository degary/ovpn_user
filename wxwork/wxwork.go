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

func SendMsg(user, msg, token string) error {
	url := "https://qyapi.weixin.qq.com/cgi-bin/message/send"
	message := fmt.Sprintf("**您的OPENVPN账号已开通**\n"+
		">用户名:<font color=\"info\">%s</font>\n"+
		">密 码: <font color=\"warning\">%s</font>\n"+
		">客户端下载地址: <font color=\"comment\">xxxxx</font>\n", user, msg)
	content := map[string]interface{}{
		"content": message,
	}

	param := req.QueryParam{
		"access_token": token,
	}

	body := map[string]interface{}{
		"touser":                   user,
		"msgtype":                  "markdown",
		"agentid":                  1000023,
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
