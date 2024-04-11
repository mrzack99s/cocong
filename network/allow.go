package network

import (
	"github.com/mrzack99s/cocong/model/inmemory_model"
	"github.com/mrzack99s/cocong/utils"
	"github.com/mrzack99s/cocong/vars"
)

func AllowAccess(ss *inmemory_model.Session) (err error) {

	err = IPT.Insert("nat", "PREROUTING", 4, "-s", ss.IPAddress, "-i", vars.Config.SecureInterface, "-j", "ACCEPT")
	if err != nil {
		return
	}

	err = IPT.Insert("filter", "FORWARD", last_accept_forward_rule_num, "-s", ss.IPAddress, "-i", vars.Config.SecureInterface, "-o", vars.Config.EgressInterface, "-j", "ACCEPT")
	if err != nil {
		return
	}

	return
}

func AllowAccessBypass(ss *inmemory_model.Session) (err error) {
	interfaceIp, _ := utils.GetSecureInterfaceIpv4Addr()

	err = IPT.Insert("nat", "PREROUTING", 4, "-s", ss.IPAddress, "-i", vars.Config.SecureInterface, "-j", "ACCEPT")
	if err != nil {
		return
	}

	err = IPT.Insert("filter", "FORWARD", last_accept_forward_rule_num, "-s", ss.IPAddress, "-i", vars.Config.SecureInterface, "-o", vars.Config.EgressInterface, "-j", "ACCEPT")
	if err != nil {
		return
	}

	err = IPT.Insert("filter", "INPUT", 1, "-s", ss.IPAddress, "-i", vars.Config.SecureInterface, "-p", "tcp", "--match", "multiport", "--dports", "443,8080,8443", "-d", interfaceIp, "-j", "DROP")
	if err != nil {
		return
	}
	// last_accept_input_rule_num += 1

	return
}
