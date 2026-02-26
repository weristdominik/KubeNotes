# Kubernetes Flatcar on libvirt/kvm

> [!IMPORTANT]  
> At the time of creation 02/2026 Terraform Provider [dmacvicar/libvirt v0.9.2 (latest)]() has unsolved problems - [flatcar sh tool](https://www.flatcar.org/docs/latest/installing/vms/libvirt/#creating-the-domain) is more stable

Install butane

```bash
wget https://github.com/coreos/butane/releases/download/v0.26.0/butane-x86_64-unknown-linux-gnu
mv -v butane-x86_64-unknown-linux-gnu butane
sudo chmod +x butane 
```

Get [Offical flatcar ](https://www.flatcar.org/docs/latest/installing/vms/libvirt/#choosing-a-channel)
(base) img

```bash
mkdir -p /var/lib/libvirt/images/flatcar-linux
cd /var/lib/libvirt/images/flatcar-linux
wget https://stable.release.flatcar-linux.net/amd64-usr/current/flatcar_production_qemu_image.img
```

Controlplane and Workernode Network: [Static IP](https://www.flatcar.org/docs/latest/installing/vms/libvirt/#static-ip)

&nbsp;

## Install ControlPlane

Create ignition file

`vi control01.yaml`

<details>

<summary>control01.yaml</summary>

```yaml
version: 1.0.0
variant: flatcar
passwd:
  users:
    - name: core
      ssh_authorized_keys:
        - "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDSMNacUETwYWSvL8XNA54QIErHJfvolD4FqG18hdkpVInubRKZVb43erYY0UlH3xQoUBKkxj1+ssAP3IXWwwsUSc7VNXobZHjuo21+ikx96UGF5A8/JQ7Wi5yhzLoqq5dMfxPfPCWPiqbagiB8hJrorjv2q3peWxTBXw4od0pbzKSRODhUdTAuDEFFiL+fcTwcdE12dUQPKWi5U3dhB6QYFSMnlYkwxpAm1Ol7KgyzIqu8D2YpKdXVZ/xdpX93h3xoR/7nKmx8P9KuHi076w+hebZUzFOoLEDDMW0jiw2GqO/flHAo3myhlv3ftHPQcDim7BC13e+QbyjfBygux42mtVQj5HzOCgS8gl3YYkDlHSzaLg3SOF1K5nrvd3FVRuSQ9OZm5aSsTukd+5mA8t7nSLKIi4hiS4/iSeskO4pH9w5/+lX09RZpZzEKHj6sEmEcVrlADuEnsmjimrJdvtZsQWYb2GZy/Ggua4FrcfnGuDHzCnNXKUmkYoDJqiQn5d7cZcolI+BULMMjaVdZrMvgwbDqI1nxvBD05JrvPdZaVMq1pO8Xv3G9i01vfnY6l1u70DeUnHPH15KOdaHT5VnOx8H0I/bt1dqw8A63xBr9p3GrKotrcGxXcilx5kUaAiiRseaV5GVGVtFnKunDkVOojE3ioNCQ4rX7OYuWEY4yUw== netadmin@brick"
storage:
  links:
    - target: /opt/extensions/kubernetes/kubernetes-v1.33.2-x86-64.raw
      path: /etc/extensions/kubernetes.raw
      hard: false
  files:
    - path: /etc/sysupdate.kubernetes.d/kubernetes-v1.33.conf
      contents:
        source: https://extensions.flatcar.org/extensions/kubernetes/kubernetes-v1.33.conf
    - path: /etc/sysupdate.d/noop.conf
      contents:
        source: https://extensions.flatcar.org/extensions/noop.conf
    - path: /opt/extensions/kubernetes/kubernetes-v1.33.2-x86-64.raw
      contents:
        source: https://extensions.flatcar.org/extensions/kubernetes-v1.33.2-x86-64.raw
    - path: /etc/hostname
      contents:
        inline: "control01"
systemd:
  units:
    - name: systemd-sysupdate.timer
      enabled: true
    - name: systemd-sysupdate.service
      dropins:
        - name: kubernetes.conf
          contents: |
            [Service]
            ExecStartPre=/usr/bin/sh -c "readlink --canonicalize /etc/extensions/kubernetes.raw > /tmp/kubernetes"
            ExecStartPre=/usr/lib/systemd/systemd-sysupdate -C kubernetes update
            ExecStartPost=/usr/bin/sh -c "readlink --canonicalize /etc/extensions/kubernetes.raw > /tmp/kubernetes-new"
            ExecStartPost=/usr/bin/sh -c "if ! cmp --silent /tmp/kubernetes /tmp/kubernetes-new; then touch /run/reboot-required; fi"
    - name: locksmithd.service
      # NOTE: To coordinate the node reboot in this context, we recommend to use Kured.
      mask: true
    - name: kubeadm.service
      enabled: true
      contents: |
        [Unit]
        Description=Kubeadm service
        Requires=containerd.service
        After=containerd.service
        After=network-online.target
        ConditionPathExists=!/etc/kubernetes/kubelet.conf
        [Service]
        ExecStartPre=/usr/bin/kubeadm init
        ExecStartPre=/usr/bin/mkdir /home/core/.kube
        ExecStartPre=/usr/bin/cp /etc/kubernetes/admin.conf /home/core/.kube/config
        ExecStart=/usr/bin/chown -R core:core /home/core/.kube
        [Install]
        WantedBy=multi-user.target
```

### Key Takeaways

- The kubernetes-v1.33.conf is a configuration file that tells the system how to handle updates for Kubernetes.

- The kubernetes-v1.33.2-x86-64.raw is the actual binary file that contains the Kubernetes version.

- *To upgrade Kubernetes* it is enough to edit kubernetes-v1.33.conf since the binaries will be automatically downloaded into /opt/extensions/kubernetes/

</details>

&nbsp;

Convert YAML to ignition file

```bash
./butane control01.yaml > control01.ign
```

Create VM control01

```bash
bash create-controlplane01.sh
# After boot you will see console
ip a
# ssh core@<ip> (authentification will be done by your ssh key)
```

&nbsp;

## Install Workernode(s)

Create as many Workernodes as needed. Each Workernode will have its own __*ign__ file.

`vi worker01.yaml`

<details>

<summary>worker01.yaml</summary>

```yaml
version: 1.0.0
variant: flatcar
passwd:
  users:
    - name: core
      ssh_authorized_keys:
        - "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDSMNacUETwYWSvL8XNA54QIErHJfvolD4FqG18hdkpVInubRKZVb43erYY0UlH3xQoUBKkxj1+ssAP3IXWwwsUSc7VNXobZHjuo21+ikx96UGF5A8/JQ7Wi5yhzLoqq5dMfxPfPCWPiqbagiB8hJrorjv2q3peWxTBXw4od0pbzKSRODhUdTAuDEFFiL+fcTwcdE12dUQPKWi5U3dhB6QYFSMnlYkwxpAm1Ol7KgyzIqu8D2YpKdXVZ/xdpX93h3xoR/7nKmx8P9KuHi076w+hebZUzFOoLEDDMW0jiw2GqO/flHAo3myhlv3ftHPQcDim7BC13e+QbyjfBygux42mtVQj5HzOCgS8gl3YYkDlHSzaLg3SOF1K5nrvd3FVRuSQ9OZm5aSsTukd+5mA8t7nSLKIi4hiS4/iSeskO4pH9w5/+lX09RZpZzEKHj6sEmEcVrlADuEnsmjimrJdvtZsQWYb2GZy/Ggua4FrcfnGuDHzCnNXKUmkYoDJqiQn5d7cZcolI+BULMMjaVdZrMvgwbDqI1nxvBD05JrvPdZaVMq1pO8Xv3G9i01vfnY6l1u70DeUnHPH15KOdaHT5VnOx8H0I/bt1dqw8A63xBr9p3GrKotrcGxXcilx5kUaAiiRseaV5GVGVtFnKunDkVOojE3ioNCQ4rX7OYuWEY4yUw== netadmin@brick"
storage:
  links:
    - target: /opt/extensions/kubernetes/kubernetes-v1.33.2-x86-64.raw
      path: /etc/extensions/kubernetes.raw
      hard: false
  files:
    - path: /etc/sysupdate.kubernetes.d/kubernetes-v1.33.conf
      contents:
        source: https://extensions.flatcar.org/extensions/kubernetes/kubernetes-v1.33.conf
    - path: /etc/sysupdate.d/noop.conf
      contents:
        source: https://extensions.flatcar.org/extensions/noop.conf
    - path: /opt/extensions/kubernetes/kubernetes-v1.33.2-x86-64.raw
      contents:
        source: https://extensions.flatcar.org/extensions/kubernetes-v1.33.2-x86-64.raw
    - path: /etc/hostname
      contents:
        inline: "worker01"
systemd:
  units:
    - name: systemd-sysupdate.timer
      enabled: true
    - name: systemd-sysupdate.service
      dropins:
        - name: kubernetes.conf
          contents: |
            [Service]
            ExecStartPre=/usr/bin/sh -c "readlink --canonicalize /etc/extensions/kubernetes.raw > /tmp/kubernetes"
            ExecStartPre=/usr/lib/systemd/systemd-sysupdate -C kubernetes update
            ExecStartPost=/usr/bin/sh -c "readlink --canonicalize /etc/extensions/kubernetes.raw > /tmp/kubernetes-new"
            ExecStartPost=/usr/bin/sh -c "if ! cmp --silent /tmp/kubernetes /tmp/kubernetes-new; then touch /run/reboot-required; fi"
    - name: locksmithd.service
      mask: true
    - name: kubeadm.service
      enabled: false
      # Kubeadm init removed for worker nodes
```

</details>

&nbsp;

Convert YAML to ignition file

```bash
./butane worker01.yaml > worker01.ign
```

Create VM worker01

```bash
bash create-worker01.sh
# After boot you will see console
ip a
# ssh core@<ip> (authentification will be done by your ssh key)
```

&nbsp;

## Join Cluster

```bash
# ControlPlane
ssh core@control01
$ sudo kubeadm token create --print-join-command

# Workernode
ssh core@worker0X
$ <paste join command>
```

&nbsp;

## Update Kubernetes Version

> Note before running updates it is necessary to [install an CNI](#install-cni) so Controlplanes and Workernodes reach __Ready__ status.

<details>

<summary>Failed pre-flight checks no CNI installed</summary>

Error Message:

```bash
core@flatcar-controlplane ~ $ sudo kubeadm upgrade plan
[preflight] Running pre-flight checks.
[upgrade/config] Reading configuration from the "kubeadm-config" ConfigMap in namespace "kube-system"...
[upgrade/config] Use 'kubeadm init phase upload-config --config your-config-file' to re-upload it.
[upgrade] Running cluster health checks
[upgrade/health] FATAL: [preflight] Some fatal errors occurred:
	[ERROR CreateJob]: Job "upgrade-health-check-1771588625883" in the namespace "kube-system" did not complete in 15s: no condition of type Complete
	[ERROR ControlPlaneNodesReady]: there are NotReady control-planes in the cluster: [flatcar-controlplane]
[preflight] If you know what you are doing, you can make a check non-fatal with `--ignore-preflight-errors=...`
To see the stack trace of this error execute with --v=5 or higher
core@flatcar-controlplane ~ $ kubectl get nodes
NAME                   STATUS     ROLES           AGE   VERSION
flatcar-controlplane   NotReady   control-plane   24m   v1.33.2
```

</details>

&nbsp;

Example shown on __control01__ but need to be done on all Hosts. Following the [original documentation](https://v1-34.docs.kubernetes.io/docs/tasks/administer-cluster/kubeadm/kubeadm-upgrade/)

Available Versions checkout: [flatcar/sysext-bakery release page](https://github.com/flatcar/sysext-bakery/releases/tag/kubernetes)

### v1.33 to v1.34

Update Version in */etc/sysupdate.kubernetes.d/kubernetes-v1.33.conf*

```bash
# ssh core@control01
sudo vi /etc/sysupdate.kubernetes.d/kubernetes-v1.33.conf
# ==> MatchPattern=kubernetes-v1.34.@v-%a.raw
```

Run update process

```bash
sudo systemctl restart systemd-sysupdate.timer
sudo systemctl start systemd-sysupdate.service
# file /run/reboot-required created
sudo reboot
```

After reboot

```bash
# core@control01 ~ $ kubectl version
# Client Version: v1.34.4
# Kustomize Version: v5.7.1
# Server Version: v1.33.8

sudo kubeadm upgrade plan
# Your output meight differ
sudo kubeadm upgrade apply v1.34.4
# done
```

<details>

<summary>Do I need to download new binaries?</summary>

By updating /etc/sysupdate.kubernetes.d/kubernetes-v1.33.conf

`MatchPattern=kubernetes-v1.34.@v-%a.raw`

Flatcar will automatically download new needed binaries and will update the /etc/extensions/kubernetes.raw systemlink

```bash
core@control01 ~ $ ls -l /etc/extensions/kubernetes.raw
lrwxrwxrwx. 1 root root 61 Feb 20 13:04 /etc/extensions/kubernetes.raw -> ../../opt/extensions/kubernetes/kubernetes-v1.34.4-x86-64.raw
```

</details>

&nbsp;

## Install CNI 

Cilium:

```bash
kubectl \
  apply -f https://raw.githubusercontent.com/projectcalico/calico/v3.24.1/manifests/calico.yaml

$ kubectl get nodes
NAME                   STATUS   ROLES           AGE   VERSION
flatcar-controlplane   Ready    control-plane   45m   v1.33.2
```
