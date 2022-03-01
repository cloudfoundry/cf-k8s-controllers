package e2e_test

import (
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	rbacv1 "k8s.io/api/rbac/v1"
)

var _ = Describe("Service Bindings", func() {
	Describe("List", func() {
		var (
			appGUID      string
			orgGUID      string
			spaceGUID    string
			instanceGUID string
			httpResp     *resty.Response
			httpError    error
			queryString  string
			result       resourceListWithInclusion
		)

		BeforeEach(func() {
			orgGUID = createOrg(generateGUID("org"))
			time.Sleep(time.Second) // this appears to reduce flakes, but should be removed once we have better logic to determine org/space readiness
			spaceGUID = createSpace(generateGUID("space1"), orgGUID)
			time.Sleep(time.Second)
			createOrgRole("organization_user", rbacv1.UserKind, certUserName, orgGUID)
			instanceGUID = createServiceInstance(spaceGUID, generateGUID("service-instance"))
			appGUID = createApp(spaceGUID, generateGUID("app"))
			createServiceBinding(appGUID, instanceGUID)

			queryString = ""
			result = resourceListWithInclusion{}
		})

		JustBeforeEach(func() {
			httpResp, httpError = certClient.R().SetResult(&result).Get("/v3/service_credential_bindings" + queryString)
		})

		AfterEach(func() {
			deleteOrg(orgGUID)
		})

		It("Returns an empty list", func() {
			Expect(httpError).NotTo(HaveOccurred())
			Expect(httpResp).To(HaveRestyStatusCode(http.StatusOK))
			Expect(result.Resources).To(HaveLen(0))
		})

		When("the user has space manager role", func() {
			BeforeEach(func() {
				createSpaceRole("space_manager", rbacv1.UserKind, certUserName, spaceGUID)
			})

			It("succeeds", func() {
				Expect(httpError).NotTo(HaveOccurred())
				Expect(httpResp).To(HaveRestyStatusCode(http.StatusOK))
				Expect(result.Resources).To(HaveLen(1))
			})
		})

		When("the user has space developer role", func() {
			BeforeEach(func() {
				createSpaceRole("space_developer", rbacv1.UserKind, certUserName, spaceGUID)
			})

			It("succeeds", func() {
				Expect(httpError).NotTo(HaveOccurred())
				Expect(httpResp).To(HaveRestyStatusCode(http.StatusOK))
				Expect(result.Resources).To(HaveLen(1))
			})

			It("doesn't return anything in the 'included' list", func() {
				Expect(result.Included).To(BeNil())
			})

			When("the 'include=app' querystring is set", func() {
				BeforeEach(func() {
					queryString = `?include=app`
				})

				It("returns an app in the 'included' list", func() {
					Expect(httpError).NotTo(HaveOccurred())
					Expect(httpResp).To(HaveRestyStatusCode(http.StatusOK))
					Expect(result.Resources).To(HaveLen(1))
					Expect(result.Included).NotTo(BeNil())
					Expect(result.Included.Apps).To(ConsistOf(
						MatchFields(IgnoreExtras, Fields{"GUID": Equal(appGUID)}),
					))
				})
			})
		})
	})
})
