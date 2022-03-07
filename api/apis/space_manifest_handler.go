package apis

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/go-http-utils/headers"
	"github.com/go-logr/logr"
	"github.com/gorilla/mux"

	"code.cloudfoundry.org/cf-k8s-controllers/api/apierrors"
	"code.cloudfoundry.org/cf-k8s-controllers/api/authorization"
	"code.cloudfoundry.org/cf-k8s-controllers/api/payloads"
	"code.cloudfoundry.org/cf-k8s-controllers/api/repositories"
)

const (
	SpaceManifestApplyEndpoint = "/v3/spaces/{spaceGUID}/actions/apply_manifest"
	SpaceManifestDiffEndpoint  = "/v3/spaces/{spaceGUID}/manifest_diff"
)

type SpaceManifestHandler struct {
	logger              logr.Logger
	serverURL           url.URL
	defaultDomainName   string
	applyManifestAction ApplyManifestAction
	spaceRepo           repositories.CFSpaceRepository
	decoderValidator    *DecoderValidator
}

//counterfeiter:generate -o fake -fake-name ApplyManifestAction . ApplyManifestAction
type ApplyManifestAction func(ctx context.Context, authInfo authorization.Info, spaceGUID string, defaultDomainName string, manifest payloads.Manifest) error

func NewSpaceManifestHandler(
	logger logr.Logger,
	serverURL url.URL,
	defaultDomainName string,
	applyManifestAction ApplyManifestAction,
	spaceRepo repositories.CFSpaceRepository,
	decoderValidator *DecoderValidator,
) *SpaceManifestHandler {
	return &SpaceManifestHandler{
		logger:              logger,
		serverURL:           serverURL,
		defaultDomainName:   defaultDomainName,
		applyManifestAction: applyManifestAction,
		spaceRepo:           spaceRepo,
		decoderValidator:    decoderValidator,
	}
}

func (h *SpaceManifestHandler) RegisterRoutes(router *mux.Router) {
	w := NewAuthAwareHandlerFuncWrapper(h.logger)
	router.Path(SpaceManifestApplyEndpoint).Methods("POST").HandlerFunc(w.Wrap(h.applyManifestHandler))
	router.Path(SpaceManifestDiffEndpoint).Methods("POST").HandlerFunc(w.Wrap(h.diffManifestHandler))
}

func (h *SpaceManifestHandler) applyManifestHandler(authInfo authorization.Info, r *http.Request) (*HandlerResponse, error) {
	vars := mux.Vars(r)
	spaceGUID := vars["spaceGUID"]
	var manifest payloads.Manifest
	if err := h.decoderValidator.DecodeAndValidateYAMLPayload(r, &manifest); err != nil {
		return nil, err
	}

	if err := h.applyManifestAction(r.Context(), authInfo, spaceGUID, h.defaultDomainName, manifest); err != nil {
		if errors.As(err, &repositories.NotFoundError{}) {
			h.logger.Info("Domain not found", "error: ", err.Error())
			return nil, apierrors.NewUnprocessableEntityError(err, "The configured default domain `"+h.defaultDomainName+"` was not found")
		}

		h.logger.Error(err, "Error applying manifest")
		return nil, err
	}

	return NewHandlerResponse(http.StatusAccepted).
		WithHeader(headers.Location, fmt.Sprintf("%s/v3/jobs/space.apply_manifest-%s", h.serverURL.String(), spaceGUID)), nil
}

func (h *SpaceManifestHandler) diffManifestHandler(authInfo authorization.Info, r *http.Request) (*HandlerResponse, error) {
	vars := mux.Vars(r)
	spaceGUID := vars["spaceGUID"]

	if _, err := h.spaceRepo.GetSpace(r.Context(), authInfo, spaceGUID); err != nil {
		var notFoundErr repositories.NotFoundError
		if errors.As(err, &notFoundErr) {
			return nil, apierrors.NewNotFoundError(err, notFoundErr.ResourceType())
		}

		return nil, err
	}

	return NewHandlerResponse(http.StatusAccepted).WithBody(map[string]interface{}{"diff": []string{}}), nil
}
