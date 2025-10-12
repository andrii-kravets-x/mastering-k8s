## Problem
When we run `sudo kubebuilder/bin/kubectl create deploy demo --image nginx` our pod cannot be scheduled and has `Pending` status.


## Troubleshooting
```bash
@andrii-kravets-x ➜ /workspaces/codespaces-blank $ sudo  kubebuilder/bin/kubectl describe pod demo
Name:             demo-677cfb9d49-d2rfb
Namespace:        default
Priority:         0
Service Account:  default
Node:             <none>
Labels:           app=demo
                  pod-template-hash=677cfb9d49
Annotations:      <none>
Status:           Pending
IP:               
IPs:              <none>
Controlled By:    ReplicaSet/demo-677cfb9d49
Containers:
  nginx:
    Image:        nginx
    Port:         <none>
    Host Port:    <none>
    Environment:  <none>
    Mounts:
      /var/run/secrets/kubernetes.io/serviceaccount from kube-api-access-wmkxf (ro)
Conditions:
  Type           Status
  PodScheduled   False 
Volumes:
  kube-api-access-wmkxf:
    Type:                    Projected (a volume that contains injected data from multiple sources)
    TokenExpirationSeconds:  3607
    ConfigMapName:           kube-root-ca.crt
    ConfigMapOptional:       <nil>
    DownwardAPI:             true
QoS Class:                   BestEffort
Node-Selectors:              <none>
Tolerations:                 node.kubernetes.io/not-ready:NoExecute op=Exists for 300s
                             node.kubernetes.io/unreachable:NoExecute op=Exists for 300s
Events:
  Type     Reason            Age    From               Message
  ----     ------            ----   ----               -------
  Warning  FailedScheduling  7m40s  default-scheduler  0/1 nodes are available: 1 node(s) had untolerated taint {node.cloudprovider.kubernetes.io/uninitialized: true}. preemption: 0/1 nodes are available: 1 Preemption is not helpful for scheduling.
  Warning  FailedScheduling  2m14s  default-scheduler  0/1 nodes are available: 1 node(s) had untolerated taint {node.cloudprovider.kubernetes.io/uninitialized: true}. preemption: 0/1 nodes are available: 1 Preemption is not helpful for scheduling.

@andrii-kravets-x ➜ /workspaces/codespaces-blank $ sudo  kubebuilder/bin/kubectl describe node
Name:               codespaces-3e7062
Roles:              master
Labels:             beta.kubernetes.io/arch=amd64
                    beta.kubernetes.io/os=linux
                    kubernetes.io/arch=amd64
                    kubernetes.io/hostname=codespaces-3e7062
                    kubernetes.io/os=linux
                    node-role.kubernetes.io/master=
Annotations:        alpha.kubernetes.io/provided-node-ip: 10.0.0.18
                    node.alpha.kubernetes.io/ttl: 0
                    volumes.kubernetes.io/controller-managed-attach-detach: true
CreationTimestamp:  Sun, 05 Oct 2025 18:01:33 +0000
Taints:             node.cloudprovider.kubernetes.io/uninitialized=true:NoSchedule
Unschedulable:      false
Lease:
...
...
...
@andrii-kravets-x ➜ /workspaces/codespaces-blank $ sudo kubebuilder/bin/kubectl events -w -A
NAMESPACE   LAST SEEN   TYPE     REASON     OBJECT                   MESSAGE
default     7m10s       Normal   Starting   Node/codespaces-3e7062   Starting kubelet.
default     7m10s       Warning   InvalidDiskCapacity   Node/codespaces-3e7062   invalid capacity 0 on image filesystem
default     7m10s (x2 over 7m10s)   Normal    NodeHasSufficientMemory   Node/codespaces-3e7062   Node codespaces-3e7062 status is now: NodeHasSufficientMemory
default     7m10s (x2 over 7m10s)   Normal    NodeHasNoDiskPressure     Node/codespaces-3e7062   Node codespaces-3e7062 status is now: NodeHasNoDiskPressure
default     7m10s (x2 over 7m10s)   Normal    NodeHasSufficientPID      Node/codespaces-3e7062   Node codespaces-3e7062 status is now: NodeHasSufficientPID
default     7m10s                   Normal    NodeAllocatableEnforced   Node/codespaces-3e7062   Updated Node Allocatable limit across pods
default     7m10s                   Normal    NodeReady                 Node/codespaces-3e7062   Node codespaces-3e7062 status is now: NodeReady
default     6m3s                    Normal    RegisteredNode            Node/codespaces-3e7062   Node codespaces-3e7062 event: Registered Node codespaces-3e7062 in Controller
default     4m7s                    Warning   FailedScheduling          Pod/demo-677cfb9d49-d2rfb   0/1 nodes are available: 1 node(s) had untolerated taint {node.cloudprovider.kubernetes.io/uninitialized: true}. preemption: 0/1 nodes are available: 1 Preemption is not helpful for scheduling.
default     4m7s                    Normal    SuccessfulCreate          ReplicaSet/demo-677cfb9d49   Created pod: demo-677cfb9d49-d2rfb
default     4m7s                    Normal    ScalingReplicaSet         Deployment/demo              Scaled up replica set demo-677cfb9d49 to 1
default     0s                      Warning   FailedScheduling          Pod/demo-677cfb9d49-d2rfb    0/1 nodes are available: 1 node(s) had untolerated taint {node.cloudprovider.kubernetes.io/uninitialized: true}. preemption: 0/1 nodes are available: 1 Preemption is not helpful for scheduling.
default     0s (x2 over 5m)         Warning   FailedScheduling          Pod/demo-677cfb9d49-d2rfb    0/1 nodes are available: 1 node(s) had untolerated taint {node.cloudprovider.kubernetes.io/uninitialized: true}. preemption: 0/1 nodes are available: 1 Preemption is not helpful for scheduling.
```


## Solution: remove the taint from the node
- https://kubernetes.io/docs/tasks/administer-cluster/running-cloud-controller/#running-cloud-controller-manager
- https://kubernetes.io/docs/concepts/scheduling-eviction/taint-and-toleration/#taint-based-evictions

If you’re not using a CCM, the kubelet should not start with any cloud provider flag.
The taint exists because kubelet registered the node with:
`--cloud-provider=external`
→ which tells Kubernetes to expect a CCM to initialize it — but none exists.

Let's remove the taint **manually** with:
```bash
sudo kubebuilder/bin/kubectl taint nodes codespaces-3e7062  node.cloudprovider.kubernetes.io/uninitialized:NoSchedule-
```


## Combined logs after we removed the taint
```log
I1005 18:37:21.668857   40956 schedule_one.go:1046] "Unable to schedule pod; no fit; waiting" pod="default/demo-677cfb9d49-d2rfb" err="0/1 nodes are available: 1 node(s) had untolerated taint {node.cloudprovider.kubernetes.io/uninitialized: true}. preemption: 0/1 nodes are available: 1 Preemption is not helpful for scheduling."
AFI1005 18:42:00.074250   41312 topology_manager.go:215] "Topology Admit Handler" podUID="9174a7ff-bbef-410d-a95a-942a8d48f010" podNamespace="default" podName="demo-677cfb9d49-d2rfb"
I1005 18:42:00.076477   41578 replica_set.go:676] "Finished syncing" logger="replicaset-controller" kind="ReplicaSet" key="default/demo-677cfb9d49" duration="137.414µs"
I1005 18:42:00.080592   40956 schedule_one.go:304] "Successfully bound pod to node" pod="default/demo-677cfb9d49-d2rfb" node="codespaces-3e7062" evaluatedNodes=1 feasibleNodes=1
I1005 18:42:00.107731   41578 replica_set.go:676] "Finished syncing" logger="replicaset-controller" kind="ReplicaSet" key="default/demo-677cfb9d49" duration="100.981µs"
I1005 18:42:00.120476   41312 reconciler_common.go:247] "operationExecutor.VerifyControllerAttachedVolume started for volume \"kube-api-access-wmkxf\" (UniqueName: \"kubernetes.io/projected/9174a7ff-bbef-410d-a95a-942a8d48f010-kube-api-access-wmkxf\") pod \"demo-677cfb9d49-d2rfb\" (UID: \"9174a7ff-bbef-410d-a95a-942a8d48f010\") " pod="default/demo-677cfb9d49-d2rfb"
I1005 18:42:00.221935   41312 reconciler_common.go:220] "operationExecutor.MountVolume started for volume \"kube-api-access-wmkxf\" (UniqueName: \"kubernetes.io/projected/9174a7ff-bbef-410d-a95a-942a8d48f010-kube-api-access-wmkxf\") pod \"demo-677cfb9d49-d2rfb\" (UID: \"9174a7ff-bbef-410d-a95a-942a8d48f010\") " pod="default/demo-677cfb9d49-d2rfb"
I1005 18:42:00.229384   41312 operation_generator.go:721] "MountVolume.SetUp succeeded for volume \"kube-api-access-wmkxf\" (UniqueName: \"kubernetes.io/projected/9174a7ff-bbef-410d-a95a-942a8d48f010-kube-api-access-wmkxf\") pod \"demo-677cfb9d49-d2rfb\" (UID: \"9174a7ff-bbef-410d-a95a-942a8d48f010\") " pod="default/demo-677cfb9d49-d2rfb"
INFO[2025-10-05T18:42:00.387307718Z] RunPodSandbox for &PodSandboxMetadata{Name:demo-677cfb9d49-d2rfb,Uid:9174a7ff-bbef-410d-a95a-942a8d48f010,Namespace:default,Attempt:0,} 
{"level":"info","ts":"2025-10-05T18:42:01.329289Z","caller":"traceutil/trace.go:171","msg":"trace[545613969] transaction","detail":"{read_only:false; response_revision:872; number_of_response:1; }","duration":"115.227576ms","start":"2025-10-05T18:42:01.214045Z","end":"2025-10-05T18:42:01.329272Z","steps":["trace[545613969] 'process raft request'  (duration: 115.146432ms)"],"step_count":1}
{"level":"info","ts":"2025-10-05T18:42:01.332421Z","caller":"mvcc/index.go:214","msg":"compact tree index","revision":778}
{"level":"info","ts":"2025-10-05T18:42:01.33785Z","caller":"mvcc/kvstore_compaction.go:68","msg":"finished scheduled compaction","compact-revision":778,"took":"4.101453ms","hash":1816176270,"current-db-size-bytes":651264,"current-db-size":"651 kB","current-db-size-in-use-bytes":479232,"current-db-size-in-use":"479 kB"}
{"level":"info","ts":"2025-10-05T18:42:01.338466Z","caller":"mvcc/hash.go:137","msg":"storing new hash","hash":1816176270,"revision":778,"compact-revision":688}
INFO[2025-10-05T18:42:01.348141772Z] ImageCreate event name:"registry.k8s.io/pause:3.10" labels:{key:"io.cri-containerd.image" value:"managed"} labels:{key:"io.cri-containerd.pinned" value:"pinned"} 
INFO[2025-10-05T18:42:01.348882401Z] stop pulling image registry.k8s.io/pause:3.10: active requests=0, bytes read=321070 
INFO[2025-10-05T18:42:01.349067832Z] ImageCreate event name:"sha256:873ed75102791e5b0b8a7fcd41606c92fcec98d56d05ead4ac5131650004c136" labels:{key:"io.cri-containerd.image" value:"managed"} labels:{key:"io.cri-containerd.pinned" value:"pinned"} 
INFO[2025-10-05T18:42:01.351622722Z] ImageCreate event name:"registry.k8s.io/pause@sha256:ee6521f290b2168b6e0935a181d4cff9be1ac3f505666ef0e3c98fae8199917a" labels:{key:"io.cri-containerd.image" value:"managed"} labels:{key:"io.cri-containerd.pinned" value:"pinned"} 
INFO[2025-10-05T18:42:01.353169360Z] Pulled image "registry.k8s.io/pause:3.10" with image id "sha256:873ed75102791e5b0b8a7fcd41606c92fcec98d56d05ead4ac5131650004c136", repo tag "registry.k8s.io/pause:3.10", repo digest "registry.k8s.io/pause@sha256:ee6521f290b2168b6e0935a181d4cff9be1ac3f505666ef0e3c98fae8199917a", size "320368" in 792.176765ms 
INFO[2025-10-05T18:42:01.374600991Z] connecting to shim f575a1729f331a0628c2d568eb3315e6ce0149ae11446bf818bb146e285b62e0  address="unix:///run/containerd/s/b76a7c08c0dc0ee7bda74b11d6675fd6c534f45e182356f2b376af6d235f66a7" namespace=k8s.io protocol=ttrpc version=2
INFO[2025-10-05T18:42:01.507061589Z] RunPodSandbox for &PodSandboxMetadata{Name:demo-677cfb9d49-d2rfb,Uid:9174a7ff-bbef-410d-a95a-942a8d48f010,Namespace:default,Attempt:0,} returns sandbox id "f575a1729f331a0628c2d568eb3315e6ce0149ae11446bf818bb146e285b62e0" 
INFO[2025-10-05T18:42:01.510071107Z] PullImage "nginx:latest"                     
```

