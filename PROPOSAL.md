# Proposal - CI / CD Pipeline

The idea is to first create a simple demo program.
It does not really matter what the program is, e.g., it could be a simple calculator, as long as it has unit tests.
The repository is hosted on GitHub.

Then [Concourse](https://github.com/concourse/concourse-docker) is deployed on a cluster.
Concourse is a CI/CD system which allows creating pipelines consisting of jobs, where each job runs in its own container.

A pipeline than runs the unit tests whenever a push to the main branch occurs.
If the unit tests run through, it is built and published.
Depending on the demo app, as artifact to GitHub or as Docker image to DockerHub.

If it is deployable, i.e. the demo app is a service or webapp, it is then also deployed on the cluster.
Likely in a different namespace.

As an extension, an [operator](https://sdk.operatorframework.io/) can then be built, which checks the image hashes of the demo app.
If the hashes of the currently running container the one on DockerHub differ, it recreates the pod, updating the application.

## Areas of concern

* Demo Application
* Concourse Deployment
* Concourse Pipeline
  * Listen to repository changes
  * Unit test
  * Build
  * Delivery
* Operator

## Responsibilities

From the above areas, the following responsibilities are assigned.

[Daniel Klösler](https://github.com/Ethlaron):
* Demo Application
* Listen to repository changes
* Run unit test

[Christoph Pargfrieder](https://github.com/ChristophPargfrieder):

* Run build
* Automate Delivery

[Michael Rynkiewicz](https://github.com/meik99): 
* Concourse Deployment
* Operator

## Milestones

The following milestones can be defined for this project.

### Milestone 1

Concourse is deployed on a cluster and is reachable either from port forwarding or the web.

### Milestone 2

A concourse pipeline handles testing, building and delivering the demo application.

### Milestone 3

An operator deploys and updates the demo application as necessary.