// nolint:dupl
package handlers

import (
	"context"
	"net/http"
	"net/url"

	"code.cloudfoundry.org/korifi/api/authorization"
	apierrors "code.cloudfoundry.org/korifi/api/errors"
	"code.cloudfoundry.org/korifi/api/payloads"
	"code.cloudfoundry.org/korifi/api/presenter"
	"code.cloudfoundry.org/korifi/api/repositories"
	"code.cloudfoundry.org/korifi/api/routing"
	"github.com/go-logr/logr"
)

const (
	ServicePlansPath          = "/v3/service_plans"
	ServicePlanVisivilityPath = "/v3/service_plans/{guid}/visibility"
)

//counterfeiter:generate -o fake -fake-name CFServicePlanRepository . CFServicePlanRepository
type CFServicePlanRepository interface {
	ListPlans(context.Context, authorization.Info, repositories.ListServicePlanMessage) ([]repositories.ServicePlanRecord, error)
	GetPlanVisibility(context.Context, authorization.Info, string) (repositories.ServicePlanVisibilityRecord, error)
}

type ServicePlan struct {
	serverURL        url.URL
	requestValidator RequestValidator
	servicePlanRepo  CFServicePlanRepository
}

func NewServicePlan(
	serverURL url.URL,
	requestValidator RequestValidator,
	servicePlanRepo CFServicePlanRepository,
) *ServicePlan {
	return &ServicePlan{
		serverURL:        serverURL,
		requestValidator: requestValidator,
		servicePlanRepo:  servicePlanRepo,
	}
}

func (h *ServicePlan) list(r *http.Request) (*routing.Response, error) {
	authInfo, _ := authorization.InfoFromContext(r.Context())
	logger := logr.FromContextOrDiscard(r.Context()).WithName("handlers.service-plan.list")

	var payload payloads.ServicePlanList
	if err := h.requestValidator.DecodeAndValidateURLValues(r, &payload); err != nil {
		return nil, apierrors.LogAndReturn(logger, err, "failed to decode json payload")
	}

	servicePlanList, err := h.servicePlanRepo.ListPlans(r.Context(), authInfo, payload.ToMessage())
	if err != nil {
		return nil, apierrors.LogAndReturn(logger, err, "failed to list service plans")
	}

	return routing.NewResponse(http.StatusOK).WithBody(presenter.ForList(presenter.ForServicePlan, servicePlanList, h.serverURL, *r.URL)), nil
}

func (h *ServicePlan) getPlanVisibility(r *http.Request) (*routing.Response, error) {
	authInfo, _ := authorization.InfoFromContext(r.Context())
	logger := logr.FromContextOrDiscard(r.Context()).WithName("handlers.service-plan.get-visibility")

	planGUID := routing.URLParam(r, "guid")
	logger = logger.WithValues("guid", planGUID)

	visibility, err := h.servicePlanRepo.GetPlanVisibility(r.Context(), authInfo, planGUID)
	if err != nil {
		return nil, apierrors.LogAndReturn(logger, err, "failed to get plan visibility")
	}

	return routing.NewResponse(http.StatusOK).WithBody(presenter.ForServicePlanVisibility(visibility, h.serverURL)), nil
}

func (h *ServicePlan) UnauthenticatedRoutes() []routing.Route {
	return nil
}

func (h *ServicePlan) AuthenticatedRoutes() []routing.Route {
	return []routing.Route{
		{Method: "GET", Pattern: ServicePlansPath, Handler: h.list},
		{Method: "GET", Pattern: ServicePlanVisivilityPath, Handler: h.getPlanVisibility},
	}
}
