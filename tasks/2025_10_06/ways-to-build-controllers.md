### Here are several ways to create a Kubernetes controller:
1. **client-go (Low-level)**

   * Use [`k8s.io/client-go`](https://github.com/kubernetes/client-go) with Informers, Listers, and Workqueues.
   * Maximum control, but verbose and manual.
   * 📘 Docs: [https://pkg.go.dev/k8s.io/client-go](https://pkg.go.dev/k8s.io/client-go)

2. **controller-runtime (Standalone)**

   * Go framework [`sigs.k8s.io/controller-runtime`](https://github.com/kubernetes-sigs/controller-runtime).
   * Provides `Manager`, `Reconciler`, and utilities for building controllers.
   * 📘 Docs: [https://pkg.go.dev/sigs.k8s.io/controller-runtime](https://pkg.go.dev/sigs.k8s.io/controller-runtime)

3. **Kubebuilder**

   * Command: `kubebuilder create api --group <group> --version <version> --kind <Kind>`
   * Built on `controller-runtime`, scaffolds full project structure.
   * Best for production-grade CRDs.
   * 📘 Docs: [https://book.kubebuilder.io](https://book.kubebuilder.io)

4. **Operator SDK**

   * Command: `operator-sdk init --domain <domain>`
   * Built on `controller-runtime`; adds Operator patterns and Helm/Ansible options.
   * 📘 Docs: [https://sdk.operatorframework.io/docs](https://sdk.operatorframework.io/docs)

5. **Metacontroller**

   * Framework using CRDs and webhooks, no Go code needed.
   * Define behavior via JSON + scripts or templates.
   * 📘 Repo: [https://github.com/metacontroller/metacontroller](https://github.com/metacontroller/metacontroller)

6. **Kopf (Kubernetes Operator Framework for Python)**

   * Python-based operator framework with decorator API.
   * Fast prototyping and simpler syntax for automation.
   * 📘 Docs: [https://kopf.readthedocs.io](https://kopf.readthedocs.io)
