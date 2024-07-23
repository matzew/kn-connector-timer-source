/*
Copyright 2020 The Knative Authors

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

// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	"context"
	time "time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
	samplesv1alpha1 "knative.dev/kamelet-source/pkg/apis/samples/v1alpha1"
	versioned "knative.dev/kamelet-source/pkg/client/clientset/versioned"
	internalinterfaces "knative.dev/kamelet-source/pkg/client/informers/externalversions/internalinterfaces"
	v1alpha1 "knative.dev/kamelet-source/pkg/client/listers/samples/v1alpha1"
)

// KameletSourceInformer provides access to a shared informer and lister for
// KameletSources.
type KameletSourceInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha1.KameletSourceLister
}

type kameletSourceInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewKameletSourceInformer constructs a new informer for KameletSource type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewKameletSourceInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredKameletSourceInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredKameletSourceInformer constructs a new informer for KameletSource type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredKameletSourceInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.SamplesV1alpha1().KameletSources(namespace).List(context.TODO(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.SamplesV1alpha1().KameletSources(namespace).Watch(context.TODO(), options)
			},
		},
		&samplesv1alpha1.KameletSource{},
		resyncPeriod,
		indexers,
	)
}

func (f *kameletSourceInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredKameletSourceInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *kameletSourceInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&samplesv1alpha1.KameletSource{}, f.defaultInformer)
}

func (f *kameletSourceInformer) Lister() v1alpha1.KameletSourceLister {
	return v1alpha1.NewKameletSourceLister(f.Informer().GetIndexer())
}
