package network

import "github.com/coreos/go-iptables/iptables"

var (
	IPT                          *iptables.IPTables
	Last_accept_input_rule_num   int = 2
	Last_accept_forward_rule_num int = 3
	Last_insert_prerouting       int = 4

	AuthorizedNetworks    = []string{}
	SSHAuthorizedNetworks = []string{}
	BypassedNetworks      = []string{}
)
