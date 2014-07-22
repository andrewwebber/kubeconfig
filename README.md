# kubeconfig

kubeconfig generates cloud-config files for testing a three node Kubernetes-CoreOS cluster with VMware Fusion 6. The cloud-config files are generated assume machines with the following specs:

* 2 network interfaces (ens33, ens34)

Tested on VMware Fusion.

## Usage

```
Usage of ./kubeconfig:
  -c="kubernetes.yml": config file to use
  -iso=false: generate config-drive iso images
```

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

### Create cloud-config files

```
kubeconfig -c config.yml
```
-

```
master.yml node1.yml node2.yml
```

### Create config-drive iso images

The following command will connect to a remote ISO creation service.

```
kubeconfig -c config.yml -iso
```

-

```
master.iso  master.yml  node1.iso   node1.yml   node2.iso   node2.yml
```
