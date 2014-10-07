package main

import (
	"errors"
	"fmt"
	"github.com/miekg/dns"
	"net"
	"os"
	"os/exec"
	"strings"
)

func main() {
	addr := "127.0.0.1:8053"
	if len(os.Args) > 1 {
		addr = os.Args[1]
	}
	dns.HandleFunc("docker.local", SimpleHandler)

	server := &dns.Server{Addr: addr, Net: "udp"}
	fmt.Printf("running server on %s\n", addr)

	err := server.ListenAndServe()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func SimpleHandler(w dns.ResponseWriter, req *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(req)
	m.Authoritative = true

	fqdn := req.Question[0].Name
	fmt.Printf("resolving %s\n", req.Question[0].Name)
	split := strings.Split(fqdn, ".")
	if len(split) != 4 {
		w.WriteMsg(m)
		return
	}

	name := split[0]
	ip, err := GetContainerIp(name)
	if err != nil {
		fmt.Print(err)
		w.WriteMsg(m)
		return
	}
	ip = strings.TrimSpace(ip)
	fmt.Printf("%s resolves to %s\n", fqdn, ip)

	a := &dns.A{
		Hdr: dns.RR_Header{
			Name:   fqdn,
			Rrtype: dns.TypeA,
			Class:  dns.ClassINET,
			Ttl:    0,
		},
		A: net.ParseIP(ip).To4(),
	}
	m.Answer = append(m.Answer, a)

	w.WriteMsg(m)
}

func GetContainerIp(name string) (string, error) {
	cmd := exec.Command("docker", "inspect", "-f", "{{ .NetworkSettings.IPAddress }}", name)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", errors.New(string(out))
	}
	return string(out), nil
}
