package relationships_test

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	korifiv1alpha1 "code.cloudfoundry.org/korifi/controllers/api/v1alpha1"
	"code.cloudfoundry.org/korifi/controllers/webhooks/relationships"
	"code.cloudfoundry.org/korifi/tests/helpers"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec" //lint:ignore ST1001 this is a test file
	admissionv1beta1 "k8s.io/api/admission/v1beta1"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	//+kubebuilder:scaffold:imports
)

var (
	stopManager     context.CancelFunc
	stopClientCache context.CancelFunc
	testEnv         *envtest.Environment
	adminClient     client.Client
	ctx             context.Context
)

func TestWorkloadsWebhooks(t *testing.T) {
	SetDefaultEventuallyTimeout(10 * time.Second)
	SetDefaultEventuallyPollingInterval(250 * time.Millisecond)

	RegisterFailHandler(Fail)
	RunSpecs(t, "Finalizer Webhook Integration Test Suite")
}

var _ = BeforeSuite(func() {
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))

	webhookManifestsPath := generateWebhookManifest()
	DeferCleanup(func() {
		Expect(os.RemoveAll(filepath.Dir(webhookManifestsPath))).To(Succeed())
	})
	testEnv = &envtest.Environment{
		CRDDirectoryPaths: []string{
			filepath.Join("..", "..", "..", "helm", "korifi", "controllers", "crds"),
		},
		ErrorIfCRDPathMissing: true,
		WebhookInstallOptions: envtest.WebhookInstallOptions{
			Paths: []string{webhookManifestsPath},
		},
	}

	_, err := testEnv.Start()
	Expect(err).NotTo(HaveOccurred())

	Expect(korifiv1alpha1.AddToScheme(scheme.Scheme)).To(Succeed())
	Expect(admissionv1beta1.AddToScheme(scheme.Scheme)).To(Succeed())

	k8sManager := helpers.NewK8sManager(testEnv, filepath.Join("helm", "korifi", "controllers", "role.yaml"))

	adminClient, stopClientCache = helpers.NewCachedClient(testEnv.Config)

	relationships.NewSpaceGUIDWebhook().SetupWebhookWithManager(k8sManager)

	stopManager = helpers.StartK8sManager(k8sManager)
})

var _ = BeforeEach(func() {
	ctx = context.Background()
})

var _ = AfterSuite(func() {
	stopClientCache()
	stopManager()
	Expect(testEnv.Stop()).To(Succeed())
})

func generateWebhookManifest() string {
	tmpDir, err := os.MkdirTemp("", "")
	Expect(err).NotTo(HaveOccurred())

	controllerGenSession, err := gexec.Start(exec.Command(
		"controller-gen",
		"paths=code.cloudfoundry.org/korifi/controllers/webhooks/relationships",
		"webhook",
		fmt.Sprintf("output:webhook:artifacts:config=%s", tmpDir),
	), GinkgoWriter, GinkgoWriter)

	Expect(err).NotTo(HaveOccurred())
	Eventually(controllerGenSession).Should(gexec.Exit(0))

	return filepath.Join(tmpDir, "manifests.yaml")
}
