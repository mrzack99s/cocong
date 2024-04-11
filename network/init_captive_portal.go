package network

import (
	"github.com/mrzack99s/cocong/model/inmemory_model"
	"github.com/mrzack99s/cocong/utils"
	"github.com/mrzack99s/cocong/vars"
)

func InitializeCaptivePortal() (err error) {
	interfaceIp, _ := utils.GetSecureInterfaceIpv4Addr()
	// Flush chain
	err = IPT.ClearAll()
	if err != nil {
		return
	}
	err = IPT.ClearChain("nat", "PREROUTING")
	if err != nil {
		return
	}

	// err = IPT.ClearChain("raw", "PREROUTING")
	// if err != nil {
	// 	return
	// }

	// Append Rules
	// err = IPT.AppendUnique("filter", "INPUT", "-p", "icmp", "-j", "DROP")
	// if err != nil {
	// 	return
	// }

	err = IPT.AppendUnique("filter", "INPUT", "-p", "tcp", "-j", "DROP")
	if err != nil {
		return
	}
	err = IPT.AppendUnique("filter", "INPUT", "-p", "udp", "-j", "DROP")
	if err != nil {
		return
	}
	err = IPT.InsertUnique("filter", "INPUT", 1, "-p", "udp", "-s", interfaceIp, "-j", "ACCEPT")
	if err != nil {
		return
	}

	err = IPT.InsertUnique("filter", "INPUT", 1, "-i", vars.Config.EgressInterface, "-m", "state", "--state", "ESTABLISHED,RELATED", "-j", "ACCEPT")
	if err != nil {
		return
	}

	err = IPT.InsertUnique("filter", "INPUT", 1, "-p", "tcp", "-s", "0.0.0.0/0", "--match", "multiport", "--dports", "3000,443", "-j", "ACCEPT")
	if err != nil {
		return
	}
	// last_accept_input_rule_num = 2

	if !vars.SYS_DEBUG {

		for _, net := range vars.Config.SSHAuthorizedNetworks {
			err = IPT.Insert("filter", "INPUT", 1, "-p", "tcp", "-s", net, "--match", "multiport", "--dports", "22,9090", "-j", "ACCEPT")
			if err != nil {
				return
			}

		}

		// last_accept_input_rule_num += len(vars.Config.SSHAuthorizedNetworks)
	} else {
		err = IPT.Insert("filter", "INPUT", 1, "-p", "tcp", "-s", "0.0.0.0/0", "--dport", "22", "-d", interfaceIp, "-j", "ACCEPT")
		if err != nil {
			return
		}
		// last_accept_input_rule_num += 1
	}

	for _, net := range vars.Config.AuthorizedNetworks {
		err = IPT.Insert("filter", "INPUT", 1, "-p", "tcp", "-s", net, "--match", "multiport", "--dports", "53,8080,8443", "-d", interfaceIp, "-j", "ACCEPT")
		if err != nil {
			return
		}
		err = IPT.Insert("filter", "INPUT", 1, "-p", "udp", "-s", net, "--dport", "53", "-d", interfaceIp, "-j", "ACCEPT")
		if err != nil {
			return
		}
	}

	// last_accept_input_rule_num += (len(vars.Config.AuthorizedNetworks) * 2)

	err = IPT.AppendUnique("filter", "FORWARD", "-i", vars.Config.SecureInterface, "-o", vars.Config.EgressInterface, "-j", "DROP")
	if err != nil {
		return
	}

	if vars.Config.ExternalPortalURL != "" {
		_, host, port, _ := utils.ParseURL(vars.Config.ExternalPortalURL)
		hostIp, _ := utils.ResolveIp(host)

		err = IPT.InsertUnique("filter", "FORWARD", 1, "-p", "tcp", "-i", vars.Config.SecureInterface, "--match", "multiport", "--dports", port, "-d", hostIp, "-j", "ACCEPT")
		if err != nil {
			return
		}
		last_accept_forward_rule_num += 1
	}

	err = IPT.InsertUnique("nat", "PREROUTING", 1, "-p", "tcp", "-i", vars.Config.SecureInterface, "--dport", "80", "-j", "DNAT", "--to-destination", interfaceIp+":8080")
	if err != nil {
		return
	}

	err = IPT.InsertUnique("nat", "PREROUTING", 1, "-p", "tcp", "-i", vars.Config.SecureInterface, "--dport", "443", "-j", "DNAT", "--to-destination", interfaceIp+":8443")
	if err != nil {
		return
	}

	err = IPT.InsertUnique("nat", "PREROUTING", 1, "-p", "tcp", "-i", vars.Config.SecureInterface, "-d", "1.1.1.1", "--dport", "80", "-j", "DNAT", "--to-destination", interfaceIp+":8080")
	if err != nil {
		return
	}

	err = IPT.InsertUnique("nat", "PREROUTING", 1, "-p", "tcp", "-i", vars.Config.SecureInterface, "-d", "1.1.1.1", "--dport", "443", "-j", "DNAT", "--to-destination", interfaceIp+":8443")
	if err != nil {
		return
	}

	err = IPT.InsertUnique("nat", "PREROUTING", 1, "-s", "0.0.0.0/0", "-p", "tcp", "-i", vars.Config.SecureInterface, "-d", interfaceIp, "-m", "tcp", "--dport", "443", "-j", "ACCEPT")
	if err != nil {
		return
	}

	// last_insert_init_prerouting = 2

	for _, net := range vars.Config.BypassNetworks {
		AllowAccessBypass(&inmemory_model.Session{
			IPAddress: net,
		})
	}

	// last_insert_init_prerouting += len(vars.Config.BypassNetworks)

	return
}
