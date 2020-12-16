package main

import (
	"bytes"
	"encoding"
	"flag"
	"github.com/mmmorris1975/simple-logger/logger"
	"net"
	"os"
	"path/filepath"
	"sync"
	"time"
	"x32-osc-bridge/osc"
	"x32-osc-bridge/x32"
)

var (
	lsnr    *net.UDPConn
	mixer   *x32.Mixer
	iaddrs  []net.Addr
	clients sync.Map

	log = logger.StdLogger

	port    int
	verbose bool

	// Version is the program version, managed at build time
	Version string
)

func init() {
	flag.IntVar(&port, "p", 0, "client port")
	flag.BoolVar(&verbose, "v", false, "enable verbose output")
	log.Infof("%s %s", filepath.Base(os.Args[0]), Version)
}

func main() {
	var err error
	flag.Parse()

	if verbose {
		log.SetLevel(logger.DEBUG)
		x32.Log.SetLevel(logger.DEBUG)
	}

	if err = initialize(); err != nil {
		log.Fatal(err)
	}
	defer lsnr.Close()
	log.Infoln(lsnr.LocalAddr())

	var n int
	var raddr *net.UDPAddr
	var buf = make([]byte, 65535)
	for {
		n, raddr, err = lsnr.ReadFromUDP(buf)
		if err != nil {
			log.Errorf("ReadFromUDP: %v", err)
			continue
		}

		if port > 0 {
			raddr.Port = port
		}

		go handleMsg(buf[:n], raddr)
	}
}

func initialize() (err error) {
	log.Debugln("Discover")
	mixer, err = x32.Discover(nil)
	if err != nil {
		return
	}

	log.Debugln("ListenUDP")
	lsnr, err = net.ListenUDP("udp", &net.UDPAddr{Port: x32.OscPort})
	if err != nil {
		return
	}

	log.Debugln("InterfaceAddrs")
	iaddrs, err = net.InterfaceAddrs()
	if err != nil {
		return
	}

	return
}

func readLoop(conn *net.UDPConn, addr *net.UDPAddr) {
	buf := make([]byte, 65535)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			if e, ok := err.(*net.OpError); ok {
				if !e.Temporary() {
					// conn is likely closed, so we should terminate the loop
					log.Debugf("mixer read: %v", e)
					return
				}
			}
			log.Errorf("readLoop:Read() - %v", err)
		}

		if _, err := lsnr.WriteToUDP(buf[:n], addr); err != nil {
			log.Errorf("readLoop:WriteToUDP() - %v", err)
			return
		}
	}
}

func handleMsg(msg []byte, raddr *net.UDPAddr) {
	var err error

	if bytes.HasPrefix(msg[:16], osc.WriteString("/xinfo")) {
		xinfoHandler(raddr)
		return
	}

	if bytes.HasPrefix(msg[:16], osc.WriteString("/status")) {
		statusHandler(raddr)
		return
	}

	c, ok := clients.Load(raddr.String())
	if !ok {
		log.Infof("new client: %v", raddr)
		var cl *x32.Client
		cl, err = x32.NewClient(mixer.Addr)
		if err != nil {
			log.Errorf("NewClient: %v", err)
			return
		}

		c, _ = clients.LoadOrStore(raddr.String(), cl)
		cl = c.(*x32.Client)
		cl.StartXremote(3 * time.Second)
		go readLoop(cl.Conn, raddr)
	}

	// if the message contains the string subscribe anywhere in the 1st 16 bytes, use it as a flag for an
	// "active" client, like X32-Edit, which manages it's own synchronization with the mixer state.  For these
	// cases, we'll stop the xremote process, which is used to keep "passive" clients informed of changes to
	// the mixer state. This should cover messages for /formatsubscribe, /batchsubscribe, and /unsubscribe.
	// NOTE - it may take up to 10 seconds for xremote to stop sending back data
	if bytes.Contains(msg[:16], []byte("subscribe")) {
		log.Debugf("%s", msg[:16])
		c.(*x32.Client).StopXremote()
	}

	if _, err = c.(*x32.Client).Conn.Write(msg); err != nil {
		log.Errorf("mixer write: %v", err)
	}

	if bytes.HasPrefix(msg, osc.WriteString("/unsubscribe")) {
		unsubscribeHandler(raddr)
	}
}

func xinfoHandler(raddr *net.UDPAddr) {
	ip := findInterface(raddr)
	msg := &osc.XInfo{
		IP:      ip,
		Name:    mixer.Info.Name,
		Model:   mixer.Info.Model,
		Version: mixer.Info.Version,
	}
	log.Debugf("xinfo: %+v", *msg)

	if err := writeMsg(msg, raddr); err != nil {
		log.Errorf("xinfo handler: %v", err)
	}
}

func statusHandler(raddr *net.UDPAddr) {
	ip := findInterface(raddr)
	msg := &osc.Status{
		State: "active",
		IP:    ip,
		Name:  mixer.Info.Name,
	}
	log.Debugf("status: %+v", *msg)

	if err := writeMsg(msg, raddr); err != nil {
		log.Errorf("status handler: %v", err)
	}
}

func unsubscribeHandler(raddr *net.UDPAddr) {
	if c, ok := clients.LoadAndDelete(raddr.String()); ok {
		log.Infof("shutting down client: %v", raddr)
		cl := c.(*x32.Client)
		cl.StopXremote()
		_ = cl.Conn.Close()
	}
}

func writeMsg(d encoding.BinaryMarshaler, raddr *net.UDPAddr) error {
	data, err := d.MarshalBinary()
	if err != nil {
		return err
	}

	if _, err = lsnr.WriteToUDP(data, raddr); err != nil {
		return err
	}

	return nil
}

func findInterface(addr *net.UDPAddr) net.IP {
	for _, n := range iaddrs {
		if ipn, ok := n.(*net.IPNet); ok {
			if ipn.Contains(addr.IP) {
				return ipn.IP
			}
		}
	}
	return nil
}
