package server_info

import (
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/net"
	"github.com/shirou/gopsutil/host"
	"github.com/progrium/go-shell"
	"agent/utils"
	"strings"
	"time"
	"io/ioutil"
	"os"
	"strconv"
)

type cpuInfo struct {
	Vcpus int `json:"cpu_vcpus"`
	ModelName string `json:"cpu_modelname"`
	CacheSize int32 `json:"cpu_cachesize"`
        Mhz float64 `json:"cpu_mhz"`
        VendorID string `json:cpu_vendorid`
	Sockets int `json:"cpu_sockets"`
	Flags string `json:"cpu_flags"`
}

type memInfo struct {
	MemTotal string `json:"mem_total"`
	SwapMemTotal string `json:"swapmem_total"`
}

type diskInfo struct {
	Device string `json:"device"`
	Mountpoint string `json:"mountpoint"`
	Fstype string `json:"fstype"`
	Opts string `json:"disk_opts"`
	TotalSize string `json:"disk_totalsize"`
}

type hostInfo struct {
    Hostname             string `json:"hostname"`
    BootTime             string `json:"bootTime"`
    OS                   string `json:"os"`              // ex: freebsd, linux
    Platform             string `json:"platform"`        // ex: ubuntu, linuxmint
    PlatformFamily       string `json:"platform_family"`  // ex: debian, rhel
    PlatformVersion      string `json:"platform_version"` // version of the complete OS
    KernelVersion        string `json:"kernel_version"`   // version of the OS kernel (if available)
    VirtualizationSystem string `json:"virtualization_system"`
    VirtualizationRole   string `json:"virtualization_role"` // guest or host
    HostID               string `json:"product_uuid"`             // ex: uuid
    BiosDate		 string `json:"bios_date"`
    ProductName	 	 string `json:"product_name"`
    ProductSerial	 string `json:"product_serial"`
    SystemVendor	 string `json:"system_vendor"`
    DefaultIPV4          string `json:"default_ipv4"`
    DefaultIPV6          string `json:"default_ipv6"`
}

type ServerInfo struct {
	cpuInfo
	Disk []diskInfo `json:"disk"` 
	memInfo
	hostInfo
	Net []net.InterfaceStat `json:"net"`
}

func GetServerInfo() ServerInfo {
	ci := getCPUInfo()
	di := getDiskInfo()
	mi := getMemInfo()
	ni := getNetInfo()
	hi := getHostInfo()
	si := ServerInfo{ci,di,mi,hi,ni}
	return si
}

func getCPUInfo() cpuInfo {
	var ci cpuInfo
	var sockets []string
	vcpus, _ := cpu.Counts(true)
	ci.Vcpus = vcpus
	cpuInfo, _ := cpu.Info()
	ci.ModelName = cpuInfo[0].ModelName
	ci.CacheSize = cpuInfo[0].CacheSize
	ci.Mhz = cpuInfo[0].Mhz
	ci.VendorID = cpuInfo[0].VendorID
	ci.Flags = strings.Join(cpuInfo[0].Flags," ")
	for _, v := range cpuInfo {
		sockets = append(sockets, v.PhysicalID)
	}
        ci.Sockets = len(utils.RemoveDuplicateElement(sockets))	
	return ci
}

func getMemInfo() memInfo {
	var mi memInfo
	swapMem, _ := mem.SwapMemory()
	mi.SwapMemTotal = strconv.FormatUint(swapMem.Total / 1024 / 1024 / 1024,10) + "GB"
	mem, _ := mem.VirtualMemory()
	mi.MemTotal = strconv.FormatUint(mem.Total/ 1024 /1024 /1024,10) + "GB"
	return mi
}

func getDiskInfo() []diskInfo {
	var diskInfoList []diskInfo
	par, _ := disk.Partitions(false)
	for _, v := range par {
		//fmt.Printf("Serial Number:%s	Device:%s	MountPoint:%s	Fstype:%s	Opts:%s\n", disk.GetDiskSerialNumber(v.Device),v.Device,v.Mountpoint,v.Fstype,v.Opts)
		var di diskInfo
		di.Device = v.Device
		di.Mountpoint = v.Mountpoint
		di.Fstype = v.Fstype
		di.Opts = v.Opts
		ui, _ := disk.Usage(v.Mountpoint)
		di.TotalSize = strconv.FormatUint(ui.Total /1024/1024/1024,10) + "GB"
		diskInfoList = append(diskInfoList,di)	
	}
	return diskInfoList
}

func getNetInfo() []net.InterfaceStat {
	ni, _ := net.Interfaces()
	return ni
}

func getHostInfo() hostInfo {
	var hi hostInfo
	hInfo, _ := host.Info()
	hi.Hostname = hInfo.Hostname
	hi.BootTime = time.Unix(int64(hInfo.BootTime),0).Format("2006-01-02 15:04:05")
	hi.OS = hInfo.OS
	hi.Platform = hInfo.Platform
	hi.PlatformFamily = hInfo.PlatformFamily
	hi.PlatformVersion = hInfo.PlatformVersion
	hi.KernelVersion = hInfo.KernelVersion
	hi.VirtualizationSystem = hInfo.VirtualizationSystem
	hi.VirtualizationRole = hInfo.VirtualizationRole
	hi.HostID = hInfo.HostID
	df := getDmiFacts()
	hi.ProductName = df["product_name"]
	hi.BiosDate = df["bios_date"]
	hi.ProductSerial = df["product_serial"]
	hi.SystemVendor = df["system_vendor"]	
	hi.DefaultIPV4, hi.DefaultIPV6 = getDefaultIP()
	return hi
}

func getDmiFacts() map[string]string {
	if _, err := os.Stat("/sys/devices/virtual/dmi/id/product_name"); err == nil {
		df := make(map[string]string)
		item := map[string]string{
			"product_name": "/sys/devices/virtual/dmi/id/product_name",
			"bios_date": "/sys/devices/virtual/dmi/id/bios_date",
			"product_serial": "/sys/devices/virtual/dmi/id/product_serial",
			"system_vendor": "/sys/devices/virtual/dmi/id/sys_vendor",
		}
		for k,v := range item {
			//fmt.Printf("%s->%s\n",k,v)
			df[k] = getFileContent(v)	
		}
		return df
	}	
	return nil
}

func getDefaultIP() (string,string) {
        ipv4 := getIP("-4","8.8.8.8")
        ipv6 := getIP("-6","2404:6800:400a:800::1012")
        return ipv4, ipv6
}

func getIP(model string, dns string) string {
        ipPath := getCommandPath("ip")
        if ipPath == "NA" {
                return "NA"
        }
        ipOut := strings.Trim(shell.Cmd(ipPath,model,"route","get",dns).Run().Stdout.String(),"\n")
        ips := strings.Fields(ipOut)
        if len(ips) >0 && ips[0] == dns {
                for i, v := range ips {
                        if v == "src" {
                                return ips[i+1]
                        }
                }
        }
        return "NA"
}

func getCommandPath(cmd string) string {
        which := shell.Cmd("which").OutputFn()
        path, err := which(cmd)
        if err != nil {
                return "NA"
        }
        return path
}

func getFileContent(path string) string {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return "NA"
	}
	return strings.TrimSuffix(string(buf),"\n")
}
