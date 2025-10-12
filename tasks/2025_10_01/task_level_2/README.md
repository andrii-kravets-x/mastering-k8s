# Task
## 2. Просунутий рівень.
- розгорніть control plane
- створіть debug privileged container з image verizondigital/kubectl-flame:v0.2.4-perf
- зробіть профілювання kube-apiserver: збір семплів з PID (perf record -F 99 -g -p ...)
- побудуйте flame graph (perf script -i /tmp/out | FlameGraph/stackcollapse-perf.pl | FlameGraph/flamegraph.pl > flame.svg)
- скопіюйте flame.svg з контейнера та збережіть у вашому репо

## Solution
- розгорніть control plane - done with kubeadm + flannel on Proxmox, [guide](/tasks/2025_10_01/task_level_2/short_k8s_install.md)
- ~~створіть debug privileged container з image verizondigital/kubectl-flame:v0.2.4-perf~~  **не вийшло**
- create an ubuntu privileged pod on the same node
- exec and install `linux-perf` package, run perf, copy to VM
- generate flamechart on VM
- copy to laptop

```bash
# create ubuntu pod
kubectl apply -f debug.yaml

kubectl exec debug -it -- bash

# inside pod
apt install -y linux-perf
pgrep -f kube-apiserver # this returned 18991 in my case

# to stop perf - pres ctrl+c
perf record -F 99 -g -p 18991
^C[ perf record: Woken up 1 times to write data ]
[ perf record: Captured and wrote 0.072 MB perf.data (223 samples) ]

# on VM
git clone --depth 1 https://github.com/brendangregg/FlameGraph.git

kubectl cp debug:/perf.data perf.data
perf script -i perf.data | FlameGraph/stackcollapse-perf.pl | FlameGraph/flamegraph.pl > flame.svg
```


<details><summary>Debug</summary>

### Try Later 
- `verizondigital/kubectl-flame:v0.2.4-perf` image [is 3y old](https://hub.docker.com/layers/verizondigital/kubectl-flame/v0.2.4-perf/images/sha256-5c098533f7731a60a73438f5aacceb402ee4ea03f693d675d87ee46244e18d01)
- source repo(not sure) [yahoo/kubectl-flame](https://github.com/yahoo/kubectl-flame) also is old
- let's try josepdcs/kubectl-prof - https://github.com/josepdcs/kubectl-prof
- https://kubernetes.io/docs/reference/kubectl/generated/kubectl_debug/#examples



```bash
kubectl debug kube-apiserver-codespaces-3e7062 -n kube-system -it --image=nicolaka/netshoot

Defaulting debug container name to debugger-26jg9.
The Pod "kube-apiserver-codespaces-3e7062" is invalid: []: Forbidden: static pods do not support ephemeral containers


# kubectl prof kube-apiserver-codespaces-3e7062 -t 1m --lang go -o flamegraph

kubectl debug kube-apiserver-codespaces-3e7062 -it --image=verizondigital/kubectl-flame:v0.2.4-perf
kubectl debug kube-apiserver-codespaces-3e7062 --privileged -it --image=verizondigital/kubectl-flame:v0.2.4-perf

kubectl exec debug -it -- bash
docker run -it --privileged nicolaka/netshoot
kubectl run tmp-shell --rm -i --tty --image nicolaka/netshoot
```


### install krew and prof 
- https://github.com/josepdcs/kubectl-prof?tab=readme-ov-file#using-krew

```
kubectl prof --help
kubectl prof mypod -t 1m --lang go -o flamegraph
kubectl prof mypod -n profiling --service-account=profiler --target-namespace=my-apps -l go

```

```
root ➜ /workspaces/codespaces-blank $ kubectl prof --help
Profiling on existing applications with low-overhead.

These commands help you identify application performance issues.

Usage:
  prof [pod-name | --selector label]

Examples:

	# Profile a pod for 5 minutes with JFR format for java language
	kubectl prof my-pod -t 5m -l java -o jfr

	# Profile an alpine based container for java language
	kubectl prof my-pod -l java --alpine 

	# Profile a pod for 5 minutes in intervals of 60 seconds for java language by giving the cpu limits, the container runtime, the agent image and the image pull policy
	kubectl my-pod -l java -o flamegraph -t 5m --interval 60s --cpu-limits=1 -r containerd --image=localhost/my-agent-image-jvm:latest --image-pull-policy=IfNotPresent

	# Profile in contprof namespace a pod running in contprof-stupid-apps namespace by using the profiler service account for go language 
	kubectl prof my-pod -n contprof --service-account=profiler --target-namespace=contprof-stupid-apps -l go

	# Set custom resource requests and limits for the agent pod (default: neither requests nor limits are set) for python language
	kubectl prof my-pod --cpu-requests 100m --cpu-limits 200m --mem-requests 100Mi --mem-limits 200Mi -l python

	# Profile the pods with the label selector "app=my-app" for 5 minutes with JFR format for java language
	kubectl prof -l java -o jfr -t 5m --selector app=my-app


Flags:
      --alpine                          TargetConfig image is based on Alpine
      --as string                       Username to impersonate for the operation. User could be a regular user or a service account in a namespace.
      --as-group stringArray            Group to impersonate for the operation, this flag can be repeated to specify multiple groups.
      --as-uid string                   UID to impersonate for the operation.
      --cache-dir string                Default cache directory (default "/root/.kube/cache")
      --capabilities strings            The capabilities to be added to the agent container. It can be used multiple times to add more than one capability (e.g. --capabilities SYS_ADMIN --capabilities SYS_PTRACE)
      --certificate-authority string    Path to a cert file for the certificate authority
      --client-certificate string       Path to a client certificate file for TLS
      --client-key string               Path to a client key file for TLS
      --cluster string                  The name of the kubeconfig cluster to use
      --context string                  The name of the kubeconfig context to use
      --cpu-limits string               CPU limits of the started profiling container
      --cpu-requests string             CPU requests of the started profiling container
      --disable-compression             If true, opt-out of response compression for all requests to the server
      --dry-run                         Simulate profiling
  -e, --event string                    Profiling event, choose one of [cpu alloc lock cache-misses wall itimer] (default "itimer")
      --grace-period-ending duration    The grace period to spend before to end the agent (default 5m0s)
      --heap-dump-split-size string     The heap dump (or snapshot, for Node.js) will be split into chunks of a specified size, following the valid format for the split command (e.g. 50M, 1G, etc.) (default "50M")
  -h, --help                            help for prof
      --image string                    Manually choose agent docker image
      --image-pull-policy string        Image pull policy, choose one of [Never Always IfNotPresent] (default "IfNotPresent")
      --image-pull-secret string        imagePullSecret for agent docker image
      --insecure-skip-tls-verify        If true, the server's certificate will not be checked for validity. This will make your HTTPS connections insecure
      --interval duration               Max scan Interval
      --kubeconfig string               Path to the kubeconfig file to use for CLI requests.
  -l, --lang string                     Programming language of the target application, choose one of [java go python ruby node clang clang++ rust]
      --local-path string               Optional local path location to store the result files. Default is current directory
      --log-level string                Log level, choose one of [info warn debug trace error panic] (default "info")
      --mem-limits string               Memory limits of the started profiling container
      --mem-requests string             Memory requests of the started profiling container
  -n, --namespace string                If present, the namespace scope for this CLI request
      --node-heap-snapshot-signal int   The signal to be sent to the target process to generate a heap snapshot for Node.js applications (default 12)
  -o, --output string                   Output type, choose one accorting tool {"perf":["flamegraph","raw"],"rbspy":["flamegraph","speedscope","callgrind","summary","summary-by-line"],"node-dummy":["heapsnapshot","heapdump"],"fake":["flamegraph"],"async-profiler":["flamegraph","jfr","flat","traces","collapsed","tree","raw"],"jcmd":["jfr","threaddump","heapdump","heaphistogram"],"pyspy":["flamegraph","speedscope","threaddump","raw"],"bpf":["flamegraph","raw"]} (default "flamegraph")
  -p, --pgrep string                    Name of the target process
      --pid string                      The PID of target process if it is known
      --pool-size-profiling-jobs int    The pool size of goroutines for launching profiling jobs when the '--selector' flag is used (default "No limit: all matching pods will be profiled simultaneously")
      --pool-size-retrieve-chunks int   The pool size of goroutines used to retrieve chunks of the obtained heap dump (or snapshot, for Node.js) from the agent (default 5)
      --print-logs                      Force agent to print the log messages type to standard output (default true)
      --privileged                      Run agent container in privileged mode (default true)
      --request-timeout string          The length of time to wait before giving up on a single server request. Non-zero values should contain a corresponding time unit (e.g. 1s, 2m, 3h). A value of zero means don't timeout requests. (default "0")
      --retrieve-file-retries int       The number of retries to retrieve a file from the remote container (default 3)
  -r, --runtime string                  The container runtime used for kubernetes, choose one of [crio containerd] (default "containerd")
      --runtime-path string             Use a different container runtime install path according to the runtime used (default "/run/containerd")
      --selector string                 Selector (label query) to filter on, supports '=', '==', and '!='.(e.g. -l key1=value1,key2=value2). Matching objects must satisfy all of the specified label constraints.
  -s, --server string                   The address and port of the Kubernetes API server
      --service-account string          serviceAccountName to be used for profiling container
      --target-container-name string    The target container name to be profiled
      --target-namespace string         namespace of target pod if different from job namespace
  -t, --time duration                   Max scan Duration
      --tls-server-name string          Server name to use for server certificate validation. If it is not provided, the hostname used to contact the server is used
      --token string                    Bearer token for authentication to the API server
      --tool string                     Profiling tool, choose one accorfing language {"node":["bpf","perf","node-dummy"],"clang":["bpf","perf"],"ruby":["rbspy"],"rust":["bpf","perf"],"python":["pyspy"],"clang++":["bpf","perf"],"fake":["fake"],"java":["jcmd","async-profiler"],"go":["bpf"]}
      --user string                     The name of the kubeconfig user to use
      --version                         Print version info
```

</details>