package software_info

import (
	"github.com/progrium/go-shell"
	"strings"
	"os"
)
var (
  rpm = shell.Cmd("rpm").OutputFn()
)

type SoftwareInfo struct {
	Name string `json:"name"`
	Version string `json:"version"`
}

var softwareInfoList []SoftwareInfo

func GetVersion(name string, exe string) string {
	_, err := rpm("-qf",exe)
	if err == nil {
		return "rpm"
	}
	if _, err := os.Stat(exe); err != nil {
		return "Deleted"
	}
	if name == "nginx" {
		version := shell.Cmd(exe,"-v").Run()
        	info := strings.Trim(version.Stderr.String(),"\n")
        	v := strings.Split(info,"/")
		return v[len(v)-1]
	}
	if name == "redis-server" {
		version := strings.Trim(shell.Cmd(exe,"-v").Pipe("awk","'{print $3}'").Pipe("awk","-F","'='","'{print $2}'").Run().Stdout.String(),"\n")
		return version
	}
	if name == "etcd" {
		version := strings.Trim(shell.Cmd(exe,"-version").Pipe("awk","'/etcd Version/ {print $NF}'").Run().Stdout.String(),"\n")
		return version
	}
	if name == "mysqld" {
		version := strings.Trim(shell.Cmd(exe,"--version").Pipe("awk","'{print $3}'").Run().Stdout.String(),"\n")
		return version
	}
	if name == "server_agent" {
                version := strings.Trim(shell.Cmd(exe,"-v").Run().Stdout.String(),"\n")
                return version
        }
	if name == "java" {
		version := shell.Cmd(exe,"-version").Run()
		info := strings.Trim(version.Stderr.String(),"\n")
		v := strings.Split(info,"\"")
		return v[1]
	}
	return "Unknown"
}

func GetSoftwareInfo() []SoftwareInfo {
	out, _ := rpm("-qa")
	s := strings.Split(out,"\n") 
	for _,v := range s {
		tempMap := make(map[string]string)
		info, _ := rpm("-qi",v)
		for _, i := range strings.Split(info,"\n") {
			tmpSlice := strings.Split(i,":")
			if len(tmpSlice) == 2 {
				key := strings.Trim(tmpSlice[0]," ")
				value := strings.Trim(tmpSlice[1]," ")
				if value != "" {
					tempMap[key] = value
				}
			}
			if len(tmpSlice) >= 3 {
				key1 := strings.Trim(tmpSlice[0]," ")
				values := strings.Trim(tmpSlice[1]," ")
				tmp := strings.Split(values," ")
				value1 := tmp[0]
				key2 := tmp[len(tmp)-1] 
				value2 := strings.Trim(tmpSlice[2]," ")
				tempMap[key1] = value1
				tempMap[key2] = value2
			}
		}
		si := SoftwareInfo{tempMap["Name"],tempMap["Version"] + "-" + tempMap["Release"]}
		softwareInfoList = append(softwareInfoList,si)
	}
	return softwareInfoList
}
