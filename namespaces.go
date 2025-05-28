package monsteraexample

import (
	"fmt"
	"io"

	"errors"

	"github.com/evrblk/monstera"
	"github.com/evrblk/monstera-example/corepb"
	monsterax "github.com/evrblk/monstera/x"
)

type NamespacesCore struct {
	badgerStore     *monstera.BadgerStore
	namespacesTable *monsterax.CompositeKeyTable[*corepb.Namespace, corepb.Namespace]
}

var _ NamespacesCoreApi = &NamespacesCore{}

func NewNamespacesCore(badgerStore *monstera.BadgerStore, shardLowerBound []byte, shardUpperBound []byte) *NamespacesCore {
	return &NamespacesCore{
		badgerStore:     badgerStore,
		namespacesTable: monsterax.NewCompositeKeyTable[*corepb.Namespace, corepb.Namespace](namespacesTableId, shardLowerBound, shardUpperBound),
	}
}

func (c *NamespacesCore) ranges() []monstera.KeyRange {
	return []monstera.KeyRange{
		c.namespacesTable.GetTableKeyRange(),
	}
}

func (c *NamespacesCore) Snapshot() monstera.ApplicationCoreSnapshot {
	return monsterax.Snapshot(c.badgerStore, c.ranges())
}

func (c *NamespacesCore) Restore(reader io.ReadCloser) error {
	return monsterax.Restore(c.badgerStore, c.ranges(), reader)
}

func (c *NamespacesCore) Close() {

}

func (c *NamespacesCore) CreateNamespace(request *corepb.CreateNamespaceRequest) (*corepb.CreateNamespaceResponse, error) {
	txn := c.badgerStore.Update()
	defer txn.Discard()

	// Validations
	if request.Name == "" {
		return nil, monsterax.NewErrorWithContext(
			monsterax.InvalidArgument,
			"Name should not be empty",
			map[string]string{})
	}

	namespaceId := &corepb.NamespaceId{
		AccountId:     request.AccountId,
		NamespaceName: request.Name,
	}

	// Checking name uniqueness
	_, err := c.getNamespace(txn, namespaceId)
	if err != nil {
		if !errors.Is(err, monstera.ErrNotFound) {
			return nil, err
		}
	} else {
		return nil, monsterax.NewErrorWithContext(
			monsterax.AlreadyExists,
			"namespace with this name already exists",
			map[string]string{"namespace_name": request.Name})
	}

	namespaces, err := c.listNamespaces(txn, request.AccountId)
	panicIfNotNil(err)

	// Checking max number of namespaces
	if int64(len(namespaces)) >= request.MaxNumberOfNamespaces {
		return nil, monsterax.NewErrorWithContext(
			monsterax.ResourceExhausted,
			"max number of namespaces reached",
			map[string]string{"limit": fmt.Sprintf("%d", request.MaxNumberOfNamespaces)})
	}

	namespace := &corepb.Namespace{
		Id:          namespaceId,
		Description: request.Description,
		CreatedAt:   request.Now,
		UpdatedAt:   request.Now,
	}

	err = c.createNamespace(txn, namespace)
	panicIfNotNil(err)

	err = txn.Commit()
	panicIfNotNil(err)

	return &corepb.CreateNamespaceResponse{
		Namespace: namespace,
	}, nil
}

func (c *NamespacesCore) UpdateNamespace(request *corepb.UpdateNamespaceRequest) (*corepb.UpdateNamespaceResponse, error) {
	txn := c.badgerStore.Update()
	defer txn.Discard()

	namespace, err := c.getNamespace(txn, request.NamespaceId)
	if err != nil {
		if errors.Is(err, monstera.ErrNotFound) {
			return nil, monsterax.NewErrorWithContext(
				monsterax.NotFound,
				"namespace not found",
				map[string]string{"namespace_name": request.NamespaceId.NamespaceName})
		} else {
			panic(err)
		}
	}

	namespace.Description = request.Description
	namespace.UpdatedAt = request.Now

	err = c.updateNamespace(txn, namespace)
	panicIfNotNil(err)

	err = txn.Commit()
	panicIfNotNil(err)

	return &corepb.UpdateNamespaceResponse{
		Namespace: namespace,
	}, nil
}

func (c *NamespacesCore) DeleteNamespace(request *corepb.DeleteNamespaceRequest) (*corepb.DeleteNamespaceResponse, error) {
	txn := c.badgerStore.Update()
	defer txn.Discard()

	namespace, err := c.getNamespace(txn, request.NamespaceId)
	if err != nil {
		if errors.Is(err, monstera.ErrNotFound) {
			return nil, monsterax.NewErrorWithContext(
				monsterax.NotFound,
				"namespace not found",
				map[string]string{"namespace_name": request.NamespaceId.NamespaceName})
		} else {
			panic(err)
		}
	}

	err = c.deleteNamespace(txn, namespace)
	panicIfNotNil(err)

	// TODO delete locks, semaphores, wgs

	err = txn.Commit()
	panicIfNotNil(err)

	return &corepb.DeleteNamespaceResponse{}, nil
}

func (c *NamespacesCore) GetNamespace(request *corepb.GetNamespaceRequest) (*corepb.GetNamespaceResponse, error) {
	txn := c.badgerStore.View()
	defer txn.Discard()

	namespace, err := c.getNamespace(txn, request.NamespaceId)
	if err != nil {
		if errors.Is(err, monstera.ErrNotFound) {
			return nil, monsterax.NewErrorWithContext(
				monsterax.NotFound,
				"namespace not found",
				map[string]string{"namespace_name": request.NamespaceId.NamespaceName})
		} else {
			panic(err)
		}
	}

	return &corepb.GetNamespaceResponse{
		Namespace: namespace,
	}, nil
}

func (c *NamespacesCore) ListNamespaces(request *corepb.ListNamespacesRequest) (*corepb.ListNamespacesResponse, error) {
	txn := c.badgerStore.View()
	defer txn.Discard()

	namespaces, err := c.listNamespaces(txn, request.AccountId)
	panicIfNotNil(err)

	return &corepb.ListNamespacesResponse{
		Namespaces: namespaces,
	}, nil
}

func (c *NamespacesCore) getNamespace(txn *monstera.Txn, namespaceId *corepb.NamespaceId) (*corepb.Namespace, error) {
	return c.namespacesTable.Get(txn, namespacesTablePK(namespaceId.AccountId), namespacesTableSK(namespaceId))
}

func (c *NamespacesCore) listNamespaces(txn *monstera.Txn, accountId uint64) ([]*corepb.Namespace, error) {
	result := make([]*corepb.Namespace, 0)

	err := c.namespacesTable.List(txn, namespacesTablePK(accountId), func(namespace *corepb.Namespace) (bool, error) {
		result = append(result, namespace)
		return true, nil
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *NamespacesCore) createNamespace(txn *monstera.Txn, namespace *corepb.Namespace) error {
	return c.namespacesTable.Set(txn, namespacesTablePK(namespace.Id.AccountId), namespacesTableSK(namespace.Id), namespace)
}

func (c *NamespacesCore) deleteNamespace(txn *monstera.Txn, namespace *corepb.Namespace) error {
	// Remove namespace from main namespacesTable
	return c.namespacesTable.Delete(txn, namespacesTablePK(namespace.Id.AccountId), namespacesTableSK(namespace.Id))
}

func (c *NamespacesCore) updateNamespace(txn *monstera.Txn, namespace *corepb.Namespace) error {
	return c.namespacesTable.Set(txn, namespacesTablePK(namespace.Id.AccountId), namespacesTableSK(namespace.Id), namespace)
}

type namespaceIdIntf interface {
	GetAccountId() uint64
	GetNamespaceName() string
}

// 1. shard key (by account id)
// 2. account id
func namespacesTablePK(accountId uint64) []byte {
	return monstera.ConcatBytes(shardByAccount(accountId), accountId)
}

// 1. namespace name
func namespacesTableSK(n namespaceIdIntf) []byte {
	return monstera.ConcatBytes(n.GetNamespaceName())
}

func panicIfNotNil(err error) {
	if err != nil {
		panic(err)
	}
}
