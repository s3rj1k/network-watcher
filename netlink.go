package main

import (
	"net"

	"github.com/vishvananda/netlink"
)

func linkHasRoutableAddressAndDefaultRoute(name string) bool {
	link, err := netlink.LinkByName(name)
	if err != nil {
		return false
	}

	attrs := link.Attrs()
	if attrs == nil {
		return false
	}

	if attrs.Flags&net.FlagUp != net.FlagUp {
		return false
	}

	addrs, err := netlink.AddrList(link, netlink.FAMILY_ALL)
	if err != nil {
		return false
	}

	if len(addrs) == 0 {
		return false
	}

	var hasAddr bool

	for _, el := range addrs {
		if el.IP.IsInterfaceLocalMulticast() ||
			el.IP.IsLinkLocalMulticast() ||
			el.IP.IsLinkLocalUnicast() ||
			el.IP.IsLoopback() ||
			el.IP.IsMulticast() {
			continue
		}

		hasAddr = true

		break
	}

	if !hasAddr {
		return false
	}

	routes, err := netlink.RouteList(link, netlink.FAMILY_ALL)
	if err != nil {
		return false
	}

	if len(routes) == 0 {
		return false
	}

	for _, el := range routes {
		if el.Src != nil &&
			el.Gw != nil &&
			el.Dst == nil { // default route
			return true
		}
	}

	return false
}
