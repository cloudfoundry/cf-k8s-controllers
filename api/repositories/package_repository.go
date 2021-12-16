package repositories

import (
	"context"
	"errors"
	"fmt"
	"sort"

	"k8s.io/apimachinery/pkg/types"

	"code.cloudfoundry.org/cf-k8s-controllers/api/authorization"
	workloadsv1alpha1 "code.cloudfoundry.org/cf-k8s-controllers/controllers/apis/workloads/v1alpha1"

	"github.com/google/uuid"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	kind = "CFPackage"

	PackageStateAwaitingUpload = "AWAITING_UPLOAD"
	PackageStateReady          = "READY"
)

//+kubebuilder:rbac:groups=workloads.cloudfoundry.org,resources=cfpackages,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=workloads.cloudfoundry.org,resources=cfpackages/status,verbs=get

//+kubebuilder:rbac:groups="",resources=serviceaccounts;secrets,verbs=get;list;watch
//+kubebuilder:rbac:groups="",resources=serviceaccounts/status;secrets/status,verbs=get

type PackageListMessage struct {
	AppGUIDs []string
}

type PackageCreateMessage struct {
	Type      string
	AppGUID   string
	SpaceGUID string
	OwnerRef  metav1.OwnerReference
}

type PackageUpdateSourceMessage struct {
	GUID               string
	SpaceGUID          string
	ImageRef           string
	RegistrySecretName string
}

type PackageRecord struct {
	GUID      string
	UID       types.UID
	Type      string
	AppGUID   string
	SpaceGUID string
	State     string
	CreatedAt string
	UpdatedAt string
}

type PackageRepo struct {
	privilegedClient client.Client
}

func NewPackageRepo(privilegedClient client.Client) *PackageRepo {
	return &PackageRepo{privilegedClient: privilegedClient}
}

func (r *PackageRepo) CreatePackage(ctx context.Context, authInfo authorization.Info, message PackageCreateMessage) (PackageRecord, error) {
	cfPackage := packageCreateToCFPackage(message)
	err := r.privilegedClient.Create(ctx, &cfPackage)
	if err != nil {
		return PackageRecord{}, err
	}
	return cfPackageToPackageRecord(cfPackage), nil
}

func (r *PackageRepo) FetchPackage(ctx context.Context, authInfo authorization.Info, guid string) (PackageRecord, error) {
	packageList := &workloadsv1alpha1.CFPackageList{}
	err := r.privilegedClient.List(ctx, packageList)
	if err != nil { // untested
		return PackageRecord{}, err
	}
	allPackages := packageList.Items
	matches := filterPackagesByMetadataName(allPackages, guid)

	return returnPackage(matches)
}

func (r *PackageRepo) FetchPackageList(ctx context.Context, authInfo authorization.Info, message PackageListMessage) ([]PackageRecord, error) {
	packageList := &workloadsv1alpha1.CFPackageList{}
	err := r.privilegedClient.List(ctx, packageList)
	if err != nil { // untested
		return []PackageRecord{}, err
	}

	filteredPackages := applyPackageFiltersAndOrder(packageList.Items, message)

	return returnPackageList(filteredPackages), nil
}

func applyPackageFiltersAndOrder(packages []workloadsv1alpha1.CFPackage, message PackageListMessage) []workloadsv1alpha1.CFPackage {
	var filtered []workloadsv1alpha1.CFPackage
	if len(message.AppGUIDs) > 0 {
		for _, currentPackage := range packages {
			for _, appGUID := range message.AppGUIDs {
				if currentPackage.Spec.AppRef.Name == appGUID {
					filtered = append(filtered, currentPackage)
					break
				}
			}
		}
	} else {
		filtered = packages
	}

	// TODO: use the future message.Order fields to reorder the list of results
	// For now, we order by created_at by default- if you really want to optimize runtime you can use bucketsort
	sort.Slice(filtered, func(i, j int) bool {
		return !filtered[i].CreationTimestamp.Before(&filtered[j].CreationTimestamp)
	})

	return filtered
}

func (r *PackageRepo) UpdatePackageSource(ctx context.Context, authInfo authorization.Info, message PackageUpdateSourceMessage) (PackageRecord, error) {
	baseCFPackage := &workloadsv1alpha1.CFPackage{
		ObjectMeta: metav1.ObjectMeta{
			Name:      message.GUID,
			Namespace: message.SpaceGUID,
		},
	}
	cfPackage := baseCFPackage.DeepCopy()
	cfPackage.Spec.Source.Registry.Image = message.ImageRef
	cfPackage.Spec.Source.Registry.ImagePullSecrets = []corev1.LocalObjectReference{{Name: message.RegistrySecretName}}

	err := r.privilegedClient.Patch(ctx, cfPackage, client.MergeFrom(baseCFPackage))
	if err != nil { // untested
		return PackageRecord{}, fmt.Errorf("err in client.Patch: %w", err)
	}

	record := cfPackageToPackageRecord(*cfPackage)
	return record, nil
}

func packageCreateToCFPackage(message PackageCreateMessage) workloadsv1alpha1.CFPackage {
	guid := uuid.NewString()
	return workloadsv1alpha1.CFPackage{
		TypeMeta: metav1.TypeMeta{
			Kind:       kind,
			APIVersion: workloadsv1alpha1.GroupVersion.Identifier(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:            guid,
			Namespace:       message.SpaceGUID,
			OwnerReferences: []metav1.OwnerReference{message.OwnerRef},
		},
		Spec: workloadsv1alpha1.CFPackageSpec{
			Type: workloadsv1alpha1.PackageType(message.Type),
			AppRef: corev1.LocalObjectReference{
				Name: message.AppGUID,
			},
		},
	}
}

func cfPackageToPackageRecord(cfPackage workloadsv1alpha1.CFPackage) PackageRecord {
	updatedAtTime, _ := getTimeLastUpdatedTimestamp(&cfPackage.ObjectMeta)
	state := PackageStateAwaitingUpload
	if cfPackage.Spec.Source.Registry.Image != "" {
		state = PackageStateReady
	}
	return PackageRecord{
		GUID:      cfPackage.Name,
		UID:       cfPackage.UID,
		SpaceGUID: cfPackage.Namespace,
		Type:      string(cfPackage.Spec.Type),
		AppGUID:   cfPackage.Spec.AppRef.Name,
		State:     state,
		CreatedAt: formatTimestamp(cfPackage.CreationTimestamp),
		UpdatedAt: updatedAtTime,
	}
}

func filterPackagesByMetadataName(packages []workloadsv1alpha1.CFPackage, name string) []workloadsv1alpha1.CFPackage {
	var filtered []workloadsv1alpha1.CFPackage
	for i, app := range packages {
		if app.Name == name {
			filtered = append(filtered, packages[i])
		}
	}
	return filtered
}

func returnPackage(packages []workloadsv1alpha1.CFPackage) (PackageRecord, error) {
	if len(packages) == 0 {
		return PackageRecord{}, NotFoundError{}
	}
	if len(packages) > 1 {
		return PackageRecord{}, errors.New("duplicate packages exist")
	}

	return cfPackageToPackageRecord(packages[0]), nil
}

func returnPackageList(packages []workloadsv1alpha1.CFPackage) []PackageRecord {
	packageRecords := make([]PackageRecord, 0, len(packages))

	for _, currentPackage := range packages {
		packageRecords = append(packageRecords, cfPackageToPackageRecord(currentPackage))
	}
	return packageRecords
}
