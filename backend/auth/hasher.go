/**
 * BLOCK_AUTH_HASHER_ARGON2ID_001
 * Purpose: Cryptographic password hashing using Argon2id with constant-time verification.
 * Inputs:  plainPassword (string), hashedPassword (string)
 * Outputs: Hash string ($argon2id$v=19$m=65536,t=3,p=4$salt$hash), Match boolean
 * Errors:  ERR_HASHING_FAILED
 */

package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

// Argon2idParams defines the security parameters for Argon2id hashing.
type Argon2idParams struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

// DefaultArgon2idParams provides recommended OWASP parameters.
var DefaultArgon2idParams = &Argon2idParams{
	Memory:      64 * 1024, // 64MB
	Iterations:  3,
	Parallelism: 4,
	SaltLength:  16,
	KeyLength:   32,
}

// Argon2idHasher implements the PasswordHasher interface.
type Argon2idHasher struct {
	params *Argon2idParams
}

// NewArgon2idHasher creates a new hasher with custom or default parameters.
func NewArgon2idHasher(params *Argon2idParams) *Argon2idHasher {
	if params == nil {
		params = DefaultArgon2idParams
	}
	return &Argon2idHasher{params: params}
}

// Hash generates an encoded Argon2id hash for the given plaintext password.
func (h *Argon2idHasher) Hash(password string) (string, error) {
	salt := make([]byte, h.params.SaltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("BLOCK_AUTH_HASHER_001: failed to generate cryptographic salt: %w", err)
	}

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		h.params.Iterations,
		h.params.Memory,
		h.params.Parallelism,
		h.params.KeyLength,
	)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encoded := fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		h.params.Memory,
		h.params.Iterations,
		h.params.Parallelism,
		b64Salt,
		b64Hash,
	)

	return encoded, nil
}

// Compare checks if a plaintext password matches an encoded Argon2id hash using constant-time comparison.
func (h *Argon2idHasher) Compare(hashedPassword, plainPassword string) bool {
	parts := strings.Split(hashedPassword, "$")
	if len(parts) != 6 {
		return false
	}

	if parts[1] != "argon2id" {
		return false
	}

	var version int
	if _, err := fmt.Sscanf(parts[2], "v=%d", &version); err != nil || version != argon2.Version {
		return false
	}

	var memory, iterations uint32
	var parallelism uint8
	if _, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &iterations, &parallelism); err != nil {
		return false
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false
	}

	expectedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false
	}

	computedHash := argon2.IDKey(
		[]byte(plainPassword),
		salt,
		iterations,
		memory,
		parallelism,
		uint32(len(expectedHash)),
	)

	return subtle.ConstantTimeCompare(expectedHash, computedHash) == 1
}

// Ensure Argon2idHasher satisfies the PasswordHasher interface.
var _ PasswordHasher = (*Argon2idHasher)(nil)
