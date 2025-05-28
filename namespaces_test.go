package monsteraexample

import (
	"math/rand/v2"
	"testing"
	"time"

	"github.com/evrblk/monstera"
	"github.com/evrblk/monstera-example/corepb"
	"github.com/stretchr/testify/require"
)

func TestCreateAndGetNamespace(t *testing.T) {
	require := require.New(t)

	namespacesCore := newNamespacesCore()

	now := time.Now()

	// Create namespace
	response1, err := namespacesCore.CreateNamespace(&corepb.CreateNamespaceRequest{
		AccountId:             rand.Uint64(),
		Name:                  "test_namespace",
		Description:           "test description",
		Now:                   now.UnixNano(),
		MaxNumberOfNamespaces: 20,
	})

	require.NoError(err)
	require.NotNil(response1.Namespace)

	// Get this newly created namespace
	response2, err := namespacesCore.GetNamespace(&corepb.GetNamespaceRequest{
		NamespaceId: response1.Namespace.Id,
	})

	require.NoError(err)
	require.NotNil(response2.Namespace)

	require.Equal("test_namespace", response2.Namespace.Id.NamespaceName)
	require.Equal("test description", response2.Namespace.Description)
	require.Equal(now.UnixNano(), response2.Namespace.CreatedAt)
	require.Equal(now.UnixNano(), response2.Namespace.UpdatedAt)

	// Get non-existent namespace
	_, err = namespacesCore.GetNamespace(&corepb.GetNamespaceRequest{
		NamespaceId: &corepb.NamespaceId{
			AccountId:     rand.Uint64(),
			NamespaceName: "random_name",
		},
	})

	require.Error(err)
}

func TestListNamespaces(t *testing.T) {
	require := require.New(t)

	namespacesCore := newNamespacesCore()

	now := time.Now()

	accountId := rand.Uint64()

	// Create namespace 1
	response1, err := namespacesCore.CreateNamespace(&corepb.CreateNamespaceRequest{
		AccountId:             accountId,
		Name:                  "test_namespace_1",
		Now:                   now.UnixNano(),
		MaxNumberOfNamespaces: 20,
	})

	require.NoError(err)
	require.NotNil(response1.Namespace)

	// Create namespace 2
	response2, err := namespacesCore.CreateNamespace(&corepb.CreateNamespaceRequest{
		AccountId:             accountId,
		Name:                  "test_namespace_2",
		Now:                   now.UnixNano(),
		MaxNumberOfNamespaces: 20,
	})

	require.NoError(err)
	require.NotNil(response2.Namespace)

	// List namespaces
	response3, err := namespacesCore.ListNamespaces(&corepb.ListNamespacesRequest{
		AccountId: accountId,
	})

	require.NoError(err)
	require.Len(response3.Namespaces, 2)
}

func TestMaxNumberOfNamespaces(t *testing.T) {
	require := require.New(t)

	namespacesCore := newNamespacesCore()

	now := time.Now()

	accountId := rand.Uint64()

	// Create namespace 1
	response1, err := namespacesCore.CreateNamespace(&corepb.CreateNamespaceRequest{
		AccountId:             accountId,
		Name:                  "test_namespace_1",
		Now:                   now.UnixNano(),
		MaxNumberOfNamespaces: 1,
	})

	require.NoError(err)
	require.NotNil(response1.Namespace)

	// Create namespace 2
	_, err = namespacesCore.CreateNamespace(&corepb.CreateNamespaceRequest{
		AccountId:             accountId,
		Name:                  "test_namespace_2",
		Now:                   now.UnixNano(),
		MaxNumberOfNamespaces: 1,
	})

	require.Error(err)
}

func newNamespacesCore() *NamespacesCore {
	return NewNamespacesCore(monstera.NewBadgerInMemoryStore(), []byte{0x00, 0x00}, []byte{0xff, 0xff})
}
