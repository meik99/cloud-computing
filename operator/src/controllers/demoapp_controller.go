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
	"github.com/go-logr/logr"
	"github.com/meik99/cloud-computing/operator/src/demoapp"
	"github.com/pkg/errors"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	cloudcomputingv1alpha1 "github.com/meik99/cloud-computing/operator/api/v1alpha1"
)

var logger logr.Logger

// DemoAppReconciler reconciles a DemoApp object
type DemoAppReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	ctx      context.Context
	req      ctrl.Request
	instance *cloudcomputingv1alpha1.DemoApp
}

// SetupWithManager sets up the controller with the Manager.
func (r *DemoAppReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&cloudcomputingv1alpha1.DemoApp{}).
		Complete(r)
}

//+kubebuilder:rbac:groups=cloud-computing.rynkbit.com,resources=demoapps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=cloud-computing.rynkbit.com,resources=demoapps/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=cloud-computing.rynkbit.com,resources=demoapps/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the DemoApp object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.10.0/pkg/reconcile
func (r *DemoAppReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger = log.FromContext(ctx)
	r.ctx = ctx
	r.req = req

	var instance cloudcomputingv1alpha1.DemoApp
	if err := r.Client.Get(ctx, req.NamespacedName, &instance); err != nil {
		if k8serrors.IsNotFound(err) {
			return ctrl.Result{RequeueAfter: 5 * time.Minute}, nil
		}
		return ctrl.Result{}, errors.WithStack(err)
	}

	r.instance = &instance

	if err := r.reconcileDemoApp(); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{RequeueAfter: 5 * time.Minute}, nil
}

func (r *DemoAppReconciler) reconcileDemoApp() error {
	logger.Info("reconciling demo app", "name", r.req.Name, "namespace", r.req.Namespace)
	if err := r.reconcileStatefulSet(); err != nil {
		return errors.WithStack(err)
	}

	if err := r.reconcileService(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r *DemoAppReconciler) reconcileStatefulSet() error {
	logger.Info("reconciling stateful set")

	var statefulSet v1.StatefulSet
	err := r.Client.Get(r.ctx, client.ObjectKey{Name: r.instance.Spec.Name, Namespace: r.instance.Namespace}, &statefulSet)

	if err != nil {
		if k8serrors.IsNotFound(err) {
			return r.createStatefulSet()
		}
		return errors.WithStack(err)
	}

	return r.updateStatefulSet(statefulSet)
}

func (r *DemoAppReconciler) createStatefulSet() error {
	logger.Info("creating stateful set")
	statefulSet := demoapp.NewDemoApp(r.instance.Spec.Name, r.req.Namespace).
		CreateDesiredStatefulSet()

	err := controllerutil.SetControllerReference(r.instance, &statefulSet, r.Scheme)
	if err != nil {
		return errors.WithStack(err)
	}

	err = r.Client.Create(r.ctx, &statefulSet)
	return errors.WithStack(err)
}

func (r *DemoAppReconciler) updateStatefulSet(currentStatefulSet v1.StatefulSet) error {
	logger.Info("updating stateful set")
	statefulSet := demoapp.NewDemoApp(r.instance.Spec.Name, r.req.Namespace).
		CreateDesiredStatefulSet()

	if statefulSet.Annotations[demoapp.AnnotationHash] == currentStatefulSet.Annotations[demoapp.AnnotationHash] {
		logger.Info("stateful set already up to date")
		return nil
	}

	err := r.Client.Update(r.ctx, &statefulSet)
	return errors.WithStack(err)
}

func (r *DemoAppReconciler) reconcileService() error {
	logger.Info("reconciling service")

	var service corev1.Service
	err := r.Client.Get(r.ctx, client.ObjectKey{Name: r.instance.Spec.Name, Namespace: r.instance.Namespace}, &service)

	if err != nil {
		if k8serrors.IsNotFound(err) {
			return r.createService()
		}
		return errors.WithStack(err)
	}

	return r.updateService(service)
}

func (r *DemoAppReconciler) createService() error {
	logger.Info("creating service")
	service := demoapp.NewDemoApp(r.instance.Spec.Name, r.req.Namespace).
		CreateDesiredService()

	err := controllerutil.SetControllerReference(r.instance, &service, r.Scheme)
	if err != nil {
		return errors.WithStack(err)
	}

	err = r.Client.Create(r.ctx, &service)
	return errors.WithStack(err)
}

func (r *DemoAppReconciler) updateService(currentService corev1.Service) error {
	logger.Info("updating service")
	service := demoapp.NewDemoApp(r.instance.Spec.Name, r.req.Namespace).
		CreateDesiredService()

	if service.Annotations[demoapp.AnnotationHash] == currentService.Annotations[demoapp.AnnotationHash] {
		logger.Info("service already up to date")
		return nil
	}

	err := r.Client.Update(r.ctx, &service)
	return errors.WithStack(err)
}
