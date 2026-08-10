/*
Copyright 2026 zhouli-95

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

package utils

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func CheckNodeResources(k8sClient kubernetes.Interface, resourceName string, num int64) error {
	nodes, err := k8sClient.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, node := range nodes.Items {
		if _, ok := node.Labels["node-role.kubernetes.io/control-plane"]; ok {
			continue
		}

		capacity, _ := node.Status.Capacity.Name(v1.ResourceName(resourceName), resource.DecimalSI).AsInt64()
		if capacity != num {
			return fmt.Errorf("capacity %d is not equal to expected %d", capacity, num)
		}

		allocatable, _ := node.Status.Allocatable.Name(v1.ResourceName(resourceName), resource.DecimalSI).AsInt64()
		if allocatable != num {
			return fmt.Errorf("allocatable %d is not equal to expected %d", allocatable, num)
		}

	}
	return nil
}

func CreatePodWithDevice(k8sClient kubernetes.Interface, ns, podName, deviceName, deviceCount string, cmds []string) (*corev1.Pod, error) {
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      podName,
			Namespace: ns,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:    "app",
					Image:   "ubuntu:24.04",
					Command: cmds,
					Resources: corev1.ResourceRequirements{
						Limits: corev1.ResourceList{
							corev1.ResourceName(deviceName): resource.MustParse(deviceCount),
						},
						Requests: corev1.ResourceList{
							corev1.ResourceName(deviceName): resource.MustParse(deviceCount),
						},
					},
				},
			},
			RestartPolicy: corev1.RestartPolicyNever,
		},
	}

	return k8sClient.CoreV1().Pods(ns).Create(context.TODO(), pod, metav1.CreateOptions{})
}
