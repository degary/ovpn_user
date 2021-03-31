package ovpn

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"os/exec"
)

type Ovpn struct {
	ConnGroup string `json:"conn_group"`
	Type      string `json:"type"`
}

//使用sacli命令,获取系统内的用户和组,并存储到reids中
func GetUserGroup(rds *redis.Client) error {
	var ctx = context.Background()
	cmd := exec.Command("/usr/local/openvpn_as/scripts/sacli", "UserPropGet")
	buf, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
		return err
	}
	users := map[string]Ovpn{}
	err = json.Unmarshal(buf, &users)
	if err != nil {
		return err
	}
	for user, v := range users {
		//如果是用户
		if v.Type == "user_compile" {
			if n, _ := rds.Exists(ctx, fmt.Sprintf("User_%s", user)).Result(); n > 0 {
				//如果用户存在
				continue
			} else {
				//如果用户不存在,则保存用户到redis
				rds.Set(ctx, fmt.Sprintf("User_%s", user), user, 0)
			}
		} else if v.Type == "group" {
			if n, _ := rds.Exists(ctx, fmt.Sprintf("Group_%s", user)).Result(); n > 0 {
				//如果组存在
				continue
			} else {
				//如果组不存在,则保存组到redis
				rds.Set(ctx, fmt.Sprintf("Group_%s", user), user, 0)
			}
		}

	}
	return nil
}

//创建用户,并存储到redis中
func CreateUser(rds *redis.Client, userName string) error {
	var ctx = context.Background()
	n, err := rds.Exists(ctx, fmt.Sprintf("User_%s", userName)).Result()
	if err != nil {
		return err
	}
	if n > 0 {
		return fmt.Errorf("用户 %s 已存在", userName)
	}
	cmd := exec.Command("/usr/local/openvpn_as/scripts/sacli", "--user", userName, "--key", "type", "--value", "user_connect", "UserPropPut")
	_, err = cmd.Output()
	if err != nil {
		return err
	}
	rds.Set(ctx, fmt.Sprintf("User_%s", userName), userName, 0)
	return nil
}

//给用户设置密码
func SetPasswd(userName, passwd string) error {
	cmd := exec.Command("/usr/local/openvpn_as/scripts/sacli", "--user", userName, "--new_pass", passwd, "SetLocalPassword")
	_, err := cmd.Output()
	if err != nil {
		return err
	}
	return nil
}

//把用户添加到组
func AddUserToGroup(userName, groupName string) error {
	cmd := exec.Command("/usr/local/openvpn_as/scripts/sacli", "--user", userName, "--key", "conn_group", "--value", groupName, "UserPropPut")
	_, err := cmd.Output()
	if err != nil {
		return err
	}
	return nil
}
