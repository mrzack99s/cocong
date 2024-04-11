package network

import (
	"github.com/mrzack99s/cocong/model/inmemory_model"
	"github.com/mrzack99s/cocong/utils"
	"github.com/mrzack99s/cocong/vars"
)

func DeleteAllowAccess(ss *inmemory_model.Session) (err error) {

	err = IPT.Delete("nat", "PREROUTING", "-s", ss.IPAddress, "-i", vars.Config.SecureInterface, "-j", "ACCEPT")
	if err != nil {
		return
	}

	err = IPT.Delete("filter", "FORWARD", "-s", ss.IPAddress, "-i", vars.Config.SecureInterface, "-o", vars.Config.EgressInterface, "-j", "ACCEPT")
	if err != nil {
		return
	}

	return
}

func DeleteAllowAccessBypass(ss *inmemory_model.Session) (err error) {
	interfaceIp, _ := utils.GetSecureInterfaceIpv4Addr()

	err = IPT.Delete("nat", "PREROUTING", "-s", ss.IPAddress, "-i", vars.Config.SecureInterface, "-j", "ACCEPT")
	if err != nil {
		return
	}

	err = IPT.Delete("filter", "FORWARD", "-s", ss.IPAddress, "-i", vars.Config.SecureInterface, "-o", vars.Config.EgressInterface, "-j", "ACCEPT")
	if err != nil {
		return
	}

	err = IPT.Delete("filter", "INPUT", "-s", ss.IPAddress, "-i", vars.Config.SecureInterface, "-p", "tcp", "--match", "multiport", "--dports", "443,8080,8443", "-d", interfaceIp, "-j", "DROP")
	if err != nil {
		return
	}

	return
}
