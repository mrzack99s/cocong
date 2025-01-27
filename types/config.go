package types

import "github.com/mrzack99s/cocong/integration"

type ConfigType struct {
	LocalNetworkName     string   `yaml:"local_network_name" json:"local_network_name"`
	ExternalPortalURL    string   `yaml:"external_portal_url" json:"external_portal_url"`
	EgressInterface      string   `yaml:"egress_interface" json:"egress_interface"`
	SecureInterface      string   `yaml:"secure_interface" json:"secure_interface"`
	SessionIdle          uint64   `yaml:"session_idle" json:"session_idle"`
	MaxConcurrentSession uint64   `yaml:"max_concurrent_session" json:"max_concurrent_session"`
	AuthorizedNetworks   []string `yaml:"authorized_networks" json:"authorized_networks"`

	BypassNetworks         []string `yaml:"bypass_networks" json:"bypass_networks"`
	SSHAuthorizedNetworks  []string `yaml:"ssh_authorized_networks" json:"ssh_authorized_networks"`
	RedirectURL            string   `yaml:"redirect_url" json:"redirect_url"`
	CustomDBDir            string   `yaml:"custom_db_dir" json:"custom_db_dir"`
	DBCacheSize            int      `yaml:"db_cache_size" json:"db_cache_size"`
	LoginFail2Ban          bool     `yaml:"login_fail2ban" json:"login_fail2ban"`
	DomainName             string   `yaml:"domain_name" json:"domain_name"`
	DisabledNetworkCapture bool     `yaml:"disabled_network_capture" json:"disabled_network_capture"`

	Radius *integration.RadiusEndpointType `yaml:"radius" json:"radius"`
	LDAP   *integration.LDAPEndpointType   `yaml:"ldap" json:"ldap"`

	TimeZone string `yaml:"timezone" json:"timezone"`
}

type EndpointType struct {
	Hostname string `yaml:"hostname" json:"hostname"`
	Port     uint64 `yaml:"port" json:"port"`
}
