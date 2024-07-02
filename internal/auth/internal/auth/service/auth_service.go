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

type AuthService struct {
	pb.UnimplementedAuthServiceServer
	userRepo    repository.UserRepository
	tokenRepo   *repository.TokenRepository
	jwtSecret   string
	tokenExpiry time.Duration
}

func NewAuthService(userRepo repository.UserRepository, tokenRepo *repository.TokenRepository, jwtSecret string, tokenExpiry time.Duration) *AuthService {
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
