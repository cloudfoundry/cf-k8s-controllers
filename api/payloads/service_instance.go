package payloads

import (
	"strings"

	"code.cloudfoundry.org/cf-k8s-controllers/api/repositories"
)

type ServiceInstanceCreate struct {
	Name          string                       `json:"name" validate:"required"`
	Type          string                       `json:"type" validate:"required,oneof=user-provided"`
	Tags          []string                     `json:"tags" validate:"serviceinstancetaglength"`
	Credentials   map[string]string            `json:"credentials"`
	Relationships ServiceInstanceRelationships `json:"relationships" validate:"required"`
	Metadata      Metadata                     `json:"metadata"`
}

type ServiceInstanceRelationships struct {
	Space Relationship `json:"space" validate:"required"`
}

func (p ServiceInstanceCreate) ToServiceInstanceCreateMessage() repositories.CreateServiceInstanceMessage {
	return repositories.CreateServiceInstanceMessage{
		Name:        p.Name,
		SpaceGUID:   p.Relationships.Space.Data.GUID,
		Credentials: p.Credentials,
		Type:        p.Type,
		Tags:        p.Tags,
		Labels:      p.Metadata.Labels,
		Annotations: p.Metadata.Annotations,
	}
}

type ServiceInstanceList struct {
	Names      *string `schema:"names"`
	SpaceGuids *string `schema:"space_guids"`
	OrderBy    string  `schema:"order_by"`
}

func (l *ServiceInstanceList) ToMessage() repositories.ListServiceInstanceMessage {
	return repositories.ListServiceInstanceMessage{
		Names:           ParseArrayParam(l.Names),
		SpaceGuids:      ParseArrayParam(l.SpaceGuids),
		OrderBy:         strings.TrimPrefix(l.OrderBy, "-"),
		DescendingOrder: strings.HasPrefix(l.OrderBy, "-"),
	}
}

func (l *ServiceInstanceList) SupportedFilterKeys() []string {
	return []string{"names", "space_guids", "fields", "order_by", "per_page"}
}
