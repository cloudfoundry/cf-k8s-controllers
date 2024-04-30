package smoke_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	korifiv1alpha1 "code.cloudfoundry.org/korifi/controllers/api/v1alpha1"
	"code.cloudfoundry.org/korifi/tests/helpers"
	"code.cloudfoundry.org/korifi/tests/helpers/fail_handler"

	"github.com/cloudfoundry/cf-test-helpers/generator"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
	"gopkg.in/yaml.v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	rest "k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

const NamePrefix = "cf-on-k8s-smoke"

var (
	appsDomain            string
	buildpackAppName      string
	cfAdmin               string
	dockerAppName         string
	orgName               string
	rootNamespace         string
	serviceAccountFactory *helpers.ServiceAccountFactory
	spaceName             string
)

func TestSmoke(t *testing.T) {
	RegisterFailHandler(fail_handler.New("Smoke Tests",
		fail_handler.Hook{
			Matcher: fail_handler.Always,
			Hook: func(config *rest.Config, failure fail_handler.TestFailure) {
				printCfApp(config)
				fail_handler.PrintKorifiLogs(config, "", failure.StartTime)
				printBuildLogs(config, spaceName)
			},
		}).Fail)

	SetDefaultEventuallyTimeout(helpers.EventuallyTimeout())
	SetDefaultEventuallyPollingInterval(helpers.EventuallyPollingInterval())
	RunSpecs(t, "Smoke Tests Suite")
}

var _ = BeforeSuite(func() {
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))

	rootNamespace = helpers.GetDefaultedEnvVar("ROOT_NAMESPACE", "cf")
	serviceAccountFactory = helpers.NewServiceAccountFactory(rootNamespace)

	cfAdmin = uuid.NewString()
	cfAdminToken := serviceAccountFactory.CreateAdminServiceAccount(cfAdmin)
	helpers.AddUserToKubeConfig(cfAdmin, cfAdminToken)

	Expect(helpers.Cf("api", helpers.GetRequiredEnvVar("API_SERVER_ROOT"), "--skip-ssl-validation")).To(Exit(0))
	Expect(helpers.Cf("auth", cfAdmin)).To(Exit(0))

	appsDomain = helpers.GetRequiredEnvVar("APP_FQDN")
	orgName = generator.PrefixedRandomName(NamePrefix, "org")
	spaceName = generator.PrefixedRandomName(NamePrefix, "space")
	buildpackAppName = generator.PrefixedRandomName(NamePrefix, "buildpackapp")
	dockerAppName = generator.PrefixedRandomName(NamePrefix, "dockerapp")

	Expect(helpers.Cf("create-org", orgName)).To(Exit(0))
	Expect(helpers.Cf("create-space", "-o", orgName, spaceName)).To(Exit(0))
	Expect(helpers.Cf("target", "-o", orgName, "-s", spaceName)).To(Exit(0))

	Expect(helpers.Cf("push", buildpackAppName, "-p", "../assets/dorifi")).To(Exit(0))
	Expect(helpers.Cf("push", dockerAppName, "-o", "eirini/dorini")).To(Exit(0))
})

var _ = AfterSuite(func() {
	Expect(helpers.Cf("delete-org", orgName, "-f").Wait()).To(Exit())

	serviceAccountFactory.DeleteServiceAccount(cfAdmin)
	helpers.RemoveUserFromKubeConfig(cfAdmin)
})

func sessionOutput(session *Session) (string, error) {
	if session.ExitCode() != 0 {
		return "", fmt.Errorf("Session %v exited with exit code %d: %s",
			session.Command,
			session.ExitCode(),
			string(session.Err.Contents()),
		)
	}
	return strings.TrimSpace(string(session.Out.Contents())), nil
}

func printCfApp(config *rest.Config) {
	utilruntime.Must(korifiv1alpha1.AddToScheme(scheme.Scheme))
	k8sClient, err := client.New(config, client.Options{Scheme: scheme.Scheme})
	if err != nil {
		fmt.Fprintf(GinkgoWriter, "failed to create k8s client: %v\n", err)
		return
	}

	cfAppNamespace, err := sessionOutput(helpers.Cf("space", spaceName, "--guid"))
	if err != nil {
		fmt.Fprintf(GinkgoWriter, "failed to run 'cf space %s --guid': %v\n", spaceName, err)
		return
	}
	cfAppGUID, err := sessionOutput(helpers.Cf("app", buildpackAppName, "--guid"))
	if err != nil {
		fmt.Fprintf(GinkgoWriter, "failed to run 'cf app %s --guid': %v\n", buildpackAppName, err)
		return
	}

	if err := printObject(k8sClient, &korifiv1alpha1.CFApp{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cfAppGUID,
			Namespace: cfAppNamespace,
		},
	}); err != nil {
		fmt.Fprintf(GinkgoWriter, "failed printing cfapp: %v\n", err)
		return
	}
}

func printObject(k8sClient client.Client, obj client.Object) error {
	if err := k8sClient.Get(context.Background(), client.ObjectKeyFromObject(obj), obj); err != nil {
		return fmt.Errorf("failed to get object %q: %v\n", client.ObjectKeyFromObject(obj), err)
	}

	fmt.Fprintf(GinkgoWriter, "\n\n========== %T %s/%s (skipping managed fields) ==========\n", obj, obj.GetNamespace(), obj.GetName())
	obj.SetManagedFields([]metav1.ManagedFieldsEntry{})
	objBytes, err := yaml.Marshal(obj)
	if err != nil {
		return fmt.Errorf("failed marshalling object %v: %v\n", obj, err)
	}
	fmt.Fprintln(GinkgoWriter, string(objBytes))
	return nil
}

func printBuildLogs(config *rest.Config, spaceName string) {
	spaceGUID, err := sessionOutput(helpers.Cf("space", spaceName, "--guid").Wait())
	if err != nil {
		fmt.Fprintf(GinkgoWriter, "failed to get space guid: %v\n", err)
		return
	}
	fail_handler.PrintAllBuildLogs(config, spaceGUID)
}
