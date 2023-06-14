package handlers_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"

	apierrors "code.cloudfoundry.org/korifi/api/errors"
	"code.cloudfoundry.org/korifi/api/handlers"
	"code.cloudfoundry.org/korifi/api/handlers/fake"
	"code.cloudfoundry.org/korifi/api/payloads"
	"code.cloudfoundry.org/korifi/api/repositories"
	. "code.cloudfoundry.org/korifi/tests/matchers"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	rbacv1 "k8s.io/api/rbac/v1"
)

const (
	rolesBase = "/v3/roles"
)

var _ = Describe("Role", func() {
	var (
		apiHandler           *handlers.Role
		roleRepo             *fake.CFRoleRepository
		requestJSONValidator *fake.RequestJSONValidator
	)

	BeforeEach(func() {
		roleRepo = new(fake.CFRoleRepository)
		requestJSONValidator = new(fake.RequestJSONValidator)

		apiHandler = handlers.NewRole(*serverURL, roleRepo, requestJSONValidator)
		routerBuilder.LoadRoutes(apiHandler)
	})

	Describe("Create Role", func() {
		var roleCreate *payloads.RoleCreate

		BeforeEach(func() {
			roleRepo.CreateRoleReturns(repositories.RoleRecord{GUID: "role-guid"}, nil)
			roleCreate = &payloads.RoleCreate{
				Type: "space_developer",
				Relationships: payloads.RoleRelationships{
					User: &payloads.UserRelationship{
						Data: payloads.UserRelationshipData{
							Username: "my-user",
						},
					},
					Space: &payloads.Relationship{
						Data: &payloads.RelationshipData{
							GUID: "my-space",
						},
					},
				},
			}

			requestJSONValidator.DecodeAndValidateJSONPayloadStub = decodeAndValidateJSONPayloadStub(roleCreate)
		})

		JustBeforeEach(func() {
			req, err := http.NewRequestWithContext(ctx, "POST", rolesBase, strings.NewReader(""))
			Expect(err).NotTo(HaveOccurred())
			routerBuilder.Build().ServeHTTP(rr, req)
		})

		It("creates the role", func() {
			Expect(roleRepo.CreateRoleCallCount()).To(Equal(1))
			_, actualAuthInfo, roleMessage := roleRepo.CreateRoleArgsForCall(0)
			Expect(actualAuthInfo).To(Equal(authInfo))
			Expect(roleMessage.Type).To(Equal("space_developer"))
			Expect(roleMessage.Space).To(Equal("my-space"))
			Expect(roleMessage.User).To(Equal("my-user"))
			Expect(roleMessage.Kind).To(Equal(rbacv1.UserKind))

			Expect(rr).To(HaveHTTPStatus(http.StatusCreated))
			Expect(rr).To(HaveHTTPHeaderWithValue("Content-Type", "application/json"))
			Expect(rr).To(HaveHTTPBody(SatisfyAll(
				MatchJSONPath("$.guid", "role-guid"),
				MatchJSONPath("$.links.self.href", "https://api.example.org/v3/roles/role-guid"),
			)))
		})

		When("username is passed in the guid field", func() {
			BeforeEach(func() {
				roleCreate.Relationships.User.Data.Username = ""
				roleCreate.Relationships.User.Data.GUID = "my-user"
			})

			It("still works as guid and username are equivalent here", func() {
				Expect(roleRepo.CreateRoleCallCount()).To(Equal(1))
				_, _, roleMessage := roleRepo.CreateRoleArgsForCall(0)
				Expect(roleMessage.User).To(Equal("my-user"))
			})
		})

		When("the role is an organisation role", func() {
			BeforeEach(func() {
				roleCreate.Type = "organization_manager"
				roleCreate.Relationships.Space = nil
				roleCreate.Relationships.Organization = &payloads.Relationship{
					Data: &payloads.RelationshipData{
						GUID: "my-org",
					},
				}
			})

			It("invokes the role repo create function with expected parameters", func() {
				Expect(roleRepo.CreateRoleCallCount()).To(Equal(1))
				_, _, roleMessage := roleRepo.CreateRoleArgsForCall(0)
				Expect(roleMessage.Type).To(Equal("organization_manager"))
				Expect(roleMessage.Org).To(Equal("my-org"))
			})
		})

		When("the kind is a service account", func() {
			BeforeEach(func() {
				roleCreate.Relationships.User.Data.GUID = "system:serviceaccount:cf:my-user"
			})

			It("creates a service account role binding", func() {
				Expect(rr).To(HaveHTTPStatus(http.StatusCreated))
				Expect(roleRepo.CreateRoleCallCount()).To(Equal(1))
				_, _, roleRecord := roleRepo.CreateRoleArgsForCall(0)
				Expect(roleRecord.User).To(Equal("my-user"))
				Expect(roleRecord.ServiceAccountNamespace).To(Equal("cf"))
				Expect(roleRecord.Kind).To(Equal(rbacv1.ServiceAccountKind))
			})
		})

		When("the payload validator returns an error", func() {
			BeforeEach(func() {
				requestJSONValidator.DecodeAndValidateJSONPayloadReturns(apierrors.NewUnprocessableEntityError(errors.New("foo"), "some error"))
			})

			It("returns an error", func() {
				expectUnprocessableEntityError("some error")
			})
		})

		When("the org repo returns another error", func() {
			BeforeEach(func() {
				roleRepo.CreateRoleReturns(repositories.RoleRecord{}, errors.New("boom"))
			})

			It("returns unknown error", func() {
				expectUnknownError()
			})
		})
	})

	Describe("List roles", func() {
		var query string

		BeforeEach(func() {
			query = ""
			roleRepo.ListRolesReturns([]repositories.RoleRecord{
				{GUID: "role-1"},
				{GUID: "role-2"},
			}, nil)
		})

		JustBeforeEach(func() {
			req, err := http.NewRequestWithContext(ctx, "GET", rolesBase+query, nil)
			Expect(err).NotTo(HaveOccurred())
			routerBuilder.Build().ServeHTTP(rr, req)
		})

		It("lists roles", func() {
			Expect(roleRepo.ListRolesCallCount()).To(Equal(1))
			_, actualAuthInfo := roleRepo.ListRolesArgsForCall(0)
			Expect(actualAuthInfo).To(Equal(authInfo))

			Expect(rr).To(HaveHTTPStatus(http.StatusOK))
			Expect(rr).To(HaveHTTPHeaderWithValue("Content-Type", "application/json"))
			Expect(rr).To(HaveHTTPBody(SatisfyAll(
				MatchJSONPath("$.pagination.total_results", BeEquivalentTo(2)),
				MatchJSONPath("$.pagination.first.href", "https://api.example.org/v3/roles"),
				MatchJSONPath("$.resources", HaveLen(2)),
				MatchJSONPath("$.resources[0].guid", "role-1"),
				MatchJSONPath("$.resources[0].links.self.href", "https://api.example.org/v3/roles/role-1"),
				MatchJSONPath("$.resources[1].guid", "role-2"),
			)))
		})

		When("include is specified", func() {
			BeforeEach(func() {
				query = "?include=user"
			})

			It("does not fail but has no effect on the result", func() {
				Expect(rr).To(HaveHTTPStatus(http.StatusOK))
				Expect(rr).To(HaveHTTPBody(MatchJSONPath("$.pagination.total_results", BeEquivalentTo(2))))
			})
		})

		Describe("filtering and ordering", func() {
			BeforeEach(func() {
				roleRepo.ListRolesReturns([]repositories.RoleRecord{
					{GUID: "1", CreatedAt: "2022-01-23T17:08:22", UpdatedAt: "2022-01-22T17:09:00", Type: "a", Space: "space1", Org: "", User: "user1"},
					{GUID: "2", CreatedAt: "2022-01-24T17:08:22", UpdatedAt: "2022-01-21T17:09:00", Type: "b", Space: "space2", Org: "", User: "user1"},
					{GUID: "3", CreatedAt: "2022-01-22T17:08:22", UpdatedAt: "2022-01-24T17:09:00", Type: "c", Space: "", Org: "org1", User: "user1"},
					{GUID: "4", CreatedAt: "2022-01-21T17:08:22", UpdatedAt: "2022-01-23T17:09:00", Type: "c", Space: "", Org: "org2", User: "user2"},
				}, nil)
			})

			DescribeTable("filtering", func(filter string, expectedGUIDs ...any) {
				req, err := http.NewRequestWithContext(ctx, "GET", rolesBase+"?"+filter, nil)
				Expect(err).NotTo(HaveOccurred())
				rr = httptest.NewRecorder()
				routerBuilder.Build().ServeHTTP(rr, req)

				Expect(rr).To(HaveHTTPStatus(http.StatusOK))
				Expect(rr).To(HaveHTTPBody(MatchJSONPath("$.resources[*].guid", expectedGUIDs)))
			},
				Entry("no filter", "", "1", "2", "3", "4"),
				Entry("guids1", "guids=4", "4"),
				Entry("guids2", "guids=1,3", "1", "3"),
				Entry("types1", "types=a", "1"),
				Entry("types2", "types=b,c", "2", "3", "4"),
				Entry("space_guids1", "space_guids=space1", "1"),
				Entry("space_guids2", "space_guids=space1,space2", "1", "2"),
				Entry("organization_guids1", "organization_guids=org1", "3"),
				Entry("organization_guids2", "organization_guids=org1,org2", "3", "4"),
				Entry("user_guids1", "user_guids=user1", "1", "2", "3"),
				Entry("user_guids2", "user_guids=user1,user2", "1", "2", "3", "4"),
			)

			DescribeTable("ordering", func(order string, expectedGUIDs ...any) {
				req, err := http.NewRequestWithContext(ctx, "GET", rolesBase+"?order_by="+order, nil)
				Expect(err).NotTo(HaveOccurred())
				rr = httptest.NewRecorder()
				routerBuilder.Build().ServeHTTP(rr, req)

				Expect(rr).To(HaveHTTPStatus(http.StatusOK))
				Expect(rr).To(HaveHTTPBody(MatchJSONPath("$.resources[*].guid", expectedGUIDs)))
			},
				Entry("created_at ASC", "created_at", "4", "3", "1", "2"),
				Entry("created_at DESC", "-created_at", "2", "1", "3", "4"),
				Entry("updated_at ASC", "updated_at", "2", "1", "4", "3"),
				Entry("updated_at DESC", "-updated_at", "3", "4", "1", "2"),
			)
		})

		When("order_by is not a valid field", func() {
			BeforeEach(func() {
				query = "?order_by=not_valid"
			})

			It("returns an Unknown key error", func() {
				expectUnknownKeyError("The query parameter is invalid: Order by can only be: .*")
			})
		})

		When("calling the repository fails", func() {
			BeforeEach(func() {
				roleRepo.ListRolesReturns(nil, errors.New("boom"))
			})

			It("returns the error", func() {
				expectUnknownError()
			})
		})
	})

	Describe("delete a role", func() {
		BeforeEach(func() {
			roleRepo.GetRoleReturns(repositories.RoleRecord{
				GUID:  "role-guid",
				Space: "my-space",
				Org:   "",
			}, nil)
		})

		JustBeforeEach(func() {
			req, err := http.NewRequestWithContext(ctx, "DELETE", rolesBase+"/role-guid", nil)
			Expect(err).NotTo(HaveOccurred())
			routerBuilder.Build().ServeHTTP(rr, req)
		})

		It("deletes the role", func() {
			Expect(roleRepo.GetRoleCallCount()).To(Equal(1))
			_, actualAuthInfo, actualRoleGuid := roleRepo.GetRoleArgsForCall(0)
			Expect(actualAuthInfo).To(Equal(authInfo))
			Expect(actualRoleGuid).To(Equal("role-guid"))

			Expect(roleRepo.DeleteRoleCallCount()).To(Equal(1))
			_, actualAuthInfo, roleDeleteMsg := roleRepo.DeleteRoleArgsForCall(0)
			Expect(actualAuthInfo).To(Equal(authInfo))
			Expect(roleDeleteMsg).To(Equal(repositories.DeleteRoleMessage{
				GUID:  "role-guid",
				Space: "my-space",
			}))

			Expect(rr).To(HaveHTTPStatus(http.StatusAccepted))
			Expect(rr).To(HaveHTTPHeaderWithValue("Location", ContainSubstring("jobs/role.delete~role-guid")))
		})

		When("getting the role is forbidden", func() {
			BeforeEach(func() {
				roleRepo.GetRoleReturns(repositories.RoleRecord{}, apierrors.NewForbiddenError(nil, "Role"))
			})

			It("returns a not found error", func() {
				expectNotFoundError("Role")
			})
		})

		When("getting the role fails", func() {
			BeforeEach(func() {
				roleRepo.GetRoleReturns(repositories.RoleRecord{}, errors.New("get-role-err"))
			})

			It("returns the error", func() {
				expectUnknownError()
			})
		})

		When("deleting the role from the repo fails", func() {
			BeforeEach(func() {
				roleRepo.DeleteRoleReturns(errors.New("delete-role-err"))
			})

			It("returns the error", func() {
				expectUnknownError()
			})
		})
	})
})
