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
	"knative.dev/pkg/client/injection/kube/informers/core/v1/service"

	reconcilersource "knative.dev/eventing/pkg/reconciler/source"

	"knative.dev/kamelet-source/pkg/apis/samples/v1alpha1"

	"github.com/kelseyhightower/envconfig"
	"k8s.io/client-go/tools/cache"

	"knative.dev/pkg/configmap"
	"knative.dev/pkg/controller"
	"knative.dev/pkg/logging"
	"knative.dev/pkg/resolver"

	eventingclient "knative.dev/eventing/pkg/client/injection/client"
	sinkbindinginformer "knative.dev/eventing/pkg/client/injection/informers/sources/v1/sinkbinding"

	"knative.dev/kamelet-source/pkg/reconciler"

	kameletsourceinformer "knative.dev/kamelet-source/pkg/client/injection/informers/samples/v1alpha1/kameletsource"
	"knative.dev/kamelet-source/pkg/client/injection/reconciler/samples/v1alpha1/kameletsource"
	kubeclient "knative.dev/pkg/client/injection/kube/client"
	deploymentinformer "knative.dev/pkg/client/injection/kube/informers/apps/v1/deployment"
)

// NewController initializes the controller and is called by the generated code
// Registers event handlers to enqueue events
func NewController(
	ctx context.Context,
	cmw configmap.Watcher,
) *controller.Impl {

	kubeClient := kubeclient.Get(ctx)
	eventingClient := eventingclient.Get(ctx)
	sinkbindingInformer := sinkbindinginformer.Get(ctx)

	serviceInformer := service.Get(ctx)
	deploymentInformer := deploymentinformer.Get(ctx)

	kameletSourceInformer := kameletsourceinformer.Get(ctx)

	r := &Reconciler{

		kubeClientSet:     kubeClient,
		eventingclientset: eventingClient,
		sinkBindingLister: sinkbindingInformer.Lister(),
		deploymentLister:  deploymentInformer.Lister(),
		serviceLister:     serviceInformer.Lister(),

		dr:             &reconciler.DeploymentReconciler{KubeClientSet: kubeclient.Get(ctx)},
		configAccessor: reconcilersource.WatchConfigurations(ctx, "kamelet-source", cmw),
	}
	if err := envconfig.Process("", r); err != nil {
		logging.FromContext(ctx).Panicf("required environment variable is not defined: %v", err)
	}

	impl := kameletsource.NewImpl(ctx, r)

	r.sinkResolver = resolver.NewURIResolverFromTracker(ctx, impl.Tracker)

	kameletSourceInformer.Informer().AddEventHandler(controller.HandleAll(impl.Enqueue))

	deploymentInformer.Informer().AddEventHandler(cache.FilteringResourceEventHandler{
		FilterFunc: controller.FilterController(&v1alpha1.KameletSource{}),
		Handler:    controller.HandleAll(impl.EnqueueControllerOf),
	})

	return impl
}
