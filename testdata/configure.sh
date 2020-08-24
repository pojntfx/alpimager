#!/bin/sh
setup-timezone -z UTC

cat <<-EOF >/etc/network/interfaces
iface lo inet loopback
iface eth0 inet dhcp
EOF

cat <<EOF >/etc/motd
Welcome to Felicitas Pojtinger's Alpine Linux Distribution!
EOF

mkdir -m 700 -p /root/.ssh
wget -O - https://github.com/pojntfx.keys | tee /root/.ssh/authorized_keys
chmod 600 /root/.ssh/authorized_keys

ln -s networking /etc/init.d/net.lo
ln -s networking /etc/init.d/net.eth0

rc-update add sshd default
rc-update add net.eth0 default
rc-update add net.lo boot

echo 'AllowTcpForwarding yes' >>/etc/ssh/sshd_config
