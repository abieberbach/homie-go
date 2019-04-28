package homie_go

import (
	"net"
	"net/url"
)

func getIpAndMacAddr(mqttURL *url.URL) (ipAddr, macAddr string) {
	targetAddrs, _ := net.LookupIP(mqttURL.Hostname())
	interfaces, _ := net.Interfaces()
	for _, i := range interfaces {

		tcpAddr, found := getAddressForTarget(i, targetAddrs)
		if found {
			ipAddr = tcpAddr.IP.String()
			macAddr = i.HardwareAddr.String()
			return
		}
	}
	return
}

func getAddressForTarget(netInterface net.Interface, ipAdresses []net.IP) (*net.IPNet, bool) {
	if netInterface.Flags&net.FlagUp == 0 {
		return nil, false
	}
	addrs, _ := netInterface.Addrs()
	for _, currentAddr := range addrs {
		ipNet, ok := currentAddr.(*net.IPNet)
		if ok {
			for _, ipAddr := range ipAdresses {
				if ipNet.Contains(ipAddr) {
					return ipNet, true
				}

			}
		}
	}
	return nil, false
}
