package monsteraexample

import (
	"context"
	"log"
	"time"

	"github.com/evrblk/monstera-example/corepb"
	"github.com/evrblk/monstera-example/gatewaypb"
	monsterax "github.com/evrblk/monstera/x"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ExampleServiceApiServer struct {
	gatewaypb.UnimplementedExampleServiceApiServer

	coreApiClient ExampleServiceCoreApi
}

func (s *ExampleServiceApiServer) Close() {
	log.Println("Stopping ApiServer...")
}

func (s *ExampleServiceApiServer) CreateNamespace(ctx context.Context, request *gatewaypb.CreateNamespaceRequest) (*gatewaypb.CreateNamespaceResponse, error) {
	now := time.Now()

	// Validation
	if err := validateCreateNamespaceRequest(request); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%s", err)
	}

	res, err := s.coreApiClient.CreateNamespace(ctx, &corepb.CreateNamespaceRequest{
		Name:        request.Name,
		Description: request.Description,
		Now:         now.UnixNano(),

		MaxNumberOfNamespaces: 25, // account.MaxNumberOfNamespaces,
	})
	if err != nil {
		return nil, monsterax.ErrorToGRPC(err)
	}

	return &gatewaypb.CreateNamespaceResponse{
		Namespace: namespaceToFront(res.Namespace),
	}, nil
}

func (s *ExampleServiceApiServer) GetNamespace(ctx context.Context, request *gatewaypb.GetNamespaceRequest) (*gatewaypb.GetNamespaceResponse, error) {
	// Validation
	if err := validateGetNamespaceRequest(request); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%s", err)
	}

	namespace, err := s.coreApiClient.GetNamespace(ctx, &corepb.GetNamespaceRequest{
		NamespaceId: &corepb.NamespaceId{
			NamespaceName: request.NamespaceName,
		},
	})
	if err != nil {
		return nil, monsterax.ErrorToGRPC(err)
	}

	return &gatewaypb.GetNamespaceResponse{
		Namespace: namespaceToFront(namespace.Namespace),
	}, nil
}

func (s *ExampleServiceApiServer) UpdateNamespace(ctx context.Context, request *gatewaypb.UpdateNamespaceRequest) (*gatewaypb.UpdateNamespaceResponse, error) {
	now := time.Now()

	// Validation
	if err := validateUpdateNamespaceRequest(request); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%s", err)
	}

	res, err := s.coreApiClient.UpdateNamespace(ctx, &corepb.UpdateNamespaceRequest{
		NamespaceId: &corepb.NamespaceId{
			NamespaceName: request.NamespaceName,
		},
		Description: request.Description,
		Now:         now.UnixNano(),
	})
	if err != nil {
		return nil, monsterax.ErrorToGRPC(err)
	}

	return &gatewaypb.UpdateNamespaceResponse{
		Namespace: namespaceToFront(res.Namespace),
	}, nil
}

func (s *ExampleServiceApiServer) DeleteNamespace(ctx context.Context, request *gatewaypb.DeleteNamespaceRequest) (*gatewaypb.DeleteNamespaceResponse, error) {
	// Validation
	if err := validateDeleteNamespaceRequest(request); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%s", err)
	}

	_, err := s.coreApiClient.DeleteNamespace(ctx, &corepb.DeleteNamespaceRequest{
		NamespaceId: &corepb.NamespaceId{
			NamespaceName: request.NamespaceName,
		},
	})
	if err != nil {
		return nil, monsterax.ErrorToGRPC(err)
	}

	return &gatewaypb.DeleteNamespaceResponse{}, nil
}

func (s *ExampleServiceApiServer) ListNamespaces(ctx context.Context, request *gatewaypb.ListNamespacesRequest) (*gatewaypb.ListNamespacesResponse, error) {
	// Validation
	if err := validateListNamespacesRequest(request); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%s", err)
	}

	res, err := s.coreApiClient.ListNamespaces(ctx, &corepb.ListNamespacesRequest{})
	if err != nil {
		return nil, monsterax.ErrorToGRPC(err)
	}

	return &gatewaypb.ListNamespacesResponse{
		Namespaces: namespacesToFront(res.Namespaces),
	}, nil
}

func (s *ExampleServiceApiServer) AcquireLock(ctx context.Context, request *gatewaypb.AcquireLockRequest) (*gatewaypb.AcquireLockResponse, error) {
	now := time.Now()

	// Validation
	if err := validateAcquireLockRequest(request); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%s", err)
	}

	res, err := s.coreApiClient.AcquireLock(ctx, &corepb.AcquireLockRequest{
		LockId: &corepb.LockId{
			NamespaceName: request.NamespaceName,
			LockName:      request.LockName,
		},
		Now:       now.UnixNano(),
		ProcessId: request.ProcessId,
	})
	if err != nil {
		return nil, monsterax.ErrorToGRPC(err)
	}

	return &gatewaypb.AcquireLockResponse{
		Lock:    lockToFront(res.Lock),
		Success: res.Success,
	}, nil
}

func (s *ExampleServiceApiServer) ReleaseLock(ctx context.Context, request *gatewaypb.ReleaseLockRequest) (*gatewaypb.ReleaseLockResponse, error) {
	now := time.Now()

	// Validation
	if err := validateReleaseLockRequest(request); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%s", err)
	}

	res, err := s.coreApiClient.ReleaseLock(ctx, &corepb.ReleaseLockRequest{
		LockId: &corepb.LockId{
			NamespaceName: request.NamespaceName,
			LockName:      request.LockName,
		},
		Now:       now.UnixNano(),
		ProcessId: request.ProcessId,
	})
	if err != nil {
		return nil, monsterax.ErrorToGRPC(err)
	}

	return &gatewaypb.ReleaseLockResponse{
		Lock: lockToFront(res.Lock),
	}, nil
}

func (s *ExampleServiceApiServer) GetLock(ctx context.Context, request *gatewaypb.GetLockRequest) (*gatewaypb.GetLockResponse, error) {
	now := time.Now()

	// Validation
	if err := validateGetLockRequest(request); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%s", err)
	}

	res, err := s.coreApiClient.GetLock(ctx, &corepb.GetLockRequest{
		LockId: &corepb.LockId{
			NamespaceName: request.NamespaceName,
			LockName:      request.LockName,
		},
		Now: now.UnixNano(),
	})
	if err != nil {
		return nil, monsterax.ErrorToGRPC(err)
	}

	return &gatewaypb.GetLockResponse{
		Lock: lockToFront(res.Lock),
	}, nil
}

func (s *ExampleServiceApiServer) DeleteLock(ctx context.Context, request *gatewaypb.DeleteLockRequest) (*gatewaypb.DeleteLockResponse, error) {
	now := time.Now()

	// Validation
	if err := validateDeleteLockRequest(request); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%s", err)
	}

	_, err := s.coreApiClient.DeleteLock(ctx, &corepb.DeleteLockRequest{
		LockId: &corepb.LockId{
			NamespaceName: request.NamespaceName,
			LockName:      request.LockName,
		},
		Now: now.UnixNano(),
	})
	if err != nil {
		return nil, monsterax.ErrorToGRPC(err)
	}

	return &gatewaypb.DeleteLockResponse{}, nil
}

func NewExampleServiceApiServer(coreApiClient ExampleServiceCoreApi) *ExampleServiceApiServer {
	return &ExampleServiceApiServer{
		coreApiClient: coreApiClient,
	}
}
