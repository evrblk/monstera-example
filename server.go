package monsteraexample

import (
	"context"
	"log"
	"time"

	"github.com/evrblk/monstera-example/corepb"
	"github.com/evrblk/monstera-example/gatewaypb"
	monsterax "github.com/evrblk/monstera/x"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
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
	accountId := ctx.Value("account-id").(uint64)

	res1, err := s.coreApiClient.GetAccount(ctx, &corepb.GetAccountRequest{
		AccountId: accountId,
	})
	if err != nil {
		return nil, monsterax.ErrorToGRPC(err)
	}
	account := res1.Account

	now := time.Now()

	// Validation
	if err := validateCreateNamespaceRequest(request); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%s", err)
	}

	res2, err := s.coreApiClient.CreateNamespace(ctx, &corepb.CreateNamespaceRequest{
		AccountId:             accountId,
		Name:                  request.Name,
		Description:           request.Description,
		Now:                   now.UnixNano(),
		MaxNumberOfNamespaces: account.MaxNumberOfNamespaces,
	})
	if err != nil {
		return nil, monsterax.ErrorToGRPC(err)
	}

	return &gatewaypb.CreateNamespaceResponse{
		Namespace: namespaceToFront(res2.Namespace),
	}, nil
}

func (s *ExampleServiceApiServer) GetNamespace(ctx context.Context, request *gatewaypb.GetNamespaceRequest) (*gatewaypb.GetNamespaceResponse, error) {
	accountId := ctx.Value("account-id").(uint64)

	// Validation
	if err := validateGetNamespaceRequest(request); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%s", err)
	}

	namespace, err := s.coreApiClient.GetNamespace(ctx, &corepb.GetNamespaceRequest{
		NamespaceId: &corepb.NamespaceId{
			AccountId:     accountId,
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
	accountId := ctx.Value("account-id").(uint64)

	now := time.Now()

	// Validation
	if err := validateUpdateNamespaceRequest(request); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%s", err)
	}

	res, err := s.coreApiClient.UpdateNamespace(ctx, &corepb.UpdateNamespaceRequest{
		NamespaceId: &corepb.NamespaceId{
			AccountId:     accountId,
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
	accountId := ctx.Value("account-id").(uint64)

	// Validation
	if err := validateDeleteNamespaceRequest(request); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%s", err)
	}

	_, err := s.coreApiClient.DeleteNamespace(ctx, &corepb.DeleteNamespaceRequest{
		NamespaceId: &corepb.NamespaceId{
			AccountId:     accountId,
			NamespaceName: request.NamespaceName,
		},
	})
	if err != nil {
		return nil, monsterax.ErrorToGRPC(err)
	}

	return &gatewaypb.DeleteNamespaceResponse{}, nil
}

func (s *ExampleServiceApiServer) ListNamespaces(ctx context.Context, request *gatewaypb.ListNamespacesRequest) (*gatewaypb.ListNamespacesResponse, error) {
	accountId := ctx.Value("account-id").(uint64)

	// Validation
	if err := validateListNamespacesRequest(request); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%s", err)
	}

	res, err := s.coreApiClient.ListNamespaces(ctx, &corepb.ListNamespacesRequest{
		AccountId: accountId,
	})
	if err != nil {
		return nil, monsterax.ErrorToGRPC(err)
	}

	return &gatewaypb.ListNamespacesResponse{
		Namespaces: namespacesToFront(res.Namespaces),
	}, nil
}

func (s *ExampleServiceApiServer) AcquireLock(ctx context.Context, request *gatewaypb.AcquireLockRequest) (*gatewaypb.AcquireLockResponse, error) {
	accountId := ctx.Value("account-id").(uint64)

	now := time.Now()

	// Validation
	if err := validateAcquireLockRequest(request); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%s", err)
	}

	res, err := s.coreApiClient.AcquireLock(ctx, &corepb.AcquireLockRequest{
		LockId: &corepb.LockId{
			AccountId:     accountId,
			NamespaceName: request.NamespaceName,
			LockName:      request.LockName,
		},
		Now:       now.UnixNano(),
		ProcessId: request.ProcessId,
		ExpiresAt: request.ExpiresAt,
		WriteLock: request.WriteLock,
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
	accountId := ctx.Value("account-id").(uint64)

	now := time.Now()

	// Validation
	if err := validateReleaseLockRequest(request); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%s", err)
	}

	res, err := s.coreApiClient.ReleaseLock(ctx, &corepb.ReleaseLockRequest{
		LockId: &corepb.LockId{
			AccountId:     accountId,
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
	accountId := ctx.Value("account-id").(uint64)

	now := time.Now()

	// Validation
	if err := validateGetLockRequest(request); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%s", err)
	}

	res, err := s.coreApiClient.GetLock(ctx, &corepb.GetLockRequest{
		LockId: &corepb.LockId{
			AccountId:     accountId,
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
	accountId := ctx.Value("account-id").(uint64)

	now := time.Now()

	// Validation
	if err := validateDeleteLockRequest(request); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%s", err)
	}

	_, err := s.coreApiClient.DeleteLock(ctx, &corepb.DeleteLockRequest{
		LockId: &corepb.LockId{
			AccountId:     accountId,
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

type AuthenticationMiddleware struct {
}

func (m *AuthenticationMiddleware) Unary(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "no metadata")
	}

	if len(md.Get("account-id")) != 1 {
		return nil, status.Errorf(codes.Unauthenticated, "no account-id in metadata")
	}
	accountIdStr := md.Get("account-id")[0]
	accountId, err := DecodeAccountId(accountIdStr)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid account-id: %s", err)
	}

	ctx = context.WithValue(ctx, "account-id", accountId)

	return handler(ctx, req)
}
