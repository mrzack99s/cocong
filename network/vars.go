package network

import "github.com/coreos/go-iptables/iptables"

var (
	IPT *iptables.IPTables
	// last_accept_input_rule_num   int = 0
	last_accept_forward_rule_num int = 3
	// last_insert_init_prerouting  int = 0
)
