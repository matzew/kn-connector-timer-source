/*
Copyright 2019 The Knative Authors

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

package sample

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	v1 "knative.dev/eventing/pkg/apis/sources/v1"
	"knative.dev/kamelet-source/pkg/reconciler/sample/resources"
	"knative.dev/pkg/controller"
	"knative.dev/pkg/logging"

	appsv1listers "k8s.io/client-go/listers/apps/v1"

	corev1listers "k8s.io/client-go/listers/core/v1"

	eventingclientset "knative.dev/eventing/pkg/client/clientset/versioned"

	pkgreconciler "knative.dev/pkg/reconciler"
	"knative.dev/pkg/resolver"

	listers "knative.dev/eventing/pkg/client/listers/sources/v1"

	reconcilersource "knative.dev/eventing/pkg/reconciler/source"

	"knative.dev/kamelet-source/pkg/apis/samples/v1alpha1"
	reconcilerkameletsource "knative.dev/kamelet-source/pkg/client/injection/reconciler/samples/v1alpha1/kameletsource"
	"knative.dev/kamelet-source/pkg/reconciler"
)

const (
	// Name of the corev1.Events emitted from the reconciliation process
	sourceReconciled   = "KameletSourceReconciled"
	deploymentCreated  = "KameletSourceDeploymentCreated"
	serciceCreated     = "KameletSourceServiceCreated"
	deploymentUpdated  = "KameletSourceDeploymentUpdated"
	sinkBindingCreated = "KameletSourceSinkBindingCreated"
	sinkBindingUpdated = "KameletSourceSinkBindingUpdated"
)

// Reconciler reconciles a KameletSource object
type Reconciler struct {
	dr *reconciler.DeploymentReconciler

	kubeClientSet     kubernetes.Interface
	eventingclientset eventingclientset.Interface

	sinkResolver *resolver.URIResolver

	sinkBindingLister listers.SinkBindingLister
	deploymentLister  appsv1listers.DeploymentLister
	serviceLister     corev1listers.ServiceLister

	configAccessor reconcilersource.ConfigAccessor
}

// Check that our Reconciler implements Interface
var _ reconcilerkameletsource.Interface = (*Reconciler)(nil)

// ReconcileKind implements Interface.ReconcileKind.
func (r *Reconciler) ReconcileKind(ctx context.Context, source *v1alpha1.KameletSource) pkgreconciler.Event {

	_, err := r.reconcileSinkBinding(ctx, source)
	if err != nil {
		logging.FromContext(ctx).Errorw("Error reconciling SinkBinding", zap.Error(err))
		return err
	}

	_, err = r.reconcileDeployment(ctx, source)
	if err != nil {
		logging.FromContext(ctx).Errorw("Error reconciling ReceiveAdapter", zap.Error(err))
		return err
	}

	_, err = r.reconcileService(ctx, source)
	if err != nil {
		logging.FromContext(ctx).Errorw("Error reconciling ReceiveAdapter", zap.Error(err))
		return err
	}

	return nil

}
func (r *Reconciler) reconcileSinkBinding(ctx context.Context, source *v1alpha1.KameletSource) (*v1.SinkBinding, error) {

	expected := resources.MakeSinkBinding(source)

	sb, err := r.sinkBindingLister.SinkBindings(source.Namespace).Get(expected.Name)
	if apierrors.IsNotFound(err) {
		sb, err = r.eventingclientset.SourcesV1().SinkBindings(source.Namespace).Create(ctx, expected, metav1.CreateOptions{})
		if err != nil {
			return nil, fmt.Errorf("creating new SinkBinding: %v", err)
		}
		controller.GetEventRecorder(ctx).Eventf(source, corev1.EventTypeNormal, sinkBindingCreated, "SinkBinding created %q", sb.Name)
	} else if err != nil {
		return nil, fmt.Errorf("getting SinkBinding: %v", err)
	} else if !metav1.IsControlledBy(sb, source) {
		return nil, fmt.Errorf("SinkBinding %q is not owned by KameletSource %q", sb.Name, source.Name)
	} else if r.sinkBindingSpecChanged(&sb.Spec, &expected.Spec) {
		sb.Spec = expected.Spec
		sb, err = r.eventingclientset.SourcesV1().SinkBindings(source.Namespace).Update(ctx, sb, metav1.UpdateOptions{})
		if err != nil {
			return nil, fmt.Errorf("updating SinkBinding: %v", err)
		}
		controller.GetEventRecorder(ctx).Eventf(source, corev1.EventTypeNormal, sinkBindingUpdated, "SinkBinding updated %q", sb.Name)
	} else {
		logging.FromContext(ctx).Debugw("Reusing existing SinkBinding", zap.Any("SinkBinding", sb))
	}

	//	source.Status.PropagateSinkBindingStatus(&sb.Status)
	return sb, nil

}

func (r *Reconciler) reconcileDeployment(ctx context.Context, source *v1alpha1.KameletSource) (*appsv1.Deployment, error) {
	args := &resources.ReceiveAdapterArgs{
		EventSource:    source.Namespace + "/" + source.Name,
		Source:         source,
		Labels:         resources.Labels(source.Name),
		AdditionalEnvs: r.configAccessor.ToEnvVars(), // Grab config envs for tracing/logging/metrics
	}

	addr, err := r.sinkResolver.AddressableFromDestinationV1(ctx, source.Spec.Sink, source)
	if err != nil {
		logging.FromContext(ctx).Errorf("Failed to get Addressable from Destination: %w", err)
		return nil, err
	}
	args.Address = addr.URL.String()

	expected := resources.MakeDeployment(args)

	ra, err := r.deploymentLister.Deployments(expected.Namespace).Get(expected.Name)
	if apierrors.IsNotFound(err) {
		ra, err = r.kubeClientSet.AppsV1().Deployments(expected.Namespace).Create(ctx, expected, metav1.CreateOptions{})
		if err != nil {
			return nil, fmt.Errorf("creating new Deployment: %v", err)
		}
		controller.GetEventRecorder(ctx).Eventf(source, corev1.EventTypeNormal, deploymentCreated, "Deployment created %q", ra.Name)
	} else if r.podSpecChanged(&ra.Spec.Template.Spec, &expected.Spec.Template.Spec) {
		ra.Spec.Template.Spec = expected.Spec.Template.Spec
		ra, err = r.kubeClientSet.AppsV1().Deployments(expected.Namespace).Update(ctx, ra, metav1.UpdateOptions{})
		if err != nil {
			return nil, fmt.Errorf("updating Deployment: %v", err)
		}
		controller.GetEventRecorder(ctx).Eventf(source, corev1.EventTypeNormal, deploymentUpdated, "Deployment updated %q", ra.Name)
	} else {
		logging.FromContext(ctx).Debugw("Reusing existing Deployment", zap.Any("Deployment", ra))
	}

	return ra, nil
}

func (r *Reconciler) reconcileService(ctx context.Context, source *v1alpha1.KameletSource) (*corev1.Service, error) {

	args := &resources.ReceiveAdapterArgs{
		EventSource:    source.Namespace + "/" + source.Name,
		Source:         source,
		Labels:         resources.Labels(source.Name),
		AdditionalEnvs: r.configAccessor.ToEnvVars(), // Grab config envs for tracing/logging/metrics
	}

	expected := resources.NewK8sService(args)
	svc, err := r.serviceLister.Services(expected.Namespace).Get(expected.Name)

	if apierrors.IsNotFound(err) {

		svc, err = r.kubeClientSet.CoreV1().Services(expected.Namespace).Create(ctx, expected, metav1.CreateOptions{})
		if err != nil {
			return nil, fmt.Errorf("creating new Service: %v", err)
		}
		controller.GetEventRecorder(ctx).Eventf(source, corev1.EventTypeNormal, serciceCreated, "Service created %q", svc.Name)

	} else {
		logging.FromContext(ctx).Debugw("Reusing existing Service", zap.Any("Service", svc))
	}

	return svc, nil

}

func (r *Reconciler) podSpecChanged(have *corev1.PodSpec, want *corev1.PodSpec) bool {
	// TODO this won't work, SinkBinding messes with this. n3wscott working on a fix.
	return !equality.Semantic.DeepDerivative(want, have)
}

func (r *Reconciler) sinkBindingSpecChanged(have *v1.SinkBindingSpec, want *v1.SinkBindingSpec) bool {
	return !equality.Semantic.DeepDerivative(want, have)
}
