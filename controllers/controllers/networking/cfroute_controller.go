/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package networking

import (
	"context"
	"fmt"

	"code.cloudfoundry.org/cf-k8s-controllers/controllers/config"

	networkingv1alpha1 "code.cloudfoundry.org/cf-k8s-controllers/controllers/apis/networking/v1alpha1"
	workloadsv1alpha1 "code.cloudfoundry.org/cf-k8s-controllers/controllers/apis/workloads/v1alpha1"

	"github.com/go-logr/logr"
	contourv1 "github.com/projectcontour/contour/apis/projectcontour/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const (
	FinalizerName = "cfRoute.networking.cloudfoundry.org"
)

// CFRouteReconciler reconciles a CFRoute object to create Contour resources
type CFRouteReconciler struct {
	Client           client.Client
	Scheme           *runtime.Scheme
	Log              logr.Logger
	ControllerConfig *config.ControllerConfig
}

//+kubebuilder:rbac:groups=networking.cloudfoundry.org,resources=cfroutes,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=networking.cloudfoundry.org,resources=cfroutes/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=networking.cloudfoundry.org,resources=cfroutes/finalizers,verbs=update

//+kubebuilder:rbac:groups=projectcontour.io,resources=httpproxies,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=projectcontour.io,resources=httpproxies/status,verbs=get
//+kubebuilder:rbac:groups=projectcontour.io,resources=httpproxies/finalizers,verbs=update

//+kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch;delete

func (r *CFRouteReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	cfRoute := new(networkingv1alpha1.CFRoute)
	err := r.Client.Get(ctx, req.NamespacedName, cfRoute)
	if err != nil {
		if !apierrors.IsNotFound(err) {
			r.Log.Error(err, "failed to get CFRoute")
		}
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	var cfDomain networkingv1alpha1.CFDomain
	err = r.Client.Get(ctx, types.NamespacedName{Name: cfRoute.Spec.DomainRef.Name}, &cfDomain)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return r.setRouteErrorStatusAndReturn(ctx, cfRoute, err, "CFDomain not found", "InvalidDomainRef")
		}
		return r.setRouteErrorStatusAndReturn(ctx, cfRoute, err, "Error fetching domain reference", "FetchDomainRef")
	}

	err = r.addFinalizer(ctx, cfRoute)
	if err != nil {
		description := "Error adding finalizer"
		r.Log.Error(err, description)
		errMsg := fmt.Sprintf("%v", err)
		if statusErr := r.setRouteStatus(ctx, cfRoute, networkingv1alpha1.InvalidStatus, description, "AddFinalizer", errMsg); statusErr != nil {
			r.Log.Error(statusErr, "Error when updating CFRoute status")
			return ctrl.Result{}, statusErr
		}
		return ctrl.Result{}, err
	}

	if isFinalizing(cfRoute) {
		return r.finalizeCFRoute(ctx, cfRoute, &cfDomain)
	}

	err = r.createOrPatchServices(ctx, cfRoute)
	if err != nil {
		return r.setRouteErrorStatusAndReturn(ctx, cfRoute, err, "Error creating/patching services", "CreatePatchServices")
	}

	err = r.createOrPatchRouteProxy(ctx, cfRoute)
	if err != nil {
		return r.setRouteErrorStatusAndReturn(ctx, cfRoute, err, "Error creating/patching Route Proxy", "CreatePatchRouteProxy")
	}

	err = r.createOrPatchFQDNProxy(ctx, cfRoute, &cfDomain)
	if err != nil {
		return r.setRouteErrorStatusAndReturn(ctx, cfRoute, err, "Error creating/patching FQDN Proxy", "CreatePatchFQDNProxy")
	}

	err = r.deleteOrphanedServices(ctx, cfRoute)
	if err != nil { // technically, failing to delete the orphaned services does not make the CFRoute invalid so we don't mess with the cfRoute status here
		return ctrl.Result{}, err
	}

	// setCFRouteCFDomainStatusFields
	cfRoute.Status.FQDN = cfRoute.Spec.Host + "." + cfDomain.Spec.Name
	cfRoute.Status.URI = cfRoute.Status.FQDN + cfRoute.Spec.Path
	if err := r.setRouteStatus(ctx, cfRoute, networkingv1alpha1.ValidStatus, "Valid CFRoute", "Valid", "Valid CFRoute"); err != nil {
		r.Log.Error(err, "Error when updating CFRoute status")
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

func (r *CFRouteReconciler) setRouteStatus(ctx context.Context, cfRoute *networkingv1alpha1.CFRoute, statusValue networkingv1alpha1.CurrentStatus, description, reason, message string) error {
	cfRoute.Status.CurrentStatus = statusValue
	cfRoute.Status.Description = description

	statusConditionValue := metav1.ConditionUnknown
	if statusValue == networkingv1alpha1.InvalidStatus {
		statusConditionValue = metav1.ConditionFalse
	} else if statusValue == networkingv1alpha1.ValidStatus {
		statusConditionValue = metav1.ConditionTrue
	}

	setStatusConditionOnLocalCopy(&cfRoute.Status.Conditions, "Valid", statusConditionValue, reason, message)

	return r.Client.Status().Update(ctx, cfRoute)
}

func (r *CFRouteReconciler) setRouteErrorStatusAndReturn(ctx context.Context, cfRoute *networkingv1alpha1.CFRoute, err error, description, reason string) (ctrl.Result, error) {
	r.Log.Error(err, description)
	errMsg := fmt.Sprintf("%v", err)
	if statusErr := r.setRouteStatus(ctx, cfRoute, networkingv1alpha1.InvalidStatus, description, reason, errMsg); statusErr != nil {
		r.Log.Error(statusErr, "Error when updating CFRoute status")
		return ctrl.Result{}, statusErr
	}
	return ctrl.Result{}, err
}

// SetupWithManager sets up the controller with the Manager.
func (r *CFRouteReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&networkingv1alpha1.CFRoute{}).
		Complete(r)
}

func (r *CFRouteReconciler) addFinalizer(ctx context.Context, cfRoute *networkingv1alpha1.CFRoute) error {
	if controllerutil.ContainsFinalizer(cfRoute, FinalizerName) {
		return nil
	}

	originalCFRoute := cfRoute.DeepCopy()
	controllerutil.AddFinalizer(cfRoute, FinalizerName)

	err := r.Client.Patch(ctx, cfRoute, client.MergeFrom(originalCFRoute))
	if err != nil {
		r.Log.Error(err, fmt.Sprintf("Error adding finalizer to CFRoute/%s", cfRoute.Name))
		return err
	}

	r.Log.Info(fmt.Sprintf("Finalizer added to CFRoute/%s", cfRoute.Name))
	return nil
}

func (r *CFRouteReconciler) finalizeCFRoute(ctx context.Context, cfRoute *networkingv1alpha1.CFRoute, cfDomain *networkingv1alpha1.CFDomain) (ctrl.Result, error) {
	r.Log.Info(fmt.Sprintf("Reconciling deletion of CFRoute/%s", cfRoute.Name))

	if !controllerutil.ContainsFinalizer(cfRoute, FinalizerName) {
		return ctrl.Result{}, nil
	}

	fqdnHTTPProxy, foundFQDNProxy, err := r.getFQDNProxy(ctx, cfRoute.Spec.Host, cfDomain.Spec.Name, cfRoute.Namespace)
	if err != nil {
		return ctrl.Result{}, err
	}

	// Cleanup the FQDN HTTPProxy on delete
	if foundFQDNProxy {
		err := r.finalizeFQDNProxy(ctx, cfRoute.Name, fqdnHTTPProxy)
		if err != nil {
			return ctrl.Result{}, err
		}
	}

	controllerutil.RemoveFinalizer(cfRoute, FinalizerName)
	if err := r.Client.Update(ctx, cfRoute); err != nil {
		r.Log.Error(err, "Failed to remove finalizer")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *CFRouteReconciler) finalizeFQDNProxy(ctx context.Context, cfRouteName string, fqdnProxy *contourv1.HTTPProxy) error {
	originalFQDNProxy := fqdnProxy.DeepCopy()
	var retainedIncludes []contourv1.Include
	for _, include := range fqdnProxy.Spec.Includes {
		if include.Name != cfRouteName {
			retainedIncludes = append(retainedIncludes, include)
		} else {
			r.Log.Info(fmt.Sprintf("Removing sub-HTTPProxy for route %s from FQDN HTTPProxy", cfRouteName))
		}
	}
	fqdnProxy.Spec.Includes = retainedIncludes
	err := r.Client.Patch(ctx, fqdnProxy, client.MergeFrom(originalFQDNProxy))
	if err != nil {
		r.Log.Error(err, "failed to patch FQDN HTTPProxy to remove sub HTTPProxy")
		return err
	}

	return nil
}

func (r *CFRouteReconciler) createOrPatchServices(ctx context.Context, cfRoute *networkingv1alpha1.CFRoute) error {
	for _, destination := range cfRoute.Spec.Destinations {
		service := &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      generateServiceName(&destination),
				Namespace: cfRoute.Namespace,
			},
		}

		result, err := controllerutil.CreateOrPatch(ctx, r.Client, service, func() error {
			service.ObjectMeta.Labels = map[string]string{
				workloadsv1alpha1.CFAppGUIDLabelKey:    destination.AppRef.Name,
				networkingv1alpha1.CFRouteGUIDLabelKey: cfRoute.Name,
			}

			err := controllerutil.SetOwnerReference(cfRoute, service, r.Scheme)
			if err != nil {
				r.Log.Error(err, "failed to set OwnerRef on Service")
				return err
			}

			service.Spec.Ports = []corev1.ServicePort{{
				Port: int32(destination.Port),
			}}
			service.Spec.Selector = map[string]string{
				workloadsv1alpha1.CFAppGUIDLabelKey:     destination.AppRef.Name,
				workloadsv1alpha1.CFProcessTypeLabelKey: destination.ProcessType,
			}

			return nil
		})
		if err != nil {
			r.Log.Error(err, fmt.Sprintf("failed to patch Service/%s", service.Name))
			return fmt.Errorf("service reconciliation failed for CFRoute/%s destinations", cfRoute.Name)
		}

		r.Log.Info(fmt.Sprintf("Service/%s %s", service.Name, result))
	}

	return nil
}

func (r *CFRouteReconciler) createOrPatchRouteProxy(ctx context.Context, cfRoute *networkingv1alpha1.CFRoute) error {
	services := make([]contourv1.Service, 0, len(cfRoute.Spec.Destinations))

	for _, destination := range cfRoute.Spec.Destinations {
		services = append(services, contourv1.Service{
			Name: generateServiceName(&destination),
			Port: destination.Port,
		})
	}

	routeHTTPProxy := &contourv1.HTTPProxy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cfRoute.Name,
			Namespace: cfRoute.Namespace,
		},
	}

	result, err := controllerutil.CreateOrPatch(ctx, r.Client, routeHTTPProxy, func() error {
		if len(services) == 0 {
			routeHTTPProxy.Spec.Routes = []contourv1.Route{}
		} else {
			routeHTTPProxy.Spec.Routes = []contourv1.Route{
				{
					Conditions: []contourv1.MatchCondition{
						{Prefix: cfRoute.Spec.Path},
					},
					Services: services,
				},
			}
		}

		err := controllerutil.SetOwnerReference(cfRoute, routeHTTPProxy, r.Scheme)
		if err != nil {
			r.Log.Error(err, "failed to set OwnerRef on route HTTPProxy")
			return err
		}

		return nil
	})
	if err != nil {
		r.Log.Error(err, "failed to patch route HTTPProxy")
		return err
	}

	r.Log.Info(fmt.Sprintf("Route HTTPProxy/%s %s", routeHTTPProxy.Name, result))
	return nil
}

func (r *CFRouteReconciler) createOrPatchFQDNProxy(ctx context.Context, cfRoute *networkingv1alpha1.CFRoute, cfDomain *networkingv1alpha1.CFDomain) error {
	fqdnHTTPProxy, foundFQDNPRoxy, err := r.getFQDNProxy(ctx, cfRoute.Spec.Host, cfDomain.Spec.Name, cfRoute.Namespace)
	if err != nil {
		return err
	}

	fqdn := fmt.Sprintf("%s.%s", cfRoute.Spec.Host, cfDomain.Spec.Name)

	if !foundFQDNPRoxy {
		fqdnHTTPProxy = &contourv1.HTTPProxy{
			ObjectMeta: metav1.ObjectMeta{
				Name:      fqdn,
				Namespace: cfRoute.Namespace,
			},
		}
	}

	result, err := controllerutil.CreateOrPatch(ctx, r.Client, fqdnHTTPProxy, func() error {
		fqdnHTTPProxy.Spec.VirtualHost = &contourv1.VirtualHost{
			Fqdn: fqdn,
		}

		if tlsSecret := r.ControllerConfig.WorkloadsTLSSecretNameWithNamespace(); tlsSecret != "" {
			fqdnHTTPProxy.Spec.VirtualHost.TLS = &contourv1.TLS{SecretName: tlsSecret}
		}

		routeAlreadyIncluded := false
		for _, include := range fqdnHTTPProxy.Spec.Includes {
			if include.Name == cfRoute.Name && include.Namespace == cfRoute.Namespace {
				routeAlreadyIncluded = true
			}
		}

		if !routeAlreadyIncluded {
			fqdnHTTPProxy.Spec.Includes = append(fqdnHTTPProxy.Spec.Includes, contourv1.Include{
				Name:      cfRoute.Name,
				Namespace: cfRoute.Namespace,
			})
		}

		err = controllerutil.SetOwnerReference(cfRoute, fqdnHTTPProxy, r.Scheme)
		if err != nil {
			r.Log.Error(err, "failed to set OwnerRef on FQDN HTTPProxy")
			return err
		}

		return nil
	})
	if err != nil {
		r.Log.Error(err, "failed to patch FQDN HTTPProxy")
		return err
	}

	r.Log.Info(fmt.Sprintf("FQDN HTTPProxy/%s %s", fqdnHTTPProxy.Name, result))
	return nil
}

func (r *CFRouteReconciler) getFQDNProxy(ctx context.Context, routeHostname, domainName, namespace string) (*contourv1.HTTPProxy, bool, error) {
	var fqdnHTTPProxy contourv1.HTTPProxy
	fqdn := fmt.Sprintf("%s.%s", routeHostname, domainName)

	// Check all namespaces for FQDN proxy with the matching FQDN label
	var proxies contourv1.HTTPProxyList
	err := r.Client.List(ctx, &proxies)
	if err != nil {
		r.Log.Error(err, "failed to list HTTPProxies")
		return nil, false, err
	}

	var found bool
	for _, proxy := range proxies.Items {
		if proxy.Spec.VirtualHost != nil && proxy.Spec.VirtualHost.Fqdn == fqdn {
			if found {
				err = fmt.Errorf("found multiple HTTPProxy with FQDN %s", fqdn)
				r.Log.Error(err, "")
				return nil, false, err
			} else if proxy.Namespace != namespace {
				err = fmt.Errorf("found existing HTTPProxy with FQDN %s in another space", fqdn)
				r.Log.Error(err, fmt.Sprintf("existing proxy found in namespace %s", proxy.Namespace))
				return nil, false, err
			}

			fqdnHTTPProxy = proxy
			found = true
		}
	}

	return &fqdnHTTPProxy, found, nil
}

func (r *CFRouteReconciler) deleteOrphanedServices(ctx context.Context, cfRoute *networkingv1alpha1.CFRoute) error {
	matchingLabelSet := map[string]string{
		networkingv1alpha1.CFRouteGUIDLabelKey: cfRoute.Name,
	}

	serviceList, err := r.fetchServicesByMatchingLabels(ctx, matchingLabelSet, cfRoute.Namespace)
	if err != nil {
		r.Log.Error(err, fmt.Sprintf("Failed to fetch services using label - %s : %s", networkingv1alpha1.CFRouteGUIDLabelKey, cfRoute.Name))
		return err
	}

	for _, service := range serviceList.Items {
		isOrphan := true
		for _, destination := range cfRoute.Spec.Destinations {
			if service.Name == generateServiceName(&destination) {
				isOrphan = false
				break
			}
		}
		if isOrphan {
			err = r.Client.Delete(ctx, &service)
			if err != nil {
				r.Log.Error(err, fmt.Sprintf("failed to delete service %s", service.Name))
				return err
			}
		}
	}

	return nil
}

func (r *CFRouteReconciler) fetchServicesByMatchingLabels(ctx context.Context, labelSet map[string]string, namespace string) (*corev1.ServiceList, error) {
	selector, err := labels.ValidatedSelectorFromSet(labelSet)
	if err != nil {
		r.Log.Error(err, "Error initializing label selector")
		return nil, err
	}

	serviceList := corev1.ServiceList{}
	err = r.Client.List(ctx, &serviceList, &client.ListOptions{LabelSelector: selector, Namespace: namespace})
	if err != nil {
		r.Log.Error(err, "Failed to list services")
		return nil, err
	}

	return &serviceList, nil
}

func isFinalizing(cfRoute *networkingv1alpha1.CFRoute) bool {
	return cfRoute.ObjectMeta.DeletionTimestamp != nil && !cfRoute.ObjectMeta.DeletionTimestamp.IsZero()
}

func generateServiceName(destination *networkingv1alpha1.Destination) string {
	return fmt.Sprintf("s-%s", destination.GUID)
}
