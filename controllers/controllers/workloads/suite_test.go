package workloads_test

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	korifiv1alpha1 "code.cloudfoundry.org/korifi/controllers/api/v1alpha1"
	"code.cloudfoundry.org/korifi/controllers/config"
	. "code.cloudfoundry.org/korifi/controllers/controllers/shared"
	. "code.cloudfoundry.org/korifi/controllers/controllers/workloads"
	"code.cloudfoundry.org/korifi/controllers/controllers/workloads/env"
	"code.cloudfoundry.org/korifi/controllers/controllers/workloads/testutils"
	"code.cloudfoundry.org/korifi/tools/k8s"

	"github.com/jonboulle/clockwork"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	servicebindingv1beta1 "github.com/servicebinding/service-binding-controller/apis/v1beta1"
	"go.uber.org/zap/zapcore"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sclient "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	//+kubebuilder:scaffold:imports
)

var (
	cancel            context.CancelFunc
	testEnv           *envtest.Environment
	k8sClient         client.Client
	cfProcessDefaults config.CFProcessDefaults
	cfRootNamespace   string
)

const (
	packageRegistrySecretName = "test-package-registry-secret"
)

func TestWorkloadsControllers(t *testing.T) {
	SetDefaultEventuallyTimeout(10 * time.Second)
	SetDefaultEventuallyPollingInterval(250 * time.Millisecond)

	RegisterFailHandler(Fail)
	RunSpecs(t, "Workloads Controllers Integration Suite")
}

var _ = BeforeSuite(func() {
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true), zap.Level(zapcore.DebugLevel)))

	ctx, cancelFunc := context.WithCancel(context.TODO())
	cancel = cancelFunc

	testEnv = &envtest.Environment{
		CRDDirectoryPaths: []string{
			filepath.Join("..", "..", "..", "helm", "controllers", "templates", "crds"),
		},
		ErrorIfCRDPathMissing: true,
	}

	cfg, err := testEnv.Start()
	Expect(err).NotTo(HaveOccurred())
	Expect(cfg).NotTo(BeNil())

	Expect(korifiv1alpha1.AddToScheme(scheme.Scheme)).To(Succeed())
	Expect(servicebindingv1beta1.AddToScheme(scheme.Scheme)).To(Succeed())
	Expect(corev1.AddToScheme(scheme.Scheme)).To(Succeed())

	//+kubebuilder:scaffold:scheme

	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
	Expect(err).NotTo(HaveOccurred())
	Expect(k8sClient).NotTo(BeNil())

	webhookInstallOptions := &testEnv.WebhookInstallOptions
	k8sManager, err := ctrl.NewManager(cfg, ctrl.Options{
		Scheme:             scheme.Scheme,
		Host:               webhookInstallOptions.LocalServingHost,
		Port:               webhookInstallOptions.LocalServingPort,
		CertDir:            webhookInstallOptions.LocalServingCertDir,
		LeaderElection:     false,
		MetricsBindAddress: "0",
	})
	Expect(err).NotTo(HaveOccurred())

	cfRootNamespace = testutils.PrefixedGUID("root-namespace")

	createNamespace(ctx, k8sClient, cfRootNamespace)

	controllerConfig := &config.ControllerConfig{
		CFProcessDefaults: config.CFProcessDefaults{
			MemoryMB:    500,
			DiskQuotaMB: 512,
		},
		CFRootNamespace:             cfRootNamespace,
		PackageRegistrySecretName:   packageRegistrySecretName,
		WorkloadsTLSSecretName:      "korifi-workloads-ingress-cert",
		WorkloadsTLSSecretNamespace: "korifi-controllers-system",
	}

	err = (NewCFAppReconciler(
		k8sManager.GetClient(),
		k8sManager.GetScheme(),
		ctrl.Log.WithName("controllers").WithName("CFApp"),
		env.NewBuilder(k8sManager.GetClient()),
	)).SetupWithManager(k8sManager)
	Expect(err).NotTo(HaveOccurred())

	registryAuthFetcherClient, err := k8sclient.NewForConfig(cfg)
	Expect(err).NotTo(HaveOccurred())
	Expect(registryAuthFetcherClient).NotTo(BeNil())
	cfBuildReconciler := NewCFBuildReconciler(
		k8sManager.GetClient(),
		k8sManager.GetScheme(),
		ctrl.Log.WithName("controllers").WithName("CFBuild"),
		controllerConfig,
		env.NewBuilder(k8sManager.GetClient()),
	)
	err = (cfBuildReconciler).SetupWithManager(k8sManager)
	Expect(err).NotTo(HaveOccurred())

	err = (NewCFProcessReconciler(
		k8sManager.GetClient(),
		k8sManager.GetScheme(),
		ctrl.Log.WithName("controllers").WithName("CFProcess"),
		controllerConfig,
		env.NewBuilder(k8sManager.GetClient()),
	)).SetupWithManager(k8sManager)
	Expect(err).NotTo(HaveOccurred())

	err = (NewCFPackageReconciler(
		k8sManager.GetClient(),
		k8sManager.GetScheme(),
		ctrl.Log.WithName("controllers").WithName("CFPackage"),
	)).SetupWithManager(k8sManager)
	Expect(err).NotTo(HaveOccurred())

	err = NewCFOrgReconciler(
		k8sManager.GetClient(),
		k8sManager.GetScheme(),
		ctrl.Log.WithName("controllers").WithName("CFOrg"),
		controllerConfig.PackageRegistrySecretName,
	).SetupWithManager(k8sManager)
	Expect(err).NotTo(HaveOccurred())

	cfProcessDefaults = config.CFProcessDefaults{
		MemoryMB:    256,
		DiskQuotaMB: 128,
	}
	err = NewCFTaskReconciler(
		k8sManager.GetClient(),
		k8sManager.GetScheme(),
		k8sManager.GetEventRecorderFor("cftask-controller"),
		ctrl.Log.WithName("controllers").WithName("CFTask"),
		NewSequenceId(clockwork.NewRealClock()),
		env.NewBuilder(k8sManager.GetClient()),
		cfProcessDefaults,
		2*time.Second,
	).SetupWithManager(k8sManager)
	Expect(err).NotTo(HaveOccurred())

	err = NewCFSpaceReconciler(
		k8sManager.GetClient(),
		k8sManager.GetScheme(),
		ctrl.Log.WithName("controllers").WithName("CFSpace"),
		controllerConfig.PackageRegistrySecretName,
		controllerConfig.CFRootNamespace,
	).SetupWithManager(k8sManager)
	Expect(err).NotTo(HaveOccurred())

	// Add new reconcilers here

	// Setup index for manager
	err = SetupIndexWithManager(k8sManager)
	Expect(err).NotTo(HaveOccurred())

	go func() {
		defer GinkgoRecover()
		err = k8sManager.Start(ctx)
		Expect(err).NotTo(HaveOccurred())
	}()
})

var _ = AfterSuite(func() {
	cancel()
	Expect(testEnv.Stop()).To(Succeed())
})

func createBuildWithDroplet(ctx context.Context, k8sClient client.Client, cfBuild *korifiv1alpha1.CFBuild, droplet *korifiv1alpha1.BuildDropletStatus) *korifiv1alpha1.CFBuild {
	Expect(
		k8sClient.Create(ctx, cfBuild),
	).To(Succeed())
	patchedCFBuild := cfBuild.DeepCopy()
	patchedCFBuild.Status.Conditions = []metav1.Condition{}
	patchedCFBuild.Status.Droplet = droplet
	Expect(
		k8sClient.Status().Patch(ctx, patchedCFBuild, client.MergeFrom(cfBuild)),
	).To(Succeed())
	return patchedCFBuild
}

func createNamespace(ctx context.Context, k8sClient client.Client, name string) *corev1.Namespace {
	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
	Expect(
		k8sClient.Create(ctx, ns)).To(Succeed())
	return ns
}

func createSecret(ctx context.Context, k8sClient client.Client, name string, namespace string) *corev1.Secret {
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		StringData: map[string]string{
			"foo": "bar",
		},
		Type: "Docker",
	}
	Expect(k8sClient.Create(ctx, secret)).To(Succeed())
	return secret
}

func createClusterRole(ctx context.Context, k8sClient client.Client, name string, rules []rbacv1.PolicyRule) *rbacv1.ClusterRole {
	role := &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Rules: rules,
	}
	Expect(k8sClient.Create(ctx, role)).To(Succeed())
	return role
}

func createRoleBinding(ctx context.Context, k8sClient client.Client, roleBindingName, subjectName, roleReference, namespace string, annotations map[string]string) *rbacv1.RoleBinding {
	roleBinding := &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:        roleBindingName,
			Namespace:   namespace,
			Annotations: annotations,
		},
		Subjects: []rbacv1.Subject{{
			Kind: rbacv1.UserKind,
			Name: subjectName,
		}},
		RoleRef: rbacv1.RoleRef{
			Kind: "ClusterRole",
			Name: roleReference,
		},
	}
	Expect(k8sClient.Create(ctx, roleBinding)).To(Succeed())
	return roleBinding
}

func createServiceAccount(ctx context.Context, k8sclient client.Client, serviceAccountName, namespace string, annotations map[string]string) *corev1.ServiceAccount {
	serviceAccount := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:        serviceAccountName,
			Namespace:   namespace,
			Annotations: annotations,
		},
		Secrets: []corev1.ObjectReference{
			{Name: serviceAccountName + "-token-someguid"},
			{Name: serviceAccountName + "-dockercfg-someguid"},
			{Name: packageRegistrySecretName},
		},
		ImagePullSecrets: []corev1.LocalObjectReference{
			{Name: serviceAccountName + "-dockercfg-someguid"},
			{Name: packageRegistrySecretName},
		},
	}
	Expect(k8sClient.Create(ctx, serviceAccount)).To(Succeed())
	return serviceAccount
}

func createNamespaceWithCleanup(ctx context.Context, k8sClient client.Client, name string) *corev1.Namespace {
	namespace := createNamespace(ctx, k8sClient, name)

	DeferCleanup(func() {
		Expect(k8sClient.Delete(ctx, namespace)).To(Succeed())
	})

	return namespace
}

func patchAppWithDroplet(ctx context.Context, k8sClient client.Client, appGUID, spaceGUID, buildGUID string) *korifiv1alpha1.CFApp {
	cfApp := &korifiv1alpha1.CFApp{
		ObjectMeta: metav1.ObjectMeta{
			Name:      appGUID,
			Namespace: spaceGUID,
		},
	}
	Expect(k8s.Patch(ctx, k8sClient, cfApp, func() {
		cfApp.Spec.CurrentDropletRef = corev1.LocalObjectReference{Name: buildGUID}
	})).To(Succeed())
	return cfApp
}
