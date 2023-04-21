/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	tutorialsv1 "github.com/hsaid4327/configwatcher-go-operator/api/v1"
	metav1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

// ConfigWatcherReconciler reconciles a ConfigWatcher object
type ConfigWatcherReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=tutorials.github.com,resources=configwatchers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=tutorials.github.com,resources=configwatchers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=tutorials.github.com,resources=configwatchers/finalizers,verbs=update
//+kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch;delete
//+kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;list;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ConfigWatcher object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *ConfigWatcherReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	log := log.FromContext(ctx)

	// Fetch the Memcached instance
	// The purpose is check if the Custom Resource for the Kind configmap
	// is present in the namespace on the cluster if not we return nil to stop the reconciliation
	configMap := &metav1.ConfigMap{}
	
	err := r.Get(ctx, req.NamespacedName, configMap)
	//configMap.
	if err != nil {
		if apierrors.IsNotFound(err) {
			// If the configmap is not found then, it usually means that it was deleted or not created
			// In this way, we will stop the reconciliation
			log.Info("confimap resource not found")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		log.Error(err, "Failed to get configmap")
		return ctrl.Result{}, err
	}
	configwatcher := &tutorialsv1.ConfigWatcher{}
	err = r.Get(ctx, req.NamespacedName, configwatcher)
	if err != nil {
		if apierrors.IsNotFound(err) {
			// If the custom resource is not found then, it usually means that it was deleted or not created
			// In this way, we will stop the reconciliation
			log.Info("configwatcher CR resource not found")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		log.Error(err, "Failed to get ConfigWatcher")
		return ctrl.Result{}, err
	}
	configMapName := configwatcher.Spec.ConfigMap
	// if this configmap is specified in CR
	if configMapName == configMap.Name {
		podSelector := configwatcher.Spec.PodSelector
		err = deletePods(podSelector, ctx, req, r)
		if err != nil {
			if apierrors.IsNotFound(err) {
				// If the configmap is not found then, it usually means that it was deleted or not created
				// In this way, we will stop the reconciliation
				log.Info("No pods in running state are found to delete")
				return ctrl.Result{}, nil
			}
			// Error reading the object - requeue the request.
			log.Error(err, "Failed to delete pods ")
			return ctrl.Result{}, err
		}

	}

	//
	//if isConfigMapModified(configMapName) {
	//	podSelector := configwatcher.Spec.PodSelector
	//	deletePods(podSelector)
	//}
	//  if isConfigMapModified(configMap) {

	//	}
	return ctrl.Result{}, nil
}

func deletePods(podSelector map[string]string, ctx context.Context, req ctrl.Request, r *ConfigWatcherReconciler) error {

	//podList := &metav1.PodList{}
	//opts := []client.ListOption{
	//	client.InNamespace(req.NamespacedName.Namespace),
	//	client.MatchingLabels{key: value},
	//	client.MatchingFields{"status.phase": "Running"},
	//}
	//err := r.List(ctx, podList, opts...)
	//if err != nil {
	//	return err
	//}
	//items := podList.Items
	//for _, pod := range items {
	//	deletePod(pod, r, ctx)
	//}
	pod := &metav1.Pod{}
	opts := []client.DeleteAllOfOption{
		client.InNamespace(req.NamespacedName.Namespace),
		client.MatchingLabels(podSelector),
		client.MatchingFields{"status.phase": "Running"},
		client.GracePeriodSeconds(9),
	}

	err := r.DeleteAllOf(ctx, pod, opts...)
	return err

}

// SetupWithManager sets up the controller with the Manager.
func (r *ConfigWatcherReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).For(&metav1.ConfigMap{}).WithEventFilter(predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool {
			return true
		},
	}).Complete(r)
	//For(&tutorialsv1.ConfigWatcher{}).Complete(r)

}
