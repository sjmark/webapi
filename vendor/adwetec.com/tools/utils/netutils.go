package utils

import (
	"strings"
	"net"
	"os"
)

func ResetIp(ipStr string, resetCount int) string { // 重新设置IP(用于和屏蔽网段匹配)

	ipArray := strings.Split(ipStr, ".")

	for i := 1; i <= resetCount; i++ {
		ipArray[len(ipArray)-i] = "*"
	}

	var result string

	for _, str := range ipArray {
		result += str + "."
	}

	rs := []rune(result)

	return string(rs[0: len(rs)-1]) // 生成: xxx.xxx.*.*、xxx.xxx.xxx.*
}

func GetInternalHostIp() string {

	addrs, err := net.InterfaceAddrs()

	if err == nil {

		for _, a := range addrs {

			if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {

				if ipnet.IP.To4() != nil {
					return ipnet.IP.String()
				}

			}
		}
	}

	return "unknown"
}

func GetHostName() string {

	host, err := os.Hostname()

	if err != nil {
		return "unknown"
	} else {
		return host
	}
}
