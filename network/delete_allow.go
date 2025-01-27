package network

import (
	"github.com/mrzack99s/cocong/types"
	"github.com/mrzack99s/cocong/utils"
	"github.com/mrzack99s/cocong/vars"
)

func DeleteAllowAccess(ss *types.SessionInfo) {

	IPT.Delete("nat", "PREROUTING", "-s", ss.IPAddress, "-i", vars.Config.SecureInterface, "-j", "ACCEPT")
	IPT.Delete("filter", "FORWARD", "-s", ss.IPAddress, "-i", vars.Config.SecureInterface, "-o", vars.Config.EgressInterface, "-j", "ACCEPT")

}

func DeleteAllowAccessBypass(ss *types.SessionInfo) (err error) {
	interfaceIp, _ := utils.GetSecureInterfaceIpv4Addr()

	err = IPT.Delete("nat", "PREROUTING", "-s", ss.IPAddress, "-i", vars.Config.SecureInterface, "-j", "ACCEPT")
	if err != nil {
		return
	}

	Last_insert_prerouting -= 1

	err = IPT.Delete("filter", "FORWARD", "-s", ss.IPAddress, "-i", vars.Config.SecureInterface, "-o", vars.Config.EgressInterface, "-j", "ACCEPT")
	if err != nil {
		return
	}

	Last_accept_forward_rule_num -= 1

	err = IPT.Delete("filter", "INPUT", "-p", "udp", "-s", ss.IPAddress, "-i", vars.Config.SecureInterface, "--dport", "53", "-d", interfaceIp, "-j", "ACCEPT")
	if err != nil {
		return
	}

	err = IPT.Delete("filter", "INPUT", "-p", "tcp", "-s", ss.IPAddress, "-i", vars.Config.SecureInterface, "--dport", "53", "-d", interfaceIp, "-j", "ACCEPT")
	if err != nil {
		return
	}

	err = IPT.Delete("filter", "INPUT", "-s", ss.IPAddress, "-i", vars.Config.SecureInterface, "-p", "tcp", "--match", "multiport", "--dports", "443,8080,8443", "-d", interfaceIp, "-j", "DROP")
	if err != nil {
		return
	}

	Last_accept_input_rule_num -= 3

	return
}
