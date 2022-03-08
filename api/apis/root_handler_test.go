package apis_test

import (
	"encoding/json"
	"net/http"

	"code.cloudfoundry.org/cf-k8s-controllers/api/apis"
	"code.cloudfoundry.org/cf-k8s-controllers/api/presenter"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gstruct"
)

var _ = Describe("RootHandler", func() {
	var req *http.Request

	BeforeEach(func() {
		apiHandler := apis.NewRootHandler(
			defaultServerURL,
		)
		apiHandler.RegisterRoutes(router)
	})

	JustBeforeEach(func() {
		router.ServeHTTP(rr, req)
	})

	Describe("GET / endpoint", func() {
		BeforeEach(func() {
			var err error
			req, err = http.NewRequest("GET", "/", nil)
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns status 200 OK", func() {
			Expect(rr.Code).To(Equal(http.StatusOK), "Matching HTTP response code:")
		})

		It("returns Content-Type as JSON in header", func() {
			contentTypeHeader := rr.Header().Get("Content-Type")
			Expect(contentTypeHeader).To(Equal(jsonHeader), "Matching Content-Type header:")
		})

		It("has a non-empty body", func() {
			Expect(rr.Body.Bytes()).NotTo(BeEmpty())
		})

		It("matches the expected response body format", func() {
			var resp presenter.RootResponse
			Expect(json.Unmarshal(rr.Body.Bytes(), &resp)).To(Succeed())

			Expect(resp).To(gstruct.MatchAllFields(gstruct.Fields{
				"Links": Equal(map[string]*presenter.APILink{
					"self": {
						Link: presenter.Link{HREF: defaultServerURL},
					},
					"bits_service":        nil,
					"cloud_controller_v2": nil,
					"cloud_controller_v3": {
						Link: presenter.Link{HREF: defaultServerURL + "/v3"},
						Meta: presenter.APILinkMeta{Version: "3.111.0+cf-k8s"},
					},
					"network_policy_v0": nil,
					"network_policy_v1": nil,
					"login":             nil,
					"uaa":               nil,
					"credhub":           nil,
					"routing":           nil,
					"logging":           nil,
					"log_cache": {
						Link: presenter.Link{HREF: defaultServerURL},
					},
					"log_stream": nil,
					"app_ssh":    nil,
				}),
				"CFOnK8s": Equal(true),
			}))
		})
	})
})
