package e2e_test

import (
	"net/http"

	"github.com/go-resty/resty/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	rbacv1 "k8s.io/api/rbac/v1"
)

var _ = Describe("Package", func() {
	var (
		orgGUID   string
		spaceGUID string
		appGUID   string
		resp      *resty.Response
		result    packageResource
		resultErr cfErrs
	)

	BeforeEach(func() {
		orgGUID = createOrg(generateGUID("org"))
		createOrgRole("organization_user", rbacv1.UserKind, certUserName, orgGUID)

		spaceGUID = createSpace(generateGUID("space"), orgGUID)
		appGUID = createApp(spaceGUID, generateGUID("app"))

		result = packageResource{}
		resultErr = cfErrs{}
	})

	AfterEach(func() {
		deleteOrg(orgGUID)
	})

	Describe("Create", func() {
		JustBeforeEach(func() {
			var err error
			resp, err = certClient.R().
				SetBody(packageResource{
					Type: "bits",
					resource: resource{
						Relationships: relationships{
							"app": relationship{Data: resource{GUID: appGUID}},
						},
					},
				}).
				SetError(&resultErr).
				SetResult(&result).
				Post("/v3/packages")
			Expect(err).NotTo(HaveOccurred())
		})

		It("fails with a resource not found error", func() {
			Expect(resp).To(HaveRestyStatusCode(http.StatusUnprocessableEntity))
			Expect(resultErr.Errors).To(HaveLen(1))
			Expect(resultErr.Errors[0].Title).To(Equal("CF-UnprocessableEntity"))
			Expect(resultErr.Errors[0].Code).To(Equal(10008))
			Expect(resultErr.Errors[0].Detail).To(Equal("App is invalid. Ensure it exists and you have access to it."))
		})

		When("the user is a SpaceDeveloper", func() {
			BeforeEach(func() {
				createSpaceRole("space_developer", rbacv1.UserKind, certUserName, spaceGUID)
			})

			It("succeeds", func() {
				Expect(resultErr.Errors).To(HaveLen(0))
				Expect(resp).To(HaveRestyStatusCode(http.StatusCreated))
				Expect(result.GUID).ToNot(BeEmpty())
			})
		})

		When("the user is a SpaceManager (i.e. can get apps but cannot create packages)", func() {
			BeforeEach(func() {
				createSpaceRole("space_manager", rbacv1.UserKind, certUserName, spaceGUID)
			})

			It("fails with a forbidden error", func() {
				Expect(resp).To(HaveRestyStatusCode(http.StatusForbidden))
				Expect(resultErr.Errors).To(HaveLen(1))
				Expect(resultErr.Errors[0].Title).To(Equal("CF-NotAuthorized"))
				Expect(resultErr.Errors[0].Code).To(Equal(10003))
				Expect(resultErr.Errors[0].Detail).To(Equal("You are not authorized to perform the requested action"))
			})
		})
	})

	Describe("Upload", func() {
		var pkgGUID string

		BeforeEach(func() {
			pkgGUID = createPackage(appGUID)
		})

		JustBeforeEach(func() {
			var err error
			resp, err = certClient.R().
				SetFile("bits", "assets/node.zip").
				SetError(&resultErr).
				SetResult(&result).
				Post("/v3/packages/" + pkgGUID + "/upload")
			Expect(err).NotTo(HaveOccurred())
		})

		When("the user is a SpaceManager (i.e. can get apps but cannot update packages)", func() {
			BeforeEach(func() {
				createSpaceRole("space_manager", rbacv1.UserKind, certUserName, spaceGUID)
			})

			It("fails with a forbidden error", func() {
				Expect(resp).To(HaveRestyStatusCode(http.StatusForbidden))
				Expect(resultErr.Errors).To(HaveLen(1))
				Expect(resultErr.Errors[0].Title).To(Equal("CF-NotAuthorized"))
				Expect(resultErr.Errors[0].Code).To(Equal(10003))
				Expect(resultErr.Errors[0].Detail).To(Equal("You are not authorized to perform the requested action"))
			})
		})

		When("the user is a SpaceDeveloper", func() {
			BeforeEach(func() {
				createSpaceRole("space_developer", rbacv1.UserKind, certUserName, spaceGUID)
			})

			It("succeeds", func() {
				Expect(resp).To(HaveRestyStatusCode(http.StatusOK))
			})
		})
	})
})
