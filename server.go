package main

import (
	"fmt"
	"github.com/miekg/dns"
	"net"
	"os/exec"
	"strings"
)

func main() {
	dns.HandleFunc("docker.local", SimpleHandler)

	server := &dns.Server{Addr: ":8053", Net: "udp"}
	fmt.Printf("running server on %s\n", "127.0.0.1:8053")

	server.ListenAndServe()
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
		fmt.Printf("error getting ip of container: %s\n", err)
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
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}
