Here are all the properties that can be set for the korifi chart.
It also serves as documentation for each individual subchart,
where the values just omit the component header ID.
Global values apply to all components.

```yaml
global:
  # The namespace where the central CF resources are created
  rootNamespace: cf
  # Enable remote debugging by running delve and opening ports
  debug: false
  # Default suffix for app domains
  defaultAppDomainName: apps.my-cf-domain.com
  # Use cert-manager to generate self-signed certificates for the API and app endpoints
  generateIngressCertificates: false
  # The secret to use when pushing source and droplet images to the package registry
  packageRegistrySecret: image-registry-credentials

# Name of admin user that will be bound to the cf admin role
adminUserName: cf-admin

# API component configuration
api:
  # Deploy the API component
  include: true
  # Namespace for the API resources
  namespace: korifi-api-system
  # Number of replicas
  replicas: 1
  # Resource requests
  resources:
    requests:
      cpu: 50m
      memory: 100Mi

  apiServer:
    # API URL
    url: api.my-cf-domain.com # externalFQDN
    # External port. Defaults to 443. To override default port, set port to a non-zero value
    port: 0
    # Container port
    internalPort: 9000
    # HTTP timeouts
    timeouts:
      read: 900
      write: 900
      idle: 900
      readHeader: 10

  # Docker image
  image: cloudfoundry/korifi-api:latest

  # Lifecycle details
  lifecycle:
    type: buildpack
    stack: cflinuxfs3
    stagingRequirements:
      memoryMB: 1024
      diskMB: 1024

  # ID of the builder to use on source packages
  builderName: kpack-image-builder
  # The container registry where app source packages will be stored
  packageRegistry: registry-org/package-repo-name
  # Warn if user cert provided for login has a long expiry
  userCertificateExpirationWarningDuration: 168h

  # If using a Cluster authentication proxy, e.g. with pinniped, (optional)
  authProxy:
    # proxy Host IP address
    # Host must be a host string, a host:port pair, or a URL to the base of the apiserver.
    host:
    # proxy's PEM-encoded CA certificate (not base64'ed)
    caCert:

# Controller component configuration
controllers:
  # Deploy the controllers component
  include: true
  # Namespace for the controllers resources
  namespace: korifi-controllers-system
  # Number of replicas
  replicas: 1
  # Resource requests and limits
  resources:
    limits:
      cpu: 1000m
      memory: 1Gi
    requests:
      cpu: 50m
      memory: 100Mi

  # Docker image
  image: cloudfoundry/korifi-controllers:latest
  reconcilers:
    # Name of the image builder to set on all BuildWorkload objects. Has to match the builder's BuilderInfo name
    build: kpack-image-builder
    # Name of the workload runner to set on all AppWorkload objects.
    app: statefulset-runner
  processDefaults:
    # Default memory limit for the web process
    memoryMB: 1024
    # Default disk quota for the web process
    diskQuotaMB: 1024
  # How long before the CFTask object is deleted after the task has completed
  taskTTL: 30d
  # The TLS secret used when setting up app route
  workloadsTLSSecret:
    name: korifi-workloads-ingress-cert
    namespace: korifi-controllers-system

job-task-runner:
  # Deploy the job-task-runner component
  include: true
  # Namespace of the job-task-runner resources
  namespace: korifi-job-task-runner-system
  # Number of replicas
  replicas: 1
  # Resource requests and limits
  resources:
    limits:
      cpu: 500m
      memory: 128Mi
    requests:
      cpu: 10m
      memory: 64Mi

  # Docker image
  image: cloudfoundry/korifi-job-task-runner:latest
  # How long before the Job backing up a task is deleted after completion
  jobTTL: 24h

kpack-image-builder:
  # Deploy the kpack-image-builder component
  include: true
  # Namespace of the kpack-image-builder resources
  namespace: korifi-kpack-build-system
  # Number of replicas
  replicas: 1
  # Resource requests and limits
  resources:
    limits:
      cpu: 1000m
      memory: 1Gi
    requests:
      cpu: 50m
      memory: 100Mi

  # Docker image
  image: cloudfoundry/korifi-kpack-image-builder:latest

  # The image registry where droplet images are pushed to
  packageRegistry: registry-org/package-repo-name
  # The name of the cluster builder kpack has been configured with
  clusterBuilderName: cf-kpack-cluster-builder
  exampleClusterBuilder:
    # create an example cluster store, stack and builder
    create: false
    # registry location to store cluster builder image
    kpackBuilderRegistry: registry-org/kpack-builder-repo-name

statefulset-runner:
  # Deploy the statefulset-runner component
  include: true
  # Namespace of the statefulset-runner resources
  namespace: korifi-statefulset-runner-system
  # Number of replicas
  replicas: 1
  # Resource requests and limits
  resources:
    limits:
      cpu: 500m
      memory: 128Mi
    requests:
      cpu: 10m
      memory: 64Mi

  # Docker image
  image: cloudfoundry/korifi-statefulset-runner:latest
```
