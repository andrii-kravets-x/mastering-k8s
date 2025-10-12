package test

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	newv1 "github.com/mastering-k8s/new-controller/api/v1alpha1"
	"github.com/stretchr/testify/require"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/scale/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/metrics/server"
)

var (
	mainTestEnv    *envtest.Environment
	mainTestConfig *rest.Config
	mainTestScheme *runtime.Scheme
	mainTestClient client.Client
	mainTestMgr    manager.Manager
)

// setupMainTestEnv creates a test environment for main.go testing
func setupMainTestEnv(t *testing.T) func() {
	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))

	// Create a new scheme that includes both core Kubernetes types and your CRDs
	mainTestScheme = runtime.NewScheme()

	// Add the core Kubernetes schemes
	err := scheme.AddToScheme(mainTestScheme)
	require.NoError(t, err)

	// Add the CRD scheme
	metav1.AddToGroupVersion(mainTestScheme, newv1.GroupVersion)
	require.NoError(t, err)

	// Add CRD scheme
	err = apiextensionsv1.AddToScheme(mainTestScheme)
	require.NoError(t, err)

	// Add your custom schemes
	err = newv1.AddToScheme(mainTestScheme)
	require.NoError(t, err)

	// Create test environment
	mainTestEnv = &envtest.Environment{
		CRDDirectoryPaths: []string{
			filepath.Join("..", "config", "crd", "bases"),
		},
		ErrorIfCRDPathMissing:    true,
		AttachControlPlaneOutput: false,
	}

	// Create a longer context timeout for environment startup
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var startErr = make(chan error)
	go func() {
		var err error
		mainTestConfig, err = mainTestEnv.Start()
		startErr <- err
	}()

	// Wait for environment to start with timeout
	select {
	case err := <-startErr:
		require.NoError(t, err, "Timeout waiting for test environment to start")
	case <-ctx.Done():
		require.NoError(t, err, "Failed to start test environment")
	}

	// Create the client with the combined scheme
	mainTestClient, err = client.New(mainTestConfig, client.Options{Scheme: mainTestScheme})
	if err != nil {
		t.Fatalf("Failed to create test client: %v", err)
	}

	// Initialize manager
	mainTestMgr, err = manager.New(mainTestConfig, manager.Options{
		Scheme:                  mainTestScheme,
		LeaderElection:          false,
		LeaderElectionNamespace: "default",
		HealthProbeBindAddress:  "0",
		Metrics:                 server.Options{BindAddress: "0"},
	})
	require.NoError(t, err)

	// Return cleanup function
	return func() {
		err := mainTestEnv.Stop()
		require.NoError(t, err)
	}
}
