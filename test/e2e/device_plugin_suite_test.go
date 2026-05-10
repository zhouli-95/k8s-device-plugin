package e2e

import (
	"context"
	"fmt"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"test/utils"
)

var DefaultPollingInterval = 5 * time.Second
var DefaultWaitTimeout = 60 * time.Second
var DefaultSpecTimeout = 120 * time.Second

var TestNamespace = "default"
var DefaultResourceName = "example.com/dev"
var DefaultResourceNum = int64(4)

var _ = Describe("Device Plugin Basic", Label("basic"), Ordered, func() {

	BeforeAll(func() {
		By("deploy device plugin")
		err := utils.Deploy()
		Expect(err, BeNil())

		By("wait for daemonset ready")
		err = utils.WaitForDaemonSetReady(kubeClient, "device-plugin", TestNamespace, DefaultWaitTimeout)
		Expect(err, BeNil())

		By("wait for pod running")
		Eventually(func() bool {
			pods, err := kubeClient.CoreV1().Pods(TestNamespace).List(context.TODO(), metav1.ListOptions{})
			Expect(err).To(BeNil())
			for _, pod := range pods.Items {
				if pod.Labels["app"] != "device-plugin" {
					continue
				}
				if pod.Status.Phase != corev1.PodRunning {
					return false
				}
			}
			return true
		}).
			WithPolling(DefaultPollingInterval).
			WithTimeout(DefaultWaitTimeout).
			Should(BeTrue())
	})

	AfterAll(func() {
		By("delete device plugin")
		err := utils.UnDeploy()
		Expect(err, BeNil())

		By("wait for node resource cleanup")
		Eventually(func() bool {
			err := utils.CheckNodeResources(kubeClient, DefaultResourceName, 0)
			return err == nil
		}).WithPolling(DefaultPollingInterval).
			WithTimeout(DefaultWaitTimeout).
			Should(BeTrue())
	})

	When("register node resource", func() {
		It("resource should existed", func(ctx SpecContext) {
			Eventually(func() bool {
				err := utils.CheckNodeResources(kubeClient, DefaultResourceName, DefaultResourceNum)
				GinkgoLogr.Info("check resource", "error", err)
				return err == nil

			}).WithPolling(DefaultPollingInterval).
				WithTimeout(DefaultWaitTimeout).
				Should(BeTrue())
		}, SpecTimeout(DefaultSpecTimeout))
	})

	When("create pod with device", func() {
		It("control device should existed", func(ctx SpecContext) {
			podName := "test-pod-dev"
			cmds := []string{"bash", "-c", "ls -l /dev/control-device"}

			_, err := utils.CreatePodWithDevice(kubeClient, TestNamespace, podName, DefaultResourceName, "1", cmds)
			Expect(err).To(BeNil())

			DeferCleanup(func() {
				utils.RemovePod(kubeClient, TestNamespace, podName)
			})

			By("check pod completed")
			Eventually(func() bool {
				return utils.PodCompleted(kubeClient, TestNamespace, podName)
			}).
				WithPolling(DefaultPollingInterval).
				WithTimeout(DefaultWaitTimeout).
				Should(BeTrue())

			logs, err := utils.GetPodLogs(kubeClient, TestNamespace, podName, "app")
			Expect(err).To(BeNil())

			GinkgoLogr.Info("pod logs", "log", logs)
			Expect(logs).To(ContainSubstring("/dev/control-device"))

		}, SpecTimeout(DefaultSpecTimeout))

		It("compute device should existed", func(ctx SpecContext) {
			podName := "test-pod-dev"
			cmds := []string{"bash", "-c", "ls -l /dev/compute-device*"}

			_, err := utils.CreatePodWithDevice(kubeClient, TestNamespace, podName, DefaultResourceName, "1", cmds)
			Expect(err).To(BeNil())

			DeferCleanup(func() {
				utils.RemovePod(kubeClient, TestNamespace, podName)
			})

			By("check pod completed")
			Eventually(func() bool {
				return utils.PodCompleted(kubeClient, TestNamespace, podName)
			}).
				WithPolling(DefaultPollingInterval).
				WithTimeout(DefaultWaitTimeout).
				Should(BeTrue())

			logs, err := utils.GetPodLogs(kubeClient, TestNamespace, podName, "app")
			Expect(err).To(BeNil())

			GinkgoLogr.Info("pod logs", "log", logs)
			Expect(logs).To(ContainSubstring("/dev/compute-device"))

		}, SpecTimeout(DefaultSpecTimeout))

		It("env should existed", func(ctx SpecContext) {
			podName := "test-pod-env"
			cmds := []string{"bash", "-c", "env"}

			_, err := utils.CreatePodWithDevice(kubeClient, TestNamespace, podName, DefaultResourceName, "1", cmds)
			Expect(err).To(BeNil())

			DeferCleanup(func() {
				utils.RemovePod(kubeClient, TestNamespace, podName)
			})

			By("check pod completed")
			Eventually(func() bool {
				return utils.PodCompleted(kubeClient, TestNamespace, podName)
			}).
				WithPolling(DefaultPollingInterval).
				WithTimeout(DefaultWaitTimeout).Should(BeTrue())

			logs, err := utils.GetPodLogs(kubeClient, TestNamespace, podName, "app")
			Expect(err).To(BeNil())
			GinkgoLogr.Info(logs)
			Expect(logs).To(ContainSubstring("DEVICE"))

		}, SpecTimeout(DefaultSpecTimeout))
	})

	When("create pod with multi-device", func() {
		It("control device should existed", func(ctx SpecContext) {
			podName := "test-pod-dev"
			cmds := []string{"bash", "-c", "ls -l /dev/control-device* | wc -l"}

			_, err := utils.CreatePodWithDevice(kubeClient, TestNamespace, podName, DefaultResourceName, "2", cmds)
			Expect(err).To(BeNil())

			DeferCleanup(func() {
				utils.RemovePod(kubeClient, TestNamespace, podName)
			})

			By("check pod completed")
			Eventually(func() bool {
				return utils.PodCompleted(kubeClient, TestNamespace, podName)
			}).
				WithPolling(DefaultPollingInterval).
				WithTimeout(DefaultWaitTimeout).
				Should(BeTrue())

			logs, err := utils.GetPodLogs(kubeClient, TestNamespace, podName, "app")
			Expect(err).To(BeNil())

			GinkgoLogr.Info("pod logs", "log", logs)
			Expect(logs).To(ContainSubstring("1"))

		}, SpecTimeout(DefaultSpecTimeout))

		It("compute device should existed", func(ctx SpecContext) {
			podName := "test-pod-dev"
			cmds := []string{"bash", "-c", "ls -l /dev/compute-device* | wc -l"}

			_, err := utils.CreatePodWithDevice(kubeClient, TestNamespace, podName, DefaultResourceName, "2", cmds)
			Expect(err).To(BeNil())

			DeferCleanup(func() {
				utils.RemovePod(kubeClient, TestNamespace, podName)
			})

			By("check pod completed")
			Eventually(func() bool {
				return utils.PodCompleted(kubeClient, TestNamespace, podName)
			}).
				WithPolling(DefaultPollingInterval).
				WithTimeout(DefaultWaitTimeout).
				Should(BeTrue())

			logs, err := utils.GetPodLogs(kubeClient, TestNamespace, podName, "app")
			Expect(err).To(BeNil())

			GinkgoLogr.Info("pod logs", "log", logs)
			Expect(logs).To(ContainSubstring("2"))

		}, SpecTimeout(DefaultSpecTimeout))

		It("env should existed", func(ctx SpecContext) {
			podName := "test-pod-env"
			cmds := []string{"bash", "-c", "env"}

			_, err := utils.CreatePodWithDevice(kubeClient, TestNamespace, podName, DefaultResourceName, "2", cmds)
			Expect(err).To(BeNil())

			DeferCleanup(func() {
				utils.RemovePod(kubeClient, TestNamespace, podName)
			})

			By("check pod completed")
			Eventually(func() bool {
				return utils.PodCompleted(kubeClient, TestNamespace, podName)
			}).
				WithPolling(DefaultPollingInterval).
				WithTimeout(DefaultWaitTimeout).Should(BeTrue())

			logs, err := utils.GetPodLogs(kubeClient, TestNamespace, podName, "app")
			Expect(err).To(BeNil())
			GinkgoLogr.Info(logs)
			Expect(logs).To(ContainSubstring("DEVICE"))

		}, SpecTimeout(DefaultSpecTimeout))
	})
})

var _ = Describe("Device Plugin Config", Label("config"), Ordered, func() {

	type Config struct {
		domain         string
		resourceName   string
		deviceStrategy string
		kustomizeDir   string
	}

	DescribeTableSubtree("device strategy", Label("device-strategy"),
		func(config Config) {

			BeforeAll(func() {
				By("deploy device plugin")
				err := utils.Deploy(utils.WithKustomizeDir(config.kustomizeDir))
				Expect(err, BeNil())

				By("wait for daemonset ready")
				err = utils.WaitForDaemonSetReady(kubeClient, "device-plugin", TestNamespace, DefaultWaitTimeout)
				Expect(err, BeNil())

				By("wait for pod running")
				Eventually(func() bool {
					pods, err := kubeClient.CoreV1().Pods(TestNamespace).List(context.TODO(), metav1.ListOptions{})
					Expect(err).To(BeNil())
					for _, pod := range pods.Items {
						if pod.Labels["app"] != "device-plugin" {
							continue
						}
						if pod.Status.Phase != corev1.PodRunning {
							return false
						}
					}
					return true
				}).
					WithPolling(DefaultPollingInterval).
					WithTimeout(DefaultWaitTimeout).
					Should(BeTrue())

				By("wait for node resource")
				Eventually(func() bool {
					err := utils.CheckNodeResources(kubeClient, DefaultResourceName, DefaultResourceNum)
					GinkgoLogr.Info("check resource", "error", err)
					return err == nil
				}).WithTimeout(DefaultWaitTimeout).
					WithPolling(DefaultPollingInterval).
					Should(BeTrue())
			})

			AfterAll(func() {
				By("delete device plugin")
				err := utils.UnDeploy(utils.WithKustomizeDir(config.kustomizeDir))
				Expect(err, BeNil())

				By("wait for node resource cleanup")
				Eventually(func() bool {
					err := utils.CheckNodeResources(kubeClient, DefaultResourceName, 0)
					return err == nil
				}).WithPolling(DefaultPollingInterval).
					WithTimeout(DefaultWaitTimeout).
					Should(BeTrue())
			})

			Context("with device", func() {
				It("device should existed", func(ctx SpecContext) {
					podName := "test-pod-dev"
					cmds := []string{"bash", "-c", "ls -l /dev/compute-device*"}

					_, err := utils.CreatePodWithDevice(kubeClient, TestNamespace, podName, DefaultResourceName, "1", cmds)
					Expect(err).To(BeNil())

					DeferCleanup(func() {
						utils.RemovePod(kubeClient, TestNamespace, podName)
					})

					By("check pod completed")
					Eventually(func() bool {
						GinkgoLogr.Info("check pod", "status", utils.BashRun("kubectl get po -owide"))
						return utils.PodCompleted(kubeClient, TestNamespace, podName)
					}).
						WithPolling(DefaultPollingInterval).
						WithTimeout(DefaultWaitTimeout).
						Should(BeTrue())

					logs, err := utils.GetPodLogs(kubeClient, TestNamespace, podName, "app")
					Expect(err).To(BeNil())

					GinkgoLogr.Info("pod logs", "log", logs)
					Expect(logs).To(ContainSubstring("/dev/compute-device"))

				}, SpecTimeout(DefaultSpecTimeout))

				It("env should existed", func(ctx SpecContext) {
					podName := "test-pod-env"
					cmds := []string{"bash", "-c", "env"}

					_, err := utils.CreatePodWithDevice(kubeClient, TestNamespace, podName, DefaultResourceName, "1", cmds)
					Expect(err).To(BeNil())

					DeferCleanup(func() {
						utils.RemovePod(kubeClient, TestNamespace, podName)
					})

					By("check pod completed")
					Eventually(func() bool {
						return utils.PodCompleted(kubeClient, TestNamespace, podName)
					}).WithPolling(DefaultPollingInterval).
						WithTimeout(DefaultWaitTimeout).
						Should(BeTrue())

					logs, err := utils.GetPodLogs(kubeClient, TestNamespace, podName, "app")
					Expect(err).To(BeNil())
					GinkgoLogr.Info(logs)
					Expect(logs).To(ContainSubstring(fmt.Sprintf("DEVICE_STRATEGY=%s", config.deviceStrategy)))

				}, SpecTimeout(DefaultSpecTimeout))
			})

		},

		Entry("native", Config{deviceStrategy: "native"}),
		Entry("cdi-annotation", Config{deviceStrategy: "cdi-annotation", kustomizeDir: "cdi-annotation"}),
		Entry("cdi-cri", Config{deviceStrategy: "cdi-cri", kustomizeDir: "cdi-cri"}),
	)

	DescribeTableSubtree("custom domain and resource name", Label("custom-domain-resource"),
		func(config Config) {

			BeforeAll(func() {
				By("deploy device plugin")
				err := utils.Deploy(utils.WithKustomizeDir(config.kustomizeDir))
				Expect(err, BeNil())

				By("wait for daemonset ready")
				err = utils.WaitForDaemonSetReady(kubeClient, "device-plugin", TestNamespace, DefaultWaitTimeout)
				Expect(err, BeNil())

				By("wait for pod running")
				Eventually(func() bool {
					pods, err := kubeClient.CoreV1().Pods(TestNamespace).List(context.TODO(), metav1.ListOptions{})
					Expect(err).To(BeNil())
					for _, pod := range pods.Items {
						if pod.Labels["app"] != "device-plugin" {
							continue
						}
						if pod.Status.Phase != corev1.PodRunning {
							return false
						}
					}
					return true
				}).WithTimeout(DefaultWaitTimeout).
					WithPolling(DefaultPollingInterval).
					Should(BeTrue())

				By("wait for node resource")
				Eventually(func() bool {
					err := utils.CheckNodeResources(kubeClient, config.resourceName, DefaultResourceNum)
					GinkgoLogr.Info("check resource", "resource name", config.resourceName, "error", err)
					return err == nil
				}).WithTimeout(DefaultWaitTimeout).
					WithPolling(DefaultPollingInterval).
					Should(BeTrue())
			})

			AfterAll(func() {
				By("delete device plugin")
				err := utils.UnDeploy(utils.WithKustomizeDir(config.kustomizeDir))
				Expect(err, BeNil())

				By("wait for node resource cleanup")
				Eventually(func() bool {
					err := utils.CheckNodeResources(kubeClient, DefaultResourceName, 0)
					return err == nil
				}).WithPolling(DefaultPollingInterval).
					WithTimeout(DefaultWaitTimeout).
					Should(BeTrue())
			})

			It("compute device should existed", func(ctx SpecContext) {
				podName := "test-pod-custom-dev"
				cmds := []string{"bash", "-c", "ls -l /dev/compute-device*"}

				_, err := utils.CreatePodWithDevice(kubeClient, TestNamespace, podName, config.resourceName, "1", cmds)
				Expect(err).To(BeNil())

				DeferCleanup(func() {
					utils.RemovePod(kubeClient, TestNamespace, podName)
				})

				By("check pod completed")
				Eventually(func() bool {
					return utils.PodCompleted(kubeClient, TestNamespace, podName)
				}).
					WithPolling(DefaultPollingInterval).
					WithTimeout(DefaultWaitTimeout).
					Should(BeTrue())

				logs, err := utils.GetPodLogs(kubeClient, TestNamespace, podName, "app")
				Expect(err).To(BeNil())

				GinkgoLogr.Info("pod logs", "log", logs)
				Expect(logs).To(ContainSubstring("/dev/compute-device"))

			}, SpecTimeout(DefaultSpecTimeout))

			It("control device should existed", func(ctx SpecContext) {
				podName := "test-pod-custom-ctrl"
				cmds := []string{"bash", "-c", "ls -l /dev/control-device"}

				_, err := utils.CreatePodWithDevice(kubeClient, TestNamespace, podName, config.resourceName, "1", cmds)
				Expect(err).To(BeNil())

				DeferCleanup(func() {
					utils.RemovePod(kubeClient, TestNamespace, podName)
				})

				By("check pod completed")
				Eventually(func() bool {
					return utils.PodCompleted(kubeClient, TestNamespace, podName)
				}).
					WithPolling(DefaultPollingInterval).
					WithTimeout(DefaultWaitTimeout).
					Should(BeTrue())

				logs, err := utils.GetPodLogs(kubeClient, TestNamespace, podName, "app")
				Expect(err).To(BeNil())

				GinkgoLogr.Info("pod logs", "log", logs)
				Expect(logs).To(ContainSubstring("/dev/control-device"))

			}, SpecTimeout(DefaultSpecTimeout))

			It("env should existed", func(ctx SpecContext) {
				podName := "test-pod-custom-env"
				cmds := []string{"bash", "-c", "env"}

				_, err := utils.CreatePodWithDevice(kubeClient, TestNamespace, podName, config.resourceName, "1", cmds)
				Expect(err).To(BeNil())

				DeferCleanup(func() {
					utils.RemovePod(kubeClient, TestNamespace, podName)
				})

				By("check pod completed")
				Eventually(func() bool {
					return utils.PodCompleted(kubeClient, TestNamespace, podName)
				}).
					WithPolling(DefaultPollingInterval).
					WithTimeout(DefaultWaitTimeout).
					Should(BeTrue())

				logs, err := utils.GetPodLogs(kubeClient, TestNamespace, podName, "app")
				Expect(err).To(BeNil())
				GinkgoLogr.Info("pod logs", "log", logs)
				Expect(logs).To(ContainSubstring("DEVICE_STRATEGY"))

			}, SpecTimeout(DefaultSpecTimeout))
		},

		Entry("custom domain and resource", Config{
			kustomizeDir: "custom-domain-resource",
			resourceName: "test.io/resource",
		}),
	)
})
