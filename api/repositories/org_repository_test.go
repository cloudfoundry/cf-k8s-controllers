package repositories_test

import (
	"context"
	"sort"
	"time"

	"code.cloudfoundry.org/cf-k8s-controllers/api/authorization"
	"code.cloudfoundry.org/cf-k8s-controllers/api/repositories"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"sigs.k8s.io/controller-runtime/pkg/client"
	hnsv1alpha2 "sigs.k8s.io/hierarchical-namespaces/api/v1alpha2"
)

var _ = Describe("OrgRepository", func() {
	var (
		ctx     context.Context
		orgRepo *repositories.OrgRepo
		info    authorization.Info
	)

	BeforeEach(func() {
		ctx = context.Background()
		Expect(k8sClient.Create(context.Background(), &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: rootNamespace}})).To(Succeed())
		orgRepo = repositories.NewOrgRepo(rootNamespace, k8sClient, time.Millisecond*500)
		info = authorization.Info{Token: "a-token"}
	})
	Describe("Create", func() {
		updateStatus := func(anchorNamespace, anchorName string) {
			defer GinkgoRecover()

			anchor := &hnsv1alpha2.SubnamespaceAnchor{}
			for {
				err := k8sClient.Get(ctx, client.ObjectKey{Namespace: anchorNamespace, Name: anchorName}, anchor)
				if err == nil {
					break
				}

				time.Sleep(time.Millisecond * 100)
				continue
			}

			newAnchor := anchor.DeepCopy()
			newAnchor.Status.State = hnsv1alpha2.Ok
			Expect(k8sClient.Patch(ctx, newAnchor, client.MergeFrom(anchor))).To(Succeed())
		}

		Describe("Org", func() {
			It("creates a subnamespace anchor in the root namespace", func() {
				go updateStatus(rootNamespace, "some-guid")
				org, err := orgRepo.CreateOrg(ctx, info, repositories.CreateOrgMessage{
					GUID: "some-guid",
					Name: "our-org",
				})
				Expect(err).NotTo(HaveOccurred())

				namesRequirement, err := labels.NewRequirement(repositories.OrgNameLabel, selection.Equals, []string{"our-org"})
				Expect(err).NotTo(HaveOccurred())
				anchorList := hnsv1alpha2.SubnamespaceAnchorList{}
				err = k8sClient.List(ctx, &anchorList, client.InNamespace(rootNamespace), client.MatchingLabelsSelector{
					Selector: labels.NewSelector().Add(*namesRequirement),
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(anchorList.Items).To(HaveLen(1))

				Expect(org.Name).To(Equal("our-org"))
				Expect(org.GUID).To(Equal("some-guid"))
				Expect(org.CreatedAt).To(BeTemporally("~", time.Now(), 2*time.Second))
				Expect(org.UpdatedAt).To(BeTemporally("~", time.Now(), 2*time.Second))
			})

			When("the org isn't ready in the timeout", func() {
				It("returns an error", func() {
					// we do not call updateStatus() to set state = ok
					_, err := orgRepo.CreateOrg(ctx, info, repositories.CreateOrgMessage{
						GUID: "some-guid",
						Name: "our-org",
					})
					Expect(err).To(MatchError(ContainSubstring("did not get state 'ok'")))
				})
			})

			When("the client fails to create the org", func() {
				It("returns an error", func() {
					_, err := orgRepo.CreateOrg(ctx, info, repositories.CreateOrgMessage{
						Name: "this-string-has-illegal-characters-ц",
					})
					Expect(err).To(HaveOccurred())
				})
			})
		})

		Describe("Space", func() {
			var org *hnsv1alpha2.SubnamespaceAnchor
			var spaceGUID string
			imageRegistryCredentials := "image-registry-credentials"

			BeforeEach(func() {
				spaceGUID = generateGUID()
				org = createOrgAnchorAndNamespace(ctx, rootNamespace, "org")
				// In the absence of HNC reconciling the SubnamespaceAnchor into a namespace, we must manually create
				// for subsequent use by the Repository createSpace function.
				_ = createNamespace(ctx, "org", spaceGUID)
			})

			It("creates a Space", func() {
				go updateStatus(org.Name, spaceGUID)

				space, err := orgRepo.CreateSpace(ctx, info, repositories.CreateSpaceMessage{
					GUID:                     spaceGUID,
					Name:                     "our-space",
					OrganizationGUID:         org.Name,
					ImageRegistryCredentials: imageRegistryCredentials,
				})
				Expect(err).NotTo(HaveOccurred())

				By("Creating a SubnamespaceAnchor in the Org namespace", func() {
					var namesRequirement *labels.Requirement
					namesRequirement, err = labels.NewRequirement(repositories.SpaceNameLabel, selection.Equals, []string{"our-space"})
					Expect(err).NotTo(HaveOccurred())
					anchorList := hnsv1alpha2.SubnamespaceAnchorList{}
					err = k8sClient.List(ctx, &anchorList, client.InNamespace(org.Name), client.MatchingLabelsSelector{
						Selector: labels.NewSelector().Add(*namesRequirement),
					})
					Expect(err).NotTo(HaveOccurred())
					Expect(anchorList.Items).To(HaveLen(1))

					Expect(space.Name).To(Equal("our-space"))
					Expect(space.GUID).To(Equal(spaceGUID))
					Expect(space.CreatedAt).To(BeTemporally("~", time.Now(), 2*time.Second))
					Expect(space.UpdatedAt).To(BeTemporally("~", time.Now(), 2*time.Second))
				})

				By("Creating ServiceAccounts in the Space namespace", func() {
					serviceAccountList := corev1.ServiceAccountList{}
					Eventually(func() []corev1.ServiceAccount {
						err = k8sClient.List(ctx, &serviceAccountList, client.InNamespace(spaceGUID))
						if err != nil {
							return []corev1.ServiceAccount{}
						}
						return serviceAccountList.Items
					}, timeCheckThreshold*time.Second, 250*time.Millisecond).Should(HaveLen(2), "could not find the service accounts created by the repo")
					Expect(err).NotTo(HaveOccurred())

					sort.Slice(serviceAccountList.Items, func(i, j int) bool {
						return serviceAccountList.Items[i].Name < serviceAccountList.Items[j].Name
					})
					serviceAccount := serviceAccountList.Items[0]
					Expect(serviceAccount.Name).To(Equal("eirini"))
					serviceAccount = serviceAccountList.Items[1]
					Expect(serviceAccount.Name).To(Equal("kpack-service-account"))
					Expect(serviceAccount.ImagePullSecrets).To(ConsistOf(corev1.LocalObjectReference{Name: imageRegistryCredentials}))
					Expect(serviceAccount.Secrets).To(ConsistOf(corev1.ObjectReference{Name: imageRegistryCredentials}))
				})
			})

			When("the space isn't ready in the timeout", func() {
				It("returns an error", func() {
					// we do not call updateStatus() to set state = ok
					_, err := orgRepo.CreateSpace(ctx, info, repositories.CreateSpaceMessage{
						GUID:             "some-guid",
						Name:             "our-org",
						OrganizationGUID: org.Name,
					})
					Expect(err).To(MatchError(ContainSubstring("did not get state 'ok'")))
				})
			})

			When("the client fails to create the space", func() {
				It("returns an error", func() {
					_, err := orgRepo.CreateSpace(ctx, info, repositories.CreateSpaceMessage{
						GUID:             "some-guid",
						Name:             "this-string-has-illegal-characters-ц",
						OrganizationGUID: org.Name,
					})
					Expect(err).To(HaveOccurred())
				})
			})
		})
	})

	Describe("List", func() {
		var (
			ctx context.Context

			org1Anchor, org2Anchor, org3Anchor *hnsv1alpha2.SubnamespaceAnchor
		)

		BeforeEach(func() {
			ctx = context.Background()

			org1Anchor = createOrgAnchorAndNamespace(ctx, rootNamespace, "org1")
			org2Anchor = createOrgAnchorAndNamespace(ctx, rootNamespace, "org2")
			org3Anchor = createOrgAnchorAndNamespace(ctx, rootNamespace, "org3")
		})

		Describe("Orgs", func() {
			It("returns the 3 orgs", func() {
				orgs, err := orgRepo.ListOrgs(ctx, info, nil)
				Expect(err).NotTo(HaveOccurred())

				Expect(orgs).To(ConsistOf(
					repositories.OrgRecord{
						Name:      "org1",
						CreatedAt: org1Anchor.CreationTimestamp.Time,
						UpdatedAt: org1Anchor.CreationTimestamp.Time,
						GUID:      org1Anchor.Name,
					},
					repositories.OrgRecord{
						Name:      "org2",
						CreatedAt: org2Anchor.CreationTimestamp.Time,
						UpdatedAt: org2Anchor.CreationTimestamp.Time,
						GUID:      org2Anchor.Name,
					},
					repositories.OrgRecord{
						Name:      "org3",
						CreatedAt: org3Anchor.CreationTimestamp.Time,
						UpdatedAt: org3Anchor.CreationTimestamp.Time,
						GUID:      org3Anchor.Name,
					},
				))
			})

			When("the org anchor is not ready", func() {
				BeforeEach(func() {
					org1AnchorCopy := org1Anchor.DeepCopy()
					org1AnchorCopy.Status.State = hnsv1alpha2.Missing
					Expect(k8sClient.Patch(ctx, org1AnchorCopy, client.MergeFrom(org1Anchor))).To(Succeed())
				})

				It("does not list it", func() {
					orgs, err := orgRepo.ListOrgs(ctx, info, nil)
					Expect(err).NotTo(HaveOccurred())

					Expect(orgs).NotTo(ContainElement(
						repositories.OrgRecord{
							Name:      "org1",
							CreatedAt: org1Anchor.CreationTimestamp.Time,
							UpdatedAt: org1Anchor.CreationTimestamp.Time,
							GUID:      org1Anchor.Name,
						},
					))
				})
			})

			When("we filter for org1 and org3", func() {
				It("returns just those", func() {
					orgs, err := orgRepo.ListOrgs(ctx, info, []string{"org1", "org3"})
					Expect(err).NotTo(HaveOccurred())

					Expect(orgs).To(ConsistOf(
						repositories.OrgRecord{
							Name:      "org1",
							CreatedAt: org1Anchor.CreationTimestamp.Time,
							UpdatedAt: org1Anchor.CreationTimestamp.Time,
							GUID:      org1Anchor.Name,
						},
						repositories.OrgRecord{
							Name:      "org3",
							CreatedAt: org3Anchor.CreationTimestamp.Time,
							UpdatedAt: org3Anchor.CreationTimestamp.Time,
							GUID:      org3Anchor.Name,
						},
					))
				})
			})
		})

		Describe("Spaces", func() {
			var space11Anchor, space12Anchor, space21Anchor, space22Anchor, space31Anchor, space32Anchor *hnsv1alpha2.SubnamespaceAnchor

			BeforeEach(func() {
				space11Anchor = createSpaceAnchorAndNamespace(ctx, org1Anchor.Name, "space1")
				space12Anchor = createSpaceAnchorAndNamespace(ctx, org1Anchor.Name, "space2")

				space21Anchor = createSpaceAnchorAndNamespace(ctx, org2Anchor.Name, "space1")
				space22Anchor = createSpaceAnchorAndNamespace(ctx, org2Anchor.Name, "space3")

				space31Anchor = createSpaceAnchorAndNamespace(ctx, org3Anchor.Name, "space1")
				space32Anchor = createSpaceAnchorAndNamespace(ctx, org3Anchor.Name, "space4")
			})

			It("returns the 6 spaces", func() {
				spaces, err := orgRepo.ListSpaces(ctx, info, []string{}, []string{})
				Expect(err).NotTo(HaveOccurred())

				Expect(spaces).To(ConsistOf(
					repositories.SpaceRecord{
						Name:             "space1",
						CreatedAt:        space11Anchor.CreationTimestamp.Time,
						UpdatedAt:        space11Anchor.CreationTimestamp.Time,
						GUID:             space11Anchor.Name,
						OrganizationGUID: org1Anchor.Name,
					},
					repositories.SpaceRecord{
						Name:             "space2",
						CreatedAt:        space12Anchor.CreationTimestamp.Time,
						UpdatedAt:        space12Anchor.CreationTimestamp.Time,
						GUID:             space12Anchor.Name,
						OrganizationGUID: org1Anchor.Name,
					},
					repositories.SpaceRecord{
						Name:             "space1",
						CreatedAt:        space21Anchor.CreationTimestamp.Time,
						UpdatedAt:        space21Anchor.CreationTimestamp.Time,
						GUID:             space21Anchor.Name,
						OrganizationGUID: org2Anchor.Name,
					},
					repositories.SpaceRecord{
						Name:             "space3",
						CreatedAt:        space22Anchor.CreationTimestamp.Time,
						UpdatedAt:        space22Anchor.CreationTimestamp.Time,
						GUID:             space22Anchor.Name,
						OrganizationGUID: org2Anchor.Name,
					},
					repositories.SpaceRecord{
						Name:             "space1",
						CreatedAt:        space31Anchor.CreationTimestamp.Time,
						UpdatedAt:        space31Anchor.CreationTimestamp.Time,
						GUID:             space31Anchor.Name,
						OrganizationGUID: org3Anchor.Name,
					},
					repositories.SpaceRecord{
						Name:             "space4",
						CreatedAt:        space32Anchor.CreationTimestamp.Time,
						UpdatedAt:        space32Anchor.CreationTimestamp.Time,
						GUID:             space32Anchor.Name,
						OrganizationGUID: org3Anchor.Name,
					},
				))
			})

			When("the space anchor is not ready", func() {
				BeforeEach(func() {
					space11AnchorCopy := space11Anchor.DeepCopy()
					space11AnchorCopy.Status.State = hnsv1alpha2.Missing
					Expect(k8sClient.Patch(ctx, space11AnchorCopy, client.MergeFrom(space11Anchor))).To(Succeed())
				})

				It("does not list it", func() {
					spaces, err := orgRepo.ListSpaces(ctx, info, []string{}, []string{})
					Expect(err).NotTo(HaveOccurred())

					Expect(spaces).NotTo(ContainElement(
						repositories.SpaceRecord{
							Name:             "space1",
							CreatedAt:        space11Anchor.CreationTimestamp.Time,
							UpdatedAt:        space11Anchor.CreationTimestamp.Time,
							GUID:             space11Anchor.Name,
							OrganizationGUID: org1Anchor.Name,
						},
					))
				})
			})

			When("filtering by org guids", func() {
				It("only retruns the spaces belonging to the specified org guids", func() {
					spaces, err := orgRepo.ListSpaces(ctx, info, []string{string(org1Anchor.Name), string(org3Anchor.Name), "does-not-exist"}, []string{})
					Expect(err).NotTo(HaveOccurred())
					Expect(spaces).To(ConsistOf(
						MatchFields(IgnoreExtras, Fields{
							"Name":             Equal("space1"),
							"OrganizationGUID": Equal(string(org1Anchor.Name)),
						}),
						MatchFields(IgnoreExtras, Fields{
							"Name":             Equal("space1"),
							"OrganizationGUID": Equal(string(org3Anchor.Name)),
						}),
						MatchFields(IgnoreExtras, Fields{"Name": Equal("space2")}),
						MatchFields(IgnoreExtras, Fields{"Name": Equal("space4")}),
					))
				})
			})

			When("filtering by space names", func() {
				It("only retruns the spaces matching the specified names", func() {
					spaces, err := orgRepo.ListSpaces(ctx, info, []string{}, []string{"space1", "space3", "does-not-exist"})
					Expect(err).NotTo(HaveOccurred())
					Expect(spaces).To(ConsistOf(
						MatchFields(IgnoreExtras, Fields{
							"Name":             Equal("space1"),
							"OrganizationGUID": Equal(string(org1Anchor.Name)),
						}),
						MatchFields(IgnoreExtras, Fields{
							"Name":             Equal("space1"),
							"OrganizationGUID": Equal(string(org2Anchor.Name)),
						}),
						MatchFields(IgnoreExtras, Fields{
							"Name":             Equal("space1"),
							"OrganizationGUID": Equal(string(org3Anchor.Name)),
						}),
						MatchFields(IgnoreExtras, Fields{"Name": Equal("space3")}),
					))
				})
			})

			When("filtering by org guids and space names", func() {
				It("only retruns the spaces matching the specified names", func() {
					spaces, err := orgRepo.ListSpaces(ctx, info, []string{string(org1Anchor.Name), string(org2Anchor.Name)}, []string{"space1", "space2", "space4"})
					Expect(err).NotTo(HaveOccurred())
					Expect(spaces).To(ConsistOf(
						MatchFields(IgnoreExtras, Fields{
							"Name":             Equal("space1"),
							"OrganizationGUID": Equal(string(org1Anchor.Name)),
						}),
						MatchFields(IgnoreExtras, Fields{
							"Name":             Equal("space1"),
							"OrganizationGUID": Equal(string(org2Anchor.Name)),
						}),
						MatchFields(IgnoreExtras, Fields{"Name": Equal("space2")}),
					))
				})
			})

			When("filtering by space names that don't exist", func() {
				It("only retruns the spaces matching the specified names", func() {
					spaces, err := orgRepo.ListSpaces(ctx, info, []string{}, []string{"does-not-exist", "still-does-not-exist"})
					Expect(err).NotTo(HaveOccurred())
					Expect(spaces).To(BeEmpty())
				})
			})

			When("filtering by org uids that don't exist", func() {
				It("only retruns the spaces matching the specified names", func() {
					spaces, err := orgRepo.ListSpaces(ctx, info, []string{"does-not-exist", "still-does-not-exist"}, []string{})
					Expect(err).NotTo(HaveOccurred())
					Expect(spaces).To(BeEmpty())
				})
			})
		})
	})
})
