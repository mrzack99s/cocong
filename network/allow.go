package network

import (
	"github.com/mrzack99s/cocong/types"
	"github.com/mrzack99s/cocong/utils"
	"github.com/mrzack99s/cocong/vars"
)

func AllowAccess(ss *types.SessionInfo) (err error) {

	err = IPT.Insert("nat", "PREROUTING", Last_insert_prerouting, "-s", ss.IPAddress, "-i", vars.Config.SecureInterface, "-j", "ACCEPT")
	if err != nil {
		return
	}

	err = IPT.Insert("filter", "FORWARD", Last_accept_forward_rule_num, "-s", ss.IPAddress, "-i", vars.Config.SecureInterface, "-o", vars.Config.EgressInterface, "-j", "ACCEPT")
	if err != nil {
		return
	}

	return
}

func AllowAccessBypass(ss *types.SessionInfo) (err error) {
	interfaceIp, _ := utils.GetSecureInterfaceIpv4Addr()

	err = IPT.Insert("nat", "PREROUTING", 4, "-s", ss.IPAddress, "-i", vars.Config.SecureInterface, "-j", "ACCEPT")
	if err != nil {
		return
	}

	Last_insert_prerouting += 1

	err = IPT.Insert("filter", "FORWARD", Last_accept_forward_rule_num, "-s", ss.IPAddress, "-i", vars.Config.SecureInterface, "-o", vars.Config.EgressInterface, "-j", "ACCEPT")
	if err != nil {
		return
	}

	Last_accept_forward_rule_num += 1

	err = IPT.Insert("filter", "INPUT", 1, "-p", "udp", "-s", ss.IPAddress, "-i", vars.Config.SecureInterface, "--dport", "53", "-d", interfaceIp, "-j", "ACCEPT")
	if err != nil {
		return
	}

	err = IPT.Insert("filter", "INPUT", 1, "-p", "tcp", "-s", ss.IPAddress, "-i", vars.Config.SecureInterface, "--dport", "53", "-d", interfaceIp, "-j", "ACCEPT")
	if err != nil {
		return
	}

	err = IPT.Insert("filter", "INPUT", 1, "-s", ss.IPAddress, "-i", vars.Config.SecureInterface, "-p", "tcp", "--match", "multiport", "--dports", "443,8080,8443", "-d", interfaceIp, "-j", "DROP")
	if err != nil {
		return
	}

	Last_accept_input_rule_num += 3

	return
}
