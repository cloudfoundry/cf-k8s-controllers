package handlers_test

import (
	"encoding/json"
	"net/http"
	"time"

	apis "code.cloudfoundry.org/korifi/api/handlers"

	"github.com/SermoDigital/jose/jws"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("OAuthToken", func() {
	var OAuthTokenHandler *apis.OAuthTokenHandler

	BeforeEach(func() {
		OAuthTokenHandler = apis.NewOAuthToken(*serverURL)
		OAuthTokenHandler.RegisterRoutes(router)
	})

	JustBeforeEach(func() {
		req, err := http.NewRequest(http.MethodPost, "/oauth/token", nil)
		Expect(err).NotTo(HaveOccurred())

		router.ServeHTTP(rr, req)
	})

	Describe("POST OAuthToken", func() {
		It("returns 201 with appropriate success JSON", func() {
			Expect(rr).To(HaveHTTPStatus(http.StatusOK))
			Expect(rr).To(HaveHTTPHeaderWithValue("Content-Type", "application/json"))
			jsonBody := map[string]string{}
			Expect(json.NewDecoder(rr.Body).Decode(&jsonBody)).To(Succeed())
			Expect(jsonBody).To(HaveKeyWithValue("token_type", "bearer"))
			Expect(jsonBody).To(HaveKeyWithValue("access_token", Not(BeEmpty())))

			tokenout, err := jws.ParseJWT([]byte(jsonBody["access_token"]))
			Expect(err).NotTo(HaveOccurred())
			expiration, ok := tokenout.Claims().Expiration()
			Expect(ok).To(BeTrue())
			Expect(expiration.Unix()).To(BeNumerically(">", time.Now().Add(time.Minute*59).Unix()))
		})
	})
})
