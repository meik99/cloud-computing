package demoapp

import (
	"github.com/meik99/cloud-computing/operator/src/docker"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"strings"
)

func GetOutdatedPods(podList corev1.PodList) (corev1.PodList, error) {
	latestDigest, err := docker.GetDigest()

	if err != nil {
		return corev1.PodList{}, errors.WithStack(err)
	}

	var outdatedPods corev1.PodList
	for _, pod := range podList.Items {
		for _, status := range pod.Status.ContainerStatuses {
			if !strings.Contains(status.ImageID, latestDigest) {
				outdatedPods.Items = append(outdatedPods.Items, pod)
			}
		}
	}

	return outdatedPods, nil
}
