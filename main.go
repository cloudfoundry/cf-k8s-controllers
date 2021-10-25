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

package main

import (
	"flag"
	"fmt"
	"os"

	networkingv1alpha1 "code.cloudfoundry.org/cf-k8s-controllers/apis/networking/v1alpha1"
	workloadsv1alpha1 "code.cloudfoundry.org/cf-k8s-controllers/apis/workloads/v1alpha1"
	config "code.cloudfoundry.org/cf-k8s-controllers/config/base"
	networkingcontrollers "code.cloudfoundry.org/cf-k8s-controllers/controllers/networking"
	workloadscontrollers "code.cloudfoundry.org/cf-k8s-controllers/controllers/workloads"
	"code.cloudfoundry.org/cf-k8s-controllers/controllers/workloads/imageprocessfetcher"
	"code.cloudfoundry.org/cf-k8s-controllers/webhooks/workloads"

	eiriniv1 "code.cloudfoundry.org/eirini-controller/pkg/apis/eirini/v1"
	buildv1alpha1 "github.com/pivotal/kpack/pkg/apis/build/v1alpha1"
	contourv1 "github.com/projectcontour/contour/apis/projectcontour/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	k8sclient "k8s.io/client-go/kubernetes"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	hnsv1alpha2 "sigs.k8s.io/hierarchical-namespaces/api/v1alpha2"
	//+kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(workloadsv1alpha1.AddToScheme(scheme))
	utilruntime.Must(networkingv1alpha1.AddToScheme(scheme))
	utilruntime.Must(buildv1alpha1.AddToScheme(scheme))
	utilruntime.Must(contourv1.AddToScheme(scheme))
	utilruntime.Must(eiriniv1.AddToScheme(scheme))
	utilruntime.Must(hnsv1alpha2.AddToScheme(scheme))
	//+kubebuilder:scaffold:scheme
}

func main() {
	var metricsAddr string
	var enableLeaderElection bool
	var probeAddr string

	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")

	opts := zap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     metricsAddr,
		Port:                   9443,
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "13c200ec.cloudfoundry.org",
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	configPath, found := os.LookupEnv("CONFIG")
	if !found {
		panic("CONFIG must be set")
	}

	controllerConfig, err := config.LoadConfigFromPath(configPath)
	if err != nil {
		errorMessage := fmt.Sprintf("Config could not be read: %v", err)
		panic(errorMessage)
	}

	k8sClientConfig := ctrl.GetConfigOrDie()
	privilegedK8sClient, err := k8sclient.NewForConfig(k8sClientConfig)
	if err != nil {
		panic(fmt.Sprintf("could not create privileged k8s client: %v", err))
	}

	// Setup with manager

	if err = (&workloadscontrollers.CFAppReconciler{
		Client:           mgr.GetClient(),
		Scheme:           mgr.GetScheme(),
		Log:              ctrl.Log.WithName("controllers").WithName("CFApp"),
		ControllerConfig: controllerConfig,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "CFApp")
		os.Exit(1)
	}

	cfBuildImageProcessFetcher := &imageprocessfetcher.ImageProcessFetcher{
		Log: ctrl.Log.WithName("controllers").WithName("CFBuildImageProcessFetcher"),
	}
	if err = (&workloadscontrollers.CFBuildReconciler{
		Client:              mgr.GetClient(),
		Scheme:              mgr.GetScheme(),
		Log:                 ctrl.Log.WithName("controllers").WithName("CFBuild"),
		ControllerConfig:    controllerConfig,
		RegistryAuthFetcher: workloadscontrollers.NewRegistryAuthFetcher(privilegedK8sClient),
		ImageProcessFetcher: cfBuildImageProcessFetcher.Fetch,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "CFBuild")
		os.Exit(1)
	}

	if err = (&networkingcontrollers.CFDomainReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "CFDomain")
		os.Exit(1)
	}

	if err = (&workloadscontrollers.CFPackageReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "CFPackage")
		os.Exit(1)
	}

	if err = (&workloadscontrollers.CFProcessReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "CFProcess")
		os.Exit(1)
	}

	if err = (&networkingcontrollers.CFRouteReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
		Log:    ctrl.Log.WithName("controllers").WithName("CFRoute"),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "CFRoute")
		os.Exit(1)
	}

	// Setup webhooks with manager

	if os.Getenv("ENABLE_WEBHOOKS") != "false" {
		if err = (&workloadsv1alpha1.CFApp{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "CFApp")
			os.Exit(1)
		}

		if err = (&workloadsv1alpha1.CFPackage{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "CFPackage")
			os.Exit(1)
		}
		if err = (&workloadsv1alpha1.CFBuild{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "CFBuild")
			os.Exit(1)
		}

		if err = (&workloadsv1alpha1.CFProcess{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "CFProcess")
			os.Exit(1)
		}

		if err = (&workloads.CFAppValidation{
			Client: mgr.GetClient(),
		}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "CFApp")
			os.Exit(1)
		}

		if err = workloads.NewSubnamespaceAnchorValidation(mgr.GetClient()).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "SubnamespaceAnchors")
			os.Exit(1)
		}
	} else {
		setupLog.Info("Skipping webhook setup because ENABLE_WEBHOOKS set to false.")
	}

	//+kubebuilder:scaffold:builder

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}

	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
