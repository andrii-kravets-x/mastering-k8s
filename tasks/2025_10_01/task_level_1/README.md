# Task
## 1. Початковий рівень. Базуючись на досвіді першого завдання:
- згенеруйте 4 маніфести для etcd, kube-apiserver, kube-scheduler та kube-controller-manager
- розгорніть control plane за допомогою kubelet staticPod
- розгорніть кастомний deployment (наприклад, nginx) на 3 реплікі.

Питання: Чому не працює? Як пофіксити?


## Solution
We know that kubeadm deploys control plane components using static Pod manifests:
- https://kubernetes.io/docs/reference/setup-tools/kubeadm/implementation-details/#generate-static-pod-manifests-for-control-plane-components

Let's generate and edit them to be able to run with kubelet from setup.sh

See edited versions [here](/tasks/2025_10_01/task_level_1/control-plane-manifests)

### On fresh codescape run setup.sh to regenerate certificates 
```bash
# run in the first term
./setup.sh start

# and kill services that we will deploy in containers
sudo pkill -f kube-controller-manager
sudo pkill kube-scheduler
sudo pkill kube-apiserver
sudo pkill etcd
```

### run only kubelet
```bash
# run in the first term
HOST_IP=$(hostname -I | awk '{print $1}')
sudo PATH=$PATH:/opt/cni/bin:/usr/sbin kubebuilder/bin/kubelet \
    --kubeconfig=/var/lib/kubelet/kubeconfig \
    --config=/var/lib/kubelet/config.yaml \
    --root-dir=/var/lib/kubelet \
    --cert-dir=/var/lib/kubelet/pki \
    --tls-cert-file=/var/lib/kubelet/pki/kubelet.crt \
    --tls-private-key-file=/var/lib/kubelet/pki/kubelet.key \
    --hostname-override=$(hostname) \
    --pod-infra-container-image=registry.k8s.io/pause:3.10 \
    --node-ip=$HOST_IP \
    --cloud-provider=external \
    --cgroup-driver=cgroupfs \
    --max-pods=50  \
    --v=1 2>&1 | ./tspin

# tspin is tailspin log formatter(optional)
# https://github.com/bensadeh/tailspin/releases/tag/5.5.0
```

### generate default manifests with kubeadm 
```bash
# run this and everything else in the second term
# download kubeadm binary
sudo curl -L "https://dl.k8s.io/v1.30.0/bin/linux/amd64/kubeadm" -o kubebuilder/bin/kubeadm
sudo chmod +x kubebuilder/bin/kubeadm
```

```bash
# to see default values that kubeadm uses for templating
# kubeadm config print init-defaults > kubeadm.yaml

codespace:/workspaces/codespaces-blank> ./kubebuilder/bin/kubeadm version
kubeadm version: &version.Info{Major:"1", Minor:"30", GitVersion:"v1.30.0", GitCommit:"7c48c2bd72b9bf5c44d21d7338cc7bea77d0ad2a", GitTreeState:"clean", BuildDate:"2024-04-17T17:34:08Z", GoVersion:"go1.22.2", Compiler:"gc", Platform:"linux/amd64"}

# here we generate 3/4 manifests
codespace:/workspaces/codespaces-blank> sudo ./kubebuilder/bin/kubeadm init phase control-plane all --dry-run
I1012 12:17:42.610593   12774 version.go:256] remote version is much newer: v1.34.1; falling back to: stable-1.30
[control-plane] Using manifest folder "/etc/kubernetes/tmp/kubeadm-init-dryrun1631701269"
[control-plane] Creating static Pod manifest for "kube-apiserver"
[control-plane] Creating static Pod manifest for "kube-controller-manager"
[control-plane] Creating static Pod manifest for "kube-scheduler"

codespace:/workspaces/codespaces-blank> sudo ls -lah /etc/kubernetes/tmp/kubeadm-init-dryrun1631701269
total 20K
drwx------ 2 root root 4.0K Oct 12 12:17 .
drwx------ 3 root root 4.0K Oct 12 12:17 ..
-rw------- 1 root root 3.8K Oct 12 12:17 kube-apiserver.yaml
-rw------- 1 root root 3.3K Oct 12 12:17 kube-controller-manager.yaml
-rw------- 1 root root 1.5K Oct 12 12:17 kube-scheduler.yaml

# kubeadm generates the etcd manifest **only** during the actual init, not in a full dry-run of all phases.
# we need to generate it separately 
@andrii-kravets-x ➜ /workspaces/codespaces-blank $ sudo ./kubebuilder/bin/kubeadm init phase etcd local --dry-run
I1012 12:44:24.039706   27437 version.go:256] remote version is much newer: v1.34.1; falling back to: stable-1.30
[etcd] Would ensure that "/var/lib/etcd" directory is present
[etcd] Creating static Pod manifest for local etcd in "/etc/kubernetes/tmp/kubeadm-init-dryrun3550435242"
```


### prepare KUBECONFIG
```bash
sudo cp /root/.kube/config /var/lib/kubelet/kubeconfig
export KUBECONFIG=~/.kube/config
cp /tmp/sa.pub /tmp/ca.crt
```


### put edited manifests from `tasks/2025_10_01/control-plane-manifests/` folder
```bash
git clone https://github.com/andrii-kravets-x/mastering-k8s
sudo cp mastering-k8s/tasks/2025_10_01/task_level_1/control-plane-manifests/* /etc/kubernetes/manifests/
```

### debug pods/containers 
```bash
# get container id's, then get logs 
sudo ./crictl ps -a
sudo ./crictl logs -f a6e190d40a6e6


# somehow before etcd and api-server are fully working containers are not updated by kubelet and we must stop/start kubelet to let it get updated manifest
# and stop/delete containers manually 

# stop all pods
sudo ./crictl pods --quiet | xargs -r ./crictl stopp

# remove pods
sudo ./crictl pods --quiet | xargs -r ./crictl rmp
```

<details><summary>Debug</summary>

```bash
codespace:/workspaces/codespaces-blank> curl -k "https://localhost:10250/healthz?verbose"
[+]ping ok
[+]log ok
[+]syncloop ok
healthz check passed
```


```bash
codespace:/workspaces/codespaces-blank> sudo ctr -n k8s.io container ls
CONTAINER                                                           IMAGE                            RUNTIME                  
72904b0ef09d4b839457c6ecd8407b9f953fc255104c30d275343800e428cf1a    registry.k8s.io/pause:3.10       io.containerd.runc.v2    
a3548a121f1ce2375c6cf12d421f358bb7ea7ea2762a95d70290a5b62f63812a    registry.k8s.io/etcd:3.5.12-0    io.containerd.runc.v2    

codespace:/workspaces/codespaces-blank> sudo ctr -n k8s.io tasks ls
TASK                                                                PID       STATUS    
72904b0ef09d4b839457c6ecd8407b9f953fc255104c30d275343800e428cf1a    127249    RUNNING

codespace:/workspaces/codespaces-blank> sudo ./crictl pods
WARN[0000] Config "/etc/crictl.yaml" does not exist, trying next: "/workspaces/codespaces-blank/crictl.yaml" 
WARN[0000] runtime connect using default endpoints: [unix:///run/containerd/containerd.sock unix:///run/crio/crio.sock unix:///var/run/cri-dockerd.sock]. As the default settings are now deprecated, you should set the endpoint instead. 
POD ID              CREATED             STATE               NAME                                        NAMESPACE           ATTEMPT             RUNTIME
72904b0ef09d4       8 minutes ago       Ready               etcd-codespaces-3e7062                      kube-system         1                   (default)
```

### runc does not display anything...
```bash
sudo runc --root /run/containerd/io.containerd.runtime.v2.task/k8s.io list
load container 72904b0ef09d4b839457c6ecd8407b9f953fc255104c30d275343800e428cf1a: container does not exist
ID          PID         STATUS      BUNDLE      CREATED     OWNER
```

### Why?
That’s expected — `runc` doesn’t manage those containers directly; it only executes them on behalf of **containerd**, then exits.

In the `io.containerd.runc.v2` model, each container has its **own dedicated shim** process that holds its state. After creation, `runc` exits, so there’s no persistent runtime state for `runc` to list.

**Explanation:**

* Containerd spawns `runc` temporarily to create/start a container.
* Then it **hands off** control to the shim (`containerd-shim-runc-v2`).
* The real state lives under containerd, not runc.

**Therefore:**
`runc list` under any root will usually show nothing for containerd-managed containers.

**If you really need runc-level info:**
You can inspect the bundle (OCI state) used by containerd:

```bash
sudo ls /run/containerd/io.containerd.runtime.v2.task/k8s.io/<container-id>/
sudo cat /run/containerd/io.containerd.runtime.v2.task/k8s.io/<container-id>/config.json
```

That’s the actual OCI spec `runc` consumed.

</details>

### run podinfo image to test if typical deployment works 
```bash
# ensure that you did not forget to untaint the node
# ensure that kubelet max-pods is high enough

kubectl apply -f mastering-k8s/tasks/2025_10_01/task_level_1/podinfo.yaml
```