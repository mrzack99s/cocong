#!/bin/bash
# !Important change this below
DOMAIN_NAME="cocong.local"
PRIVATE_IP="Changeme"
OS_BASED="rhel" #(rhel/debian)
###############################

command_exists() {
	command -v "$@" > /dev/null 2>&1
}


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
rm -rf /etc/cocong
rm -rf /etc/cocong/certs
rm -rf /etc/cocong/dns
rm -rf /usr/share/cocong
rm -rf /var/log/cocong

rm -f /usr/local/bin/cocong
rm -f /usr/local/bin/coredns
rm -f /usr/local/bin/tcdel
rm -f /usr/local/bin/tcset
rm -f /usr/local/bin/tcshow

systemctl stop cocong
systemctl stop cocong-admin
systemctl stop coredns
systemctl disable cocong
systemctl disable cocong-admin
systemctl disable coredns

rm -f /etc/sysctl.d/99-cocong.conf
rm -f /etc/systemd/system/cocong.service
rm -f /etc/systemd/system/cocong-admin.service
rm -f /etc/systemd/system/coredns.service

systemctl daemon-reload