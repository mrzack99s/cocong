package utils

import (
	"net"
	"net/url"
)

func ParseURL(link string) (schema, host, port string, err error) {

	u, err := url.Parse(link)
	if err != nil {
		return
	}

	schema = u.Scheme
	host, port, err = net.SplitHostPort(u.Host)
	if err != nil {
		if schema == "https" {
			u.Host = u.Host + ":443"
		} else {
			u.Host = u.Host + ":80"
		}
		host, port, _ = net.SplitHostPort(u.Host)
	}

	return
}
