package apis_test

import (
	"net/http"

	"code.cloudfoundry.org/cf-k8s-controllers/api/apis"
	"github.com/go-http-utils/headers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("RootV3Handler", func() {
	var req *http.Request

	BeforeEach(func() {
		handler := apis.NewRootV3Handler(defaultServerURL)
		handler.RegisterRoutes(router)
	})

	JustBeforeEach(func() {
		router.ServeHTTP(rr, req)
	})

	Describe("the GET /v3 endpoint", func() {
		BeforeEach(func() {
			var err error
			req, err = http.NewRequest("GET", "/v3", nil)
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns status 200 OK", func() {
			Expect(rr.Code).To(Equal(http.StatusOK))
		})

		It("returns Content-Type as JSON in header", func() {
			Expect(rr).To(HaveHTTPHeaderWithValue(headers.ContentType, jsonHeader))
		})

		It("matches the expected response body format", func() {
			expectedBody := `{"links":{"self":{"href":"` + defaultServerURL + `/v3"}}}`
			Expect(rr.Body).To(MatchJSON(expectedBody))
		})
	})
})
