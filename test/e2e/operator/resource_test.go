// test/e2e/operator/resources_test.go
package operator_test

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Resource Management", func() {
	var (
		testNamespace    string
		managedNamespace string
	)

	BeforeEach(func() {
		testNamespace = createTestNamespace("test-ns-", nil)
		managedNamespace = createTestNamespace("managed-ns-", map[string]string{
			"kueue.openshift.io/managed": "true",
		})
	})

	AfterEach(func() {
		cleanupNamespace(testNamespace)
		cleanupNamespace(managedNamespace)
	})

	Context("When managing jobs", func() {
		It("Should handle jobs in managed namespaces", func() {
			job := createTestJob(managedNamespace)
			_, err := f.KubeClient.BatchV1().Jobs(managedNamespace).Create(
				context.TODO(), job, metav1.CreateOptions{})
			Expect(err).NotTo(HaveOccurred())

			Eventually(func() bool {
				return verifyJobManaged(managedNamespace, job.Name)
			}, operatorReadyTimeout, operatorPollInterval).Should(BeTrue())
		})

		It("Should ignore jobs in unmanaged namespaces", func() {
			job := createTestJob(testNamespace)
			_, err := f.KubeClient.BatchV1().Jobs(testNamespace).Create(
				context.TODO(), job, metav1.CreateOptions{})
			Expect(err).NotTo(HaveOccurred())

			Consistently(func() bool {
				return verifyJobManaged(testNamespace, job.Name)
			}, 30*time.Second, 5*time.Second).Should(BeFalse())
		})
	})
})

func createTestNamespace(prefix string, labels map[string]string) string {
	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: prefix,
			Labels:       labels,
		},
	}
	created, err := f.KubeClient.CoreV1().Namespaces().Create(
		context.TODO(), ns, metav1.CreateOptions{})
	Expect(err).NotTo(HaveOccurred())
	return created.Name
}

func cleanupNamespace(name string) {
	err := f.KubeClient.CoreV1().Namespaces().Delete(
		context.TODO(), name, metav1.DeleteOptions{})
	Expect(err).NotTo(HaveOccurred())
}
