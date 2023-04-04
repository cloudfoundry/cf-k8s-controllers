package handlers_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"code.cloudfoundry.org/korifi/api/authorization"
	"code.cloudfoundry.org/korifi/api/routing"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

const (
	defaultServerURL = "https://api.example.org"
)

var (
	rr            *httptest.ResponseRecorder
	routerBuilder *routing.RouterBuilder
	serverURL     *url.URL
	ctx           context.Context
	authInfo      authorization.Info
)

func TestHandlers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Handlers Suite")
}

var _ = BeforeSuite(func() {
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter)))
})

var _ = BeforeEach(func() {
	authInfo = authorization.Info{Token: "a-token"}
	ctx = authorization.NewContext(context.Background(), &authInfo)
	rr = httptest.NewRecorder()
	routerBuilder = routing.NewRouterBuilder()

	var err error
	serverURL, err = url.Parse(defaultServerURL)
	Expect(err).NotTo(HaveOccurred())
})

func expectErrorResponse(status int, title, detail string, code int) {
	ExpectWithOffset(2, rr).To(HaveHTTPStatus(status))
	ExpectWithOffset(2, rr).To(HaveHTTPHeaderWithValue("Content-Type", "application/json"))
	ExpectWithOffset(2, rr).To(HaveHTTPBody(MatchJSON(fmt.Sprintf(`{
		"errors": [{
			"title": %q,
			"detail": %q,
			"code": %d
		}]
	}`, title, detail, code))))
}

func expectUnknownError() {
	expectErrorResponse(http.StatusInternalServerError, "UnknownError", "An unknown error occurred.", 10001)
}

func expectNotAuthorizedError() {
	expectErrorResponse(http.StatusForbidden, "CF-NotAuthorized", "You are not authorized to perform the requested action", 10003)
}

func expectNotFoundError(resourceType string) {
	expectErrorResponse(http.StatusNotFound, "CF-ResourceNotFound", resourceType+" not found. Ensure it exists and you have access to it.", 10010)
}

func expectUnprocessableEntityError(detail string) {
	expectErrorResponse(http.StatusUnprocessableEntity, "CF-UnprocessableEntity", detail, 10008)
}

func expectBadRequestError() {
	expectErrorResponse(http.StatusBadRequest, "CF-MessageParseError", "Request invalid due to parse error: invalid request body", 1001)
}

func expectBlobstoreUnavailableError() {
	expectErrorResponse(http.StatusBadGateway, "CF-BlobstoreUnavailable", "Error uploading source package to the container registry", 150006)
}

func expectUnknownKeyError(detail string) {
	expectErrorResponse(http.StatusBadRequest, "CF-BadQueryParameter", detail, 10005)
}

func generateGUID(prefix string) string {
	guid := uuid.NewString()

	return fmt.Sprintf("%s-%s", prefix, guid[:13])
}
