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

	networkingv1alpha1 "code.cloudfoundry.org/cf-k8s-controllers/apis/networking/v1alpha1"

	"github.com/go-logr/logr"
	contourv1 "github.com/projectcontour/contour/apis/projectcontour/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	Client client.Client
	Scheme *runtime.Scheme
	Log    logr.Logger
}

//+kubebuilder:rbac:groups=networking.cloudfoundry.org,resources=cfroutes,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=networking.cloudfoundry.org,resources=cfroutes/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=networking.cloudfoundry.org,resources=cfroutes/finalizers,verbs=update

//+kubebuilder:rbac:groups=projectcontour.io,resources=httpproxies,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=projectcontour.io,resources=httpproxies/status,verbs=get
//+kubebuilder:rbac:groups=projectcontour.io,resources=httpproxies/finalizers,verbs=update

//+kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch;delete

func (r *CFRouteReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var cfRoute networkingv1alpha1.CFRoute
	err := r.Client.Get(ctx, req.NamespacedName, &cfRoute)
	if err != nil {
		if !apierrors.IsNotFound(err) {
			r.Log.Error(err, "failed to get CFRoute")
		}
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	var cfDomain networkingv1alpha1.CFDomain
	err = r.Client.Get(ctx, types.NamespacedName{Name: cfRoute.Spec.DomainRef.Name}, &cfDomain)
	if err != nil {
		r.Log.Error(err, "failed to get CFDomain")

		// TODO: General status management in follow up story, possibly set CFRoute to status invalid?
		return ctrl.Result{}, err
	}

	// Check all namespaces for FQDN proxy with the matching FQDN label
	var proxies contourv1.HTTPProxyList
	err = r.Client.List(ctx, &proxies)
	if err != nil {
		r.Log.Error(err, "failed to list HTTPProxies")
		return ctrl.Result{}, err
	}

	var fqdnHTTPProxy contourv1.HTTPProxy
	fqdn := fmt.Sprintf("%s.%s", cfRoute.Spec.Host, cfDomain.Spec.Name)

	found := false
	for _, proxy := range proxies.Items {
		if proxy.Spec.VirtualHost != nil && proxy.Spec.VirtualHost.Fqdn == fqdn {
			if found {
				err = fmt.Errorf("found multiple HTTPProxy with FQDN %s", fqdn)
				r.Log.Error(err, "")
				return ctrl.Result{}, err
			} else if proxy.Namespace != cfRoute.Namespace {
				err = fmt.Errorf("found existing HTTPProxy with FQDN %s in another space", fqdn)
				r.Log.Error(err, fmt.Sprintf("existing proxy found in namespace %s", proxy.Namespace))
				return ctrl.Result{}, err
			}

			fqdnHTTPProxy = proxy
			found = true
		}
	}

	// If proxy with desired FQDN not found, create in current namespace
	if !found {
		fqdnHTTPProxy = contourv1.HTTPProxy{
			ObjectMeta: metav1.ObjectMeta{
				Name:      fqdn,
				Namespace: cfRoute.Namespace,
			},
		}
	}

	if cfRoute.ObjectMeta.DeletionTimestamp != nil && cfRoute.ObjectMeta.DeletionTimestamp.IsZero() == false {
		r.Log.Info(fmt.Sprintf("Reconciling deletion of CFRoute/%s", cfRoute.Name))
		// Cleanup the FQDN HTTPProxy on delete
		if hasFinalizer(&cfRoute, FinalizerName) {
			err = r.finalizeCFRouteForDeletion(ctx, req, &cfRoute, &fqdnHTTPProxy, found)
			if err != nil {
				return ctrl.Result{}, err
			}

			controllerutil.RemoveFinalizer(&cfRoute, FinalizerName)
			if err := r.Client.Update(ctx, &cfRoute); err != nil {
				r.Log.Error(err, "Failed to remove finalizer")
				return ctrl.Result{}, err
			}
		}

		return ctrl.Result{}, nil
	}

	// Update FQDN proxy with include for new sub-proxy
	result, err := controllerutil.CreateOrPatch(ctx, r.Client, &fqdnHTTPProxy, func() error {
		fqdnHTTPProxy.Spec.VirtualHost = &contourv1.VirtualHost{
			Fqdn: fqdn,
		}

		found := false
		for _, include := range fqdnHTTPProxy.Spec.Includes {
			if include.Name == cfRoute.Name && include.Namespace == cfRoute.Namespace {
				found = true
			}
		}

		if !found {
			fqdnHTTPProxy.Spec.Includes = append(fqdnHTTPProxy.Spec.Includes, contourv1.Include{
				Name:      cfRoute.Name,
				Namespace: cfRoute.Namespace,
			})
		}

		err = controllerutil.SetOwnerReference(&cfRoute, &fqdnHTTPProxy, r.Scheme)
		if err != nil {
			r.Log.Error(err, "failed to set OwnerRef on FQDN HTTPProxy")
			return err
		}

		return nil
	})
	if err != nil {
		r.Log.Error(err, "failed to patch FQDN HTTPProxy")
		return ctrl.Result{}, err
	}
	r.Log.Info(fmt.Sprintf("FQDN HTTPProxy/%s %s", fqdnHTTPProxy.Name, result))

	routeHTTPProxy := &contourv1.HTTPProxy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cfRoute.Name,
			Namespace: cfRoute.Namespace,
		},
	}

	result, err = controllerutil.CreateOrPatch(ctx, r.Client, routeHTTPProxy, func() error {
		desiredRoutes := make([]contourv1.Route, 0, len(cfRoute.Spec.Destinations))

		for _, destination := range cfRoute.Spec.Destinations {
			desiredRoute := contourv1.Route{
				Conditions: []contourv1.MatchCondition{
					{
						Prefix: cfRoute.Spec.Path,
					},
				},
				Services: []contourv1.Service{
					{
						Name: fmt.Sprintf("s-%s-%s", destination.AppRef.Name, destination.ProcessType),
						Port: destination.Port,
					},
				},
			}
			desiredRoutes = append(desiredRoutes, desiredRoute)
		}

		routeHTTPProxy.Spec.Routes = desiredRoutes

		err = controllerutil.SetOwnerReference(&cfRoute, routeHTTPProxy, r.Scheme)
		if err != nil {
			r.Log.Error(err, "failed to set OwnerRef on route HTTPProxy")
			return err
		}

		return nil
	})
	if err != nil {
		r.Log.Error(err, "failed to patch route HTTPProxy")
		return ctrl.Result{}, err
	}
	r.Log.Info(fmt.Sprintf("Route HTTPProxy/%s %s", routeHTTPProxy.Name, result))

	var serviceReconcileErr error
	for _, destination := range cfRoute.Spec.Destinations {
		service := &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				// TODO: Make name GUID to avoid complexity on delete trying to check for other proxies referring to the service
				Name:      fmt.Sprintf("s-%s-%s", destination.AppRef.Name, destination.ProcessType),
				Namespace: cfRoute.Namespace,
			},
		}

		result, err = controllerutil.CreateOrPatch(ctx, r.Client, service, func() error {
			service.ObjectMeta.Labels = map[string]string{
				"workloads.cloudfoundry.org/app-guid": destination.AppRef.Name,
			}

			err = controllerutil.SetOwnerReference(&cfRoute, service, r.Scheme)
			if err != nil {
				r.Log.Error(err, "failed to set OwnerRef on Service")
				return err
			}

			service.Spec.Ports = []corev1.ServicePort{{
				Port: int32(destination.Port),
			}}
			service.Spec.Selector = map[string]string{
				"workloads.cloudfoundry.org/app-guid":     destination.AppRef.Name,
				"workloads.cloudfoundry.org/process-type": destination.ProcessType,
			}

			return nil
		})
		if err != nil {
			r.Log.Error(err, fmt.Sprintf("failed to patch Service/%s", service.Name))
			serviceReconcileErr = fmt.Errorf("service reconciliation failed for CFRoute/%s destinations", cfRoute.Name)
		} else {
			r.Log.Info(fmt.Sprintf("Service/%s %s", service.Name, result))
		}
	}

	if serviceReconcileErr != nil {
		return ctrl.Result{}, serviceReconcileErr
	}

	// Add the finalizer
	result, err = controllerutil.CreateOrPatch(ctx, r.Client, &cfRoute, func() error {
		controllerutil.AddFinalizer(&cfRoute, FinalizerName)
		return nil
	})
	if err != nil {
		r.Log.Error(err, fmt.Sprintf("Error updating CFRoute/%s", cfRoute.Name))
	} else {
		r.Log.Info(fmt.Sprintf("CFRoute/%s %s with finalizer", cfRoute.Name, result))
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *CFRouteReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&networkingv1alpha1.CFRoute{}).
		Complete(r)
}

func hasFinalizer(o metav1.Object, finalizerName string) bool {
	for _, f := range o.GetFinalizers() {
		if f == finalizerName {
			return true
		}
	}
	return false
}

func (r *CFRouteReconciler) finalizeCFRouteForDeletion(ctx context.Context, req ctrl.Request, cfRoute *networkingv1alpha1.CFRoute, fqdnHTTPProxy *contourv1.HTTPProxy, found bool) error {
	// If the FQDN HTTPProxy was not found, there is nothing to do
	if !found {
		return nil
	}

	// Remove the sub-HTTPProxy (name equal to the CFRoute name) from the list of includes
	_, err := controllerutil.CreateOrPatch(ctx, r.Client, fqdnHTTPProxy, func() error {
		for idx, include := range fqdnHTTPProxy.Spec.Includes {
			if include.Name == cfRoute.Name {
				r.Log.Info(fmt.Sprintf("Removing sub-HTTPProxy for route %s from FQDN HTTPProxy", cfRoute.Name))
				fqdnHTTPProxy.Spec.Includes[idx] = fqdnHTTPProxy.Spec.Includes[len(fqdnHTTPProxy.Spec.Includes)-1]
				fqdnHTTPProxy.Spec.Includes = fqdnHTTPProxy.Spec.Includes[:len(fqdnHTTPProxy.Spec.Includes)-1]
			}
		}

		return nil
	})
	if err != nil {
		r.Log.Error(err, "failed to patch FQDN HTTPProxy to remove sub HTTPProxy")
		return err
	}

	return nil
}
