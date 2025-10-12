## How to install k8s 
Links:
- https://kubernetes.io/docs/setup/production-environment/container-runtimes/
- https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/install-kubeadm/
- https://kubernetes.io/docs/concepts/cluster-administration/addons/
- https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/troubleshooting-kubeadm/#coredns-is-stuck-in-the-pending-state
- https://github.com/flannel-io/flannel?tab=readme-ov-file#deploying-flannel-with-kubectl

### root user
```bash
sudo -i

# prerequisites
swapoff -a
sed -i '/ swap / s/^/#/' /etc/fstab

cat <<EOF | tee /etc/modules-load.d/k8s.conf
overlay
br_netfilter
EOF
modprobe overlay
modprobe br_netfilter
# Flannel requires the br_netfilter module to start and from version 1.30 kubeadm doesn't check if the module is installed and Flannel will not rightly start in case the module is missing.

cat <<EOF | tee /etc/sysctl.d/99-k8s.conf
net.bridge.bridge-nf-call-ip6tables = 1
net.bridge.bridge-nf-call-iptables  = 1
net.ipv4.ip_forward                 = 1
EOF
sysctl --system

# CRI install
apt install -y curl apt-transport-https ca-certificates software-properties-common
apt install -y containerd
mkdir -p /etc/containerd

containerd config default | tee /etc/containerd/config.toml
sed -i 's/SystemdCgroup = false/SystemdCgroup = true/' /etc/containerd/config.toml
systemctl restart containerd
systemctl enable containerd
systemctl status containerd

# k8s components install
echo 'deb [signed-by=/etc/apt/keyrings/kubernetes-apt-keyring.gpg] https://pkgs.k8s.io/core:/stable:/v1.34/deb/ /' | sudo tee /etc/apt/sources.list.d/kubernetes.list
curl -fsSL https://pkgs.k8s.io/core:/stable:/v1.34/deb/Release.key | sudo gpg --dearmor -o /etc/apt/keyrings/kubernetes-apt-keyring.gpg
apt update
apt install -y kubeadm kubelet kubectl
kubeadm version

kubeadm init --pod-network-cidr=10.0.20.0/24 --ignore-preflight-errors=all
```

### non-root user
```bash
# without root !!!

mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config

# Deploying Flannel with kubectl
# If you use custom podCIDR (not 10.244.0.0/16) you first need to download the above manifest and modify the network to match your one.
wget https://github.com/flannel-io/flannel/releases/latest/download/kube-flannel.yml
nano kube-flannel.yml
kubectl apply -f kube-flannel.yml

# untaint to run payload on master node
kubectl taint nodes --all node-role.kubernetes.io/control-plane-

# now you can use it :)
kubectl get nodes -owide
kubectl get pods -A -owide
kubectl get events -w -A
watch "kubectl get pods -A -owide"

kubectl create deployment nginx --image nginx
kubectl port-forward svc/nginx-service 8080:80 and then curl http://localhost:8080
curl localhost:8080
```


#### if you run in LXC - you may later need
```bash
conntrack -L
cat /proc/sys/net/netfilter/nf_conntrack_max
sudo udevadm control --reload-rules
sudo nano /etc/udev/rules.d/99-kmsg-symlink.rules
```
#### nice to have if not already present
```bash
# fix editor var
apt install -y nano micro
set EDITOR nano
set VISUAL nano
export EDITOR nano
export VISUAL nano
```
