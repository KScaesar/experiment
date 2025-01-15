#!/bin/bash
set -euo pipefail

sudo apt update

sudo timedatectl set-timezone Asia/Taipei

# fluent pre-installation Guide
# https://docs.fluentd.org/installation/install-by-deb#installing-fluent-package

sudo apt install ntp -y

Fluent_Limit_Config="
root soft nofile 65536
root hard nofile 65536
* soft nofile 65536
* hard nofile 65536
"
echo "$Fluent_Limit_Config" | sudo tee -a /etc/security/limits.conf > /dev/null

Fluent_Net_Config="
net.core.somaxconn = 1024
net.core.netdev_max_backlog = 5000
net.core.rmem_max = 16777216
net.core.wmem_max = 16777216
net.ipv4.tcp_wmem = 4096 12582912 16777216
net.ipv4.tcp_rmem = 4096 12582912 16777216
net.ipv4.tcp_max_syn_backlog = 8096
net.ipv4.tcp_slow_start_after_idle = 0
net.ipv4.tcp_tw_reuse = 1
net.ipv4.ip_local_port_range = 10240 65535
net.ipv4.ip_local_reserved_ports = 24224
"
echo "$Fluent_Net_Config" | sudo tee -a /etc/sysctl.conf > /dev/null

# Fluent-package-v5-lts
curl -fsSL https://toolbelt.treasuredata.com/sh/install-ubuntu-noble-fluent-package5-lts.sh | sh
sudo systemctl enable fluentd.service
sudo systemctl start fluentd.service

sudo usermod -a -G ec2-user _fluentd

# https://docs.fluentd.org/installation/post-installation-guide#connect-to-the-other-services
# https://docs.fluentd.org/deployment/plugin-management
# https://www.fluentd.org/plugins
sudo fluent-gem install fluent-plugin-record-modifier fluent-plugin-bigquery
