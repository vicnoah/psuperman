package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/wuwengang/psuperman/core/parser"

	"github.com/wuwengang/psuperman/core/packet"

	"github.com/wuwengang/psuperman/core/listen"

	"golang.org/x/net/ipv4"
)

const (
	// I use TUN interface, so only plain IP packet, no ethernet header + mtu is set to 1300
	BUFFERSIZE = 1500
	MTU        = "1300"
)

var (
	localIP  = flag.String("local", "", "Local tun interface IP/MASK like 192.168.3.3⁄24")
	remoteIP = flag.String("remote", "", "Remote server (external) IP like 8.8.8.8")
	port     = flag.Int("port", 4321, "UDP port for communication")
	tun      = &listen.Tun{}
	local    = &listen.Udp{}
	prs      = &parser.Parser{}
)

func main() {
	flag.Parse()
	// check if we have anything
	if "" == *localIP {
		flag.Usage()
		log.Fatalln("\nlocal ip is not specified")
	}
	if "" == *remoteIP {
		flag.Usage()
		log.Fatalln("\nremote server is not specified")
	}
	err := tun.Listen(*localIP, MTU, BUFFERSIZE)
	fmt.Println(err)

	err = local.Listen(fmt.Sprintf(":%v", *port), BUFFERSIZE)
	fmt.Println(err)

	defer func() {
		_ = local.Close()
	}()

	// recv udp packet to tun
	go func() {
		for {
			n, addr, payload, err := local.Read()
			// just debug
			header, _ := ipv4.ParseHeader(payload)
			fmt.Printf("Received %d bytes from %v: %+v\n", n, addr, header)
			prs.SetHeader(*header)
			lan := prs.IsLan()
			if lan {
				fmt.Println("局域网ip")
			} else {
				fmt.Println("广域网ip")
			}
			protocol, err := packet.Ipv4Packet{}.Protocol(header.Protocol)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(protocol)
			}
			if err != nil || n == 0 {
				fmt.Println("Error: ", err)
				continue
			}
			// write to TUN interface
			_, _ = tun.Write(payload)
		}
	}()

	// recv tun packet to udp
	for {
		plen, packet, err := tun.Read()
		if err != nil {
			break
		}
		// debug :)
		header, _ := ipv4.ParseHeader(packet[:plen])
		fmt.Printf("Sending to remote: %+v (%+v)\n", header, err)
		// real send
		addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%v", *remoteIP, *port))
		if err != nil {
			return
		}
		_, _ = local.Write(packet, addr)
	}
}
