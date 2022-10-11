> **Warning**
> Make sure you are using the correct version of these instructions by using the link in the release notes for the version you're trying to install. If you're not sure, check our [latest release](https://github.com/cloudfoundry/korifi/releases/latest).

# Korifi installation guide

The following lines will guide you through the process of deploying a [released version](https://github.com/cloudfoundry/korifi/releases) of [Korifi](https://github.com/cloudfoundry/korifi). This document is written with the intent to act both as a runbook as well as a starting point in understanding basic concepts of Korifi and its dependencies.

## Prerequisites

-   Tools:
    -   [kubectl](https://kubernetes.io/docs/tasks/tools/#kubectl);
    -   [cf](https://docs.cloudfoundry.org/cf-cli/install-go-cli.html) CLI version 8.5 or greater;
    -   [Helm](https://helm.sh/docs/intro/install/).
-   Resources:
    -   Kubernetes cluster of one of the [upstream releases](https://kubernetes.io/releases/);
    -   Container Registry on which you have write permissions.

This document was tested on:

-   [EKS](https://aws.amazon.com/eks/), using GCP's [Artifact Registry](https://cloud.google.com/artifact-registry);
-   [GKE](https://cloud.google.com/kubernetes-engine), using GCP's [Artifact Registry](https://cloud.google.com/artifact-registry);
-   [kind](https://kind.sigs.k8s.io/), using [DockerHub](https://hub.docker.com/) (see [_Install Korifi on kind_](./INSTALL_kind.md)).

## Initial setup

The following environment variables will be needed throughout this guide:

-   `ROOT_NAMESPACE`: the namespace at the root of the Korifi org and space hierarchy. The default value is `cf`.
-   `ADMIN_USERNAME`: the name of the Kubernetes user who will have CF admin privileges on the Korifi installation. For security reasons, you should choose or create a user that is different from your cluster admin user. To provision new users, follow the user management instructions specific for your cluster's [authentication configuration](https://kubernetes.io/docs/reference/access-authn-authz/authentication/) or create a [new (short-lived) client certificate for user authentication](https://kubernetes.io/docs/reference/access-authn-authz/certificate-signing-requests/#normal-user).
-   `BASE_DOMAIN`: the base domain used by both the Korifi API and, by default, all apps running on Korifi.

Here are the example values we'll use in this guide:

```sh
export ROOT_NAMESPACE="cf"
export ADMIN_USERNAME="cf-admin"
export BASE_DOMAIN="korifi.example.org"
```

### Registries with Custom CA

See [_Using container registry signed by custom CA_](docs/using-container-registry-signed-by-custom-ca.md).

## Dependencies

### cert-Manager

[cert-Manager](https://cert-manager.io) allows us to automatically create internal certificates within the cluster. Follow the [instructions](https://cert-manager.io/docs/installation/) to install the latest version.

### Kpack

[Kpack](https://github.com/pivotal/kpack) is used to build runnable applications from source code using [Cloud Native Buildpacks](https://buildpacks.io/). Follow the [instructions](https://github.com/pivotal/kpack/blob/main/docs/install.md) to install the latest version.

### Contour

[Contour](https://projectcontour.io/) is our [ingress](https://kubernetes.io/docs/concepts/services-networking/ingress/) controller. Follow the [instructions](https://projectcontour.io/getting-started/#install-contour-and-envoy) from the getting started guide to install the latest version.

### Metrics Server

We use the [Kubernetes Metrics Server](https://github.com/kubernetes-sigs/metrics-server) to implement [process stats](https://v3-apidocs.cloudfoundry.org/#get-stats-for-a-process).
Most Kubernetes distributions will come with `metrics-server` already installed. If yours does not, you should follow the [instructions](https://github.com/kubernetes-sigs/metrics-server#installation) to install it.

### Optional: Service Bindings Controller

We use the [Service Binding Specification for Kubernetes](https://github.com/servicebinding/spec) and its [controller reference implementation](https://github.com/servicebinding/runtime) to implement [Cloud Foundry service bindings](https://docs.cloudfoundry.org/devguide/services/application-binding.html) ([see this issue](https://github.com/cloudfoundry/cf-k8s-controllers/issues/462)). Follow the [instructions](https://github.com/servicebinding/runtime/releases/latest) to install the latest version.

## Deploy Korifi

As of v0.4.0 Korifi is distributed as a helm chart. You can set the required configuration inline as below, or use a values file to store the settings.

```
helm install korifi https://github.com/cloudfoundry/korifi/releases/download/v<version>/korifi-<version>.tgz \
    --set=global.generateIngressCertificates=true \
    --set=global.rootNamespace=$ROOT_NAMESPACE \
    --set=adminUserName=$ADMIN_USERNAME \
    --set=api.apiServer.url=api.$BASE_DOMAIN \
    --set=global.defaultAppDomainName=apps.$BASE_DOMAIN \
    --set=api.packageRegistry=us-east4-docker.pkg.dev/vigilant-card-347116/korifi/packages \
    --set=kpack-image-builder.exampleClusterBuilder.create=true \
    --set=kpack-image-builder.exampleClusterBuilder.kpackBuilderRegistry=us-east4-docker.pkg.dev/vigilant-card-347116/korifi/kpack \
    --set=kpack-image-builder.packageRegistry=us-east4-docker.pkg.dev/vigilant-card-347116/korifi/droplets
```

### Description of helm values

- `global.generateIngressCertificates` when set to `true` generates self-signed certificates for the applications and API HTTP endpoint
- `global.rootNamespace` is the name of the CF root namespace containing base CF resources, like CFOrgs.
- `adminUserName` is the username that will be bound to the cf admin role.
- `api.apiServer.url` is the domain name that will be used by the Korifi API, and is usually of the format `api.$BASE_DOMAIN`.
- `global.defaultAppDomainName` is the default base domain name for the apps deployed by Korifi, and is usually of the format `apps.$BASE_DOMAIN`.
- `api.packageRegistry` specifies the tag prefix used for the source packages uploaded to Korifi. Its hostname should point to your container registry and its path should be valid for the registry.
  - If using **DockerHub**, `api.packageRegistry` should be `index.docker.io/<username>`.
  - If using **GCR**, `api.packageRegistry` should be `gcr.io/<project-id>/packages`.
- `kpack-image-builder.exampleClusterBuilder.create` activates creation of the example kpack cluster builder, store and stack resources.
- `kpack-image-builder.exampleClusterBuilder.kpackBuilderRegistry` is the registry location for the kpack builder image (pushed by kpack).
- `kpack-image-builder.packageRegistry` specifies the tag prefix used for the images built by Korifi. Its hostname should point to your container registry and its path should be valid for the registry.
  - If using **DockerHub**, `kpack-image-builder.packageRegistry` should be `index.docker.io/<username>`.
  - If using **GCR**, `kpack-image-builder.packageRegistry` should be `gcr.io/<project-id>/droplets`.

The chart provides various other values that can be set. See [helm/README.values.md](./helm/README.values.md) for details.

### Configure an Authentication Proxy (optional)

If you are using an authentication proxy with your cluster to enable SSO, you must alter the above `helm install` command to set the following values:

-   Set the `api.authProxy.host` helm value to the IP address of your cluster's auth proxy.
-   Set the `api.authProxy.caCert` helm value to the CA certificate of your cluster's auth proxy.

## Post-install Configuration

### Kpack Configuration

The korifi helm chart will create an example kpack configuration (cluster builder, cluster store and cluster stack) if the `kpack-image-builder.exampleClusterBuilder.create` helm property has been set to `true`.
You can opt out of doing that by setting the property to `false` (that's the default value).
In that case you have to configure those yourself.

#### `ClusterStore`

Follow the [documentation](https://github.com/pivotal/kpack/blob/main/docs/store.md) to create a `ClusterStore` for your cluster.

#### `ClusterStack`:

Follow the [documentation](https://github.com/pivotal/kpack/blob/main/docs/stack.md) to create a `ClusterStack` for your cluster.

#### `ClusterBuilder`

Follow the [documentation](https://github.com/pivotal/kpack/blob/main/docs/builders.md#cluster-builders) to create a `ClusterBuilder` for your cluster. Make sure that:

-   `metadata.name` matches the `korifi-image-builder.clusterBuilderName` helm value (default is `cf-kpack-cluster-builder`)
-   `spec.tag` points to your container registry:
    -   if using **DockerHub**, it should be `index.docker.io/<username>/korifi-cluster-builder`;
    -   if using **GCP**, it should be `gcr.io/<project-id>/korifi-cluster-builder`;
-   `spec.stack` references to the previously created `ClusterStack`;
-   `spec.store` references to the previously created `ClusterStore`;
-   `spec.serviceAccountRef` should be `kpack-service-account`.

### Container registry credentials `Secret`

Use the following command to create a `Secret` that Korifi and kpack will use to connect to your container registry:

```sh
kubectl --namespace "$ROOT_NAMESPACE" create secret docker-registry image-registry-credentials \
    --docker-username="<your-container-registry-username>" \
    --docker-password="<your-container-registry-password>" \
    --docker-server="<your-container-registry-hostname-and-port>"
```

Make sure the value of `--docker-server` is a valid [URI authority](https://datatracker.ietf.org/doc/html/rfc3986#section-3.2).

-   If using **DockerHub**:
    -   `--docker-server` should be `https://index.docker.io/v1/`;
    -   `--docker-username` should be your DockerHub user;
    -   `--docker-password` can be either your DockerHub password or a [generated personal access token](https://hub.docker.com/settings/security?generateToken=true).
-   If using **GCR**:
    -   `--docker-server` should be `gcr.io`;
    -   `--docker-username` should be `_json_key`;
    -   `--docker-password` should be the JSON-formatted access token for a service account that has permission to manage images in GCR.

### TLS certificates

Self-signed TLS certificates are generated automatically by the installation if `global.generateIngressCertificates` has been set to `true`.

If you want to generate certificates yourself, you should not set the `global.generateIngressCertificates` value, and instead provide your certificates to Korifi by creating two TLS secrets:

1. `korifi-api-ingress-cert` in the `korifi-api-system` namespace, and
2. `korifi-workloads-ingress-cert` in the `korifi-controllers-system` namespace

with the appropriate values.

### DNS

Create DNS entries for the Korifi API and for the apps running on Korifi. They should match the halm values when [deploying korifi](#deploy korifi):

-   The Korifi API entry should match the `api.apiServer.url` helm value. In our example, that would be `api.korifi.example.org`.
-   The apps entry should be a wildcard matching the `global.defaultAppDomainName` helm value: In our example, `*.apps.korifi.example.org`.

The DNS entries should point to the load balancer endpoint created by Contour when installed. To discover your endpoint, run:

```sh
kubectl get service envoy -n projectcontour -ojsonpath='{.status.loadBalancer.ingress[0]}'
```

It may take some time before the address is available. Retry this until you see a result.

The type of DNS records to create will differ based on the type of the endpoint: `ip` endpoints (e.g. the ones created by GKE) will need an `A` record, while `hostname` endpoints (e.g. on EKS) a `CNAME` record.

## Test Korifi

```sh
cf api https://api.$BASE_DOMAIN --skip-ssl-validation
cf auth $ADMIN_USERNAME
cf create-org org1
cf create-space -o org1 space1
cf target -o org1
cd <directory of a test cf app>
cf push test-app
```
