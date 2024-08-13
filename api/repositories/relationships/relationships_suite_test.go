package relationships_test

import (
	"context"
	"testing"

	"code.cloudfoundry.org/korifi/api/authorization"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	ctx      context.Context
	authInfo authorization.Info
)

var _ = BeforeSuite(func() {
	ctx = context.Background()
	authInfo = authorization.Info{Token: "my-token"}
})

func TestRelationships(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Relationships Suite")
}
