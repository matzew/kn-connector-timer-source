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

// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
	v1alpha1 "knative.dev/kamelet-source/pkg/apis/samples/v1alpha1"
)

// KameletSourceLister helps list KameletSources.
// All objects returned here must be treated as read-only.
type KameletSourceLister interface {
	// List lists all KameletSources in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.KameletSource, err error)
	// KameletSources returns an object that can list and get KameletSources.
	KameletSources(namespace string) KameletSourceNamespaceLister
	KameletSourceListerExpansion
}

// kameletSourceLister implements the KameletSourceLister interface.
type kameletSourceLister struct {
	indexer cache.Indexer
}

// NewKameletSourceLister returns a new KameletSourceLister.
func NewKameletSourceLister(indexer cache.Indexer) KameletSourceLister {
	return &kameletSourceLister{indexer: indexer}
}

// List lists all KameletSources in the indexer.
func (s *kameletSourceLister) List(selector labels.Selector) (ret []*v1alpha1.KameletSource, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.KameletSource))
	})
	return ret, err
}

// KameletSources returns an object that can list and get KameletSources.
func (s *kameletSourceLister) KameletSources(namespace string) KameletSourceNamespaceLister {
	return kameletSourceNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// KameletSourceNamespaceLister helps list and get KameletSources.
// All objects returned here must be treated as read-only.
type KameletSourceNamespaceLister interface {
	// List lists all KameletSources in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.KameletSource, err error)
	// Get retrieves the KameletSource from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1alpha1.KameletSource, error)
	KameletSourceNamespaceListerExpansion
}

// kameletSourceNamespaceLister implements the KameletSourceNamespaceLister
// interface.
type kameletSourceNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all KameletSources in the indexer for a given namespace.
func (s kameletSourceNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.KameletSource, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.KameletSource))
	})
	return ret, err
}

// Get retrieves the KameletSource from the indexer for a given namespace and name.
func (s kameletSourceNamespaceLister) Get(name string) (*v1alpha1.KameletSource, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("kameletsource"), name)
	}
	return obj.(*v1alpha1.KameletSource), nil
}