package handlers

import (
	"context"
	"net/http"
	"net/url"
	"slices"

	"code.cloudfoundry.org/korifi/api/authorization"
	apierrors "code.cloudfoundry.org/korifi/api/errors"
	"code.cloudfoundry.org/korifi/api/payloads"
	"code.cloudfoundry.org/korifi/api/presenter"
	"code.cloudfoundry.org/korifi/api/repositories"
	"code.cloudfoundry.org/korifi/api/routing"
	korifiv1alpha1 "code.cloudfoundry.org/korifi/controllers/api/v1alpha1"
	"github.com/go-logr/logr"
)

const (
	ServiceOfferingsPath      = "/v3/service_offerings"
	ServicePlansPath          = "/v3/service_plans"
	ServicePlanPath           = "/v3/service_plans/{guid}"
	ServiceOfferingPath       = "/v3/service_offerings/{guid}"
	ServicePlanVisivilityPath = "/v3/service_plans/{guid}/visibility"
)

type ServiceCatalogRepo interface {
	ListServiceOfferings(ctx context.Context, authInfo authorization.Info, message repositories.ListServiceOfferingMessage) ([]korifiv1alpha1.ServiceOfferingResource, error)
	ListServicePlans(ctx context.Context, authInfo authorization.Info, message repositories.ListServicePlanMessage) ([]korifiv1alpha1.ServicePlanResource, error)
	GetServicePlan(ctx context.Context, authInfo authorization.Info, guid string) (korifiv1alpha1.ServicePlanResource, error)
	ApplyPlanVisibility(context.Context, authorization.Info, repositories.PlanVisibilityApplyMessage) (korifiv1alpha1.ServicePlanVisibilityResource, error)
	GetServiceOffering(ctx context.Context, authInfo authorization.Info, guid string) (korifiv1alpha1.ServiceOfferingResource, error)
}

type ServiceCatalog struct {
	serverURL          url.URL
	serviceCatalogRepo ServiceCatalogRepo
	requestValidator   RequestValidator
}

func NewServiceCatalog(
	serverURL url.URL,
	serviceCatalogRepo ServiceCatalogRepo,
	requestValidator RequestValidator,
) *ServiceCatalog {
	return &ServiceCatalog{
		serverURL:          serverURL,
		serviceCatalogRepo: serviceCatalogRepo,
		requestValidator:   requestValidator,
	}
}

func (h *ServiceCatalog) listOfferings(r *http.Request) (*routing.Response, error) {
	authInfo, _ := authorization.InfoFromContext(r.Context())
	logger := logr.FromContextOrDiscard(r.Context()).WithName("handlers.service-instance.list")

	if err := r.ParseForm(); err != nil {
		return nil, apierrors.LogAndReturn(logger, err, "Unable to parse request query parameters")
	}

	listFilter := new(payloads.ServiceOfferingList)
	err := h.requestValidator.DecodeAndValidateURLValues(r, listFilter)
	if err != nil {
		return nil, apierrors.LogAndReturn(logger, err, "Unable to decode request query parameters")
	}

	serviceOfferingList, err := h.serviceCatalogRepo.ListServiceOfferings(r.Context(), authInfo, listFilter.ToMessage())
	if err != nil {
		return nil, apierrors.LogAndReturn(logger, err, "Failed to list service instance")
	}

	return routing.NewResponse(http.StatusOK).WithBody(presenter.ForServiceOfferingList(serviceOfferingList, h.serverURL, *r.URL)), nil
}

func (h *ServiceCatalog) listPlans(r *http.Request) (*routing.Response, error) {
	authInfo, _ := authorization.InfoFromContext(r.Context())
	logger := logr.FromContextOrDiscard(r.Context()).WithName("handlers.service-instance.list")

	if err := r.ParseForm(); err != nil {
		return nil, apierrors.LogAndReturn(logger, err, "Unable to parse request query parameters")
	}

	listFilter := payloads.ServicePlanList{}
	err := h.requestValidator.DecodeAndValidateURLValues(r, &listFilter)
	if err != nil {
		return nil, apierrors.LogAndReturn(logger, err, "Unable to decode request query parameters")
	}

	servicePlanList, err := h.serviceCatalogRepo.ListServicePlans(r.Context(), authInfo, listFilter.ToMessage())
	if err != nil {
		return nil, apierrors.LogAndReturn(logger, err, "Failed to list service plans")
	}

	includedResources := []presenter.IncludedResources{}
	if slices.Contains(listFilter.Include, "service_offering") {
		var offerings []any
		for _, plan := range servicePlanList {
			offeringGUID := plan.Relationships.Service_offering.Data.GUID
			offering, err := h.serviceCatalogRepo.GetServiceOffering(r.Context(), authInfo, offeringGUID)
			if err != nil {
				return nil, apierrors.LogAndReturn(logger, err, "Failed to get service offering", "guid", offeringGUID)
			}
			offerings = append(offerings, presenter.ForServiceOffering(offering, h.serverURL))
		}
		includedResources = append(includedResources, presenter.IncludedResources{
			Type:      "service_offerings",
			Resources: offerings,
		})
	}

	return routing.NewResponse(http.StatusOK).WithBody(presenter.ForServicePlanList(servicePlanList, h.serverURL, *r.URL, includedResources...)), nil
}

func (h *ServiceCatalog) getPlan(r *http.Request) (*routing.Response, error) {
	authInfo, _ := authorization.InfoFromContext(r.Context())
	logger := logr.FromContextOrDiscard(r.Context()).WithName("handlers.service-catalog.get-plan")

	if err := r.ParseForm(); err != nil {
		return nil, apierrors.LogAndReturn(logger, err, "Unable to parse request query parameters")
	}

	planGUID := routing.URLParam(r, "guid")

	servicePlan, err := h.serviceCatalogRepo.GetServicePlan(r.Context(), authInfo, planGUID)
	if err != nil {
		return nil, apierrors.LogAndReturn(logger, err, "Failed to get service plan", "planGUID", planGUID)
	}

	return routing.NewResponse(http.StatusOK).WithBody(presenter.ForServicePlan(servicePlan, h.serverURL)), nil
}

func (h *ServiceCatalog) applyPlanVisibility(r *http.Request) (*routing.Response, error) {
	authInfo, _ := authorization.InfoFromContext(r.Context())
	logger := logr.FromContextOrDiscard(r.Context()).WithName("handlers.service-catalog.get-plan")

	if err := r.ParseForm(); err != nil {
		return nil, apierrors.LogAndReturn(logger, err, "Unable to parse request query parameters")
	}

	planGUID := routing.URLParam(r, "guid")
	payload := payloads.PlanVisiblityApply{}
	if err := h.requestValidator.DecodeAndValidateJSONPayload(r, &payload); err != nil {
		return nil, apierrors.LogAndReturn(logger, err, "failed to decode payload")
	}

	visibilityRes, err := h.serviceCatalogRepo.ApplyPlanVisibility(r.Context(), authInfo, payload.ToMessage(planGUID))
	if err != nil {
		return nil, apierrors.LogAndReturn(logger, err, "failed to create plan visibility resource")
	}
	return routing.NewResponse(http.StatusOK).WithBody(presenter.ForServicePlanVisibility(visibilityRes)), nil
}

func (h *ServiceCatalog) getOffering(r *http.Request) (*routing.Response, error) {
	authInfo, _ := authorization.InfoFromContext(r.Context())
	logger := logr.FromContextOrDiscard(r.Context()).WithName("handlers.service-catalog.get-offering")

	if err := r.ParseForm(); err != nil {
		return nil, apierrors.LogAndReturn(logger, err, "Unable to parse request query parameters")
	}

	offeringGUID := routing.URLParam(r, "guid")

	serviceOffering, err := h.serviceCatalogRepo.GetServiceOffering(r.Context(), authInfo, offeringGUID)
	if err != nil {
		return nil, apierrors.LogAndReturn(logger, err, "Failed to get service offering", "offeringGUID", offeringGUID)
	}

	return routing.NewResponse(http.StatusOK).WithBody(presenter.ForServiceOffering(serviceOffering, h.serverURL)), nil
}

func (h *ServiceCatalog) UnauthenticatedRoutes() []routing.Route {
	return nil
}

func (h *ServiceCatalog) AuthenticatedRoutes() []routing.Route {
	return []routing.Route{
		{Method: "GET", Pattern: ServiceOfferingsPath, Handler: h.listOfferings},
		{Method: "GET", Pattern: ServicePlansPath, Handler: h.listPlans},
		{Method: "GET", Pattern: ServicePlanPath, Handler: h.getPlan},
		{Method: "POST", Pattern: ServicePlanVisivilityPath, Handler: h.applyPlanVisibility},
		{Method: "GET", Pattern: ServiceOfferingPath, Handler: h.getOffering},
	}
}
