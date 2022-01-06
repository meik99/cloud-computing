package demoapp

import (
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
)

const (
	defaultName            = "demo-app"
	defaultImage           = "meik99/cloud-computing:dev"
	defaultImagePullPolicy = "Always"
	defaultReplicas        = 3
	defaultPort            = 80
	defaultPortName        = "web"

	AnnotationHash = "cloudcomputing.demoapp.hash"
)

var selectorLabels = map[string]string{"app": "demo-app"}

func (demoApp *DemoApp) CreateDesiredStatefulSet() appsv1.StatefulSet {
	result := appsv1.StatefulSet{
		ObjectMeta: demoApp.createObjectMeta(),
		Spec:       demoApp.createSpec(),
	}

	result.Annotations = map[string]string{
		AnnotationHash: calculateHash(result),
	}

	return result
}

func (demoApp *DemoApp) createObjectMeta() metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:      demoApp.Name,
		Namespace: demoApp.Namespace,
	}
}

func (demoApp *DemoApp) createSpec() appsv1.StatefulSetSpec {
	return appsv1.StatefulSetSpec{
		Selector:    demoApp.createSelector(),
		ServiceName: demoApp.Name,
		Replicas:    pointer.Int32(defaultReplicas),
		Template:    demoApp.createTemplate(),
	}
}

func (demoApp *DemoApp) createSelector() *metav1.LabelSelector {
	return &metav1.LabelSelector{
		MatchLabels: selectorLabels,
	}
}

func (demoApp *DemoApp) createTemplate() v1.PodTemplateSpec {
	return v1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Labels: selectorLabels,
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				demoApp.createContainer(),
			},
		},
	}
}

func (demoApp *DemoApp) createContainer() v1.Container {
	return v1.Container{
		Name:            demoApp.Name,
		Image:           defaultImage,
		ImagePullPolicy: defaultImagePullPolicy,
		Ports: []v1.ContainerPort{
			{ContainerPort: defaultPort, Name: defaultPortName},
		},
	}
}
