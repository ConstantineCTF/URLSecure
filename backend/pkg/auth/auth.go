package auth

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5" // JWT library for token creation and parsing
	"github.com/joho/godotenv"     // Load .env environment variables
	"golang.org/x/crypto/bcrypt"   // Password hashing and comparison
)

var jwtKey []byte

func init() {
	// Load .env file to have access to environment variables like JWT_SECRET
	_ = godotenv.Load()

	// Get JWT secret key from environment
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		// Panic if secret is not set, as JWT cannot function without it
		panic("JWT_SECRET environment variable not set")
	}

	// Convert secret to byte slice for signing JWTs
	jwtKey = []byte(secret)
}

// Claims represents the payload stored inside JWT token
type Claims struct {
	UserID uint64 `json:"userId"` // User ID stored in token claims
	jwt.RegisteredClaims            // Standard JWT claims (expires, issued at, etc.)
}

// HashPassword hashes a plaintext password using bcrypt algorithm
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword compares a hashed password with a plaintext password
func CheckPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// CreateJWT creates a signed JWT token for a given user ID valid for 24 hours
func CreateJWT(userID uint64) (string, error) {
	expiration := time.Now().Add(24 * time.Hour)

	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiration), // Expiration time
			IssuedAt:  jwt.NewNumericDate(time.Now()), // Issue time
		},
	}

	// Create token with claims using HMAC SHA256 signing method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	
	// Sign token string with secret key
	return token.SignedString(jwtKey)
}

// ParseJWT parses and validates a JWT token string and returns claims
func ParseJWT(tokenStr string) (*Claims, error) {
	claims := &Claims{}

	// Parse token with claims, validating signature with jwtKey
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		// Enforce expected signing method
		if token.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected signing method")
		}
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	// Check token validity (expiration, signature, etc.)
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
