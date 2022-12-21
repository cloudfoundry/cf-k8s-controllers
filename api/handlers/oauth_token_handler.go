package handlers

import (
	"context"
	"net/http"
	"net/url"
	"time"

	"code.cloudfoundry.org/korifi/api/authorization"

	"github.com/go-chi/chi"
	"github.com/go-logr/logr"
	"github.com/golang-jwt/jwt"

	ctrl "sigs.k8s.io/controller-runtime"
)

const (
	OAuthTokenPath = "/oauth/token"
)

type OAuthTokenHandler struct {
	handlerWrapper *AuthAwareHandlerFuncWrapper
	apiBaseURL     url.URL
}

func NewOAuthToken(apiBaseURL url.URL) *OAuthTokenHandler {
	return &OAuthTokenHandler{
		handlerWrapper: NewUnauthenticatedHandlerFuncWrapper(ctrl.Log.WithName("OAuthTokenHandler")),
		apiBaseURL:     apiBaseURL,
	}
}

func (h *OAuthTokenHandler) oauthTokenHandler(ctx context.Context, logger logr.Logger, authInfo authorization.Info, r *http.Request) (*HandlerResponse, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte("not-a-real-secret"))
	if err != nil {
		// we don't expect to hit this ever, given that the string above is hard coded.
		panic(err.Error())
	}

	return NewHandlerResponse(http.StatusOK).WithBody(map[string]string{
		"token_type":   "bearer",
		"access_token": tokenString,
	}), nil
}

func (h *OAuthTokenHandler) RegisterRoutes(router *chi.Mux) {
	router.Post(OAuthTokenPath, h.handlerWrapper.Wrap(h.oauthTokenHandler))
}
