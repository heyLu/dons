package main

import (
	"fmt"
	"github.com/miekg/dns"
	"net"
	// "os"
)

func main() {
	dns.HandleFunc("example.com.", SimpleHandler)

	server := &dns.Server{Addr: ":8053", Net: "udp"}
	fmt.Printf("running server on %s\n", "127.0.0.1:8053")

	server.ListenAndServe()
}

func SimpleHandler(w dns.ResponseWriter, req *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(req)
	m.Authoritative = true

	a := &dns.A{
		Hdr: dns.RR_Header{
			Name:   req.Question[0].Name,
			Rrtype: dns.TypeA,
			Class:  dns.ClassINET,
			Ttl:    0,
		},
		A: net.ParseIP("42.42.42.42").To4(),
	}
	m.Answer = append(m.Answer, a)

	w.WriteMsg(m)
}
