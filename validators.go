package monsteraexample

import (
	"fmt"

	"github.com/evrblk/monstera-example/gatewaypb"
)

const (
	maxNamespaceNameLength = 128
	maxLockNameLength      = 128
	maxSemaphoreNameLength = 128
	maxProcessIdLength     = 128
	maxDescriptionLength   = 1024
)

func validateCreateNamespaceRequest(request *gatewaypb.CreateNamespaceRequest) error {
	if request.Name == "" {
		return fmt.Errorf("invalid CreateNamespaceRequest.Name")
	}
	if len(request.Name) > maxNamespaceNameLength {
		return fmt.Errorf("invalid CreateNamespaceRequest.Name")
	}

	if len(request.Description) > maxDescriptionLength {
		return fmt.Errorf("invalid CreateNamespaceRequest.Description")
	}

	return nil
}

func validateGetNamespaceRequest(request *gatewaypb.GetNamespaceRequest) error {
	if request.NamespaceName == "" {
		return fmt.Errorf("invalid GetNamespaceRequest.NamespaceName")
	}
	if len(request.NamespaceName) > maxNamespaceNameLength {
		return fmt.Errorf("invalid GetNamespaceRequest.NamespaceName")
	}

	return nil
}

func validateUpdateNamespaceRequest(request *gatewaypb.UpdateNamespaceRequest) error {
	if request.NamespaceName == "" {
		return fmt.Errorf("invalid UpdateNamespaceRequest.NamespaceName")
	}
	if len(request.NamespaceName) > maxNamespaceNameLength {
		return fmt.Errorf("invalid UpdateNamespaceRequest.NamespaceName")
	}

	if len(request.Description) > maxDescriptionLength {
		return fmt.Errorf("invalid UpdateNamespaceRequest.Description")
	}

	return nil
}

func validateDeleteNamespaceRequest(request *gatewaypb.DeleteNamespaceRequest) error {
	if request.NamespaceName == "" {
		return fmt.Errorf("invalid DeleteNamespaceRequest.NamespaceName")
	}
	if len(request.NamespaceName) > maxNamespaceNameLength {
		return fmt.Errorf("invalid DeleteNamespaceRequest.NamespaceName")
	}

	return nil
}

func validateListNamespacesRequest(request *gatewaypb.ListNamespacesRequest) error {
	return nil
}

func validateDeleteLockRequest(request *gatewaypb.DeleteLockRequest) error {
	if request.NamespaceName == "" {
		return fmt.Errorf("invalid DeleteLockRequest.NamespaceName")
	}
	if len(request.NamespaceName) > maxNamespaceNameLength {
		return fmt.Errorf("invalid DeleteLockRequest.NamespaceName")
	}

	if request.LockName == "" {
		return fmt.Errorf("invalid DeleteLockRequest.LockName")
	}
	if len(request.LockName) > maxLockNameLength {
		return fmt.Errorf("invalid DeleteLockRequest.LockName")
	}

	return nil
}

func validateGetLockRequest(request *gatewaypb.GetLockRequest) error {
	if request.NamespaceName == "" {
		return fmt.Errorf("invalid GetLockRequest.NamespaceName")
	}
	if len(request.NamespaceName) > maxNamespaceNameLength {
		return fmt.Errorf("invalid GetLockRequest.NamespaceName")
	}

	if request.LockName == "" {
		return fmt.Errorf("invalid GetLockRequest.LockName")
	}
	if len(request.LockName) > maxLockNameLength {
		return fmt.Errorf("invalid GetLockRequest.LockName")
	}

	return nil
}

func validateReleaseLockRequest(request *gatewaypb.ReleaseLockRequest) error {
	if request.NamespaceName == "" {
		return fmt.Errorf("invalid ReleaseLockRequest.NamespaceName")
	}
	if len(request.NamespaceName) > maxNamespaceNameLength {
		return fmt.Errorf("invalid ReleaseLockRequest.NamespaceName")
	}

	if request.LockName == "" {
		return fmt.Errorf("invalid ReleaseLockRequest.LockName")
	}
	if len(request.LockName) > maxLockNameLength {
		return fmt.Errorf("invalid ReleaseLockRequest.LockName")
	}

	if request.ProcessId == "" {
		return fmt.Errorf("invalid ReleaseLockRequest.ProcessId")
	}
	if len(request.ProcessId) > maxProcessIdLength {
		return fmt.Errorf("invalid ReleaseLockRequest.ProcessId")
	}

	return nil
}

func validateAcquireLockRequest(request *gatewaypb.AcquireLockRequest) error {
	if request.NamespaceName == "" {
		return fmt.Errorf("invalid AcquireLockRequest.NamespaceName")
	}
	if len(request.NamespaceName) > maxNamespaceNameLength {
		return fmt.Errorf("invalid AcquireLockRequest.NamespaceName")
	}

	if request.LockName == "" {
		return fmt.Errorf("invalid AcquireLockRequest.LockName")
	}
	if len(request.LockName) > maxLockNameLength {
		return fmt.Errorf("invalid AcquireLockRequest.LockName")
	}

	if request.ProcessId == "" {
		return fmt.Errorf("invalid AcquireLockRequest.ProcessId")
	}
	if len(request.ProcessId) > maxProcessIdLength {
		return fmt.Errorf("invalid AcquireLockRequest.ProcessId")
	}

	return nil
}

func validateCreateSemaphoreRequest(request *gatewaypb.CreateSemaphoreRequest) error {
	if request.NamespaceName == "" {
		return fmt.Errorf("invalid CreateSemaphoreRequest.NamespaceName")
	}
	if len(request.NamespaceName) > maxNamespaceNameLength {
		return fmt.Errorf("invalid CreateSemaphoreRequest.NamespaceName")
	}

	if request.SemaphoreName == "" {
		return fmt.Errorf("invalid CreateSemaphoreRequest.SemaphoreName")
	}
	if len(request.SemaphoreName) > maxSemaphoreNameLength {
		return fmt.Errorf("invalid CreateSemaphoreRequest.SemaphoreName")
	}

	return nil
}

func validateGetSemaphoreRequest(request *gatewaypb.GetSemaphoreRequest) error {
	if request.NamespaceName == "" {
		return fmt.Errorf("invalid GetSemaphoreRequest.NamespaceName")
	}
	if len(request.NamespaceName) > maxNamespaceNameLength {
		return fmt.Errorf("invalid GetSemaphoreRequest.NamespaceName")
	}

	if request.SemaphoreName == "" {
		return fmt.Errorf("invalid GetSemaphoreRequest.SemaphoreName")
	}
	if len(request.SemaphoreName) > maxSemaphoreNameLength {
		return fmt.Errorf("invalid GetSemaphoreRequest.SemaphoreName")
	}

	return nil
}

func validateReleaseSemaphoreRequest(request *gatewaypb.ReleaseSemaphoreRequest) error {
	if request.NamespaceName == "" {
		return fmt.Errorf("invalid ReleaseSemaphoreRequest.NamespaceName")
	}
	if len(request.NamespaceName) > maxNamespaceNameLength {
		return fmt.Errorf("invalid ReleaseSemaphoreRequest.NamespaceName")
	}

	if request.SemaphoreName == "" {
		return fmt.Errorf("invalid ReleaseSemaphoreRequest.SemaphoreName")
	}
	if len(request.SemaphoreName) > maxSemaphoreNameLength {
		return fmt.Errorf("invalid ReleaseSemaphoreRequest.SemaphoreName")
	}

	if request.ProcessId == "" {
		return fmt.Errorf("invalid ReleaseSemaphoreRequest.ProcessId")
	}
	if len(request.ProcessId) > maxProcessIdLength {
		return fmt.Errorf("invalid ReleaseSemaphoreRequest.ProcessId")
	}

	return nil
}

func validateUpdateSemaphoreRequest(request *gatewaypb.UpdateSemaphoreRequest) error {
	if request.NamespaceName == "" {
		return fmt.Errorf("invalid UpdateSemaphoreRequest.NamespaceName")
	}
	if len(request.NamespaceName) > maxNamespaceNameLength {
		return fmt.Errorf("invalid UpdateSemaphoreRequest.NamespaceName")
	}

	if request.SemaphoreName == "" {
		return fmt.Errorf("invalid UpdateSemaphoreRequest.SemaphoreName")
	}
	if len(request.SemaphoreName) > maxSemaphoreNameLength {
		return fmt.Errorf("invalid UpdateSemaphoreRequest.SemaphoreName")
	}

	return nil
}

func validateDeleteSemaphoreRequest(request *gatewaypb.DeleteSemaphoreRequest) error {
	if request.NamespaceName == "" {
		return fmt.Errorf("invalid DeleteSemaphoreRequest.NamespaceName")
	}
	if len(request.NamespaceName) > maxNamespaceNameLength {
		return fmt.Errorf("invalid DeleteSemaphoreRequest.NamespaceName")
	}

	if request.SemaphoreName == "" {
		return fmt.Errorf("invalid DeleteSemaphoreRequest.SemaphoreName")
	}
	if len(request.SemaphoreName) > maxSemaphoreNameLength {
		return fmt.Errorf("invalid DeleteSemaphoreRequest.SemaphoreName")
	}

	return nil
}

func validateListSemaphoresRequest(request *gatewaypb.ListSemaphoresRequest) error {
	if request.NamespaceName == "" {
		return fmt.Errorf("invalid ListSemaphoresRequest.NamespaceName")
	}
	if len(request.NamespaceName) > maxNamespaceNameLength {
		return fmt.Errorf("invalid ListSemaphoresRequest.NamespaceName")
	}

	return nil
}

func validateAcquireSemaphoreRequest(request *gatewaypb.AcquireSemaphoreRequest) error {
	if request.NamespaceName == "" {
		return fmt.Errorf("invalid AcquireSemaphoreRequest.NamespaceName")
	}
	if len(request.NamespaceName) > maxNamespaceNameLength {
		return fmt.Errorf("invalid AcquireSemaphoreRequest.NamespaceName")
	}

	if request.SemaphoreName == "" {
		return fmt.Errorf("invalid AcquireSemaphoreRequest.SemaphoreName")
	}
	if len(request.SemaphoreName) > maxSemaphoreNameLength {
		return fmt.Errorf("invalid AcquireSemaphoreRequest.SemaphoreName")
	}

	if request.ProcessId == "" {
		return fmt.Errorf("invalid AcquireSemaphoreRequest.ProcessId")
	}
	if len(request.ProcessId) > maxProcessIdLength {
		return fmt.Errorf("invalid AcquireSemaphoreRequest.ProcessId")
	}

	return nil
}
