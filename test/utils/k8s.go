package utils

import (
	"context"
	"time"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
)

func WaitForDaemonSetReady(k8sClient kubernetes.Interface, name, ns string, timeout time.Duration) error {
	return wait.PollUntilContextTimeout(context.TODO(), 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {
		ds, err := k8sClient.AppsV1().DaemonSets(ns).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return false, err
		}

		if ds.Status.NumberReady == ds.Status.DesiredNumberScheduled {
			return true, nil
		}
		return false, nil
	})
}

func RemovePod(k8sClient kubernetes.Interface, ns, podName string) error {
	err := k8sClient.CoreV1().Pods(ns).Delete(context.TODO(), podName, metav1.DeleteOptions{})
	if err != nil && apierrors.IsNotFound(err) {
		return nil
	}
	return err
}

func PodExisted(k8sClient kubernetes.Interface, ns, podName string) bool {
	_, err := k8sClient.CoreV1().Pods(ns).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil && apierrors.IsNotFound(err) {
		return false
	}
	return true
}

func PodCompleted(client kubernetes.Interface, ns, podName string) bool {
	pod, err := client.CoreV1().Pods(ns).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		return false
	}
	return pod.Status.Phase == corev1.PodSucceeded
}

func GetPodLogs(c kubernetes.Interface, namespace, podName, containerName string) (string, error) {
	request := c.CoreV1().RESTClient().Get().
		Resource("pods").
		Namespace(namespace).
		Name(podName).SubResource("log").
		Param("container", containerName)

	logs, err := request.Do(context.TODO()).Raw()
	if err != nil {
		return "", err
	}

	return string(logs), err
}
