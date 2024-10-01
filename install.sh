#!/bin/bash
# !Important change this below
DOMAIN_NAME="cocong.local"
PRIVATE_IP="Changeme"
OS_BASED="rhel" #(rhel/debian)
###############################

command_exists() {
	command -v "$@" > /dev/null 2>&1
}

cat >&2 <<-'EOF'
First please setup your network interface and any routing before install

EOF

if [[ "${OS_BASED}" != "rhel" && "${OS_BASED}" != "debian" ]]; then
	cat >&2 <<-'EOF'
		OS_BASED support only rhel and debian
	EOF
	exit 1
fi

if ! command_exists node; then
		cat >&2 <<-'EOF'
			NodeJS is required, version 20.x

			RedHat Based:
				dnf module install nodejs:20

			Debian Based:
				curl -fsSL https://deb.nodesource.com/setup_lts.x | bash - &&\
				apt-get install -y nodejs

		EOF
		exit 1
fi

if ! command_exists conntrack; then
		cat >&2 <<-'EOF'
			Conntrack-Tools is required
			Please install Conntrack-Tools

			RedHat Based:
				yum install -y conntrack-tools
				--------- Or ---------
				dnf install -y conntrack-tools

			Debian Based:
				apt-get install conntrack-tools


		EOF
		exit 1
fi

if [ -z "$(ldconfig -p | grep libpcap)" ]; then
		cat >&2 <<-'EOF'
			libpcap is required
			Please install libpcap

			RedHat Based:
				yum install -y libpcap
				--------- Or ---------
				dnf install -y libpcap

			Debian Based:
				apt-get install libpcap0.8


		EOF
		exit 1
fi


user="$(id -un 2>/dev/null || true)"
if [ "$user" != 'root' ]; then
		cat >&2 <<-'EOF'
		Error: this installer needs the ability to run commands as root.
		EOF
		exit 1
fi

if [[ -n $SUDO_USER ]]; then 
    cat >&2 <<-'EOF'
	Error: not support sudo command. required root user.
	EOF
	exit 1
fi

# Prepare workspace
mkdir -p /etc/cocong
mkdir -p /etc/cocong/certs
mkdir -p /etc/cocong/dns
mkdir -p /usr/share/cocong
mkdir -p /var/log/cocong

# Build and download
cd cocong-admin
echo "NEXT_PUBLIC_API_URL=\"https://${DOMAIN_NAME}\"" > .env.production

npm install --production
npm run build

cp -r .next/static .next/standalone/.next/
cp -r .next/standalone /usr/share/cocong/cocong-admin

cd ../

# Install package

cp cocong_${OS_BASED} /usr/local/bin/cocong
chmod +x /usr/local/bin/cocong

cp coredns /usr/local/bin/coredns
chmod +x /usr/local/bin/coredns

cp tcdel /usr/local/bin/tcdel
chmod +x /usr/local/bin/tcdel

cp tcset /usr/local/bin/tcset
chmod +x /usr/local/bin/tcset

cp tcshow /usr/local/bin/tcshow
chmod +x /usr/local/bin/tcshow

# Gencert
cocong gencert ${DOMAIN_NAME}
cp -r ./certs/* /etc/cocong/certs

# Copy template
cp -r templates /usr/share/cocong

# Default system config
cp cocong.yaml /etc/cocong/cocong.yaml
cat <<EOF > /etc/sysctl.d/99-cocong.conf
net.ipv4.ip_forward=1
net.ipv6.conf.all.disable_ipv6=1

net.netfilter.nf_conntrack_generic_timeout=300
net.netfilter.nf_conntrack_icmp_timeout=15
net.netfilter.nf_conntrack_tcp_timeout_established=86400
net.netfilter.nf_conntrack_tcp_timeout_close = 10
net.netfilter.nf_conntrack_tcp_timeout_close_wait = 30
net.netfilter.nf_conntrack_tcp_timeout_fin_wait = 30
net.netfilter.nf_conntrack_tcp_timeout_syn_recv = 30
net.netfilter.nf_conntrack_tcp_timeout_syn_sent = 60
net.netfilter.nf_conntrack_tcp_timeout_time_wait = 30
net.netfilter.nf_conntrack_udp_timeout_stream = 30
net.core.rmem_max=2097152
EOF



# Services
cat <<EOF > /etc/systemd/system/cocong.service
[Unit]
Description=Run the CoCoNG service
After=network.target

[Service]
User=root
Group=root
# Environment="COCONG_API_KEY_HASHED=<SHA512>"
WorkingDirectory=/etc/cocong
ExecStart=cocong run
Restart=on-failure

[Install]
WantedBy=multi-user.target

EOF

cat <<EOF > /etc/systemd/system/cocong-admin.service
[Unit]
Description=Run the CoCoNG-Admin-UI service
After=network.target

[Service]
User=root
Group=root
Environment="HOSTNAME=0.0.0.0"
WorkingDirectory=/usr/share/cocong/cocong-admin
ExecStart=node server.js
Restart=on-failure

[Install]
WantedBy=multi-user.target

EOF


cat <<EOF > /etc/systemd/system/coredns.service
[Unit]
Description=Run the CoreDNS service
After=network.target

[Service]
User=root
Group=root
WorkingDirectory=/etc/cocong/dns
ExecStart=coredns -dns.port=53 -conf=/etc/cocong/dns/Corefile
Restart=on-failure

[Install]
WantedBy=multi-user.target

EOF

cat <<EOF > /etc/cocong/dns/Corefile
${DOMAIN_NAME} {
	file /etc/cocong/dns/db.${DOMAIN_NAME}
}

logout.net {
	file /etc/cocong/dns/db.logout.net
}

login.net {
	file /etc/cocong/dns/db.login.net
}

. {
	forward . 8.8.8.8
	log
	cache
}

EOF



cat <<EOF > /etc/cocong/dns/db.${DOMAIN_NAME}
\$ORIGIN ${DOMAIN_NAME}.
@	3600 IN	SOA sns.dns._cocong.local. noc.dns._cocong.local. 2017042745 7200 3600 1209600 3600
	3600 IN NS ns._cocong.local.
	
	IN A     ${PRIVATE_IP}

EOF



cat <<EOF > /etc/cocong/dns/db.login.net
\$ORIGIN login.net.
@	3600 IN	SOA sns.dns._cocong.local. noc.dns._cocong.local. 2017042745 7200 3600 1209600 3600
	3600 IN NS ns._cocong.local.

        IN A     ${PRIVATE_IP}

EOF



cat <<EOF > /etc/cocong/dns/db.logout.net
\$ORIGIN logout.net.
@	3600 IN	SOA sns.dns._cocong.local. noc.dns._cocong.local. 2017042745 7200 3600 1209600 3600
	3600 IN NS ns._cocong.local.

        IN A     ${PRIVATE_IP}

EOF



systemctl daemon-reload
systemctl enable cocong
systemctl enable cocong-admin
systemctl enable coredns



cat >&2 <<-'EOF'
	If your os system is RHEL based, please do this below

	firewall-cmd --zone=public --permanent --add-port=53/udp
	firewall-cmd --zone=public --permanent --add-port=443/tcp
	firewall-cmd --zone=public --permanent --add-port=3000/tcp
	firewall-cmd --zone=public --permanent --add-port=8080/tcp
	firewall-cmd --zone=public --permanent --add-port=8443/tcp
	firewall-cmd --reload

	Recommended
	########################################################################################
	Use your trusted certificate for your domain

	-> sysctl settings

	# Add below config to /etc/sysctl.d/99-cocong.conf

	(How to calculate buckets and max conntrack)
	max conntrack = ( Memory size in MB * 1048576)/16384/2
	buckets conntrack = max conntrack/4

	net.netfilter.nf_conntrack_max=(your max conntrack from your calculate)
	net.netfilter.nf_conntrack_buckets=(your buckets conntrack from your calculate)

	########################################################################################

	Please reboot server

EOF