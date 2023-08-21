package main

import (
	"fmt"
	"net"
)

type IPinfo struct {
	Myip   string
	IPlist []string
	Port   int
}

func NewIPinfo() (bool, *IPinfo) {
	fmt.Println("[readme]if needed, change the ip info from the 'ipinfo.go' file")

	// ============ change node ip and ip list here ============
	iplist := []string{"192.168.163.129", "192.168.163.130", "192.168.163.131"}
	myip := getMyIP() // manually input if needed
	fmt.Println("[please check]my ip:", myip)
	fmt.Println("[please check]iplist:", iplist)

	ipinfo := &IPinfo{
		Myip:   myip,
		IPlist: iplist,
		Port:   10006,
	}

	// check myip in iplist
	if !checkip(iplist, myip) {
		fmt.Println("[error]my ip is not in iplist, please check the 'ipinfo.go' file")
		return false, ipinfo
	}
	return true, ipinfo
}

// automatically get node first ipv4 address
func getMyIP() string {
	interfaces, _ := net.Interfaces()
	for _, iface := range interfaces {
		addresses, _ := iface.Addrs()
		for _, addr := range addresses {
			ipNet, ok := addr.(*net.IPNet)
			if ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
				return ipNet.IP.String()
			}
		}
	}
	return "ip not found"
}

// check ip in iplist
func checkip(iplist []string, ip string) bool {
	for _, value := range iplist {
		if value == ip {
			return true
		}
	}
	return false
}