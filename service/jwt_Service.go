package service

import (
	"errors"
	"fmt"
	"gatherly-app/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)


type JwtService interface {
	GenerateToken(userID int, email string, role string, latitude, longitude float64) (string, error)
	ValidateToken(tokenString string) (*Claims, error)
}


type jwtService struct {
	cfg config.TokenConfig
}

type Claims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	Latitude float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	jwt.RegisteredClaims // Embeds standard claims like expiration (ExpiresAt), issued at (IssuedAt)
}

func NewJwtService(cfg config.TokenConfig) JwtService {
	return &jwtService{
		cfg: cfg,
	}
}


func (s *jwtService) GenerateToken(userID int, email string, role string, latitude, longitude float64) (string, error) {
	claims := &Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		Latitude: latitude,
		Longitude: longitude,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(s.cfg.AccessTokenLifeTime) * time.Hour)),
			IssuedAt: jwt.NewNumericDate(time.Now()),
			Issuer: s.cfg.ApplicationName,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(s.cfg.JwtSignatureKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}

// --- ValidateToken Method ---
// Parses and validates a JWT token string, returning the custom claims if valid
func (s *jwtService) ValidateToken(tokenString string) (*Claims, error) {
	// Create an empty claims struct to be populated by ParseWithClaims
	claims := &Claims{}

	// Parse the token string.
	// ParseWithClaims requires the token string, the claims struct instance to populate,
	// and a function to provide the key for verification.
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		// Verify that the signing method is HMAC (HS256 in this case)
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			// Return an error if the signing method is unexpected (e.g., 'none' or RSA when expecting HMAC)
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		// Return the secret key used for signing (from config) for the library to verify the signature
		return []byte(s.cfg.JwtSignatureKey), nil
	})

	// Handle potential errors during parsing and validation
	if err != nil {
		// Check for specific validation errors provided by the jwt library
		if errors.Is(err, jwt.ErrTokenExpired) {
			// Return a specific error if the token has expired
			return nil, errors.New("token has expired")
		}
		if errors.Is(err, jwt.ErrTokenNotValidYet) {
			// Return a specific error if the token is not yet valid (based on 'nbf' claim if used)
			return nil, errors.New("token not active yet")
		}
		// For other errors (e.g., signature mismatch, malformed token), return a wrapped error
		return nil, fmt.Errorf("token parsing/validation error: %w", err)
	}

	// Although ParseWithClaims handles validation, double-check the token's Valid flag
	if !token.Valid {
		// This might be redundant if 'err' catches all invalid cases, but acts as a safeguard
		return nil, errors.New("token is invalid")
	}

	// If parsing and validation were successful, the 'claims' variable is populated
	// and the token is valid. Return the populated claims struct.
	return claims, nil
}

