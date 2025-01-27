package services

import (
	"os"
	"os/exec"

	"github.com/mrzack99s/cocong/constants"
	"github.com/mrzack99s/cocong/network"
	"github.com/mrzack99s/cocong/types"
	"github.com/mrzack99s/cocong/utils"
	"github.com/mrzack99s/cocong/vars"
	"gopkg.in/yaml.v3"
)

func ConfigUpdate(updatedConfig types.ConfigType) {

	interfaceIp, _ := utils.GetSecureInterfaceIpv4Addr()

	if ok, _ := utils.CheckDifference(vars.Config.AuthorizedNetworks, updatedConfig.AuthorizedNetworks); ok {
		addedAuthorizedNetworks, deletedAuthorizedNetworks, _ := utils.GetDifferenceSlice(network.AuthorizedNetworks, updatedConfig.AuthorizedNetworks)
		for _, net := range addedAuthorizedNetworks {
			network.IPT.Insert("filter", "INPUT", network.Last_accept_input_rule_num, "-p", "tcp", "-s", net, "--match", "multiport", "--dports", "53,8080,8443", "-d", interfaceIp, "-j", "ACCEPT")
			network.IPT.Insert("filter", "INPUT", network.Last_accept_input_rule_num, "-p", "udp", "-s", net, "--dport", "53", "-d", interfaceIp, "-j", "ACCEPT")

			network.AuthorizedNetworks = append(network.AuthorizedNetworks, net)
		}

		for _, net := range deletedAuthorizedNetworks {
			network.IPT.Delete("filter", "INPUT", "-p", "tcp", "-s", net, "--match", "multiport", "--dports", "53,8080,8443", "-d", interfaceIp, "-j", "ACCEPT")
			network.IPT.Delete("filter", "INPUT", "-p", "udp", "-s", net, "--dport", "53", "-d", interfaceIp, "-j", "ACCEPT")

			network.AuthorizedNetworks, _ = utils.DeleteSliceElement(network.AuthorizedNetworks, net)
		}

		vars.Config.AuthorizedNetworks = updatedConfig.AuthorizedNetworks
	}

	if ok, _ := utils.CheckDifference(vars.Config.BypassNetworks, updatedConfig.BypassNetworks); ok {
		addedBypassNetworks, deletedBypassNetworks, _ := utils.GetDifferenceSlice(network.BypassedNetworks, updatedConfig.BypassNetworks)
		for _, net := range addedBypassNetworks {
			network.AllowAccessBypass(&types.SessionInfo{
				IPAddress: net,
			})

			network.BypassedNetworks = append(network.BypassedNetworks, net)
		}

		for _, net := range deletedBypassNetworks {
			network.DeleteAllowAccessBypass(&types.SessionInfo{
				IPAddress: net,
			})

			network.BypassedNetworks, _ = utils.DeleteSliceElement(network.BypassedNetworks, net)
		}

		vars.Config.BypassNetworks = updatedConfig.BypassNetworks
	}

	if ok, _ := utils.CheckDifference(vars.Config.SSHAuthorizedNetworks, updatedConfig.SSHAuthorizedNetworks); ok {
		addedSSHAuthorizedNetworks, deletedSSHAuthorizedNetworks, _ := utils.GetDifferenceSlice(network.SSHAuthorizedNetworks, updatedConfig.SSHAuthorizedNetworks)
		for _, net := range addedSSHAuthorizedNetworks {
			network.IPT.Insert("filter", "INPUT", network.Last_accept_input_rule_num, "-p", "tcp", "-s", net, "--match", "multiport", "--dports", "22,9090", "-j", "ACCEPT")
			network.SSHAuthorizedNetworks = append(network.SSHAuthorizedNetworks, net)
		}

		for _, net := range deletedSSHAuthorizedNetworks {
			network.IPT.Delete("filter", "INPUT", "-p", "tcp", "-s", net, "-d", interfaceIp, "--match", "multiport", "--dports", "22,9090", "-j", "ACCEPT")
			network.SSHAuthorizedNetworks, _ = utils.DeleteSliceElement(network.SSHAuthorizedNetworks, net)
		}

		vars.Config.SSHAuthorizedNetworks = updatedConfig.SSHAuthorizedNetworks
	}

	if ok, _ := utils.CheckDifference(vars.Config.LDAP, updatedConfig.LDAP); ok {
		vars.Config.LDAP = updatedConfig.LDAP
	}

	if ok, _ := utils.CheckDifference(vars.Config.Radius, updatedConfig.Radius); ok {
		vars.Config.Radius = updatedConfig.Radius
	}

	if ok, _ := utils.CheckDifference(vars.Config.SessionIdle, updatedConfig.SessionIdle); ok {
		vars.Config.SessionIdle = updatedConfig.SessionIdle
	}

	if ok, _ := utils.CheckDifference(vars.Config.MaxConcurrentSession, updatedConfig.MaxConcurrentSession); ok {
		vars.Config.MaxConcurrentSession = updatedConfig.MaxConcurrentSession
	}

	if ok, _ := utils.CheckDifference(vars.Config.RedirectURL, updatedConfig.RedirectURL); ok {
		vars.Config.RedirectURL = updatedConfig.RedirectURL
	}

}

func WriteConfigUpdate() {

	yamlBytes, _ := yaml.Marshal(vars.Config)

	if !vars.SYS_DEBUG {
		os.WriteFile(constants.CONFIG_DIR+"/cocong.yaml", yamlBytes, 0644)

	} else {
		os.WriteFile("./cocong.yaml", yamlBytes, 0644)
	}

}

func RestartService() {
	exec.Command("systemctl", "restart", "cocong").Run()
}

func RestartDNSServer() {
	exec.Command("systemctl", "restart", "coredns").Run()
}
