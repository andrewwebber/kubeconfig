# kubeconfig

kubeconfig generates cloud-config files for testing a three node Kubernetes-CoreOS cluster with VMware Fusion 6. The cloud-config files are generated assume machines with the following specs:

* 2 network interfaces (ens33, ens34)

Tested on VMware Fusion.

## Usage

config.yml
```
token: 35888e98a8a633296d3b53b2bf9a87fc
dns: 192.168.12.1
gateway: 192.168.12.1
master_ip: 192.168.12.10
node1_ip: 192.168.12.11
node2_ip: 192.168.12.12
sshkey: ssh-rsa AAAAB3NzaC1yc2...
```

```
kubeconfig -c config.yml
```

Should result in the following cloud-config files:

* master.yml
* node1.yml
* node2.yml
