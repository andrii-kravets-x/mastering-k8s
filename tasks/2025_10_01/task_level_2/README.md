# Task
## 2. Просунутий рівень.
- розгорніть control plane
- створіть debug privileged container з image verizondigital/kubectl-flame:v0.2.4-perf
- зробіть профілювання kube-apiserver: збір семплів з PID (perf record -F 99 -g -p ...)
- побудуйте flame graph (perf script -i /tmp/out | FlameGraph/stackcollapse-perf.pl | FlameGraph/flamegraph.pl > flame.svg)
- скопіюйте flame.svg з контейнера та збережіть у вашому репо

## Solution
- розгорніть control plane

done with kubeadm + flannel on Proxmox, [guide](/tasks/2025_10_01/task_level_2/short_k8s_install.md)

--- 

### For kube-apiserver profiling - four options were found:
1. Use https://github.com/yahoo/kubectl-flame or https://github.com/josepdcs/kubectl-prof, see spoiler
    <details><summary>kubectl-prof</summary>

    You will require **krew plugin manager** for kubectl and **prof** plugin
    ```bash
    ubuntu@kubeadm-01 ~> kubectl prof podinfo-544cf96bd5-66j6c -t 1m --lang clang --tool perf  -o flamegraph --local-path=./
    Verified target pod ... ✔
    Launched profiler ... 🚀
    Profiling ... 🔬
    Remote profiling file downloaded in 0.037513 seconds. ✔
    The profiling result file [podinfo-544cf96bd5-66j6c-agent-flamegraph-2598-1-2025-10-13T19_08_49Z.svg] was obtained in 60.950172 seconds. 🔥
    ```


    This was needed on VM where kubeadm-01 runs using bpf (when `--tool` flag is **not** set)
    ```
    # See: https://kernel.ubuntu.com/mainline/
    wget https://kernel.ubuntu.com/mainline/v6.17/amd64/linux-image-unsigned-6.17.0-061700-generic_6.17.0-061700.202509282239_amd64.deb
    wget https://kernel.ubuntu.com/mainline/v6.17/amd64/linux-modules-6.17.0-061700-generic_6.17.0-061700.202509282239_amd64.deb

    wget https://kernel.ubuntu.com/mainline/v6.17/amd64/linux-headers-6.17.0-061700-generic_6.17.0-061700.202509282239_amd64.deb
    wget https://kernel.ubuntu.com/mainline/v6.17/amd64/linux-headers-6.17.0-061700_6.17.0-061700.202509282239_all.deb 

    sudo dpkg -i linux-*.deb
    sudo reboot
    ```


    ```bash
    ubuntu@kubeadm-01 ~> kubectl prof kube-apiserver-kubeadm-01 -n kube-system -t 1m --lang clang -o flamegraph --local-path=./ --log-level trace
    Default profiling tool bpf will be used ... 🧐
    Verified target pod ... ✔
    Launched profiler ... 🚀
    Profiling ... 🔬
    DEBU[2025-10-13T19:46:22Z] File /tmp/agent-flamegraph-1399-1.svg.gz downloaded (local: 7f1768f9ccce8e0cc0231102a6e3b1bf (5829 bytes) | remote: 7f1768f9ccce8e0cc0231102a6e3b1bf (5829 bytes)) 
    Remote profiling file downloaded in 0.035047 seconds. ✔
    The profiling result file [kube-apiserver-kubeadm-01-agent-flamegraph-1399-1-2025-10-13T19_46_22Z.svg] was obtained in 62.178256 seconds. 🔥
    ```


    `kubectl-prof` will run k8s **Job**:
    ```log
    default        0s                   Normal    SuccessfulCreate          Job/kubectl-prof-bpf-3a7549db-df13-40eb-9235-9805a7d7a48b      Created pod: kubectl-prof-bpf-3a7549db-df13-40eb-9235-9805a7d7a48b-8w8l2
    default        0s                   Normal    Pulling                   Pod/kubectl-prof-bpf-3a7549db-df13-40eb-9235-9805a7d7a48b-8w8l2   Pulling image "josepdcs/kubectl-prof:1.6.0-bpf"
    default        0s                   Normal    Pulled                    Pod/kubectl-prof-bpf-3a7549db-df13-40eb-9235-9805a7d7a48b-8w8l2   Successfully pulled image "josepdcs/kubectl-prof:1.6.0-bpf" in 31.805s (31.805s including waiting). Image size: 133266682 bytes.
    default        0s                   Normal    Created                   Pod/kubectl-prof-bpf-3a7549db-df13-40eb-9235-9805a7d7a48b-8w8l2   Created container: kubectl-prof
    default        0s                   Normal    Started                   Pod/kubectl-prof-bpf-3a7549db-df13-40eb-9235-9805a7d7a48b-8w8l2   Started container kubectl-prof
    default        0s                   Normal    Killing                   Pod/kubectl-prof-bpf-3a7549db-df13-40eb-9235-9805a7d7a48b-8w8l2   Stopping container kubectl-prof
    default        0s                   Normal    SuccessfulCreate          Job/kubectl-prof-bpf-271b0461-faa2-47cf-b231-3cb52ac5e52a         Created pod: kubectl-prof-bpf-271b0461-faa2-47cf-b231-3cb52ac5e52a-qbdhc
    default        0s                   Normal    Pulled                    Pod/kubectl-prof-bpf-271b0461-faa2-47cf-b231-3cb52ac5e52a-qbdhc   Container image "josepdcs/kubectl-prof:1.6.0-bpf" already present on machine
    default        0s                   Normal    Created                   Pod/kubectl-prof-bpf-271b0461-faa2-47cf-b231-3cb52ac5e52a-qbdhc   Created container: kubectl-prof
    default        0s                   Normal    Started                   Pod/kubectl-prof-bpf-271b0461-faa2-47cf-b231-3cb52ac5e52a-qbdhc   Started container kubectl-prof
    default        0s                   Normal    Killing                   Pod/kubectl-prof-bpf-271b0461-faa2-47cf-b231-3cb52ac5e52a-qbdhc   Stopping container kubectl-prof
    default        0s                   Normal    SuccessfulCreate          Job/kubectl-prof-perf-9da5ddda-929d-41d2-94ee-60bba9713ccf        Created pod: kubectl-prof-perf-9da5ddda-929d-41d2-94ee-60bba9713ccf-kw7c4
    default        0s                   Normal    Pulling                   Pod/kubectl-prof-perf-9da5ddda-929d-41d2-94ee-60bba9713ccf-kw7c4   Pulling image "josepdcs/kubectl-prof:1.6.0-perf"
    default        0s                   Normal    Pulled                    Pod/kubectl-prof-perf-9da5ddda-929d-41d2-94ee-60bba9713ccf-kw7c4   Successfully pulled image "josepdcs/kubectl-prof:1.6.0-perf" in 41.603s (41.603s including waiting). Image size: 109434443 bytes.
    default        0s                   Normal    Created                   Pod/kubectl-prof-perf-9da5ddda-929d-41d2-94ee-60bba9713ccf-kw7c4   Created container: kubectl-prof
    default        0s                   Normal    Started                   Pod/kubectl-prof-perf-9da5ddda-929d-41d2-94ee-60bba9713ccf-kw7c4   Started container kubectl-prof
    default        0s                   Normal    Killing                   Pod/kubectl-prof-perf-9da5ddda-929d-41d2-94ee-60bba9713ccf-kw7c4   Stopping container kubectl-prof
    ```

    </details>
2. Run **debug pod on the node**
    > The container will run in the host namespaces and the host's filesystem will be mounted at /host
    ```bash
    kubectl debug node/kubeadm-01 -it --image=verizondigital/kubectl-flame:v0.2.4-perf --profile=general -- /app/perf record -F 99 -g -p 1611 -o /host/perf.data
    ```
    this will record profiling file in `/perf.data` on node (VM in our case)
3. Create standalone **ubuntu** privileged pod **on the same node via manifest**, see spoiler
    <details><summary>kubectl-prof</summary>

    Steps:
    - create an ubuntu privileged pod on the same node from [manifest](./ubuntu.yml)
    - exec and install `linux-perf` package, run perf, copy to VM
    - generate flamechart on VM
    - copy to laptop


    ```bash
    # create ubuntu pod
    kubectl apply -f ubuntu.yaml

    kubectl exec ubuntu -it -- bash

    # inside pod
    apt update
    apt install -y linux-perf
    pgrep -f kube-apiserver # this returned 18991 in my case

    # to stop perf - pres ctrl+c
    perf record -F 99 -g -p 18991
    ^C[ perf record: Woken up 1 times to write data ]
    [ perf record: Captured and wrote 0.072 MB perf.data (223 samples) ]

    # on VM
    git clone --depth 1 https://github.com/brendangregg/FlameGraph.git

    kubectl cp ubuntu:/perf.data perf.data
    perf script -i perf.data | FlameGraph/stackcollapse-perf.pl | FlameGraph/flamegraph.pl > flame.svg
    ```

    </details>

4. Run `kubectl debug kube-apiserver-kubeadm-01 -n kube-system -it --image=busybox --profile=general`. 
    ```bash
    Defaulting debug container name to debugger-s2gzd.
    The Pod "kube-apiserver-kubeadm-01" is invalid: []: Forbidden: static pods do not support ephemeral containers
    ```

### Flamegraphs
First 3 options are working fine and can be used to create flamegraphs as it is. Last option will work only for your usual pods, **not** kube-apiserver that kubelet runs via __static pods__.

See prepared flamecharts [here](./flamegraphs/)
