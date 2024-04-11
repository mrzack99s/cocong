package network

import "github.com/coreos/go-iptables/iptables"

func init() {
	var err error
	IPT, err = iptables.New()
	if err != nil {
		return
	}
}
