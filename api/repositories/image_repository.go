package repositories

import (
	"context"
	"errors"
	"fmt"
	"io"

	"code.cloudfoundry.org/korifi/api/authorization"
	apierrors "code.cloudfoundry.org/korifi/api/errors"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/authn/k8schain"
	registryv1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/remote"

	authv1 "k8s.io/api/authorization/v1"
	k8sclient "k8s.io/client-go/kubernetes"
)

//+kubebuilder:rbac:groups="",resources=serviceaccounts,verbs=get,namespace=ROOT_NAMESPACE
//+kubebuilder:rbac:groups="",resources=secrets,verbs=get,namespace=ROOT_NAMESPACE

const SourceImageResourceType = "SourceImage"

//counterfeiter:generate -o fake -fake-name ImageBuilder . ImageBuilder
//counterfeiter:generate -o fake -fake-name ImagePusher . ImagePusher

type ImageBuilder interface {
	Build(ctx context.Context, srcReader io.Reader) (registryv1.Image, error)
}

type ImagePusher interface {
	Push(repoRef string, image registryv1.Image, credentials remote.Option) (string, error)
}

type ImageRepository struct {
	privilegedK8sClient k8sclient.Interface
	userClientFactory   authorization.UserK8sClientFactory
	rootNamespace       string
	registrySecretName  string

	builder ImageBuilder
	pusher  ImagePusher
}

func NewImageRepository(
	privilegedK8sClient k8sclient.Interface,
	userClientFactory authorization.UserK8sClientFactory,
	rootNamespace,
	registrySecretName string,
	builder ImageBuilder,
	pusher ImagePusher,
) *ImageRepository {
	return &ImageRepository{
		privilegedK8sClient: privilegedK8sClient,
		userClientFactory:   userClientFactory,
		rootNamespace:       rootNamespace,
		registrySecretName:  registrySecretName,
		builder:             builder,
		pusher:              pusher,
	}
}

func (r *ImageRepository) UploadSourceImage(ctx context.Context, authInfo authorization.Info, imageRef string, srcReader io.Reader, spaceGUID string) (string, error) {
	authorized, err := r.canIPatchCFPackage(ctx, authInfo, spaceGUID)
	if err != nil {
		return "", fmt.Errorf("checking auth to upload source image for failed: %w", err)
	}

	if !authorized {
		return "", apierrors.NewForbiddenError(errors.New("not authorized to patch cfpackage"), PackageResourceType)
	}

	image, err := r.builder.Build(ctx, srcReader)
	if err != nil {
		return "", fmt.Errorf("image build for ref '%s' failed: %w", imageRef, err)
	}

	credentials, err := r.getCredentials(ctx)
	if err != nil {
		return "", fmt.Errorf("getting push credentials for image ref '%s' failed: %w", imageRef, err)
	}

	pushedRef, err := r.pusher.Push(imageRef, image, credentials)
	if err != nil {
		return "", apierrors.NewBlobstoreUnavailableError(fmt.Errorf("pushing image ref '%s' failed: %w", imageRef, err))
	}

	return pushedRef, nil
}

func (r *ImageRepository) canIPatchCFPackage(ctx context.Context, authInfo authorization.Info, spaceGUID string) (bool, error) {
	userClient, err := r.userClientFactory.BuildClient(authInfo)
	if err != nil {
		return false, fmt.Errorf("canIPatchCFPackage: failed to create user k8s client: %w", err)
	}

	review := authv1.SelfSubjectAccessReview{
		Spec: authv1.SelfSubjectAccessReviewSpec{
			ResourceAttributes: &authv1.ResourceAttributes{
				Namespace: spaceGUID,
				Verb:      "patch",
				Group:     "korifi.cloudfoundry.org",
				Resource:  "cfpackages",
			},
		},
	}
	if err := userClient.Create(ctx, &review); err != nil {
		return false, fmt.Errorf("canIPatchCFPackage: failed to create self subject access review: %w", apierrors.FromK8sError(err, PackageResourceType))
	}

	return review.Status.Allowed, nil
}

func (r *ImageRepository) getCredentials(ctx context.Context) (remote.Option, error) {
	var keychain authn.Keychain
	var err error

	if r.registrySecretName == "" {
		keychain, err = k8schain.NewNoClient(ctx)
	} else {
		keychain, err = k8schain.New(ctx, r.privilegedK8sClient, k8schain.Options{
			Namespace:        r.rootNamespace,
			ImagePullSecrets: []string{r.registrySecretName},
		})
	}
	if err != nil {
		return nil, fmt.Errorf("failed creating registry keychain: %w", apierrors.FromK8sError(err, SourceImageResourceType))
	}

	return remote.WithAuthFromKeychain(keychain), nil
}
