package resources

import (
	appsv1 "k8s.io/api/apps/v1"
	"knative.dev/kamelet-source/pkg/apis/samples/v1alpha1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "knative.dev/eventing/pkg/apis/sources/v1"
	duckv1 "knative.dev/pkg/apis/duck/v1"
	"knative.dev/pkg/kmeta"
	"knative.dev/pkg/tracker"
)

var subjectGVK = appsv1.SchemeGroupVersion.WithKind("Deployment")

func MakeSinkBinding(source *v1alpha1.KameletSource) *v1.SinkBinding {
	subjectAPIVersion, subjectKind := subjectGVK.ToAPIVersionAndKind()

	sb := &v1.SinkBinding{
		ObjectMeta: metav1.ObjectMeta{
			OwnerReferences: []metav1.OwnerReference{
				*kmeta.NewControllerRef(source),
			},
			Name:      "kn-connector-source-timer",
			Namespace: source.Namespace,
		},
		Spec: v1.SinkBindingSpec{
			SourceSpec: source.Spec.SourceSpec,
			BindingSpec: duckv1.BindingSpec{
				Subject: tracker.Reference{
					APIVersion: subjectAPIVersion,
					Kind:       subjectKind,
					Namespace:  source.Namespace,
					Name:       source.Name,
				},
			},
		},
	}
	return sb
}
