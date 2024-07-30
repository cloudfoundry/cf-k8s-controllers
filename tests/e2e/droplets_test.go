package e2e_test

import (
	"net/http"

	"github.com/go-resty/resty/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Droplets", func() {
	var (
		spaceGUID string
		buildGUID string
		result    responseResource
	)

	BeforeEach(func() {
		spaceGUID = createSpace(generateGUID("space1"), commonTestOrgGUID)
		appGUID := createBuildpackApp(spaceGUID, generateGUID("app"))
		pkgGUID := createBitsPackage(appGUID)
		uploadTestApp(pkgGUID, defaultAppBitsFile)
		buildGUID = createBuild(pkgGUID)
	})

	AfterEach(func() {
		deleteSpace(spaceGUID)
	})

	Describe("get", func() {
		JustBeforeEach(func() {
			Eventually(func() (*resty.Response, error) {
				return adminClient.R().
					SetResult(&result).
					Get("/v3/droplets/" + buildGUID)
			}).Should(HaveRestyStatusCode(http.StatusOK))
		})

		It("returns the droplet", func() {
			Expect(result.GUID).To(Equal(buildGUID))
		})

		When("the lifecycle type is docker", func() {
			BeforeEach(func() {
				appGUID := createApp(appResource{
					resource: resource{
						Name:          generateGUID("app"),
						Relationships: relationships{"space": {Data: resource{GUID: spaceGUID}}},
					},
					Lifecycle: &lifecycle{
						Type: "docker",
					},
				})
				pkgGUID := createPackage(packageResource{
					typedResource: typedResource{
						Type: "docker",
						resource: resource{
							Relationships: relationships{
								"app": relationship{Data: resource{GUID: appGUID}},
							},
						},
					},
					Data: &packageData{
						Image: "eirini/dorini",
					},
				})
				buildGUID = createBuild(pkgGUID)
			})

			It("returns the droplet", func() {
				Expect(result.GUID).To(Equal(buildGUID))
			})
		})
	})

	Describe("update", func() {
		JustBeforeEach(func() {
			Eventually(func() (*resty.Response, error) {
				return adminClient.R().
					SetBody(metadataResource{
						Metadata: &metadataPatch{
							Annotations: &map[string]string{"foo": "bar"},
							Labels:      &map[string]string{"baz": "bar"},
						},
					}).
					SetResult(&result).
					Patch("/v3/droplets/" + buildGUID)
			}).Should(HaveRestyStatusCode(http.StatusOK))
		})

		It("updates droplet with labels and annotations", func() {
			Expect(result.GUID).To(Equal(buildGUID))
			Expect(result.Metadata.Annotations).To(HaveKeyWithValue("foo", "bar"))
			Expect(result.Metadata.Labels).To(HaveKeyWithValue("baz", "bar"))
		})
	})
})
