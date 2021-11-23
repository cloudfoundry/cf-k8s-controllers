package apis

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"code.cloudfoundry.org/cf-k8s-controllers/api/payloads"
	"code.cloudfoundry.org/cf-k8s-controllers/api/presenter"
	"code.cloudfoundry.org/cf-k8s-controllers/api/repositories"
	"code.cloudfoundry.org/cf-k8s-controllers/api/repositories/authorization"
	"code.cloudfoundry.org/cf-k8s-controllers/controllers/webhooks/workloads"
	"github.com/go-logr/logr"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	controllerruntime "sigs.k8s.io/controller-runtime"
)

const (
	SpacesEndpoint = "/v3/spaces"
)

//counterfeiter:generate -o fake -fake-name SpaceRepositoryProvider . SpaceRepositoryProvider

type SpaceRepositoryProvider interface {
	SpaceRepoForRequest(request *http.Request) (repositories.CFSpaceRepository, error)
}

type SpaceHandler struct {
	spaceRepoProvider SpaceRepositoryProvider
	logger            logr.Logger
	apiBaseURL        url.URL
}

func NewSpaceHandler(apiBaseURL url.URL, spaceRepoProvider SpaceRepositoryProvider) *SpaceHandler {
	return &SpaceHandler{
		apiBaseURL:        apiBaseURL,
		spaceRepoProvider: spaceRepoProvider,
		logger:            controllerruntime.Log.WithName("Org Handler"),
	}
}

func (h *SpaceHandler) SpaceCreateHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	var payload payloads.SpaceCreate
	rme := decodeAndValidateJSONPayload(r, &payload)
	if rme != nil {
		h.logger.Error(rme, "Failed to decode and validate payload")
		writeErrorResponse(w, rme)
		return
	}

	space := payload.ToRecord()
	space.GUID = uuid.NewString()

	spaceRepo, err := h.spaceRepoProvider.SpaceRepoForRequest(r)
	if err != nil {
		if authorization.IsUnauthorized(err) {
			h.logger.Error(err, "unauthorized to create spaces")
			writeUnauthorizedErrorResponse(w)

			return
		}

		h.logger.Error(err, "Failed to provide space repo")
		writeUnknownErrorResponse(w)

		return
	}

	record, err := spaceRepo.CreateSpace(ctx, space)
	if err != nil {
		if workloads.HasErrorCode(err, workloads.DuplicateSpaceNameError) {
			errorDetail := fmt.Sprintf("Space '%s' already exists.", space.Name)
			h.logger.Info(errorDetail)
			writeUnprocessableEntityError(w, errorDetail)
			return
		}

		h.logger.Error(err, "Failed to create space", "Space Name", space.Name)
		writeUnknownErrorResponse(w)
		return
	}

	w.WriteHeader(http.StatusCreated)
	spaceResponse := presenter.ForCreateSpace(record, h.apiBaseURL)

	err = json.NewEncoder(w).Encode(spaceResponse)
	if err != nil {
		h.logger.Error(err, "Failed to write response")
	}
}

func (h *SpaceHandler) SpaceListHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")

	orgUIDs := parseCommaSeparatedList(r.URL.Query().Get("organization_guids"))
	names := parseCommaSeparatedList(r.URL.Query().Get("names"))

	spaceRepo, err := h.spaceRepoProvider.SpaceRepoForRequest(r)
	if err != nil {
		if authorization.IsUnauthorized(err) {
			h.logger.Error(err, "unauthorized to list spaces")
			writeUnauthorizedErrorResponse(w)

			return
		}

		h.logger.Error(err, "Failed to provide space repo")
		writeUnknownErrorResponse(w)

		return
	}

	spaces, err := spaceRepo.FetchSpaces(ctx, orgUIDs, names)
	if err != nil {
		h.logger.Error(err, "Failed to fetch spaces")
		writeUnknownErrorResponse(w)

		return
	}

	spaceList := presenter.ForSpaceList(spaces, h.apiBaseURL)
	json.NewEncoder(w).Encode(spaceList)
}

func (h *SpaceHandler) RegisterRoutes(router *mux.Router) {
	router.Path(SpacesEndpoint).Methods("GET").HandlerFunc(h.SpaceListHandler)
	router.Path(SpacesEndpoint).Methods("POST").HandlerFunc(h.SpaceCreateHandler)
}

func parseCommaSeparatedList(list string) []string {
	var elements []string
	for _, element := range strings.Split(list, ",") {
		if element != "" {
			elements = append(elements, element)
		}
	}

	return elements
}
