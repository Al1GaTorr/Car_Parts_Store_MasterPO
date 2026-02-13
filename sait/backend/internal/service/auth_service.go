package service

import (
	"context"
	"crypto/subtle"
	"errors"
	"strings"
	"time"

	"bazarpo-backend/internal/model"
	"bazarpo-backend/internal/repo"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmailPasswordInvalid = errors.New("email/password invalid")
	ErrEmailAlreadyExists   = errors.New("email already exists")
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrUnauthorized         = errors.New("unauthorized")
)

type AuthService struct {
	repo *repo.Repository
	env  model.Env
}

func NewAuthService(r *repo.Repository, env model.Env) *AuthService {
	return &AuthService{repo: r, env: env}
}

func (s *AuthService) SignToken(userID primitive.ObjectID, role string) (string, error) {
	claims := model.Claims{
		UserID: userID.Hex(),
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(s.env.JWTSecret))
}

func (s *AuthService) ParseToken(authHeader string) (*model.Claims, error) {
	if authHeader == "" {
		return nil, errors.New("missing auth")
	}
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return nil, errors.New("invalid auth header")
	}
	tokenStr := parts[1]
	tok, err := jwt.ParseWithClaims(tokenStr, &model.Claims{}, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.env.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := tok.Claims.(*model.Claims)
	if !ok || !tok.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

func (s *AuthService) EnsureAdmin(ctx context.Context) error {
	if s.env.AdminEmail == "" || s.env.AdminPassword == "" {
		return nil
	}
	email := strings.ToLower(strings.TrimSpace(s.env.AdminEmail))

	user, err := s.repo.FindUserByEmail(ctx, email)
	if err == nil {
		if user.Role != "admin" {
			return s.repo.UpdateUserRoleByID(ctx, user.ID, "admin")
		}
		return nil
	}
	if !errors.Is(err, mongo.ErrNoDocuments) {
		return err
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(s.env.AdminPassword), 10)
	_, err = s.repo.InsertUser(ctx, model.UserDoc{
		FirstName:    "Admin",
		LastName:     "",
		Email:        email,
		PasswordHash: string(hash),
		Role:         "admin",
		CreatedAt:    time.Now(),
	})
	return err
}

func (s *AuthService) Register(ctx context.Context, req model.RegisterRequest) (string, string, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))
	if email == "" || len(req.Password) < 6 {
		return "", "", ErrEmailPasswordInvalid
	}

	count, err := s.repo.CountUsersByEmail(ctx, email)
	if err != nil {
		return "", "", err
	}
	if count > 0 {
		return "", "", ErrEmailAlreadyExists
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	role := "user"
	if subtle.ConstantTimeCompare([]byte(email), []byte(s.env.AdminEmail)) == 1 &&
		subtle.ConstantTimeCompare([]byte(req.Password), []byte(s.env.AdminPassword)) == 1 {
		role = "admin"
	}

	oid, err := s.repo.InsertUser(ctx, model.UserDoc{
		FirstName:    strings.TrimSpace(req.FirstName),
		LastName:     strings.TrimSpace(req.LastName),
		Email:        email,
		PasswordHash: string(hash),
		Role:         role,
		CreatedAt:    time.Now(),
	})
	if err != nil {
		return "", "", err
	}

	token, err := s.SignToken(oid, role)
	if err != nil {
		return "", "", err
	}
	return token, role, nil
}

func (s *AuthService) Login(ctx context.Context, req model.LoginRequest) (string, string, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))
	pass := req.Password

	if subtle.ConstantTimeCompare([]byte(email), []byte(s.env.AdminEmail)) == 1 &&
		subtle.ConstantTimeCompare([]byte(pass), []byte(s.env.AdminPassword)) == 1 {
		_ = s.EnsureAdmin(ctx)
		u, err := s.repo.FindUserByEmail(ctx, email)
		if err != nil {
			return "", "", ErrInvalidCredentials
		}
		token, err := s.SignToken(u.ID, "admin")
		if err != nil {
			return "", "", err
		}
		return token, "admin", nil
	}

	u, err := s.repo.FindUserByEmail(ctx, email)
	if err != nil {
		return "", "", ErrInvalidCredentials
	}
	if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(pass)) != nil {
		return "", "", ErrInvalidCredentials
	}

	token, err := s.SignToken(u.ID, u.Role)
	if err != nil {
		return "", "", err
	}
	return token, u.Role, nil
}

func (s *AuthService) Me(ctx context.Context, userHex string) (*model.UserDoc, error) {
	oid, err := primitive.ObjectIDFromHex(userHex)
	if err != nil {
		return nil, ErrUnauthorized
	}
	u, err := s.repo.FindUserByID(ctx, oid)
	if err != nil {
		return nil, ErrUnauthorized
	}
	return u, nil
}

