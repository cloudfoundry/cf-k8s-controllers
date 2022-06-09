package actions

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"code.cloudfoundry.org/korifi/api/apierrors"
	"code.cloudfoundry.org/korifi/api/authorization"
	"code.cloudfoundry.org/korifi/api/payloads"
	"code.cloudfoundry.org/korifi/api/repositories"
)

const (
	processTypeWeb = "web"
)

type Manifest struct {
	appRepo           CFAppRepository
	domainRepo        CFDomainRepository
	processRepo       CFProcessRepository
	routeRepo         CFRouteRepository
	defaultDomainName string
}

func NewManifest(appRepo CFAppRepository, domainRepo CFDomainRepository, processRepo CFProcessRepository, routeRepo CFRouteRepository, defaultDomainName string) *Manifest {
	return &Manifest{
		appRepo:           appRepo,
		domainRepo:        domainRepo,
		processRepo:       processRepo,
		routeRepo:         routeRepo,
		defaultDomainName: defaultDomainName,
	}
}

func (a *Manifest) Apply(ctx context.Context, authInfo authorization.Info, spaceGUID string, manifest payloads.Manifest) error {
	appInfo := manifest.Applications[0]
	exists := true
	appRecord, err := a.appRepo.GetAppByNameAndSpace(ctx, authInfo, appInfo.Name, spaceGUID)
	if err != nil {
		if errors.As(err, new(apierrors.NotFoundError)) {
			exists = false
		} else {
			return apierrors.ForbiddenAsNotFound(err)
		}
	}

	if appInfo.Memory != nil {
		found := false
		for _, process := range appInfo.Processes {
			if process.Type == processTypeWeb {
				found = true
			}
		}

		if !found {
			appInfo.Processes = append(appInfo.Processes, payloads.ManifestApplicationProcess{
				Type:   processTypeWeb,
				Memory: appInfo.Memory,
			})
		}
	}

	if exists {
		err = a.updateApp(ctx, authInfo, spaceGUID, appRecord, appInfo)
	} else {
		appRecord, err = a.createApp(ctx, authInfo, spaceGUID, appInfo)
	}

	if err != nil {
		return err
	}

	err = a.checkAndUpdateDefaultRoute(ctx, authInfo, appRecord, a.defaultDomainName, &appInfo)
	if err != nil {
		return err
	}

	return a.createOrUpdateRoutes(ctx, authInfo, appRecord, appInfo.Routes)
}

// checkAndUpdateDefaultRoute may set the default route on the manifest when DefaultRoute is true
func (a *Manifest) checkAndUpdateDefaultRoute(ctx context.Context, authInfo authorization.Info, appRecord repositories.AppRecord, defaultDomainName string, appInfo *payloads.ManifestApplication) error {
	if !appInfo.DefaultRoute || len(appInfo.Routes) > 0 {
		return nil
	}

	existingRoutes, err := a.routeRepo.ListRoutesForApp(ctx, authInfo, appRecord.GUID, appRecord.SpaceGUID)
	if err != nil {
		return err
	}
	if len(existingRoutes) > 0 {
		return nil
	}

	_, err = a.domainRepo.GetDomainByName(ctx, authInfo, defaultDomainName)
	if err != nil {
		return apierrors.AsUnprocessableEntity(
			err,
			fmt.Sprintf("The configured default domain %q was not found", defaultDomainName),
			apierrors.NotFoundError{},
		)
	}
	defaultRouteString := appInfo.Name + "." + defaultDomainName
	defaultRoute := payloads.ManifestRoute{
		Route: &defaultRouteString,
	}
	// set the route field of the manifest with app-name . default domain
	appInfo.Routes = append(appInfo.Routes, defaultRoute)

	return nil
}

func (a *Manifest) updateApp(ctx context.Context, authInfo authorization.Info, spaceGUID string, appRecord repositories.AppRecord, appInfo payloads.ManifestApplication) error {
	_, err := a.appRepo.CreateOrPatchAppEnvVars(ctx, authInfo, repositories.CreateOrPatchAppEnvVarsMessage{
		AppGUID:              appRecord.GUID,
		AppEtcdUID:           appRecord.EtcdUID,
		SpaceGUID:            appRecord.SpaceGUID,
		EnvironmentVariables: appInfo.Env,
	})
	if err != nil {
		return err
	}

	for _, processInfo := range appInfo.Processes {
		exists := true

		var process repositories.ProcessRecord
		process, err = a.processRepo.GetProcessByAppTypeAndSpace(ctx, authInfo, appRecord.GUID, processInfo.Type, spaceGUID)
		if err != nil {
			if errors.As(err, new(apierrors.NotFoundError)) {
				exists = false
			} else {
				return err
			}
		}

		if exists {
			_, err = a.processRepo.PatchProcess(ctx, authInfo, processInfo.ToProcessPatchMessage(process.GUID, spaceGUID))
		} else {
			err = a.processRepo.CreateProcess(ctx, authInfo, processInfo.ToProcessCreateMessage(appRecord.GUID, spaceGUID))
		}
		if err != nil {
			return err
		}
	}

	return err
}

func (a *Manifest) createApp(ctx context.Context, authInfo authorization.Info, spaceGUID string, appInfo payloads.ManifestApplication) (repositories.AppRecord, error) {
	appRecord, err := a.appRepo.CreateApp(ctx, authInfo, appInfo.ToAppCreateMessage(spaceGUID))
	if err != nil {
		return appRecord, err
	}

	for _, processInfo := range appInfo.Processes {
		message := processInfo.ToProcessCreateMessage(appRecord.GUID, spaceGUID)
		err = a.processRepo.CreateProcess(ctx, authInfo, message)
		if err != nil {
			return appRecord, err
		}
	}

	return appRecord, nil
}

func (a *Manifest) createOrUpdateRoutes(ctx context.Context, authInfo authorization.Info, appRecord repositories.AppRecord, routes []payloads.ManifestRoute) error {
	if len(routes) == 0 {
		return nil
	}

	routeString := *routes[0].Route
	hostName, domainName, path := splitRoute(routeString)

	domainRecord, err := a.domainRepo.GetDomainByName(ctx, authInfo, domainName)
	if err != nil {
		return fmt.Errorf("createOrUpdateRoutes: %w", err)
	}

	routeRecord, err := a.routeRepo.GetOrCreateRoute(
		ctx,
		authInfo,
		repositories.CreateRouteMessage{
			Host:            hostName,
			Path:            path,
			SpaceGUID:       appRecord.SpaceGUID,
			DomainGUID:      domainRecord.GUID,
			DomainNamespace: domainRecord.Namespace,
			DomainName:      domainRecord.Name,
		})
	if err != nil {
		return fmt.Errorf("createOrUpdateRoutes: %w", err)
	}

	routeRecord, err = a.routeRepo.AddDestinationsToRoute(ctx, authInfo, repositories.AddDestinationsToRouteMessage{
		RouteGUID:            routeRecord.GUID,
		SpaceGUID:            routeRecord.SpaceGUID,
		ExistingDestinations: routeRecord.Destinations,
		NewDestinations: []repositories.DestinationMessage{
			{
				AppGUID:     appRecord.GUID,
				ProcessType: "web",
				Port:        8080,
				Protocol:    "http1",
			},
		},
	})

	return err
}

func splitRoute(route string) (string, string, string) {
	parts := strings.SplitN(route, ".", 2)
	hostName := parts[0]
	domainAndPath := parts[1]

	parts = strings.SplitN(domainAndPath, "/", 2)
	domain := parts[0]
	var path string
	if len(parts) > 1 {
		path = "/" + parts[1]
	}
	return hostName, domain, path
}
