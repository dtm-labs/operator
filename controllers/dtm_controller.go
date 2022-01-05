/*
Copyright 2022.

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
	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	dtmappv1 "github.com/dtm-labs/operator/api/v1"
)

// DtmReconciler reconciles a Dtm object
type DtmReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=app.dtm.hub,resources=dtms,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=app.dtm.hub,resources=dtms/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=app.dtm.hub,resources=dtms/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// the Dtm object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.10.0/pkg/reconcile
func (r *DtmReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// 检查 CR 是否存在
	dtm := &dtmappv1.Dtm{}
	if err := r.Get(ctx, req.NamespacedName, dtm); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	logger.Info("ensure dtm configmap")
	cm := r.GetDtmConfigMap(dtm)
	result, err := r.ensureConfigMap(ctx, cm)
	if err != nil {
		return *result, err
	}

	logger.Info("ensure dtm deployment")
	deploy := r.GetDtmDeployment(dtm)
	result, err = r.ensureDeployment(ctx, deploy)
	if err != nil {
		return *result, err
	}

	logger.Info("ensure dtm service")
	svc := r.GetDtmService(dtm)
	result, err = r.ensureService(ctx, svc)
	if err != nil {
		return *result, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DtmReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&dtmappv1.Dtm{}).
		Owns(&appv1.Deployment{}).
		Owns(&corev1.Service{}).
		Owns(&corev1.ConfigMap{}).
		Complete(r)
}
