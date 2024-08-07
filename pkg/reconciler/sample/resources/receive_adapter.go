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

package resources

import (
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knative.dev/pkg/kmeta"

	"knative.dev/kamelet-source/pkg/apis/samples/v1alpha1"
)

// ReceiveAdapterArgs are the arguments needed to create a Sample Source Receive Adapter.
// Every field is required.
type ReceiveAdapterArgs struct {
	//Image          string
	Labels         map[string]string
	Source         *v1alpha1.KameletSource
	EventSource    string
	AdditionalEnvs []corev1.EnvVar
	Address        string
}

// MakeDeployment generates (but does not insert into K8s) the Receive Adapter Deployment for
// Sample sources.
func MakeDeployment(args *ReceiveAdapterArgs) *v1.Deployment {
	replicas := int32(1)
	return &v1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: args.Source.Namespace,
			Name:      args.Source.Name,
			Labels:    args.Labels,
			OwnerReferences: []metav1.OwnerReference{
				*kmeta.NewControllerRef(args.Source),
			},
		},
		Spec: v1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: args.Labels,
			},
			Replicas: &replicas,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: args.Labels,
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: args.Source.Spec.ServiceAccountName,
					Containers: []corev1.Container{
						{
							Name:  args.Source.Name,
							Image: "quay.io/openshift-knative/kn-connector-source-timer:1.0-SNAPSHOT",
							Ports: []corev1.ContainerPort{{
								Name:          "http",
								ContainerPort: 8080,
							}},
							Env: append(
								makeEnv(args.Address, &args.Source.Spec),
								args.AdditionalEnvs...,
							),
						},
					},
				},
			},
		},
	}
}
func makeEnv(address string, spec *v1alpha1.KameletSourceSpec) []corev1.EnvVar {
	return []corev1.EnvVar{{
		Name:  "K_SINK",
		Value: address, //"http://broker-ingress.knative-eventing.svc.cluster.local/knative-samples/default",
	}, {
		Name:  "CAMEL_KAMELET_TIMER_SOURCE_MESSAGE",
		Value: spec.Text,
	}, {
		Name: "K_CE_OVERRIDES",
	}}
}
