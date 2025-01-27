package integration

type RadiusEndpointType struct {
	Hostname string `yaml:"hostname" json:"hostname"`
	Port     uint64 `yaml:"port" json:"port"`
	Secret   string `yaml:"secret" json:"secret"`
}

type LDAPEndpointType struct {
	Hostname  string `yaml:"hostname" json:"hostname"`
	Port      uint64 `yaml:"port" json:"port"`
	TLSEnable bool   `yaml:"tls_enable" json:"tls_enable"`
	Domain    string `yaml:"domain" json:"domain"`
	PoolSize  uint   `yaml:"pool_size" json:"pool_size"`
	// instance  *ldap.Conn
	pool *LDAPConnectionPool
}
