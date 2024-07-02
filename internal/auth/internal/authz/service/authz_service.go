package service

import (
	"context"

	"auth/ent"
	"auth/internal/repository"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "auth/pb"
)

type AuthzService struct {
	pb.UnimplementedAuthorizationServiceServer
	roleRepo repository.RoleRepository
	userRepo repository.UserRepository
}

func NewAuthzService(roleRepo repository.RoleRepository, userRepo repository.UserRepository) *AuthzService {
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

	if err := s.roleRepo.Create(ctx, role); err != nil {
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

	if err := s.roleRepo.Update(ctx, role); err != nil {
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
	roles, nextPageToken, err := s.roleRepo.List(ctx, int(req.PageSize), req.PageToken)
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
		NextPageToken: nextPageToken,
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
