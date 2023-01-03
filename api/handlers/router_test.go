package handlers_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"

	"code.cloudfoundry.org/korifi/api/authorization"
	"code.cloudfoundry.org/korifi/api/handlers"
	"github.com/go-logr/logr"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

type handler struct{}

func handlerFunc(ctx context.Context, logger logr.Logger, authInfo authorization.Info, r *http.Request) (*handlers.HandlerResponse, error) {
	return handlers.NewHandlerResponse(http.StatusTeapot), nil
}

func (h handler) AuthenticatedRoutes() []handlers.Route {
	return []handlers.Route{
		{Method: "GET", Pattern: "/authenticated", HandlerFunc: handlerFunc},
	}
}

func (h handler) UnauthenticatedRoutes() []handlers.Route {
	return []handlers.Route{
		{Method: "GET", Pattern: "/unauthenticated", HandlerFunc: handlerFunc},
	}
}

func middleware(h http.Handler) http.Handler {
	return h
}

var _ = FDescribe("Router", func() {
	var (
		router *handlers.Router
		server *httptest.Server
	)

	BeforeEach(func() {
		router = handlers.NewRouter(logf.Log.WithName("test"))
	})

	JustBeforeEach(func() {
		server = httptest.NewServer(router)
	})

	AfterEach(func() {
		server.Close()
	})

	When("registering a handler", func() {
		BeforeEach(func() {
			router.RegisterHandler("handler", handler{})
		})

		It("wraps the unauthenticated handlers with an unauthenticated wrapper", func() {
			res, err := http.Get(server.URL + "/unauthenticated")
			Expect(err).NotTo(HaveOccurred())
			Expect(res.StatusCode).To(Equal(http.StatusTeapot))
		})

		When("using the authentication middleware", func() {
			It("wraps the unauthenticated handlers with an authenticated wrapper", func() {
			})
		})
	})

	When("trying to use a common handler after registering an authenticated handler", func() {
		BeforeEach(func() {
			// router.UseCommonMiddleware(middleware)
		})
	})
	When("trying to use a common handler after registering an unauthenticated handler", func() {})
	When("trying to use an auth handler after registering an authenticated handler", func() {})
	When("trying to use an auth handler after registering an unauthenticated handler", func() {})
})
