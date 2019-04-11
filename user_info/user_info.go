package user_info

import (
	"fmt"
	"github.com/golang/go/src/os/user"
	"bufio"
	"os"
	"io"
	"strings"
	"strconv"
	"time"
)

type UserInfo struct {
	Name string `json:"user"`
	Uid string `json:"uid"`
	Group string `json:"group"`
        Expire string `json:"expire"`
}

var UserList []UserInfo

const passwdFile = "/etc/passwd"
const shadowFile = "/etc/shadow"

func GetUserInfo() ([]UserInfo,error) {
	sm := shadow()
	inputFile, inputError := os.Open(passwdFile)	
	if inputError != nil {
		return nil, inputError
	}
	defer inputFile.Close()

	inputReader := bufio.NewReader(inputFile)
	for {
		inputString, readError := inputReader.ReadString('\n')
		if readError == io.EOF {
			break
		}
		name := strings.Split(inputString,":")[0]
		userInfo, err := user.Lookup(name)
		if err != nil {
                        continue
                }
		groupInfo, err := user.LookupGroupId(userInfo.Gid)
		//fmt.Printf("%s %s %s %s\n", userInfo.Uid, userInfo.Gid, userInfo.Username, groupInfo.Name)
		ui := UserInfo{userInfo.Username,userInfo.Uid,groupInfo.Name,sm[userInfo.Username]}
		UserList = append(UserList,ui)
	}
	return UserList, nil
}

func shadow() map[string]string {
	var shadowMap = make(map[string]string)
	inputFile, inputError := os.Open(shadowFile)
	if inputError != nil {
		fmt.Print("%v",inputError)
		return nil
	}
	defer inputFile.Close()

	inputReader := bufio.NewReader(inputFile)
	for {
		inputString, readError := inputReader.ReadString('\n')
		if readError == io.EOF {
			break
		}
		s := strings.Split(inputString,":")
		if len(s[7]) == 0 {
			shadowMap[s[0]] = "Never"
		} else {
			start := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)	
			days, _ := strconv.Atoi(s[7])
                	end := start.AddDate(0,0,days).Format("2006-01-02 15:04:05")
			shadowMap[s[0]] = end
		}
	}
	return shadowMap
} 
