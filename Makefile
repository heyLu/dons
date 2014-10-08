PREFIX ?= /usr

all: dons

dons: server.go
	go get
	go build

install: dons
	cp dons ${PREFIX}/bin
	cp dons.service /etc/systemd/system

	grep -q docker.local /etc/dnsmasq.conf || echo 'server=/docker.local/127.0.0.1\#8053' >> /etc/dnsmasq.conf
	cp /etc/dnsmasq.conf /etc/NetworkManager/dnsmasq.d/conf

uninstall:
	rm -f /etc/systemd/system/dons.service
	rm -f ${PREFIX}/bin/dons
