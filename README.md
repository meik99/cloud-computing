# Project for Cloud Computing

This repository contains all files related to a project for a course about cloud computing.
See [PROPOSAL.md](PROPOSAL.md) for the proposal for this project. 

# Todo

- [x] Deloy Concourse
- [x] Create demo application
- [x] Create test pipeline
- [x] Run pipeline on repository change
- [x] Build demo-app Docker image
- [x] Deploy demo-app
- [x] Create Operator project
- [ ] Redeploy on pipeline success

# Documentation

## Operator

An operator is a piece of software, designed to run inside a Kubernetes cluster.
It can be used to run arbitrary code on certain events, for example, when a certain kind of resource is created, updated or deleted.
Furthermore, it can interact with a cluster and, e.g., create resources or modify them.

The operator included in this project has two objectives.
First, making sure the demo-app and its service is deployed.
Secondly, keeping the application up to date with the latest image available for it.

### Using the Operator

To use the operator, it must first be deployed and configured.
Deploying the operator can be done by applying its manifest by using the command below.

```shell
kubectl apply -f https://raw.githubusercontent.com/meik99/cloud-computing/main/operator/config/deploy/deployment.yaml
```

For the operator to recognize that it should deploy the demo-app, a custom resource must be configured and applied.
Note that the custom resource definition is already included in the deployment file used above.
The sample below creates a DemoApp resource in the demo-app namespace, called demo-app.
```shell
kubectl apply -f https://raw.githubusercontent.com/meik99/cloud-computing/main/operator/config/samples/sample.yaml
```

### Reconciliation

After the operator and the custom resource are deployed, the operator starts a routine called a "Reconcile" or a reconciliation.
First retrieves the instance of the DemoApp that triggered the event.
A client instance is used to interact with the cluster the operator runs in.
The namespace the request comes from can be found in the request object that is passed to the reconciliation.
If the instance cannot be found, be it due to the reconciliation being fired on a delete event or due to caching issues, the operator returns no error and restarts the loop after five minutes.
Should it encounter a different error, it is returned.
It is then automatically logged by the operator's framework.
The framework then triggers a new reconciliation immediately.

```go
// operator/src/controller/dempapp_controller.go:77
var instance cloudcomputingv1alpha1.DemoApp
if err := r.Client.Get(ctx, req.NamespacedName, &instance); err != nil {
    if k8serrors.IsNotFound(err) {
        return ctrl.Result{RequeueAfter: 5 * time.Minute}, nil
    }
    return ctrl.Result{}, errors.WithStack(err)
}
```

The DemoApp type is defined in code as well.
Its specification can be found in `operator/api/v1alpha1/demoapp_types.go`.
Most of the code was generated by the [Operator SDK](https://sdk.operatorframework.io/).
The part that was added was the `Name` property.
This property can be used to define the name of the stateful set and service created for the demo-app.
It is not necessary for the deployment of the application, but has been added to show how an operator can be configured.

```go
// operator/api/v1alpha1/demoapp_types.go:26

// DemoAppSpec defines the desired state of DemoApp
type DemoAppSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of DemoApp. Edit demoapp_types.go to remove/update
	//Foo string `json:"foo,omitempty"`

	Name string `json:"name,omitempty"`
}
```

After the instance is retrieved, the state of the cluster is compared to the state that is required.
To do so, the operator generates the stateful set that is required.
For easier comparison later, a hash is produced from the stateful set and stored as an annotation in it.
Then, the operator tries to get the stateful set of the demo-app from the cluster.
If it does not find one, the generated stateful set is applied.

If a stateful set exists, the hash in the annotation is compared to the generated one.
If they do not match, the stateful set is updated using the generated stateful set.
All of these steps are then repeated for the service that is needed for demo-app to be accessible from external clients.

```go
// operator/src/demoapp/statefulset.go:23
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

// operator/src/controller/demoapp_controller.go:177
if statefulSet.Annotations[demoapp.AnnotationHash] == currentStatefulSet.Annotations[demoapp.AnnotationHash] {
    logger.Info("stateful set already up to date")
    return nil
}
```

Furthermore, the images used by the pods that were spawned by the stateful set are compared to the latest available.
Before a container in a pod is created, the image needed is downloaded if it does not exist in the cluster.
The version of the image that should be used can be defined using either tags or labels, which is the most common way, or a hash value commonly referenced as "image digest" .
For example, an image can be specified by `image-name:tag`, where the tag references a specific version of the image, or with `image-name@sha256:<digest>`.

When the image is applied to a container, the image ID is also stored within the container status.
This image ID can then be used to compare images with the same tag, because although two different images may share a tag, e.g. "latest", they have different digests.
The assumption then is, that the image on DockerHub is the image that should be used by the demo-app deployment.
Finally, in order to compare the images, the image digest is requested from the registry api.
This digest is then compared to the on in the container status of the pods of the stateful set.
If the image ID of a pod does not contain the digest, the pod is considered outdated.

Outdated pods are collected and then deleted.
The stateful set makes sure its set replicas of pods exists and recreates them.
Since the image pull policy is set to `Always`, Kubernetes checks whether a new image is available.
The latest image is then downloaded and applied to the new pod.
In summary, deleting the pod triggers an automatic update of the image by the cluster.

