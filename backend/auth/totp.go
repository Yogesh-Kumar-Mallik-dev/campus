/**
 * BLOCK_AUTH_TOTP_RFC6238_001
 * Purpose: RFC 6238 Time-based One-Time Password (TOTP) generator and verifier.
 * Domain:  Central Auth & Permissions System (IAM)
 */

package auth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"crypto/subtle"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"net/url"
	"strings"
	"time"
)

// StandardTOTPProvider implements the TOTPProvider interface.
type StandardTOTPProvider struct {
	timeStep   time.Duration
	digits     int
	skewWindow int // Number of time-steps to check before/after current time
}

// NewStandardTOTPProvider creates a new RFC 6238 TOTP provider.
func NewStandardTOTPProvider() *StandardTOTPProvider {
	return &StandardTOTPProvider{
		timeStep:   30 * time.Second,
		digits:     6,
		skewWindow: 1, // +/- 30 seconds clock drift tolerance
	}
}

// GenerateSecret creates a 160-bit (20-byte) base32 secret and standard otpauth:// URL.
func (p *StandardTOTPProvider) GenerateSecret(accountEmail string, issuer string) (secret string, otpAuthURL string, err error) {
	rawSecret := make([]byte, 20)
	if _, err := rand.Read(rawSecret); err != nil {
		return "", "", fmt.Errorf("BLOCK_AUTH_TOTP_001: failed to generate secret bytes: %w", err)
	}

	secret = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(rawSecret)

	encodedIssuer := url.QueryEscape(issuer)
	encodedAccount := url.QueryEscape(accountEmail)
	otpAuthURL = fmt.Sprintf("otpauth://totp/%s:%s?secret=%s&issuer=%s&algorithm=SHA1&digits=%d&period=30",
		encodedIssuer, encodedAccount, secret, encodedIssuer, p.digits)

	return secret, otpAuthURL, nil
}

// ValidatePasscode checks if the given 6-digit passcode matches the secret within the skew window.
func (p *StandardTOTPProvider) ValidatePasscode(passcode string, secret string) bool {
	cleanedSecret := strings.ToUpper(strings.TrimSpace(secret))
	key, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(cleanedSecret)
	if err != nil {
		return false
	}

	cleanedPasscode := strings.TrimSpace(passcode)
	if len(cleanedPasscode) != p.digits {
		return false
	}

	now := time.Now().Unix()
	step := int64(p.timeStep.Seconds())
	currentInterval := now / step

	for i := -p.skewWindow; i <= p.skewWindow; i++ {
		testInterval := uint64(currentInterval + int64(i))
		expectedCode := p.computeCode(key, testInterval)
		if subtle.ConstantTimeCompare([]byte(cleanedPasscode), []byte(expectedCode)) == 1 {
			return true
		}
	}

	return false
}

func (p *StandardTOTPProvider) computeCode(key []byte, interval uint64) string {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, interval)

	mac := hmac.New(sha1.New, key)
	mac.Write(buf)
	hash := mac.Sum(nil)

	offset := hash[len(hash)-1] & 0x0f
	truncatedHash := binary.BigEndian.Uint32(hash[offset : offset+4])
	truncatedHash &= 0x7fffffff

	code := truncatedHash % 1000000
	return fmt.Sprintf("%06d", code)
}

// Ensure StandardTOTPProvider satisfies TOTPProvider interface.
var _ TOTPProvider = (*StandardTOTPProvider)(nil)
