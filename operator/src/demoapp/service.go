package demoapp

import v1 "k8s.io/api/core/v1"

func (demoApp *DemoApp) CreateDesiredService() v1.Service {
	result := v1.Service{
		ObjectMeta: demoApp.createObjectMeta(),
		Spec:       demoApp.createServiceSpec(),
	}

	result.Annotations = map[string]string{
		AnnotationHash: calculateHash(result),
	}

	return result
}

func (demoApp *DemoApp) createServiceSpec() v1.ServiceSpec {
	return v1.ServiceSpec{
		Selector: selectorLabels,
		Ports: []v1.ServicePort{
			{Port: defaultPort},
		},
		Type: v1.ServiceTypeNodePort,
	}
}
