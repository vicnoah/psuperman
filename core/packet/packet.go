package packet

import "errors"

type Ipv4Packet struct {

}

func (i Ipv4Packet) Protocol(protocol int) (string, error) {
	switch protocol {
	case 0x01:
		return "ICMP", nil
	case 0x06:
		return "TCP", nil
	case 0x11:
		return "UDP", nil
	default:
		return "", errors.New("not supported for this protocol")
	}
}