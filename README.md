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

## Deploy Operator

First, deploy the operator itself.
```shell
kubectl apply -f https://raw.githubusercontent.com/meik99/cloud-computing/main/operator/config/deploy/deployment.yaml
```

Then, deploy a CustomResource.
For testing purposes, you can use this sample.
```shell
kubectl apply -f https://raw.githubusercontent.com/meik99/cloud-computing/main/operator/config/samples/sample.yaml
```