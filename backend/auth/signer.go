/**
 * BLOCK_AUTH_SIGNER_JWT_001
 * Purpose: JWT access token issuance/validation and refresh token family generation.
 * Domain:  Central Auth & Permissions System (IAM)
 */

package auth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// JWTSigner implements the TokenSigner interface using HMAC-SHA256.
type JWTSigner struct {
	secretKey     []byte
	issuer        string
	tokenLifespan time.Duration
}

// NewJWTSigner creates a new JWT signer instance.
func NewJWTSigner(secretKey string, issuer string, tokenLifespan time.Duration) *JWTSigner {
	if tokenLifespan <= 0 {
		tokenLifespan = 15 * time.Minute
	}
	if issuer == "" {
		issuer = "campus.institution.edu"
	}
	return &JWTSigner{
		secretKey:     []byte(secretKey),
		issuer:        issuer,
		tokenLifespan: tokenLifespan,
	}
}

type jwtHeader struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

// GenerateAccessToken produces a signed HS256 JWT access token string.
func (s *JWTSigner) GenerateAccessToken(claims Claims) (string, error) {
	now := time.Now().Unix()
	if claims.IssuedAt == 0 {
		claims.IssuedAt = now
	}
	if claims.ExpiresAt == 0 {
		claims.ExpiresAt = now + int64(s.tokenLifespan.Seconds())
	}

	header := jwtHeader{
		Alg: "HS256",
		Typ: "JWT",
	}

	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", fmt.Errorf("BLOCK_AUTH_SIGNER_001: failed to marshal header: %w", err)
	}

	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		return "", fmt.Errorf("BLOCK_AUTH_SIGNER_002: failed to marshal claims: %w", err)
	}

	encodedHeader := base64.RawURLEncoding.EncodeToString(headerJSON)
	encodedClaims := base64.RawURLEncoding.EncodeToString(claimsJSON)

	signingInput := encodedHeader + "." + encodedClaims

	mac := hmac.New(sha256.New, s.secretKey)
	mac.Write([]byte(signingInput))
	signature := mac.Sum(nil)
	encodedSig := base64.RawURLEncoding.EncodeToString(signature)

	return signingInput + "." + encodedSig, nil
}

// ValidateAccessToken verifies the token signature and returns decoded Claims.
func (s *JWTSigner) ValidateAccessToken(tokenString string) (*Claims, error) {
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return nil, ErrUnauthorized
	}

	signingInput := parts[0] + "." + parts[1]
	receivedSig, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return nil, ErrUnauthorized
	}

	mac := hmac.New(sha256.New, s.secretKey)
	mac.Write([]byte(signingInput))
	expectedSig := mac.Sum(nil)

	if subtle.ConstantTimeCompare(receivedSig, expectedSig) != 1 {
		return nil, ErrUnauthorized
	}

	claimsBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, ErrUnauthorized
	}

	var claims Claims
	if err := json.Unmarshal(claimsBytes, &claims); err != nil {
		return nil, ErrUnauthorized
	}

	now := time.Now().Unix()
	if claims.ExpiresAt <= now {
		return nil, ErrTokenExpired
	}

	return &claims, nil
}

// GenerateRefreshToken generates a secure 256-bit random hex string and its SHA-256 hash.
func (s *JWTSigner) GenerateRefreshToken() (plainToken string, tokenHash string, err error) {
	randomBytes := make([]byte, 32)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", "", fmt.Errorf("BLOCK_AUTH_SIGNER_003: failed to generate random token bytes: %w", err)
	}

	plainToken = hex.EncodeToString(randomBytes)
	tokenHash = s.HashRefreshToken(plainToken)
	return plainToken, tokenHash, nil
}

// HashRefreshToken calculates the SHA-256 hash of a plaintext refresh token.
func (s *JWTSigner) HashRefreshToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}

// Ensure JWTSigner satisfies TokenSigner interface.
var _ TokenSigner = (*JWTSigner)(nil)
