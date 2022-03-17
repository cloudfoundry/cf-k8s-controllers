package integration_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"code.cloudfoundry.org/cf-k8s-controllers/api/apis"
	"code.cloudfoundry.org/cf-k8s-controllers/api/config"
	"code.cloudfoundry.org/cf-k8s-controllers/api/repositories"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	hnsv1alpha2 "sigs.k8s.io/hierarchical-namespaces/api/v1alpha2"
)

var _ = Describe("Role", func() {
	var (
		apiHandler *apis.RoleHandler
		org, space *hnsv1alpha2.SubnamespaceAnchor
	)

	BeforeEach(func() {
		// NOTE: This creates an arbitrary mapping that is not representative of production mapping
		roleMappings := map[string]config.Role{
			"space_developer":      {Name: "cf-space-developer"},
			"organization_manager": {Name: "cf-organization-manager"},
			"cf_user":              {Name: "cf-user"},
		}
		roleRepo := repositories.NewRoleRepo(k8sClient, clientFactory, nsPermissions, rootNamespace, roleMappings)
		decoderValidator, err := apis.NewDefaultDecoderValidator()
		Expect(err).NotTo(HaveOccurred())

		apiHandler = apis.NewRoleHandler(*serverURL, roleRepo, decoderValidator)
		apiHandler.RegisterRoutes(router)

		org = createOrgAnchorAndNamespace(ctx, rootNamespace, generateGUID())
		space = createSpaceAnchorAndNamespace(ctx, org.Name, "spacename-"+generateGUID())
	})

	Describe("creation", func() {
		var (
			bodyTemplate  string
			roleName      string
			userGUID      string
			orgSpaceLabel string
			orgSpaceGUID  string
		)

		BeforeEach(func() {
			roleName = ""
			userGUID = generateGUID()
			orgSpaceLabel = ""
			orgSpaceGUID = ""

			bodyTemplate = `{
              "type": %q,
              "relationships": {
                "user": {"data": {"guid": %q}},
                %q: {"data": {"guid": %q}}
              }
            }`
		})

		createTheRole := func(respRecorder *httptest.ResponseRecorder) {
			requestBody := fmt.Sprintf(bodyTemplate, roleName, userGUID, orgSpaceLabel, orgSpaceGUID)

			createRoleReq, err := http.NewRequestWithContext(ctx, "POST", serverURI("/v3/roles"), strings.NewReader(requestBody))
			Expect(err).NotTo(HaveOccurred())

			createRoleReq.Header.Add("Content-type", "application/json")
			router.ServeHTTP(respRecorder, createRoleReq)
		}

		JustBeforeEach(func() {
			createTheRole(rr)
		})

		Describe("creating an org role", func() {
			BeforeEach(func() {
				roleName = "organization_manager"
				orgSpaceLabel = "organization"
				orgSpaceGUID = org.Name
			})

			It("fails when the user is not allowed to create org roles", func() {
				Expect(rr).To(HaveHTTPStatus(http.StatusForbidden))
			})

			When("the user is admin", func() {
				BeforeEach(func() {
					createRoleBinding(ctx, userName, adminRole.Name, org.Name)
				})

				It("succeeds", func() {
					Expect(rr).To(HaveHTTPStatus(http.StatusCreated))
				})

				When("the role already exists for the user", func() {
					BeforeEach(func() {
						duplicateRr := httptest.NewRecorder()
						createTheRole(duplicateRr)
						Expect(duplicateRr).To(HaveHTTPStatus(http.StatusCreated))
					})

					It("returns unprocessable entity", func() {
						Expect(rr).To(HaveHTTPStatus(http.StatusUnprocessableEntity))
					})
				})
			})
		})

		Describe("creating a space role", func() {
			BeforeEach(func() {
				roleName = "space_developer"
				orgSpaceLabel = "space"
				orgSpaceGUID = space.Name
			})

			It("fails when the user doesn't have a role binding in the space's org", func() {
				Expect(rr).To(HaveHTTPStatus(http.StatusUnprocessableEntity))
			})

			When("the user is an org user", func() {
				BeforeEach(func() {
					createRoleBinding(ctx, userGUID, orgUserRole.Name, org.Name)
				})

				It("fails when the user is not allowed to create space roles", func() {
					Expect(rr).To(HaveHTTPStatus(http.StatusForbidden))
				})

				When("the user is admin", func() {
					BeforeEach(func() {
						createRoleBinding(ctx, userName, adminRole.Name, space.Name)
					})

					It("succeeds", func() {
						Expect(rr).To(HaveHTTPStatus(http.StatusCreated))
					})

					When("the role already exists for the user", func() {
						BeforeEach(func() {
							duplicateRr := httptest.NewRecorder()
							createTheRole(duplicateRr)
							Expect(duplicateRr).To(HaveHTTPStatus(http.StatusCreated))
						})

						It("returns unprocessable entity", func() {
							Expect(rr).To(HaveHTTPStatus(http.StatusUnprocessableEntity))
						})
					})
				})
			})
		})
	})
})
