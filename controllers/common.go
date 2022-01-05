package controllers

import (
	"context"
	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
)

func (r *DtmReconciler) ensureDeployment(ctx context.Context, deploy *appv1.Deployment) (*ctrl.Result, error) {
	found := &appv1.Deployment{}
	err := r.Client.Get(ctx, types.NamespacedName{
		Namespace: deploy.Namespace,
		Name:      deploy.Name,
	}, found)
	if err != nil {
		if errors.IsNotFound(err) {
			err := r.Client.Create(ctx, deploy)
			if err != nil {
				return &ctrl.Result{}, err
			}
		}
		return &ctrl.Result{}, err
	}
	return nil, nil
}

func (r *DtmReconciler) ensureService(ctx context.Context, svc *corev1.Service) (*ctrl.Result, error) {
	found := &corev1.Service{}
	err := r.Client.Get(ctx, types.NamespacedName{
		Namespace: svc.Namespace,
		Name:      svc.Name,
	}, found)
	if err != nil {
		if errors.IsNotFound(err) {
			err = r.Client.Create(ctx, svc)
			if err != nil {
				return &ctrl.Result{}, err
			}
		}
		return &ctrl.Result{}, err
	}

	return nil, nil
}

func (r *DtmReconciler) ensureConfigMap(ctx context.Context, cm *corev1.ConfigMap) (*ctrl.Result, error) {
	found := &corev1.ConfigMap{}
	err := r.Client.Get(ctx, types.NamespacedName{
		Namespace: cm.Namespace,
		Name:      cm.Name,
	}, found)
	if err != nil {
		if errors.IsNotFound(err) {
			err = r.Client.Create(ctx, cm)
			if err != nil {
				return &ctrl.Result{}, err
			}
		}
		return &ctrl.Result{}, err
	}

	return nil, nil
}
