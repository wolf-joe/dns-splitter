package common

import (
	"fmt"
	"github.com/miekg/dns"
	"net"
	"strings"
)

// ExtractA 提取dns响应中的A记录
func ExtractA(r *dns.Msg) (records []*dns.A) {
	if r == nil {
		return
	}
	for _, answer := range r.Answer {
		switch answer.(type) {
		case *dns.A:
			records = append(records, answer.(*dns.A))
		}
	}
	return
}

// ParseSubnet 将字符串（IP/CIDR）转换为EDNS CLIENT SUBNET对象
func ParseSubnet(s string) (ecs *dns.EDNS0_SUBNET, err error) {
	if s == "" {
		return nil, nil
	}
	if strings.Contains(s, "/") { // 解析网段
		ipAddr, ipNet, err := net.ParseCIDR(s)
		if err != nil {
			return nil, err
		}
		mask, _ := ipNet.Mask.Size()
		ecs = &dns.EDNS0_SUBNET{Address: ipAddr, SourceNetmask: uint8(mask)}
	} else { // 解析ip
		addr, mask := net.ParseIP(s), uint8(0)
		if addr.To4() != nil {
			mask = uint8(net.IPv4len * 8)
		} else if addr.To16() != nil {
			mask = uint8(net.IPv6len * 8)
		} else {
			return nil, fmt.Errorf("wrong ip address: %s", s)
		}
		ecs = &dns.EDNS0_SUBNET{Address: addr, SourceNetmask: mask}
	}
	if ecs.Address.To4() != nil {
		ecs.Family = uint16(1)
	} else {
		ecs.Family = uint16(2)
	}
	return ecs, nil
}

// FormatSubnet 将DNS请求/响应里的EDNS CLIENT SUBNET对象格式化为字符串
func FormatSubnet(r *dns.Msg) string {
	if r == nil {
		return ""
	}
	for _, extra := range r.Extra {
		switch extra.(type) {
		case *dns.OPT:
			for _, opt := range extra.(*dns.OPT).Option {
				switch opt.(type) {
				case *dns.EDNS0_SUBNET:
					ecs := opt.(*dns.EDNS0_SUBNET)
					return fmt.Sprintf("%s/%d", ecs.Address, ecs.SourceNetmask)
				}
			}
		}
	}
	return ""
}
