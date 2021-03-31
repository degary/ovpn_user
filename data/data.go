package data

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

type UserObj struct {
	UserName    string
	UserGroup   string
	CreatedTime time.Time
}

func GetUserObjs(filepath string) ([]UserObj, error) {
	var userObjs = []UserObj{}
	file, err := os.Open(filepath)
	defer file.Close()
	if err != nil {
		return nil, err
	}
	br := bufio.NewReader(file)
	for {
		s, _, err := br.ReadLine()
		if err == io.EOF {
			break
		}
		split := strings.Split(string(s), ",")
		if len(split) != 2 && len(split) != 1 {
			return nil, fmt.Errorf("文件:%s格式错误,请以','分隔", filepath)
		}
		if len(split) == 1 {
			userObjs = append(userObjs, UserObj{UserName: split[0], UserGroup: "default_group"})
			continue
		}
		userObjs = append(userObjs, UserObj{UserName: split[0], UserGroup: split[1]})
	}
	return userObjs, nil
}

func SaveToFile(user UserObj, filepath string) error {
	f, err := os.OpenFile(filepath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	defer f.Close()
	if err != nil {
		return err
	}
	data, _ := json.Marshal(&user)
	data = append(data, []byte("\n")...)
	_, err = f.Write(data)
	if err != nil {
		return err
	}
	return nil
}
