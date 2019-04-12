package main
import (
	"fmt"
	"flag"
	"encoding/json"
	"agent/user_info"
	"agent/process_info"
	"agent/server_info"
	"agent/software_info"
	"agent/elastic_send"
	"os"
	"time"
	"strings"
)

type serverInfo struct {
	server_info.ServerInfo
	TimeStamp string `json:"@timestamp"`
}

type baseInfo struct {
        HostID string `json:"product_uuid"`
        Hostname string `json:"hostname"`
        TimeStamp string `json:"@timestamp"`
        DefaultIPV4 string `json:"default_ipv4"`
        DefaultIPv6 string `json:"default_ipv6"`
}

type userInfo struct {
	user_info.UserInfo
	baseInfo
} 

type processInfo struct {
	process_info.ProcessInfo
	baseInfo
}

type softwareInfo struct {
	software_info.SoftwareInfo
	baseInfo
}

type EsMeta struct {
        EsIndex string `json:"_index"`
        EsType string `json:"_type"`
}

type Meta struct {
        Index EsMeta `json:"index"`
}


var (
	help bool
	isTest bool
	version bool
	esHost string
	esPort string
)

func init() {
	flag.BoolVar(&help,"h",false,"This help")	
	flag.BoolVar(&isTest,"t",false,"test mode, not send to elastic")
	flag.BoolVar(&version,"v",false,"show version & exit")
	flag.StringVar(&esHost,"e","","`esHost`")
	flag.StringVar(&esPort,"p","","`esPort`")
	flag.Usage = usage
}

const v = "0.0.1"

func main() {
	flag.Parse()
	if help {
		flag.Usage()
	}
	if version {
		fmt.Printf("%s\n",v)
		os.Exit(0)
	}

	if !isTest && (esHost == "" || esPort == "") {
		flag.Usage()
		os.Exit(1)
	}

	t := time.Now().Format(time.RFC3339)
	sInfo := server_info.GetServerInfo()
	product_uuid := sInfo.HostID
	hostname := sInfo.Hostname
	default_ipv4 := sInfo.DefaultIPV4
        default_ipv6 := sInfo.DefaultIPV6
        bi := baseInfo{product_uuid,hostname,t,default_ipv4,default_ipv6}

	tmpServerInfo := serverInfo{sInfo,t}
        serverInfoSend(esHost, esPort, "server","server",tmpServerInfo)

	ui, err := user_info.GetUserInfo()
	if err != nil {
		fmt.Printf("%v\n",err)
		return
	}
	var tmpUi []userInfo
	for _, v := range ui {
		tmpUi = append(tmpUi,userInfo{v,bi})
	}
	userInfoSend(esHost, esPort, "user","user",tmpUi)

	pi, si, err := process_info.GetProcessInfo()
	if err != nil {
		fmt.Printf("%v\n",err)
		return
	}
	var tmpPi []processInfo
	var tmpSi []softwareInfo
	for _, v := range pi {
		tmpPi = append(tmpPi, processInfo{v,bi})
	} 
	processInfoSend(esHost, esPort, "process","process",tmpPi)

	for _, v := range si {
		tmpSi = append(tmpSi, softwareInfo{v,bi})
	}
	softwareInfoSend(esHost, esPort, "software","software",tmpSi)
}

func usage() {
    	flag.PrintDefaults()
	os.Exit(0)
}

func userInfoSend(esHost string, esPort string, index string, docType string, ui []userInfo) {
	esMeta := EsMeta{index,docType}
	meta, _:= json.Marshal(Meta{esMeta})	
	meta_str := string(meta)
	var body []string
	for _,v := range ui {
		u, _ := json.Marshal(v)
		u_str := string(u)
		body = append(body,meta_str)
		body = append(body,u_str)
	}
	esBody := strings.Join(body,"\n") + "\n"
	if !isTest {
		status := elastic_send.Bulk(esBody,esHost,esPort)
		fmt.Println("index userinfo: ",status)
	} else {
		fmt.Printf("%s",esBody)
                fmt.Printf("--------------------------\n")
	}	
}

func processInfoSend(esHost string, esPort string, index string, docType string, pi []processInfo) {
	esMeta := EsMeta{index,docType}
        meta, _:= json.Marshal(Meta{esMeta})
        meta_str := string(meta)
        var body []string
        for _,v := range pi {
                u, _ := json.Marshal(v)
                u_str := string(u)
                body = append(body,meta_str)
                body = append(body,u_str)
        }
        esBody := strings.Join(body,"\n") + "\n"
	if !isTest {
                status := elastic_send.Bulk(esBody,esHost,esPort)
                fmt.Println("index processinfo: ",status)
        } else {
                fmt.Printf("%s",esBody)
                fmt.Printf("--------------------------\n")
        }
}

func softwareInfoSend(esHost string, esPort string, index string, docType string, si []softwareInfo) {
        esMeta := EsMeta{index,docType}
        meta, _:= json.Marshal(Meta{esMeta})
        meta_str := string(meta)
        var body []string
        for _,v := range si {
                u, _ := json.Marshal(v)
                u_str := string(u)
                body = append(body,meta_str)
                body = append(body,u_str)
        }
        esBody := strings.Join(body,"\n") + "\n"
	if !isTest {
                status := elastic_send.Bulk(esBody,esHost,esPort)
                fmt.Println("index softwarinfo: ",status)
        } else {
                fmt.Printf("%s",esBody)
                fmt.Printf("--------------------------\n")
        }
}

func serverInfoSend(esHost string, esPort string, index string, docType string, sInfo serverInfo) {
        esMeta := EsMeta{index,docType}
        meta, _:= json.Marshal(Meta{esMeta})
        meta_str := string(meta)
        var body []string
        u, _ := json.Marshal(sInfo)
        u_str := string(u)
        body = append(body,meta_str)
        body = append(body,u_str)
        esBody := strings.Join(body,"\n") + "\n"
	if !isTest {
                status := elastic_send.Bulk(esBody,esHost,esPort)
                fmt.Println("index serverinfo: ",status)
        } else {
                fmt.Printf("%s",esBody)
                fmt.Printf("--------------------------\n")
        }
}
