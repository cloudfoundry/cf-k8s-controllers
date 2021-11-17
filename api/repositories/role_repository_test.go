package repositories_test

import (
	"context"
	"errors"
	"time"

	"code.cloudfoundry.org/cf-k8s-controllers/api/repositories"
	"code.cloudfoundry.org/cf-k8s-controllers/api/repositories/fake"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	hnsv1alpha2 "sigs.k8s.io/hierarchical-namespaces/api/v1alpha2"
)

var _ = Describe("RoleRepository", func() {
	var (
		ctx                 context.Context
		rootNamespace       string
		roleRepo            *repositories.RoleRepo
		roleRecord          repositories.RoleRecord
		orgAnchor           *hnsv1alpha2.SubnamespaceAnchor
		createdRole         repositories.RoleRecord
		authorizedInChecker *fake.AuthorizedInChecker
		createErr           error
	)

	BeforeEach(func() {
		rootNamespace = uuid.NewString()
		ctx = context.Background()
		authorizedInChecker = new(fake.AuthorizedInChecker)
		Expect(k8sClient.Create(context.Background(), &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: rootNamespace}})).To(Succeed())
		roleRepo = repositories.NewRoleRepo(k8sClient, authorizedInChecker, map[string]string{
			"space_developer":      "cf-space-dev-role",
			"organization_manager": "cf-org-mgr-role",
		})

		roleRecord = repositories.RoleRecord{}
		orgAnchor = createOrgAnchorAndNamespace(ctx, rootNamespace, uuid.NewString())
	})

	Describe("Create Org Role", func() {
		var (
			createdRole repositories.RoleRecord
			createErr   error
		)

		BeforeEach(func() {
			roleRecord = repositories.RoleRecord{
				GUID: uuid.NewString(),
				Type: "organization_manager",
				User: "my-user",
				Kind: rbacv1.UserKind,
				Org:  orgAnchor.Name,
			}
		})

		JustBeforeEach(func() {
			createdRole, createErr = roleRepo.CreateRole(ctx, roleRecord)
		})

		It("succeeds", func() {
			Expect(createErr).NotTo(HaveOccurred())
		})

		It("creates a role binding in the org namespace", func() {
			roleBindingList := rbacv1.RoleBindingList{}
			Expect(k8sClient.List(ctx, &roleBindingList, client.InNamespace(orgAnchor.Name))).To(Succeed())
			Expect(roleBindingList.Items).To(HaveLen(1))

			roleBinding := roleBindingList.Items[0]

			// Sha256 sum of "organization_manager::my-user"
			Expect(roleBinding.Name).To(Equal("cf-d024ad51b9896f27fab865db894beb14992af05fbbc785bbf90d8706bc95b21b"))
			Expect(roleBinding.Labels).To(HaveKeyWithValue(repositories.RoleGuidLabel, roleRecord.GUID))
			Expect(roleBinding.Labels).To(HaveKeyWithValue(repositories.RoleUserLabel, roleRecord.User))
			Expect(roleBinding.Labels).To(HaveKeyWithValue(repositories.RoleTypeLabel, roleRecord.Type))
			Expect(roleBinding.RoleRef.Kind).To(Equal("ClusterRole"))
			Expect(roleBinding.RoleRef.Name).To(Equal("cf-org-mgr-role"))
			Expect(roleBinding.Subjects).To(HaveLen(1))
			Expect(roleBinding.Subjects[0].Kind).To(Equal(rbacv1.UserKind))
			Expect(roleBinding.Subjects[0].Name).To(Equal("my-user"))
		})

		It("updated the create/updated timestamps", func() {
			Expect(createdRole.CreatedAt).To(BeTemporally("~", time.Now(), 2*time.Second))
			Expect(createdRole.UpdatedAt).To(BeTemporally("~", time.Now(), 2*time.Second))
			Expect(createdRole.CreatedAt).To(Equal(createdRole.UpdatedAt))
		})

		When("using a service account identity", func() {
			BeforeEach(func() {
				roleRecord.Kind = rbacv1.ServiceAccountKind
			})

			It("succeeds and uses a service account subject kind", func() {
				Expect(createErr).NotTo(HaveOccurred())

				roleBindingList := rbacv1.RoleBindingList{}
				Expect(k8sClient.List(ctx, &roleBindingList, client.InNamespace(orgAnchor.Name))).To(Succeed())
				Expect(roleBindingList.Items).To(HaveLen(1))

				roleBinding := roleBindingList.Items[0]

				Expect(roleBinding.Subjects).To(HaveLen(1))
				Expect(roleBinding.Subjects[0].Kind).To(Equal(rbacv1.ServiceAccountKind))
			})
		})

		When("the org does not exist", func() {
			BeforeEach(func() {
				roleRecord.Org = "i-do-not-exist"
			})

			It("returns an error", func() {
				Expect(k8serrors.IsNotFound(createErr)).To(BeTrue())
			})
		})

		When("the role type is invalid", func() {
			BeforeEach(func() {
				roleRecord.Type = "i-am-invalid"
			})

			It("returns an error", func() {
				Expect(createErr).To(MatchError(ContainSubstring("invalid role type")))
			})
		})

		When("the user is already bound to that role", func() {
			It("returns an error", func() {
				anotherRoleRecord := repositories.RoleRecord{
					GUID: uuid.NewString(),
					Type: "organization_manager",
					User: "my-user",
					Kind: rbacv1.UserKind,
					Org:  roleRecord.Org,
				}
				_, createErr = roleRepo.CreateRole(ctx, anotherRoleRecord)
				Expect(createErr).To(Equal(repositories.ErrorDuplicateRoleBinding))
			})
		})
	})

	Describe("Create Space Role", func() {
		var spaceAnchor *hnsv1alpha2.SubnamespaceAnchor

		BeforeEach(func() {
			authorizedInChecker.AuthorizedInReturns(true, nil)
			spaceAnchor = createSpaceAnchorAndNamespace(ctx, orgAnchor.Name, uuid.NewString())

			Expect(k8sClient.Create(context.Background(), &rbacv1.RoleBinding{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: orgAnchor.Name,
				},
				Subjects: []rbacv1.Subject{
					{
						Kind: rbacv1.UserKind,
						Name: "my-user",
					},
				},
				RoleRef: rbacv1.RoleRef{
					Kind: "ClusterRole",
					Name: "org_user",
				},
			})).To(Succeed())

			roleRecord = repositories.RoleRecord{
				GUID:  uuid.NewString(),
				Type:  "space_developer",
				User:  "my-user",
				Kind:  rbacv1.UserKind,
				Space: spaceAnchor.Name,
			}
		})

		JustBeforeEach(func() {
			createdRole, createErr = roleRepo.CreateRole(ctx, roleRecord)
		})

		It("succeeds", func() {
			Expect(createErr).NotTo(HaveOccurred())
		})

		It("creates a role binding in the space namespace", func() {
			roleBindingList := rbacv1.RoleBindingList{}
			Expect(k8sClient.List(ctx, &roleBindingList, client.InNamespace(spaceAnchor.Name))).To(Succeed())
			Expect(roleBindingList.Items).To(HaveLen(1))

			roleBinding := roleBindingList.Items[0]

			// Sha256 sum of "space_developer::my-user"
			Expect(roleBinding.Name).To(Equal("cf-1b2399803c0978bcf9669095590b5f423215e053200e67d7d517db76fdedf197"))
			Expect(roleBinding.Labels).To(HaveKeyWithValue(repositories.RoleGuidLabel, roleRecord.GUID))
			Expect(roleBinding.Labels).To(HaveKeyWithValue(repositories.RoleUserLabel, roleRecord.User))
			Expect(roleBinding.Labels).To(HaveKeyWithValue(repositories.RoleTypeLabel, roleRecord.Type))
			Expect(roleBinding.RoleRef.Kind).To(Equal("ClusterRole"))
			Expect(roleBinding.RoleRef.Name).To(Equal("cf-space-dev-role"))
			Expect(roleBinding.Subjects).To(HaveLen(1))
			Expect(roleBinding.Subjects[0].Kind).To(Equal("User"))
			Expect(roleBinding.Subjects[0].Name).To(Equal("my-user"))
		})

		It("verifies that the user has a role in the parent org", func() {
			Expect(authorizedInChecker.AuthorizedInCallCount()).To(Equal(1))
			_, userIdentity, org := authorizedInChecker.AuthorizedInArgsForCall(0)
			Expect(userIdentity.Name).To(Equal("my-user"))
			Expect(userIdentity.Kind).To(Equal(rbacv1.UserKind))
			Expect(org).To(Equal(orgAnchor.Name))
		})

		It("updated the create/updated timestamps", func() {
			Expect(createdRole.CreatedAt).To(BeTemporally("~", time.Now(), 2*time.Second))
			Expect(createdRole.UpdatedAt).To(BeTemporally("~", time.Now(), 2*time.Second))
			Expect(createdRole.CreatedAt).To(Equal(createdRole.UpdatedAt))
		})

		When("getting the parent org fails", func() {
			BeforeEach(func() {
				namespace := &corev1.Namespace{
					ObjectMeta: metav1.ObjectMeta{
						Name: spaceAnchor.Name,
						Annotations: map[string]string{
							hnsv1alpha2.SubnamespaceOf: orgAnchor.Name,
						},
					},
				}
				nsCopy := namespace.DeepCopy()
				nsCopy.Annotations[hnsv1alpha2.SubnamespaceOf] = ""

				Expect(k8sClient.Patch(ctx, nsCopy, client.MergeFrom(namespace))).To(Succeed())
			})

			It("returns an error", func() {
				Expect(createErr).To(MatchError(ContainSubstring("does not have a parent")))
			})
		})

		When("checking an org role exists fails", func() {
			BeforeEach(func() {
				authorizedInChecker.AuthorizedInReturns(false, errors.New("boom!"))
			})

			It("returns an error", func() {
				Expect(createErr).To(MatchError(ContainSubstring("failed to check for role in parent org")))
			})
		})

		When("the space does not exist", func() {
			BeforeEach(func() {
				roleRecord.Space = "i-do-not-exist"
			})

			It("returns an error", func() {
				Expect(k8serrors.IsNotFound(createErr)).To(BeTrue())
			})
		})

		When("the role type is invalid", func() {
			BeforeEach(func() {
				roleRecord.Type = "i-am-invalid"
			})

			It("returns an error", func() {
				Expect(createErr).To(MatchError(ContainSubstring("invalid role type")))
			})
		})

		When("the user is already bound to that role", func() {
			It("returns an error", func() {
				anotherRoleRecord := repositories.RoleRecord{
					GUID:  uuid.NewString(),
					Type:  "space_developer",
					User:  "my-user",
					Kind:  rbacv1.UserKind,
					Space: roleRecord.Space,
				}
				_, createErr = roleRepo.CreateRole(ctx, anotherRoleRecord)
				Expect(createErr).To(Equal(repositories.ErrorDuplicateRoleBinding))
			})
		})

		When("the user does not have a role in the parent organization", func() {
			BeforeEach(func() {
				authorizedInChecker.AuthorizedInReturns(false, nil)
			})

			It("returns an error", func() {
				Expect(createErr).To(Equal(repositories.ErrorMissingRoleBindingInParentOrg))
			})
		})
	})
})
