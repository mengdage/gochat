package helper

import (
	"net"
)

func GetServerIp() (ip string) {
	// isLocal := viper.GetBool("app.local_test")
	// if isLocal {
	// 	ip = "127.0.0.1"
	// 	return
	// }

	addrs, err := net.InterfaceAddrs()

	if err != nil {
		return ""
	}

	for _, address := range addrs {
		if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ip = ipNet.IP.String()
			}
		}
	}
	return
}
