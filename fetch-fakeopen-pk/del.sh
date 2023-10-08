function delPortHoppingNat() {
	# $1 portHoppingStart
	# $2 portHoppingEnd
	# $3 portHoppingTarget
	if systemctl status firewalld 2>/dev/null | grep -q "active (running)"; then
		firewall-cmd --permanent --remove-forward-port=port=47000-48000:proto=udp:toport=37740
		firewall-cmd --reload
	else
		iptables -t nat -F PREROUTING  2>/dev/null
		ip6tables -t nat -F PREROUTING  2>/dev/null
		if systemctl status netfilter-persistent 2>/dev/null | grep -q "active (exited)"; then
			netfilter-persistent save 2> /dev/null
		fi

	fi
}
delPortHoppingNat