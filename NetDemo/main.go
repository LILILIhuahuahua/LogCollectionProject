package main

import (
	"fmt"
	"log"
	"net"
	"strings"
)

/**
 * @author xhli
 * @date 2021/6/26 20:06
 * @version 1.0
 * @description: TODO
 */

//获取本地IP的方式
func main(){
	//GetLocalIP()
	GetOutboundIP()
}

func GetLocalIP() (ip string, err error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return
	}
	for _, addr := range addrs {
		ipAddr, ok := addr.(*net.IPNet) // 类型断言
		if !ok {
			continue
		}

		if ipAddr.IP.IsLoopback() {
			continue
		}

		if !ipAddr.IP.IsGlobalUnicast() {
			continue
		}
		fmt.Println(ipAddr)
		return ipAddr.IP.String(), nil
	}
	return
}


// Get preferred outbound ip of this machine
func GetOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	fmt.Println(localAddr.String())
	return strings.Split(localAddr.IP.String(),":")[0]
}