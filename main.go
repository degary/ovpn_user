package main

import (
	"context"
	"fmt"
	"github.com/degary/ovpn_user/data"
	"github.com/degary/ovpn_user/ovpn"
	"github.com/degary/ovpn_user/wxwork"
	"github.com/go-redis/redis/v8"
	"gopkg.in/ini.v1"
	"time"
)

func main() {
	//初始化配置文件
	cfg, err := initConfig("./config.ini")
	if err != nil {
		fmt.Println(err)
		return
	}

	//初始化redis
	addr := cfg.Section("Redis").Key("Host").String()
	passwd := cfg.Section("Redis").Key("Passwd").String()
	db, err := cfg.Section("Redis").Key("Db").Int()
	if err != nil {
		fmt.Println(err)
		return
	}
	rds, err := initRedis(addr, passwd, db)
	if err != nil {
		fmt.Println(err)
		return
	}
	//获取用户及组信息并同步到redis
	fmt.Println("==========查看user,Group==========")
	err = ovpn.GetUserGroup(rds)
	if err != nil {
		fmt.Println(err)
		return
	}

	//获取token
	cid := cfg.Section("Wxwork").Key("CorpId").String()
	ckey := cfg.Section("Wxwork").Key("CorpSecret").String()
	agentid, _ := cfg.Section("Wxwork").Key("AgentId").Int()
	token, err := wxwork.GetToken(rds, cid, ckey)
	if err != nil {
		fmt.Println(err)
		return
	}

	//从文件中获取需要创建的用户
	objs, err := data.GetUserObjs("users.txt")
	if err != nil {
		fmt.Println("GetUserObjs", err)

	}
	for _, user := range objs {
		fmt.Printf("创建用户:%s\n", user.UserName)
		err = ovpn.CreateUser(rds, user.UserName)
		if err != nil {
			fmt.Printf("创建用户%s失败,err:%s\n", user.UserName, err.Error())
			continue
		}
		//设置用户
		passwordforuser := wxwork.GetPasswd(10)
		err = ovpn.SetPasswd(user.UserName, passwordforuser)
		if err != nil {
			fmt.Printf("设置用户%s密码失败,err:%s\n", user.UserName, err.Error())
			continue
		}
		//加入组
		err = ovpn.AddUserToGroup(user.UserName, user.UserGroup)
		if err != nil {
			fmt.Printf("用户%s加入到组%s失败,err:%s\n", user.UserName, user.UserGroup, err.Error())
			continue
		}
		//发送信息
		err = wxwork.SendMsg(user.UserName, passwordforuser, token, agentid)
		if err != nil {
			fmt.Printf("发送用户%s信息到企业微信失败,err:%s\n", user.UserName, err.Error())
		}

		//记录信息
		user.CreatedTime = time.Now()
		err = data.SaveToFile(user, cfg.Section("Data").Key("DbPath").String())
		if err != nil {
			fmt.Printf("保存用户信息到日志文件失败")
		}
	}
}

//func main()  {
//	user := data.UserObj{
//		UserName: "denghui",
//		UserGroup: "default",
//		CreatedTime: time.Now(),
//	}
//	err := data.SaveToFile(user, "user.log")
//	fmt.Println(err)
//}

//初始化配置
func initConfig(configPath string) (*ini.File, error) {
	cfg, err := ini.Load(configPath)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

//初始化redis客户端
func initRedis(addr, passwd string, db int) (*redis.Client, error) {
	var ctx = context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: passwd,
		DB:       db,
	})
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}
	return rdb, nil
}
