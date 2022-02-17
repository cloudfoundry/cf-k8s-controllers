package repositories

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	hnsv1alpha2 "sigs.k8s.io/hierarchical-namespaces/api/v1alpha2"

	"code.cloudfoundry.org/cf-k8s-controllers/api/authorization"
	"code.cloudfoundry.org/cf-k8s-controllers/api/config"
)

var (
	ErrorDuplicateRoleBinding          = errors.New("RoleBinding with that name already exists")
	ErrorMissingRoleBindingInParentOrg = errors.New("no RoleBinding found in parent org")
)

const (
	RoleGuidLabel         = "cloudfoundry.org/role-guid"
	roleBindingNamePrefix = "cf"

	RoleResourceType = "Role"
)

//counterfeiter:generate -o fake -fake-name AuthorizedInChecker . AuthorizedInChecker

type AuthorizedInChecker interface {
	AuthorizedIn(ctx context.Context, identity authorization.Identity, namespace string) (bool, error)
}

//+kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=rolebindings,verbs=create

type CreateRoleMessage struct {
	GUID  string
	Type  string
	Space string
	Org   string
	User  string
	Kind  string
}

type RoleRecord struct {
	GUID      string
	CreatedAt time.Time
	UpdatedAt time.Time
	Type      string
	Space     string
	Org       string
	User      string
	Kind      string
}

type RoleRepo struct {
	privilegedClient    client.Client
	roleMappings        map[string]config.Role
	authorizedInChecker AuthorizedInChecker
	userClientFactory   UserK8sClientFactory
}

func NewRoleRepo(privilegedClient client.Client, userClientFactory UserK8sClientFactory, authorizedInChecker AuthorizedInChecker, roleMappings map[string]config.Role) *RoleRepo {
	return &RoleRepo{
		privilegedClient:    privilegedClient,
		userClientFactory:   userClientFactory,
		roleMappings:        roleMappings,
		authorizedInChecker: authorizedInChecker,
	}
}

func (r *RoleRepo) CreateRole(ctx context.Context, authInfo authorization.Info, role CreateRoleMessage) (RoleRecord, error) {
	userClient, err := r.userClientFactory.BuildClient(authInfo)
	if err != nil {
		return RoleRecord{}, fmt.Errorf("failed to build user client: %w", err)
	}

	k8sRoleConfig, ok := r.roleMappings[role.Type]
	if !ok {
		return RoleRecord{}, fmt.Errorf("invalid role type: %q", role.Type)
	}

	userIdentity := authorization.Identity{
		Name: role.User,
		Kind: role.Kind,
	}

	if role.Space != "" {
		if err = r.validateOrgRequirements(ctx, role, userIdentity); err != nil {
			return RoleRecord{}, err
		}
	}

	ns := role.Space
	if ns == "" {
		ns = role.Org
	}

	annotations := map[string]string{}
	if !k8sRoleConfig.Propagate {
		annotations[hnsv1alpha2.AnnotationNoneSelector] = "true"
	}

	roleBinding := rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: ns,
			Name:      calculateRoleBindingName(role),
			Labels: map[string]string{
				RoleGuidLabel: role.GUID,
			},
			Annotations: annotations,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind: role.Kind,
				Name: role.User,
			},
		},
		RoleRef: rbacv1.RoleRef{
			Kind: "ClusterRole",
			Name: k8sRoleConfig.Name,
		},
	}

	err = userClient.Create(ctx, &roleBinding)
	if err != nil {
		if k8serrors.IsAlreadyExists(err) {
			return RoleRecord{}, ErrorDuplicateRoleBinding
		}
		if k8serrors.IsForbidden(err) {
			return RoleRecord{}, NewForbiddenError(RoleResourceType, err)
		}
		return RoleRecord{}, fmt.Errorf("failed to assign user %q to role %q: %w", role.User, role.Type, err)
	}

	roleRecord := RoleRecord{
		GUID:      role.GUID,
		CreatedAt: roleBinding.CreationTimestamp.Time,
		UpdatedAt: roleBinding.CreationTimestamp.Time,
		Type:      role.Type,
		Space:     role.Space,
		Org:       role.Org,
		User:      role.User,
		Kind:      role.Kind,
	}

	return roleRecord, nil
}

func (r *RoleRepo) validateOrgRequirements(ctx context.Context, role CreateRoleMessage, userIdentity authorization.Identity) error {
	orgName, err := r.getOrgName(ctx, role.Space)
	if err != nil {
		return err
	}

	hasOrgBinding, err := r.authorizedInChecker.AuthorizedIn(ctx, userIdentity, orgName)
	if err != nil {
		return fmt.Errorf("failed to check for role in parent org: %w", err)
	}

	if !hasOrgBinding {
		return ErrorMissingRoleBindingInParentOrg
	}
	return nil
}

func (r *RoleRepo) getOrgName(ctx context.Context, spaceGUID string) (string, error) {
	namespace := corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: spaceGUID,
		},
	}

	err := r.privilegedClient.Get(ctx, client.ObjectKeyFromObject(&namespace), &namespace)
	if err != nil {
		return "", fmt.Errorf("failed to get namespace with name %q: %w", spaceGUID, err)
	}

	orgName := namespace.Annotations[hnsv1alpha2.SubnamespaceOf]
	if orgName == "" {
		return "", fmt.Errorf("namespace %s does not have a parent", spaceGUID)
	}

	return orgName, nil
}

func calculateRoleBindingName(role CreateRoleMessage) string {
	plain := []byte(role.Type + "::" + role.User)
	sum := sha256.Sum256(plain)

	return fmt.Sprintf("%s-%x", roleBindingNamePrefix, sum)
}
