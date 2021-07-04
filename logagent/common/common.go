package common

import (
	"fmt"
	"log"
	"net"
	"strings"
)

/**
 * @author xhli
 * @date 2021/6/26 20:17
 * @version 1.0
 * @description: TODO
 */

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