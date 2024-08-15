package main

import (
	"context"
	"log"
	"sigs.k8s.io/controller-runtime/pkg/client/config"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	"github.com/iraklikairakli/appconfig-controller/v1alpha1" // Correct import path
)

type ReconcileAppConfig struct {
	client.Client
}

func (r *ReconcileAppConfig) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	// Fetch the AppConfig instance
	appConfig := &v1alpha1.AppConfig{}
	err := r.Get(ctx, request.NamespacedName, appConfig)
	if err != nil {
		return reconcile.Result{}, client.IgnoreNotFound(err)
	}

	// Create or update a ConfigMap based on the AppConfig
	configMap := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      appConfig.Name + "-config",
			Namespace: appConfig.Namespace,
		},
		Data: appConfig.Spec.Settings,
	}

	if err := r.Client.Create(ctx, configMap); err != nil {
		return reconcile.Result{}, err
	}

	// Update AppConfig status to indicate that it has been applied
	appConfig.Status.Applied = true
	if err := r.Status().Update(ctx, appConfig); err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}

func main() {
	// Set up the manager and controller
	mgr, err := manager.New(config.GetConfigOrDie(), manager.Options{})
	if err != nil {
		log.Fatalf("Failed to create manager: %v", err)
	}

	// Add the controller to the manager
	c, err := controller.New("appconfig-controller", mgr, controller.Options{
		Reconciler: &ReconcileAppConfig{
			Client: mgr.GetClient(),
		},
	})
	if err != nil {
		log.Fatalf("Failed to create controller: %v", err)
	}

	// Watch for changes to AppConfig resources
	err = c.Watch(&source.Kind{Type: &v1alpha1.AppConfig{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		log.Fatalf("Failed to watch AppConfig: %v", err)
	}

	// Start the manager
	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		log.Fatalf("Failed to start manager: %v", err)
	}
}
