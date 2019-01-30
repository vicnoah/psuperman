package listen

import (
	"log"
	"net"
	"os"
	"os/exec"

	"github.com/songgao/water"
)

func Tap() {

}

type Tun struct {
	iface     *water.Interface
	bufferLen int
}

func (t *Tun) Listen(ip string, mtu string, bufferLen int) (err error) {
	config := water.Config{
		DeviceType: water.TUN,
	}
	config.Name = "O_O"

	t.iface, err = water.New(config)
	if err != nil {
		return
	}
	t.bufferLen = bufferLen

	runIP("link", "set", "dev", t.iface.Name(), "mtu", mtu)
	runIP("addr", "add", ip, "dev", t.iface.Name())
	runIP("link", "set", "dev", t.iface.Name(), "up")
	return
}

func (t *Tun) Read() (n int, payload []byte, err error) {
	packet := make([]byte, t.bufferLen)
	n, err = t.iface.Read(packet)
	payload = packet
	return
}

func (t *Tun) Write(p []byte) (n int, err error) {
	return t.iface.Write(p)
}

func Tcp() {

}

type Udp struct {
	conn      *net.UDPConn
	Addr      *net.UDPAddr
	BufferLen int
}

func (u *Udp) Listen(address string, bufferLen int) (err error) {
	u.BufferLen = bufferLen
	u.Addr, err = net.ResolveUDPAddr("udp", address)
	if err != nil {
		return
	}
	u.conn, err = net.ListenUDP("udp", u.Addr)
	return
}

func (u *Udp) Read() (n int, addr *net.UDPAddr, payload []byte, err error) {
	packet := make([]byte, u.BufferLen)
	n, addr, err = u.conn.ReadFromUDP(packet)
	payload = packet
	return
}

func (u *Udp) Write(p []byte, addr *net.UDPAddr) (n int, err error) {
	n, err = u.conn.WriteToUDP(p, addr)
	return
}

func (u *Udp) Close() error {
	return u.conn.Close()
}

func runIP(args ...string) {
	cmd := exec.Command("/sbin/ip", args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if nil != err {
		log.Fatalln("Error running /sbin/ip:", err)
	}
}
