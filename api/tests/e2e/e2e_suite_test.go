package e2e_test

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"net/http"
	"os"
	"sync"
	"testing"
	"time"

	"code.cloudfoundry.org/cf-k8s-controllers/api/apis"
	"code.cloudfoundry.org/cf-k8s-controllers/api/tests/e2e/helpers"

	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	certsv1 "k8s.io/api/certificates/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	controllerruntime "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	hnsv1alpha2 "sigs.k8s.io/hierarchical-namespaces/api/v1alpha2"
)

var (
	k8sClient           client.WithWatch
	adminClient         *resty.Client
	certClient          *resty.Client
	tokenClient         *resty.Client
	clientset           *kubernetes.Clientset
	rootNamespace       string
	apiServerRoot       string
	serviceAccountName  string
	serviceAccountToken string
	tokenAuthHeader     string
	certUserName        string
	certSigningReq      *certsv1.CertificateSigningRequest
	certAuthHeader      string
	adminAuthHeader     string
	certPEM             string
)

type resource struct {
	Name          string        `json:"name,omitempty"`
	GUID          string        `json:"guid,omitempty"`
	Relationships relationships `json:"relationships,omitempty"`
	CreatedAt     string        `json:"created_at,omitempty"`
	UpdatedAt     string        `json:"updated_at,omitempty"`
}

type relationships map[string]relationship

type relationship struct {
	Data resource `json:"data"`
}

type resourceList struct {
	Resources []resource `json:"resources"`
}

type appResource struct {
	resource `json:",inline"`
	State    string `json:"state,omitempty"`
}

type roleResource struct {
	resource `json:",inline"`
	Type     string `json:"type,omitempty"`
}

type packageResource struct {
	resource `json:",inline"`
	Type     string `json:"type,omitempty"`
}

type buildResource struct {
	resource `json:",inline"`
	Package  resource `json:"package"`
}

type dropletResource struct {
	Data resource `json:"data"`
}

type statsResourceList struct {
	Resources []statsResource `json:"resources"`
}

type statsResource struct {
	Type  string `json:"type"`
	State string `json:"state"`
}

type cfErrs struct {
	Errors []cfErr
}

type cfErr struct {
	Detail string `json:"detail"`
	Title  string `json:"title"`
	Code   int    `json:"code"`
}

func TestE2E(t *testing.T) {
	RegisterFailHandler(helpers.E2EFailHandler)
	RunSpecs(t, "E2E Suite")
}

var _ = BeforeSuite(func() {
	SetDefaultEventuallyTimeout(240 * time.Second)
	SetDefaultEventuallyPollingInterval(2 * time.Second)

	apiServerRoot = mustHaveEnv("API_SERVER_ROOT")

	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))

	Expect(hnsv1alpha2.AddToScheme(scheme.Scheme)).To(Succeed())

	config, err := controllerruntime.GetConfig()
	Expect(err).NotTo(HaveOccurred())

	adminAuthHeader = "ClientCert " + obtainAdminUserCert()

	k8sClient, err = client.NewWithWatch(config, client.Options{Scheme: scheme.Scheme})
	Expect(err).NotTo(HaveOccurred())

	clientset, err = kubernetes.NewForConfig(config)
	Expect(err).NotTo(HaveOccurred())

	rootNamespace = mustHaveEnv("ROOT_NAMESPACE")
	ensureServerIsUp()

	serviceAccountName = generateGUID("token-user")
	serviceAccountToken = obtainServiceAccountToken(serviceAccountName)

	certUserName = generateGUID("cert-user")
	certSigningReq, certPEM = obtainClientCert(certUserName)
	certAuthHeader = "ClientCert " + certPEM
})

var _ = BeforeEach(func() {
	tokenAuthHeader = fmt.Sprintf("Bearer %s", serviceAccountToken)
	adminClient = resty.New().SetBaseURL(apiServerRoot).SetAuthScheme("ClientCert").SetAuthToken(obtainAdminUserCert())
	certClient = resty.New().SetBaseURL(apiServerRoot).SetAuthScheme("ClientCert").SetAuthToken(certPEM)
	tokenClient = resty.New().SetBaseURL(apiServerRoot).SetAuthToken(serviceAccountToken)
})

var _ = AfterSuite(func() {
	deleteServiceAccount(serviceAccountName)
	deleteCSR(certSigningReq)
})

func mustHaveEnv(key string) string {
	val, ok := os.LookupEnv(key)
	Expect(ok).To(BeTrue(), "must set env var %q", key)

	return val
}

func ensureServerIsUp() {
	Eventually(func() (int, error) {
		resp, err := http.Get(apiServerRoot)
		if err != nil {
			return 0, err
		}

		resp.Body.Close()

		return resp.StatusCode, nil
	}, "5m").Should(Equal(http.StatusOK), "API Server at %s was not running after 5 minutes", apiServerRoot)
}

func generateGUID(prefix string) string {
	guid := uuid.NewString()

	return fmt.Sprintf("%s-%s", prefix, guid[:13])
}

func deleteOrg(name string) {
	if name == "" {
		return
	}
	deleteSubnamespace(rootNamespace, name)
}

func asyncDeleteOrg(orgID string, wg *sync.WaitGroup) {
	go func() {
		defer wg.Done()
		defer GinkgoRecover()

		deleteOrg(orgID)
	}()
}

func deleteSubnamespace(parent, name string) {
	ctx := context.Background()

	subnsList := &hnsv1alpha2.SubnamespaceAnchorList{}
	Expect(k8sClient.List(ctx, subnsList, client.InNamespace(name))).To(Succeed())

	var wg sync.WaitGroup
	wg.Add(len(subnsList.Items))
	for _, subns := range subnsList.Items {
		go func(subns string) {
			defer wg.Done()
			defer GinkgoRecover()

			deleteSubnamespace(name, subns)
		}(subns.Name)
	}
	wg.Wait()

	anchor := hnsv1alpha2.SubnamespaceAnchor{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: parent,
			Name:      name,
		},
	}
	err := k8sClient.Delete(ctx, &anchor)
	if errors.IsNotFound(err) {
		return
	}
	Expect(err).NotTo(HaveOccurred())

	Eventually(func() bool {
		err := k8sClient.Get(ctx, client.ObjectKeyFromObject(&anchor), &anchor)

		return errors.IsNotFound(err)
	}).Should(BeTrue())
}

func createOrgRaw(orgName string) (string, error) {
	var org resource
	resp, err := adminClient.R().
		SetBody(resource{Name: orgName}).
		SetResult(&org).
		Post("/v3/organizations")
	if err != nil {
		return "", err
	}
	if resp.StatusCode() != http.StatusCreated {
		return "", fmt.Errorf("expected status code %d, got %d", http.StatusCreated, resp.StatusCode())
	}

	return org.GUID, nil
}

func createOrg(orgName string) string {
	orgGUID, err := createOrgRaw(orgName)
	Expect(err).NotTo(HaveOccurred())
	Expect(waitForAdminRoleBinding(orgGUID)).To(Succeed())

	return orgGUID
}

func asyncCreateOrg(orgName string, createdOrgGUID *string, wg *sync.WaitGroup, errChan chan error) {
	go func() {
		defer wg.Done()
		defer GinkgoRecover()

		var err error
		*createdOrgGUID, err = createOrgRaw(orgName)
		if err != nil {
			errChan <- err
			return
		}

		err = waitForAdminRoleBinding(*createdOrgGUID)
		if err != nil {
			errChan <- err
			return
		}
	}()
}

func createSpaceRaw(spaceName, orgGUID string) (string, error) {
	var space resource
	resp, err := adminClient.R().
		SetBody(resource{
			Name: spaceName,
			Relationships: relationships{
				"organization": relationship{Data: resource{GUID: orgGUID}},
			},
		}).
		SetResult(&space).
		Post("/v3/spaces")
	if err != nil {
		return "", err
	}

	if resp.StatusCode() != http.StatusCreated {
		return "", fmt.Errorf("expected status code %d, got %d", http.StatusCreated, resp.StatusCode())
	}

	return space.GUID, nil
}

func createSpace(spaceName, orgGUID string) string {
	spaceGUID, err := createSpaceRaw(spaceName, orgGUID)
	Expect(err).NotTo(HaveOccurred())
	Expect(waitForAdminRoleBinding(spaceGUID)).To(Succeed())

	return spaceGUID
}

func asyncCreateSpace(spaceName, orgGUID string, createdSpaceGUID *string, wg *sync.WaitGroup, errChan chan error) {
	go func() {
		defer wg.Done()
		defer GinkgoRecover()

		var err error
		*createdSpaceGUID, err = createSpaceRaw(spaceName, orgGUID)
		if err != nil {
			errChan <- err
			return
		}

		err = waitForAdminRoleBinding(*createdSpaceGUID)
		if err != nil {
			errChan <- err
			return
		}
	}()
}

// createRole creates an org or space role
// You should probably invoke this via createOrgRole or createSpaceRole
func createRole(roleName, kind, orgSpaceType, userName, orgSpaceGUID string) {
	rolesURL := apiServerRoot + apis.RolesEndpoint

	userOrServiceAccount := "user"
	if kind == rbacv1.ServiceAccountKind {
		userOrServiceAccount = "kubernetesServiceAccount"
	}

	payload := roleResource{
		Type: roleName,
		resource: resource{
			Relationships: relationships{
				userOrServiceAccount: relationship{Data: resource{GUID: userName}},
				orgSpaceType:         relationship{Data: resource{GUID: orgSpaceGUID}},
			},
		},
	}

	resp, err := adminClient.R().
		SetBody(payload).
		Post(rolesURL)

	ExpectWithOffset(2, err).NotTo(HaveOccurred())
	ExpectWithOffset(2, resp).To(HaveRestyStatusCode(http.StatusCreated))
}

func createOrgRole(roleName, kind, userName, orgGUID string) {
	createRole(roleName, kind, "organization", userName, orgGUID)
}

func createSpaceRole(roleName, kind, userName, spaceGUID string) {
	createRole(roleName, kind, "space", userName, spaceGUID)
}

func obtainServiceAccountToken(name string) string {
	var err error

	serviceAccount := corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: rootNamespace,
		},
	}
	err = k8sClient.Create(context.Background(), &serviceAccount)
	Expect(err).NotTo(HaveOccurred())

	Eventually(func() error {
		if err := k8sClient.Get(context.Background(), client.ObjectKeyFromObject(&serviceAccount), &serviceAccount); err != nil {
			return err
		}

		if len(serviceAccount.Secrets) != 1 {
			return fmt.Errorf("expected exactly 1 secret, got %d", len(serviceAccount.Secrets))
		}

		return nil
	}, "120s").Should(Succeed())

	tokenSecret := corev1.Secret{}
	Eventually(func() error {
		return k8sClient.Get(context.Background(), client.ObjectKey{Name: serviceAccount.Secrets[0].Name, Namespace: rootNamespace}, &tokenSecret)
	}).Should(Succeed())

	return string(tokenSecret.Data["token"])
}

func deleteServiceAccount(name string) {
	serviceAccount := corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: rootNamespace,
		},
	}

	Expect(k8sClient.Delete(context.Background(), &serviceAccount)).To(Succeed())
}

func obtainClientCert(name string) (*certsv1.CertificateSigningRequest, string) {
	privKey, err := rsa.GenerateKey(rand.Reader, 1024)
	Expect(err).NotTo(HaveOccurred())

	template := x509.CertificateRequest{
		Subject:            pkix.Name{CommonName: name},
		SignatureAlgorithm: x509.SHA256WithRSA,
	}

	csrBytes, err := x509.CreateCertificateRequest(rand.Reader, &template, privKey)
	Expect(err).NotTo(HaveOccurred())

	k8sCSR := &certsv1.CertificateSigningRequest{
		ObjectMeta: metav1.ObjectMeta{
			Name: uuid.NewString(),
		},
		Spec: certsv1.CertificateSigningRequestSpec{
			Request:    pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE REQUEST", Bytes: csrBytes}),
			SignerName: "kubernetes.io/kube-apiserver-client",
			Usages:     []certsv1.KeyUsage{certsv1.UsageClientAuth},
		},
	}

	Expect(k8sClient.Create(context.Background(), k8sCSR)).To(Succeed())

	k8sCSR.Status.Conditions = append(k8sCSR.Status.Conditions, certsv1.CertificateSigningRequestCondition{
		Type:   certsv1.CertificateApproved,
		Status: "True",
	})

	k8sCSR, err = clientset.CertificatesV1().CertificateSigningRequests().UpdateApproval(context.Background(), k8sCSR.Name, k8sCSR, metav1.UpdateOptions{})
	Expect(err).NotTo(HaveOccurred())

	var certPEM []byte
	Eventually(func() ([]byte, error) {
		err := k8sClient.Get(context.Background(), client.ObjectKeyFromObject(k8sCSR), k8sCSR)
		if err != nil {
			return nil, err
		}

		if len(k8sCSR.Status.Certificate) == 0 {
			return nil, nil
		}

		certPEM = k8sCSR.Status.Certificate

		return certPEM, nil
	}).ShouldNot(BeEmpty())

	buf := bytes.NewBuffer(certPEM)
	Expect(pem.Encode(buf, &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privKey),
	})).To(Succeed())

	return k8sCSR, base64.StdEncoding.EncodeToString(buf.Bytes())
}

func obtainAdminUserCert() string {
	crtBytes, err := base64.StdEncoding.DecodeString(mustHaveEnv("CF_ADMIN_CERT"))
	Expect(err).NotTo(HaveOccurred())
	keyBytes, err := base64.StdEncoding.DecodeString(mustHaveEnv("CF_ADMIN_KEY"))
	Expect(err).NotTo(HaveOccurred())

	return base64.StdEncoding.EncodeToString(append(crtBytes, keyBytes...))
}

func deleteCSR(csr *certsv1.CertificateSigningRequest) {
	Expect(k8sClient.Delete(context.Background(), csr)).To(Succeed())
}

func createApp(spaceGUID, name string) string {
	var app resource

	resp, err := adminClient.R().
		SetBody(appResource{
			resource: resource{
				Name:          name,
				Relationships: relationships{"space": {Data: resource{GUID: spaceGUID}}},
			},
		}).
		SetResult(&app).
		Post("/v3/apps")

	Expect(err).NotTo(HaveOccurred())
	Expect(resp).To(HaveRestyStatusCode(http.StatusCreated))

	return app.GUID
}

func getProcess(appGUID, processType string) string {
	var processList resourceList

	resp, err := adminClient.R().
		SetResult(&processList).
		Get("/v3/processes?app_guids=" + appGUID)

	Expect(err).NotTo(HaveOccurred())
	Expect(resp).To(HaveRestyStatusCode(http.StatusOK))
	Expect(processList.Resources).To(HaveLen(1))

	return processList.Resources[0].GUID
}

func createPackage(appGUID string) string {
	var pkg resource
	resp, err := adminClient.R().
		SetBody(packageResource{
			Type: "bits",
			resource: resource{
				Relationships: relationships{
					"app": relationship{Data: resource{GUID: appGUID}},
				},
			},
		}).
		SetResult(&pkg).
		Post("/v3/packages")

	Expect(err).NotTo(HaveOccurred())
	Expect(resp).To(HaveRestyStatusCode(http.StatusCreated))

	return pkg.GUID
}

func createBuild(packageGUID string) string {
	var build resource

	resp, err := adminClient.R().
		SetBody(buildResource{Package: resource{GUID: packageGUID}}).
		SetResult(&build).
		Post("/v3/builds")

	Expect(err).NotTo(HaveOccurred())
	Expect(resp).To(HaveRestyStatusCode(http.StatusCreated))

	return build.GUID
}

func waitForDroplet(buildGUID string) {
	Eventually(func() (*resty.Response, error) {
		resp, err := adminClient.R().
			Get("/v3/droplets/" + buildGUID)
		return resp, err
	}).Should(HaveRestyStatusCode(http.StatusOK))
}

func setCurrentDroplet(appGUID, dropletGUID string) {
	resp, err := adminClient.R().
		SetBody(dropletResource{Data: resource{GUID: dropletGUID}}).
		Patch("/v3/apps/" + appGUID + "/relationships/current_droplet")

	Expect(err).NotTo(HaveOccurred())
	Expect(resp).To(HaveRestyStatusCode(http.StatusOK))
}

func startApp(appGUID string) {
	resp, err := adminClient.R().
		Post("/v3/apps/" + appGUID + "/actions/start")

	Expect(err).NotTo(HaveOccurred())
	Expect(resp).To(HaveRestyStatusCode(http.StatusOK))
}

func uploadNodeApp(pkgGUID string) {
	resp, err := adminClient.R().
		SetFiles(map[string]string{
			"bits": "assets/node.zip",
		}).Post("/v3/packages/" + pkgGUID + "/upload")
	Expect(err).NotTo(HaveOccurred())
	Expect(resp).To(HaveRestyStatusCode(http.StatusOK))
}

// pushNodeApp creates a running node app in the given space
func pushNodeApp(spaceGUID string) string {
	appGUID := createApp(spaceGUID, generateGUID("app"))
	pkgGUID := createPackage(appGUID)
	uploadNodeApp(pkgGUID)
	buildGUID := createBuild(pkgGUID)
	waitForDroplet(buildGUID)
	setCurrentDroplet(appGUID, buildGUID)
	startApp(appGUID)

	return appGUID
}

func waitForAdminRoleBinding(namespace string) error {
	timeout := 10 * time.Second
	timeoutCtx, cancelFn := context.WithTimeout(context.Background(), timeout)
	defer cancelFn()

	watch, err := k8sClient.Watch(timeoutCtx, &rbacv1.RoleBindingList{}, client.InNamespace(namespace))
	if err != nil {
		return fmt.Errorf("failed to create a rolebindings watch on namespace %s: %v", namespace, err)
	}

	adminRolebindingPropagated := false
	for res := range watch.ResultChan() {
		roleBinding, ok := res.Object.(*rbacv1.RoleBinding)
		if !ok {
			// should never happen, but avoids panic above
			continue
		}
		if roleBinding.RoleRef.Name == "cf-k8s-controllers-admin" {
			watch.Stop()
			adminRolebindingPropagated = true
			break
		}

	}

	if !adminRolebindingPropagated {
		return fmt.Errorf("role binding to role 'cf-k8s-controllers-admin' has not been propagated within timeout period %d ms", timeout.Milliseconds())
	}

	return nil
}
