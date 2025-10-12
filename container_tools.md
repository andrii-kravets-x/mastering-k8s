**Summary:**

| Tool        | Purpose                                             | Layer             | Typical Use                                                                   |
| ----------- | --------------------------------------------------- | ----------------- | ----------------------------------------------------------------------------- |
| **unshare** | Create new Linux namespaces manually                | Kernel / syscalls | Debugging, testing isolation (e.g. `unshare --uts --ipc --net --pid --mount`) |
| **runc**    | Low-level container runtime (implements OCI spec)   | Runtime           | Used by Docker, containerd, CRI-O to spawn containers                         |
| **crictl**  | CLI for CRI-compatible runtimes (containerd, CRI-O) | K8s interface     | Inspect/manage pods/containers in Kubernetes (CRI layer)                      |

**Relations:**
`crictl` → talks to CRI → (containerd/CRI-O) → uses `runc` → uses namespaces (`unshare` syscalls).

**Alternative Tools:**

* `nsenter` — enter existing namespaces
* `podman` — daemonless container tool using `runc`
* `ctr` — containerd’s direct CLI
* `nerdctl` - Docker-compatible CLI for containerd