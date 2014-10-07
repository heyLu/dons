# dons

Resolve `<container-id>.docker.local` to the ip of that container via DNS.

## usage

	$ ./dons
	running server on 127.0.0.1:8053
	# somewhere else ...
	$ dig +nocmd @127.0.0.1 -p 8053 53a98c2ac.docker.local
	[...]
	;; ANSWER SECTION:
	53a98c2ac.docker.local. 0       IN      A       172.17.0.11

To use it directly, you have to run it as `./dons 127.0.0.1:53` and configure
it as your DNS server:

	# /etc/resolv.conf
	nameserver 127.0.0.1

However, you need root privileges to do that, so you might want to install
`dnsmasq` and configure it as follows in addition to the above:

	# /etc/dnsmasq.conf
	server=/docker.local/127.0.0.1#8053

The you can run `./dons` as a regular user.

## installation

use a [binary release](https://github.com/heyLu/dons/releases), or build
it yourself using `go get . && go build`.