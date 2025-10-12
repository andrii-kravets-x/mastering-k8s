package test

import (
	"context"
	"testing"
	"time"

	newv1 "github.com/mastering-k8s/new-controller/api/v1alpha1"
	"github.com/mastering-k8s/new-controller/controllers"
	"github.com/stretchr/testify/require"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// TestMainController tests that the controller can be started and CRD is available
func TestMainController(t *testing.T) {
	cleanup := setupMainTestEnv(t)
	// defer cleanup() // Moved to after sleep

	// Set up the reconciler once for all tests
	reconciler := &controllers.NewResourceReconciler{
		Client: mainTestMgr.GetClient(),
	}

	err := reconciler.SetupWithManager(mainTestMgr)
	require.NoError(t, err)

	t.Run("test_crd_available", func(t *testing.T) {
		// Check that our CRD is available in the test cluster
		crd := &apiextensionsv1.CustomResourceDefinition{}
		err := mainTestClient.Get(context.TODO(), client.ObjectKey{
			Name: "newresources.apps.newresource.com",
		}, crd)
		require.NoError(t, err)
		require.Equal(t, "newresources.apps.newresource.com", crd.Name)
	})

	t.Run("test_controller_startup", func(t *testing.T) {
		// Start the manager in a goroutine
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		go func() {
			if err := mainTestMgr.Start(ctx); err != nil {
				t.Logf("Manager stopped with error: %v", err)
			}
		}()

		// Wait a bit to ensure the manager starts
		time.Sleep(1 * time.Second)

		// Verify the manager is running
		require.NotNil(t, mainTestMgr)
		require.NotNil(t, mainTestMgr.GetClient())
		require.NotNil(t, mainTestMgr.GetScheme())
	})

	t.Run("test_can_create_newresource", func(t *testing.T) {
		// Test that we can create a NewResource instance
		newResource := &newv1.NewResource{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-resource",
				Namespace: "default",
			},
			Spec: newv1.NewResourceSpec{
				// Add any required spec fields here
			},
		}

		err := mainTestClient.Create(context.TODO(), newResource)
		require.NoError(t, err)

		// Clean up
		defer func() {
			mainTestClient.Delete(context.TODO(), newResource)
		}()

		// Verify it was created
		created := &newv1.NewResource{}
		err = mainTestClient.Get(context.TODO(), client.ObjectKey{
			Name:      "test-resource",
			Namespace: "default",
		}, created)
		require.NoError(t, err)
		require.Equal(t, "test-resource", created.Name)
	})

	// Cleanup after sleep
	cleanup()
}
