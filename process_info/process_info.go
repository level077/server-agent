package process_info

import (
	"github.com/shirou/gopsutil/process"
	"time"
	"strconv"
	"agent/utils"
	"path"
	"agent/software_info"
)

type ProcessInfo  struct {
	Name string `json:"name"`
	Exe string `json:"exe"`
	CreateTime string `json:"createtime"`
	Pid int32 `json:"pid"`
	Parent string `json:"parent"`
	Cmdline string	`json:"cmdline"`
	Ports []string `json:"ports"`
	Listen []string `json:"listen"`
}

var ProcessList []ProcessInfo

func GetProcessInfo() ([]ProcessInfo, []software_info.SoftwareInfo, error) {
	si := software_info.GetSoftwareInfo()
	processes, err := process.Processes()
	if err != nil {
		return nil, si, err
	}
	for _, v := range processes {
		exe, err := v.Exe()
                if err != nil {
                        continue
                }
		name, err := v.Name()
                if err != nil {
                        return nil, si, err
                }
		if path.Base(exe) != name {
			continue
		} 
		pi := ProcessInfo{}
		version := software_info.GetVersion(name, exe)
		if version != "rpm" {
			if contains(si,software_info.SoftwareInfo{name,version}) != true {
				si = append(si,software_info.SoftwareInfo{name,version})
			}
		}
		pi.Exe = exe
		pi.Name = name
		pi.Pid = v.Pid
		create_time, err := v.CreateTime()
  		if err != nil {
			return nil, si, err
		}	
		pi.CreateTime =time.Unix(create_time/1000,0).Format("2006-01-02 15:04:05")
		cmdline, err := v.Cmdline()
		if err != nil {
			return nil, si, err
                }
		pi.Cmdline = cmdline
		conns, err := v.Connections()
		if err != nil {
			return nil, si, err
                }
		var tempPort []string
		var tempListen []string
		for _, c := range conns {
			if c.Status == "LISTEN" {
				port := strconv.Itoa(int(c.Laddr.Port))
				tempPort = append(tempPort,port)
				tempListen = append(tempListen,c.Laddr.IP)
			}
		}
		pi.Ports = utils.RemoveDuplicateElement(tempPort)	
		pi.Listen = utils.RemoveDuplicateElement(tempListen)
		ppid, err := v.Ppid()
		if err != nil {
			return nil, si, err
                }
		parent_name := "None"
		if ppid >0 {
			parent, err := v.Parent()
			if err != nil {
				return nil, si, err
                        }
			parent_name, err = parent.Name()
			if err != nil {
				return nil, si, err
                	}	
		} 
		pi.Parent = parent_name
		ProcessList = append(ProcessList,pi)
	}
	return ProcessList, si, nil
}

func contains(ele []software_info.SoftwareInfo, s software_info.SoftwareInfo) bool {
	for _, v := range ele {
		if v == s {
			return true
		}
	}	
	return false
}
