package service

import (
	"context"
	"strconv"

	"auth/ent"
	"auth/internal/repository"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "auth/pb"
)

// IAuthzService defines the interface for authorization and role-based access control (RBAC) operations.
// It provides methods for managing permissions, roles, and user-role associations.
type IAuthzService interface {
	// CheckPermission verifies if a user has a specific permission.
	//
	// Parameters:
	//   - ctx: A context.Context for handling deadlines, cancellations, and request-scoped values.
	//   - req: A pointer to pb.CheckPermissionRequest containing the user ID and the permission to check.
	//
	// Returns:
	//   - *pb.CheckPermissionResponse: A response indicating whether the user has the specified permission.
	//   - error: An error if the check fails due to invalid input or internal server issues.
	//
	// This method should efficiently check the user's roles and the permissions associated with those roles.
	// It may implement caching mechanisms to improve performance for frequent permission checks.
	CheckPermission(ctx context.Context, req *pb.CheckPermissionRequest) (*pb.CheckPermissionResponse, error)

	// GetUserRoles retrieves all roles assigned to a specific user.
	//
	// Parameters:
	//   - ctx: A context.Context for handling deadlines, cancellations, and request-scoped values.
	//   - req: A pointer to pb.GetUserRolesRequest containing the user ID.
	//
	// Returns:
	//   - *pb.GetUserRolesResponse: A response containing a list of roles assigned to the user.
	//   - error: An error if the retrieval fails, such as user not found or internal server issues.
	//
	// This method should handle cases where a user might have multiple roles and ensure all are returned.
	GetUserRoles(ctx context.Context, req *pb.GetUserRolesRequest) (*pb.GetUserRolesResponse, error)

	// AssignRoleToUser assigns a specific role to a user.
	//
	// Parameters:
	//   - ctx: A context.Context for handling deadlines, cancellations, and request-scoped values.
	//   - req: A pointer to pb.AssignRoleToUserRequest containing the user ID and role ID to be assigned.
	//
	// Returns:
	//   - *emptypb.Empty: An empty response indicating successful assignment.
	//   - error: An error if the assignment fails, such as invalid user/role ID or internal server issues.
	//
	// This method should check if the role already exists for the user to avoid duplicate assignments.
	// It may also trigger any necessary cache invalidations or notifications.
	AssignRoleToUser(ctx context.Context, req *pb.AssignRoleToUserRequest) (*emptypb.Empty, error)

	// RemoveRoleFromUser removes a specific role from a user.
	//
	// Parameters:
	//   - ctx: A context.Context for handling deadlines, cancellations, and request-scoped values.
	//   - req: A pointer to pb.RemoveRoleFromUserRequest containing the user ID and role ID to be removed.
	//
	// Returns:
	//   - *emptypb.Empty: An empty response indicating successful removal.
	//   - error: An error if the removal fails, such as role not assigned to the user or internal server issues.
	//
	// This method should handle cases where the role might not be assigned to the user gracefully.
	// It may also trigger any necessary cache invalidations or notifications.
	RemoveRoleFromUser(ctx context.Context, req *pb.RemoveRoleFromUserRequest) (*emptypb.Empty, error)

	// CreateRole creates a new role in the system.
	//
	// Parameters:
	//   - ctx: A context.Context for handling deadlines, cancellations, and request-scoped values.
	//   - req: A pointer to pb.CreateRoleRequest containing the details of the new role.
	//
	// Returns:
	//   - *pb.CreateRoleResponse: A response containing the ID of the newly created role.
	//   - error: An error if creation fails, such as duplicate role name or internal server issues.
	//
	// This method should validate the role details and ensure uniqueness of the role name.
	// It may also set up any default permissions associated with the new role.
	CreateRole(ctx context.Context, req *pb.CreateRoleRequest) (*pb.CreateRoleResponse, error)

	// UpdateRole modifies an existing role in the system.
	//
	// Parameters:
	//   - ctx: A context.Context for handling deadlines, cancellations, and request-scoped values.
	//   - req: A pointer to pb.UpdateRoleRequest containing the role ID and updated details.
	//
	// Returns:
	//   - *pb.UpdateRoleResponse: A response confirming the update and containing the updated role information.
	//   - error: An error if the update fails, such as role not found or internal server issues.
	//
	// This method should validate the updated role details and handle any conflicts with existing roles.
	// It may also trigger updates to user permissions if the role's permissions have changed.
	UpdateRole(ctx context.Context, req *pb.UpdateRoleRequest) (*pb.UpdateRoleResponse, error)

	// DeleteRole removes a role from the system.
	//
	// Parameters:
	//   - ctx: A context.Context for handling deadlines, cancellations, and request-scoped values.
	//   - req: A pointer to pb.DeleteRoleRequest containing the ID of the role to be deleted.
	//
	// Returns:
	//   - *emptypb.Empty: An empty response indicating successful deletion.
	//   - error: An error if the deletion fails, such as role not found or internal server issues.
	//
	// This method should handle the removal of the role from all users it was assigned to.
	// It should also consider the implications of deleting a role and may implement safeguards against deleting critical roles.
	DeleteRole(ctx context.Context, req *pb.DeleteRoleRequest) (*emptypb.Empty, error)

	// GetRole retrieves detailed information about a specific role.
	//
	// Parameters:
	//   - ctx: A context.Context for handling deadlines, cancellations, and request-scoped values.
	//   - req: A pointer to pb.GetRoleRequest containing the ID of the role to retrieve.
	//
	// Returns:
	//   - *pb.GetRoleResponse: A response containing detailed information about the requested role.
	//   - error: An error if the retrieval fails, such as role not found or internal server issues.
	//
	// This method should return comprehensive information about the role, including its permissions and any metadata.
	GetRole(ctx context.Context, req *pb.GetRoleRequest) (*pb.GetRoleResponse, error)

	// ListRoles retrieves a list of roles based on specified criteria.
	//
	// Parameters:
	//   - ctx: A context.Context for handling deadlines, cancellations, and request-scoped values.
	//   - req: A pointer to pb.ListRolesRequest containing any filtering or pagination parameters.
	//
	// Returns:
	//   - *pb.ListRolesResponse: A response containing a list of roles matching the specified criteria.
	//   - error: An error if the listing fails due to invalid parameters or internal server issues.
	//
	// This method should support pagination and filtering to handle large numbers of roles efficiently.
	// It may also implement sorting options for the returned list.
	ListRoles(ctx context.Context, req *pb.ListRolesRequest) (*pb.ListRolesResponse, error)

	// AddPermissionToRole adds a specific permission to a role.
	//
	// Parameters:
	//   - ctx: A context.Context for handling deadlines, cancellations, and request-scoped values.
	//   - req: A pointer to pb.AddPermissionToRoleRequest containing the role ID and the permission to be added.
	//
	// Returns:
	//   - *emptypb.Empty: An empty response indicating successful addition of the permission.
	//   - error: An error if the addition fails, such as invalid role/permission or internal server issues.
	//
	// This method should check if the permission already exists for the role to avoid duplicates.
	// It may also trigger updates to user permissions for all users with this role.
	AddPermissionToRole(ctx context.Context, req *pb.AddPermissionToRoleRequest) (*emptypb.Empty, error)

	// RemovePermissionFromRole removes a specific permission from a role.
	//
	// Parameters:
	//   - ctx: A context.Context for handling deadlines, cancellations, and request-scoped values.
	//   - req: A pointer to pb.RemovePermissionFromRoleRequest containing the role ID and the permission to be removed.
	//
	// Returns:
	//   - *emptypb.Empty: An empty response indicating successful removal of the permission.
	//   - error: An error if the removal fails, such as permission not found in role or internal server issues.
	//
	// This method should handle cases where the permission might not be assigned to the role gracefully.
	// It may also trigger updates to user permissions for all users with this role.
	RemovePermissionFromRole(ctx context.Context, req *pb.RemovePermissionFromRoleRequest) (*emptypb.Empty, error)
}

type AuthzService struct {
	pb.UnimplementedAuthorizationServiceServer
	roleRepo repository.IRoleRepository
	userRepo repository.IUserRepository
}

func NewAuthzService(roleRepo repository.IRoleRepository, userRepo repository.IUserRepository) *AuthzService {
	return &AuthzService{
		roleRepo: roleRepo,
		userRepo: userRepo,
	}
}

func (s *AuthzService) CheckPermission(ctx context.Context, req *pb.CheckPermissionRequest) (*pb.CheckPermissionResponse, error) {
	roles, err := s.roleRepo.GetRolesByUserID(ctx, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to fetch roles: %v", err)
	}

	for _, role := range roles {
		if hasPermission(role, req.Permission.Resource, req.Permission.Action) {
			return &pb.CheckPermissionResponse{HasPermission: true}, nil
		}
	}

	return &pb.CheckPermissionResponse{HasPermission: false}, nil
}

func (s *AuthzService) GetUserRoles(ctx context.Context, req *pb.GetUserRolesRequest) (*pb.GetUserRolesResponse, error) {
	roles, err := s.roleRepo.GetRolesByUserID(ctx, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to fetch roles: %v", err)
	}

	pbRoles := make([]*pb.Role, len(roles))
	for i, role := range roles {
		pbRoles[i] = &pb.Role{
			Id:          role.ID,
			Name:        role.Name,
			Permissions: role.Permissions,
		}
	}

	return &pb.GetUserRolesResponse{Roles: pbRoles}, nil
}

func (s *AuthzService) AssignRoleToUser(ctx context.Context, req *pb.AssignRoleToUserRequest) (*emptypb.Empty, error) {
	if err := s.roleRepo.AssignRoleToUser(ctx, req.UserId, req.RoleId); err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to assign role: %v", err)
	}
	return &emptypb.Empty{}, nil
}

func (s *AuthzService) RemoveRoleFromUser(ctx context.Context, req *pb.RemoveRoleFromUserRequest) (*emptypb.Empty, error) {
	if err := s.roleRepo.RemoveRoleFromUser(ctx, req.UserId, req.RoleId); err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to remove role: %v", err)
	}
	return &emptypb.Empty{}, nil
}

func (s *AuthzService) CreateRole(ctx context.Context, req *pb.CreateRoleRequest) (*pb.CreateRoleResponse, error) {
	role := &ent.Role{
		Name:        req.Name,
		Permissions: make([]string, len(req.Permissions)),
	}
	for i, perm := range req.Permissions {
		role.Permissions[i] = perm.Resource + ":" + perm.Action
	}

	var err error

	if role, err = s.roleRepo.Create(ctx, role); err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to create role: %v", err)
	}

	return &pb.CreateRoleResponse{
		Role: &pb.Role{
			Id:          role.ID,
			Name:        role.Name,
			Permissions: role.Permissions,
		},
	}, nil
}

func (s *AuthzService) UpdateRole(ctx context.Context, req *pb.UpdateRoleRequest) (*pb.UpdateRoleResponse, error) {
	role, err := s.roleRepo.GetByID(ctx, req.RoleId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Role not found: %v", err)
	}

	if req.Name != nil {
		role.Name = *req.Name
	}
	if req.Permissions != nil {
		role.Permissions = make([]string, len(req.Permissions))
		for i, perm := range req.Permissions {
			role.Permissions[i] = perm.Resource + ":" + perm.Action
		}
	}

	if role, err = s.roleRepo.Update(ctx, role); err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to update role: %v", err)
	}

	return &pb.UpdateRoleResponse{
		Role: &pb.Role{
			Id:          role.ID,
			Name:        role.Name,
			Permissions: role.Permissions,
		},
	}, nil
}

func (s *AuthzService) DeleteRole(ctx context.Context, req *pb.DeleteRoleRequest) (*emptypb.Empty, error) {
	if err := s.roleRepo.Delete(ctx, req.RoleId); err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to delete role: %v", err)

	}

	return &emptypb.Empty{}, nil
}

func (s *AuthzService) GetRole(ctx context.Context, req *pb.GetRoleRequest) (*pb.GetRoleResponse, error) {
	role, err := s.roleRepo.GetByID(ctx, req.RoleId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Role not found: %v", err)
	}

	return &pb.GetRoleResponse{
		Role: &pb.Role{
			Id:          role.ID,
			Name:        role.Name,
			Permissions: role.Permissions,
		},
	}, nil
}

func (s *AuthzService) ListRoles(ctx context.Context, req *pb.ListRolesRequest) (*pb.ListRolesResponse, error) {
	nextPageToken, err := strconv.Atoi(req.PageToken)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid page token: %v", err)
	}
	nextPageToken += int(req.PageSize)
	roles, err := s.roleRepo.List(ctx, int(req.PageSize), nextPageToken)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to list roles: %v", err)
	}

	pbRoles := make([]*pb.Role, len(roles))
	for i, role := range roles {
		pbRoles[i] = &pb.Role{
			Id:          role.ID,
			Name:        role.Name,
			Permissions: role.Permissions,
		}
	}

	return &pb.ListRolesResponse{
		Roles:         pbRoles,
		NextPageToken: strconv.Itoa(nextPageToken),
	}, nil
}

func (s *AuthzService) AddPermissionToRole(ctx context.Context, req *pb.AddPermissionToRoleRequest) (*emptypb.Empty, error) {
	permission := req.Permission.Resource + ":" + req.Permission.Action
	if err := s.roleRepo.AddPermission(ctx, req.RoleId, permission); err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to add permission: %v", err)
	}
	return &emptypb.Empty{}, nil
}

func (s *AuthzService) RemovePermissionFromRole(ctx context.Context, req *pb.RemovePermissionFromRoleRequest) (*emptypb.Empty, error) {
	permission := req.Permission.Resource + ":" + req.Permission.Action
	if err := s.roleRepo.RemovePermission(ctx, req.RoleId, permission); err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to remove permission: %v", err)
	}
	return &emptypb.Empty{}, nil
}

func hasPermission(role *ent.Role, resource, action string) bool {
	permission := resource + ":" + action
	for _, perm := range role.Permissions {
		if perm == permission {
			return true
		}
	}
	return false
}
