# New Controller

A simple Kubernetes controller example built with [controller-runtime](https://github.com/kubernetes-sigs/controller-runtime). This project demonstrates how to create a custom controller that manages a custom resource type.

## Overview

This controller manages `NewResource` custom resources and sets their status to `Ready: true` when reconciled. It's a minimal example that shows the basic structure of a Kubernetes controller.

## Custom Resource Definition

The controller manages a custom resource called `NewResource` with the following structure:

```yaml
apiVersion: apps.newresource.com/v1alpha1
kind: NewResource
metadata:
  name: example-resource
spec:
  foo: "bar"
status:
  ready: true
```

## Quick Start

### Prerequisites

- Go 1.24.2 or later
- Kubernetes cluster (local with kind/minikube or remote)
- kubectl configured to access your cluster
- controller-gen installed: `go install sigs.k8s.io/controller-tools/cmd/controller-gen@latest`

### Step 1: Start with API Structure

Create the basic API folder structure with two essential files:

```
api/v1alpha1/
├── newresource_types.go    # Your CRD definition
└── groupversion.go         # Group and version registration
```

**newresource_types.go** - Define your custom resource:
```go
// +groupName=apps.newresource.com
package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
type NewResource struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NewResourceSpec   `json:"spec,omitempty"`
	Status NewResourceStatus `json:"status,omitempty"`
}

type NewResourceSpec struct {
	Foo string `json:"foo,omitempty"`
}

type NewResourceStatus struct {
	Ready bool `json:"ready,omitempty"`
}

// +kubebuilder:object:root=true
type NewResourceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NewResource `json:"items"`
}
```

**groupversion.go** - Register your API group:
```go
// +groupName=apps.newresource.com
package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	GroupVersion  = schema.GroupVersion{Group: "apps.newresource.com", Version: "v1alpha1"}
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme   = SchemeBuilder.AddToScheme
)

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(GroupVersion,
		&NewResource{},
		&NewResourceList{},
	)
	metav1.AddToGroupVersion(scheme, GroupVersion)
	return nil
}
```

### Step 2: Generate Objects and APIs

Generate deep copy methods and CRDs:

```bash
# Generate deep copy methods
controller-gen object paths="./api/..."

# Generate CRDs
controller-gen crd:crdVersions=v1 paths=./... output:crd:artifacts:config=config/crd/bases
```

### Step 3: Create Main Application

Create the main application entry point:

```go
// main.go
package main

import (
	"flag"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/metrics/server"

	newv1 "github.com/den-vasyliev/new-controller/api/v1alpha1"
	"github.com/den-vasyliev/new-controller/controllers"
)

func main() {
	var (
		metricsAddr          string
		enableLeaderElection bool
	)

	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false, "Enable leader election for controller manager.")
	flag.Parse()

	scheme := runtime.NewScheme()
	utilruntime.Must(newv1.AddToScheme(scheme))

	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))

	mgr, err := manager.New(config.GetConfigOrDie(), manager.Options{
		Scheme:           scheme,
		Metrics:          server.Options{BindAddress: metricsAddr},
		LeaderElection:   enableLeaderElection,
		LeaderElectionID: "newresource-controller",
	})
	if err != nil {
		panic(err)
	}

	if err := (&controllers.NewResourceReconciler{
		Client: mgr.GetClient(),
	}).SetupWithManager(mgr); err != nil {
		panic(err)
	}

	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		panic(err)
	}
}
```

### Step 4: Create Controller

Create the controller logic:

```go
// controllers/resource_controller.go
package controllers

import (
	"context"
	newv1 "github.com/den-vasyliev/new-controller/api/v1alpha1"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type NewResourceReconciler struct {
	client.Client
}

func (r *NewResourceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	var resource newv1.NewResource
	if err := r.Get(ctx, req.NamespacedName, &resource); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	logger.Info("Reconciling", "name", resource.Name)

	resource.Status.Ready = true
	if err := r.Status().Update(ctx, &resource); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *NewResourceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&newv1.NewResource{}).
		Complete(r)
}
```

### Step 5: Build and Run

```bash
# Initialize Go module
go mod init github.com/den-vasyliev/mastering-k8s/new-controller

# Build the controller
go mod tidy
go build -o bin/manager main.go

# Install CRD
kubectl apply -f config/crd/bases/apps.newresource.com_newresources.yaml

# Run the controller
./bin/manager
```

### Step 6: Test the Controller

Create a test resource:

```bash
cat <<EOF | kubectl apply -f -
apiVersion: apps.newresource.com/v1alpha1
kind: NewResource
metadata:
  name: example-resource
spec:
  foo: "hello world"
EOF
```

Verify the controller is working:

```bash
kubectl get newresource example-resource -o yaml
```

You should see `status.ready: true` in the output.

## Development

### Running Tests

```bash
KUBEBUILDER_ASSETS="<PATH_TO_TESTENV_BIN>" go test ./...
```

### Generating CRDs

If you modify the API types, regenerate the CRD:

```bash
# Generate CRDs
controller-gen crd paths=./api/... output:crd:dir=config/crd/bases
```

### Project Structure

```
.
├── api/v1alpha1/           # API definitions
│   ├── groupversion.go     # Group version configuration
│   └── newresource_types.go # NewResource CRD definition
├── config/crd/bases/       # Generated CRD manifests
├── controllers/            # Controller logic
│   └── resource_controller.go
├── test/                   # Test utilities and examples
├── main.go                 # Application entry point
├── go.mod                  # Go module dependencies
└── README.md              # This file
```

## How It Works

1. The controller watches for `NewResource` objects in the cluster
2. When a resource is created or updated, the `Reconcile` method is called
3. The controller logs the reconciliation and sets `status.ready = true`
4. The status update is persisted back to the cluster

## Customization

To extend this controller for your own use case:

1. Modify `api/v1alpha1/newresource_types.go` to define your custom resource spec and status
2. Update the `Reconcile` method in `controllers/resource_controller.go` to implement your business logic
3. Regenerate CRDs if you changed the API types
4. Add any additional dependencies to `go.mod`

## License

This project is licensed under the MIT License - see the LICENSE file for details.
