package parser

import (
	"fmt"

	"golang.org/x/net/ipv4"
)

type Parser struct {
	header ipv4.Header
}

func (p *Parser) SetHeader(header ipv4.Header) {
	p.header = header
}

func (p *Parser) IsLan() bool {
	ip := []byte(p.header.Dst)
	if len(ip) == 4 {
		switch ip[1] {
		case 10:
			return true
		case 172:
			if ip[2] >= 16 && ip[2] <= 31 {
				return true
			}
			return false
		case 192:
			if ip[2] == 168 {
				return true
			}
			return false
		default:
			return false
		}
	} else {
		switch ip[12] {
		case 10:
			return true
		case 172:
			if ip[13] >= 16 && ip[13] <= 31 {
				return true
			}
			return false
		case 192:
			if ip[13] == 168 {
				return true
			}
			return false
		default:
			return false
		}
	}
}

func (p *Parser) Parser() {
	fmt.Println(p.header.Src)
	fmt.Println(p.header.Dst)
	fmt.Println(p.header.TTL)
	fmt.Println(p.header.Version)
}
