package service

import (
	"context"
	"time"

	"auth/ent"
	"auth/internal/repository"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "auth/pb"
)

// IAuthService defines the interface for authentication and user management operations.
// It provides methods for user authentication, token management, and user CRUD operations.
type IAuthService interface {
	// Login authenticates a user and returns a login response with tokens.
	//
	// Parameters:
	//   - ctx: A context.Context for handling deadlines, cancellations, and request-scoped values.
	//   - req: A pointer to pb.LoginRequest containing user credentials and any additional login parameters.
	//
	// Returns:
	//   - *pb.LoginResponse: A response containing authentication tokens (e.g., access token, refresh token) and any additional user information.
	//   - error: An error if authentication fails, such as invalid credentials, account lockout, or internal server issues.
	//
	// The method should implement proper security measures, such as rate limiting and account lockout mechanisms.
	// It should also log authentication attempts for security auditing purposes.
	Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error)

	// Logout terminates a user's active session, invalidating their current authentication tokens.
	//
	// Parameters:
	//   - ctx: A context.Context for handling deadlines, cancellations, and request-scoped values.
	//   - req: A pointer to pb.LogoutRequest containing the user's session information or tokens to be invalidated.
	//
	// Returns:
	//   - *emptypb.Empty: An empty response indicating successful logout.
	//   - error: An error if the logout process fails, such as invalid session, or internal server issues.
	//
	// This method should ensure all related session data is cleared and any distributed caches are updated.
	// It should also log the logout event for security auditing purposes.
	Logout(ctx context.Context, req *pb.LogoutRequest) (*emptypb.Empty, error)

	// RefreshToken extends a user's session by providing a new access token.
	//
	// Parameters:
	//   - ctx: A context.Context for handling deadlines, cancellations, and request-scoped values.
	//   - req: A pointer to pb.RefreshTokenRequest containing the current refresh token.
	//
	// Returns:
	//   - *pb.RefreshTokenResponse: A response containing a new access token and optionally a new refresh token.
	//   - error: An error if token refresh fails, such as expired refresh token, token reuse, or internal server issues.
	//
	// This method should implement proper security checks, including refresh token rotation if applicable.
	// It should also validate the refresh token's expiration and ensure it hasn't been revoked.
	RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error)

	// ValidateToken checks the validity and integrity of a given token.
	//
	// Parameters:
	//   - ctx: A context.Context for handling deadlines, cancellations, and request-scoped values.
	//   - req: A pointer to pb.ValidateTokenRequest containing the token to be validated.
	//
	// Returns:
	//   - *pb.ValidateTokenResponse: A response indicating the token's validity and any associated metadata.
	//   - error: An error if the validation process fails due to internal server issues.
	//
	// This method should check the token's signature, expiration, and ensure it hasn't been revoked.
	// It may also return additional information about the token, such as associated user ID or permissions.
	ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error)

	// RegisterUser creates a new user account in the system.
	//
	// Parameters:
	//   - ctx: A context.Context for handling deadlines, cancellations, and request-scoped values.
	//   - req: A pointer to pb.RegisterUserRequest containing the new user's information.
	//
	// Returns:
	//   - *pb.RegisterUserResponse: A response containing the newly created user's ID and any additional information.
	//   - error: An error if user registration fails, such as duplicate username/email, invalid data, or internal server issues.
	//
	// This method should implement proper data validation, secure password hashing, and any necessary unique constraint checks.
	// It may also trigger additional processes like sending a verification email or assigning default roles.
	RegisterUser(ctx context.Context, req *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error)

	// UpdateUser modifies existing user information.
	//
	// Parameters:
	//   - ctx: A context.Context for handling deadlines, cancellations, and request-scoped values.
	//   - req: A pointer to pb.UpdateUserRequest containing the user ID and fields to be updated.
	//
	// Returns:
	//   - *pb.UpdateUserResponse: A response confirming the update and containing the updated user information.
	//   - error: An error if the update fails, such as user not found, invalid data, or internal server issues.
	//
	// This method should implement proper authorization checks to ensure the requester has permission to update the user.
	// It should also validate the updated data and handle any unique constraint violations.
	UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error)

	// DeleteUser removes a user account from the system.
	//
	// Parameters:
	//   - ctx: A context.Context for handling deadlines, cancellations, and request-scoped values.
	//   - req: A pointer to pb.DeleteUserRequest containing the ID of the user to be deleted.
	//
	// Returns:
	//   - *emptypb.Empty: An empty response indicating successful deletion.
	//   - error: An error if the deletion fails, such as user not found, unauthorized deletion, or internal server issues.
	//
	// This method should implement proper authorization checks to ensure the requester has permission to delete the user.
	// It should also handle related data cleanups, such as revoking all active sessions and handling foreign key constraints.
	DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*emptypb.Empty, error)

	// GetUser retrieves user information based on provided criteria.
	//
	// Parameters:
	//   - ctx: A context.Context for handling deadlines, cancellations, and request-scoped values.
	//   - req: A pointer to pb.GetUserRequest containing search criteria (e.g., user ID, username, email).
	//
	// Returns:
	//   - *pb.GetUserResponse: A response containing the requested user information.
	//   - error: An error if the retrieval fails, such as user not found, unauthorized access, or internal server issues.
	//
	// This method should implement proper authorization checks to ensure the requester has permission to access the user information.
	// It should also consider privacy settings and may return different levels of detail based on the requester's permissions.
	GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error)
}

type AuthService struct {
	pb.UnimplementedAuthServiceServer
	userRepo    repository.IUserRepository
	tokenRepo   repository.ITokenRepository
	jwtSecret   string
	tokenExpiry time.Duration
}

func NewAuthService(userRepo repository.IUserRepository, tokenRepo repository.ITokenRepository, jwtSecret string, tokenExpiry time.Duration) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		tokenRepo:   tokenRepo,
		jwtSecret:   jwtSecret,
		tokenExpiry: tokenExpiry,
	}
}

func (s *AuthService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	user, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to fetch user: %v", err)
	}

	if user == nil || !s.userRepo.CheckPassword(ctx, req.Password) {
		return nil, status.Errorf(codes.Unauthenticated, "Invalid username or password")
	}

	accessToken, err := s.generateJWT(user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to generate access token: %v", err)
	}

	refreshToken, err := s.generateRefreshToken(user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to generate refresh token: %v", err)
	}

	// Store tokens in the database
	_, err = s.tokenRepo.Create(ctx, user.ID, accessToken, "access", time.Now().Add(s.tokenExpiry))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to store access token: %v", err)
	}

	_, err = s.tokenRepo.Create(ctx, user.ID, refreshToken, "refresh", time.Now().Add(24*time.Hour))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to store refresh token: %v", err)
	}

	return &pb.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(s.tokenExpiry.Seconds()),
	}, nil
}

func (s *AuthService) Logout(ctx context.Context, req *pb.LogoutRequest) (*emptypb.Empty, error) {
	if err := s.tokenRepo.RevokeToken(ctx, req.AccessToken); err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to revoke token: %v", err)
	}
	return &emptypb.Empty{}, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	token, err := s.tokenRepo.GetByToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "Invalid refresh token: %v", err)
	}

	if token.Type != "refresh" || token.Revoked || token.ExpiresAt.Before(time.Now()) {
		return nil, status.Error(codes.Unauthenticated, "Invalid or expired refresh token")
	}

	user, err := s.userRepo.GetByID(ctx, token.Edges.User.ID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to fetch user: %v", err)
	}

	accessToken, err := s.generateJWT(user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to generate access token: %v", err)
	}

	refreshToken, err := s.generateRefreshToken(user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to generate refresh token: %v", err)
	}

	// Revoke the old refresh token
	if err := s.tokenRepo.RevokeToken(ctx, req.RefreshToken); err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to revoke old refresh token: %v", err)
	}

	// Store new tokens in the database
	_, err = s.tokenRepo.Create(ctx, user.ID, accessToken, "access", time.Now().Add(s.tokenExpiry))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to store access token: %v", err)
	}

	_, err = s.tokenRepo.Create(ctx, user.ID, refreshToken, "refresh", time.Now().Add(24*time.Hour))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to store refresh token: %v", err)
	}

	return &pb.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(s.tokenExpiry.Seconds()),
	}, nil
}

func (s *AuthService) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	token, err := s.tokenRepo.GetByToken(ctx, req.Token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "Invalid token: %v", err)
	}

	if token.Revoked || token.ExpiresAt.Before(time.Now()) {
		return nil, status.Error(codes.Unauthenticated, "Token is revoked or expired")
	}

	user, err := s.userRepo.GetByID(ctx, token.Edges.User.ID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to fetch user: %v", err)
	}

	return &pb.ValidateTokenResponse{
		User: &pb.User{
			Id:       user.ID,
			Username: user.Username,
			Email:    user.Email,
		},
	}, nil
}

func (s *AuthService) RegisterUser(ctx context.Context, req *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error) {
	existingUser, _ := s.userRepo.GetByUsername(ctx, req.Username)
	if existingUser != nil {
		return nil, status.Error(codes.AlreadyExists, "Username already exists")
	}

	existingUser, _ = s.userRepo.GetByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, status.Error(codes.AlreadyExists, "Email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to hash password: %v", err)
	}

	user := &ent.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	if user, err = s.userRepo.Create(ctx, user); err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to create user: %v", err)
	}

	return &pb.RegisterUserResponse{
		User: &pb.User{
			Id:       user.ID,
			Username: user.Username,
			Email:    user.Email,
		},
	}, nil
}

func (s *AuthService) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "User not found: %v", err)
	}

	if req.Username != nil {
		user.Username = *req.Username
	}
	if req.Email != nil {
		user.Email = *req.Email
	}
	if req.Password != nil {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "Failed to hash password: %v", err)
		}
		user.Password = string(hashedPassword)
	}

	if user, err = s.userRepo.Update(ctx, user); err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to update user: %v", err)
	}

	// Revoke all existing tokens for the user
	if err := s.tokenRepo.RevokeAllUserTokens(ctx, user.ID); err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to revoke user tokens: %v", err)
	}

	return &pb.UpdateUserResponse{
		User: &pb.User{
			Id:       user.ID,
			Username: user.Username,
			Email:    user.Email,
		},
	}, nil
}

func (s *AuthService) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*emptypb.Empty, error) {
	// Revoke all tokens for the user
	if err := s.tokenRepo.RevokeAllUserTokens(ctx, req.UserId); err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to revoke user tokens: %v", err)
	}

	if err := s.userRepo.Delete(ctx, req.UserId); err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to delete user: %v", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *AuthService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "User not found: %v", err)
	}

	return &pb.GetUserResponse{
		User: &pb.User{
			Id:       user.ID,
			Username: user.Username,
			Email:    user.Email,
		},
	}, nil
}

func (s *AuthService) generateJWT(user *ent.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(s.tokenExpiry).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

func (s *AuthService) generateRefreshToken(user *ent.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(), // Refresh token valid for 24 hours
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}
