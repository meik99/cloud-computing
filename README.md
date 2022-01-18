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
- [x] Redeploy on pipeline success

# Documentation

## Overview

What if you could just push a change into a git repository and have everything from tests to deployment in a cluster run automatically without any needed supervision? Our project shows exactly how to do that.

Using an Angular demo app, called demo-app, that simply displays something and includes some tests, we use Concourse (an open source CI/CD system) to run those tests.
If they succeed, a job uploads a new docker image of the demo app to Docker Hub. In addition, a Discord notification informs us about certain events (e.g. test failure).
Finally, an operator, watching over the demo app deployment and its service, updates it if a new version is available.

All of the above (demo app, Concourse and the operator) run on a Kubernetes cluster. 
All images are hosted on and uploaded to Docker Hub.

## Demo App

This is just a demo app built with Angular. 
It displays some text and contains basic unit tests.

The [demo-app-statefulset.yaml](demo-app-statefulset.yaml) configuration defines a stateful set with three pods and a service with type NodePort.

We use a stateful set instead of a replica set because of the nice predictable and readable names (e.g. demo-app-0, demo-app-1 and demo-app-2).
The service type is `NodePort`, since the cluster we use is deployed on a self-managed root server.
Therefore, it has no cloud provided load balancer, but is made available externally using an Nginx reverse proxy.

## Concourse

Concourse is a CI/CD system which allows creating pipelines consisting of jobs.
Each job runs tasks which run, e.g., scripts, in their own containers. 
We chose Concourse because it is a common CI/CD system, uses YAML to configure its pipelines, and does what it should: Running automated jobs (which includes tests, of course).

### Deployment Customization

The starting point for the Concourse deployment is found in [kustomization.yaml](kustomization.yaml). 
The configuration is split into manageable parts, and therefore this file contains just the reference to the folders postgresql and concourse, which, again, contain kustomization files, declaring the files needed for the deployment.

```yaml
# kustomization.yaml
resources:
  - postgresql
  - concourse
```

This pattern of splitting the configuration files is continued in those folders. The postgresql folder contains the definition of the service and stateful set for the database Concourse uses. The service gets a ClusterIP, because it is only accessed from inside the cluster. The stateful set contains one interesting part: the definition for a persistent volume.

```yaml
# postgresql/statefulset.yaml:31
resources:
  volumeClaimTemplates:
  - metadata:
      name: data
      namespace: concourse
    spec:
      accessModes: [ "ReadWriteOnce" ]
      resources:
        requests:
          storage: 5Gi
```

The secrets template (as all other templates in this repository) contains a list of the needed secrets to run this projct.

In the concourse folder, the definition for the actual Concourse deployment can be found. There are services and stateful sets for both the web page and the workers. There is one web pod and there are three worker pods. The web service gets a NodePort to enable access through the Nginx reverse proxy, the workers with their ClusterIP are only accessable from inside the cluster.

One detail worth discussing is in the definition of the stateful set for the worker. 
The workers are privileged to allow them access to the Docker process of the host system.
They need this access to start new containers in their pod, since Concourse runs tasks in containers and would not be allowed to if not privileged.

```yaml
# concourse/statefulset-worker.yaml:36
        securityContext:
          privileged: true
```

### Pipeline

The pipeline folder contains the definition for the pipeline and the two jobs we want to run. 
The file [pipeline.yaml](pipeline/pipeline.yaml) contains the definition for three resources used by the jobs: our GitHub repository, the Docker image of the demo app on Docker Hub and a webhook-notification.
The two jobs are testing the demo app, and building and pushing the docker image to Docker Hub. In addition, notifications will be sent to a Discord channel on some events (e.g. test failure or docker push was successful).
The test is triggered on a change in the repository, and the build and push job only runs after a successful test. 
The jobs themselves are found in their respective sub-folders.

The test task is described in [test-demo-app/task.yaml](pipeline/tasks/test-demo-app/task.yaml).

```yaml
# pipeline/tasks/test-demo-app/task.yaml
platform: linux
image_resource: 
  type: registry-image
  source: 
    repository: node
    tag: latest


inputs:
  - name: repository

run:
  path: repository/pipeline/tasks/test-demo-app/task.sh
```

The task.sh installs the demo app and runs the tests.

The build and push task described in [build-and-push/task.yaml](pipeline/tasks/build-and-push/task.yaml) runs a Docker image, that builds and pushes the demo app to Docker Hub.
This image is already available and is specifically made to allow the use of Docker commands inside a Docker container.

After the image has been pushed to Docker Hub successfully, a notification is sent to a Discord channel. For this purpose we added a webhook-notification resource to the pipeline configuration. Since the required resource type is not directly available in Concourse, a new resource type is declared at the beginning:

```yaml
resource_types:
  - name: webhook-notification
    type: registry-image
    source:
      repository: flavorjones/webhook-notification-resource
      tag: latest
```

This declaration refers to an image on Docker Hub that is used for this resource type (we use this implementation: https://github.com/flavorjones/webhook-notification-resource). In the pipeline configuration, we declare then a Discord resource that references to this type. Here, the webhook from Discord is defined which has to be provided in the form of a secret. To trigger a notification on success of the "build-and-push" job, the Discord resource is used as follows:

```yaml
  - name: build-and-push
    plan:
# ...
      - put: demo-app
        params:
          image: image/image.tar
    on_success:
      put: discord
      params:
        status: "succeeded"
```

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

# Lessons Learned

Concourse is really versatile and useful. 
Running tests in their own containers guarantees high isolation.
The configuration is all in files and not through some UI, which allows for more collaborative work.
It does its own internal resource management, which allows, e.g., automatic polling for changes. 
This means some initial investment and overhead on one hand, but increased maintainability and usability on the other.

This project shows how to do everything in a cluster, even the management of said cluster. 
It shows how cloud computing can be complex, but also versatile and rewarding.

Releases of updates can be painless and almost invisible. 
With good test coverage and the right tooling an update can be just a push to a git repository. 
Removing the need for manual tests, builds and deployments.

Operators can aid with custom tooling. 
They offer a lot of freedom in customizing the workflows in a cluster, at the expense of time investment and needed domain knowledge.

The downside is, of course, the high initial investment of time and the domain knowledge needed to set up such a system. 
Having the system working is wonderful, but setting it up is its own little project in itself.
This indicates the need for specific DevOps team in enterprise environments.
