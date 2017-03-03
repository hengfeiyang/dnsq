package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/miekg/dns"
	"github.com/safeie/daemon"
)

var conf *Configer

func main() {
	var err error
	conf, err = NewConfig("config.yml")
	if err != nil {
		panic(err)
	}

	if conf.Daemon {
		daemon.Daemon(0, 0)
	}

	udpServer := &dns.Server{Addr: conf.Server.UDPAddr, Net: "udp"}
	tcpServer := &dns.Server{Addr: conf.Server.TCPAddr, Net: "tcp"}
	dns.HandleFunc(".", route)
	go func() {
		if err := udpServer.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()
	go func() {
		if err := tcpServer.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	log.Printf("start ok, listen: UDP %v TCP %v\n", conf.Server.UDPAddr, conf.Server.TCPAddr)

	// Wait for SIGINT or SIGTERM
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	udpServer.Shutdown()
	tcpServer.Shutdown()
}

func route(w dns.ResponseWriter, req *dns.Msg) {
	if len(req.Question) == 0 {
		dns.HandleFailed(w, req)
		return
	}

	var err error
	for _, name := range conf.Rule {
		if strings.HasSuffix("."+req.Question[0].Name, "."+name+".") {
			for _, addr := range conf.DNS.OUT {
				fmt.Printf("hit: %v -> %v\n", name, addr)
				if err = proxy(addr, w, req); err == nil {
					return
				}
				fmt.Println("proxy hit error: ", err)
			}
			dns.HandleFailed(w, req)
			return
		}
	}

	for _, addr := range conf.DNS.IN {
		fmt.Printf("default: %v -> %v\n", req.Question[0].Name, addr)
		if err = proxy(addr, w, req); err == nil {
			return
		}
		fmt.Println("proxy default error: ", err)
	}
	dns.HandleFailed(w, req)
}

func isTransfer(req *dns.Msg) bool {
	for _, q := range req.Question {
		switch q.Qtype {
		case dns.TypeIXFR, dns.TypeAXFR:
			return true
		}
	}
	return false
}

func proxy(addr string, w dns.ResponseWriter, req *dns.Msg) error {
	transport := "udp"
	if _, ok := w.RemoteAddr().(*net.TCPAddr); ok {
		transport = "tcp"
	}
	if isTransfer(req) {
		if transport != "tcp" {
			return errors.New("unsupported protocol")
		}
		t := new(dns.Transfer)
		c, err := t.In(req, addr)
		if err != nil {
			return err
		}
		if err = t.Out(w, req, c); err != nil {
			return err
		}
		return nil
	}
	c := &dns.Client{Net: transport}
	resp, _, err := c.Exchange(req, addr)
	if err != nil {
		return err
	}
	w.WriteMsg(resp)
	return nil
}
