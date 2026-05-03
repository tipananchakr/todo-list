package security

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"strings"
	"time"

	"github.com/tipananchakr/todo-list/internal/core/domain"
)

type HMACTokenManager struct {
	secret []byte
	ttl    time.Duration
	now    func() time.Time
}

type jwtHeader struct {
	Algorithm string `json:"alg"`
	Type      string `json:"typ"`
}

type jwtClaims struct {
	Subject   string `json:"sub"`
	ExpiresAt int64  `json:"exp"`
	IssuedAt  int64  `json:"iat"`
}

func NewHMACTokenManager(secret string, ttl time.Duration) *HMACTokenManager {
	return &HMACTokenManager{
		secret: []byte(secret),
		ttl:    ttl,
		now:    time.Now,
	}
}

func (m *HMACTokenManager) Generate(userID string) (string, error) {
	header := jwtHeader{Algorithm: "HS256", Type: "JWT"}
	now := m.now().UTC()
	claims := jwtClaims{
		Subject:   userID,
		IssuedAt:  now.Unix(),
		ExpiresAt: now.Add(m.ttl).Unix(),
	}

	headerSegment, err := encodeJSONSegment(header)
	if err != nil {
		return "", err
	}

	claimsSegment, err := encodeJSONSegment(claims)
	if err != nil {
		return "", err
	}

	unsignedToken := headerSegment + "." + claimsSegment
	return unsignedToken + "." + m.sign(unsignedToken), nil
}

func (m *HMACTokenManager) Validate(token string) (string, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return "", domain.ErrInvalidToken
	}

	unsignedToken := parts[0] + "." + parts[1]
	if !hmac.Equal([]byte(m.sign(unsignedToken)), []byte(parts[2])) {
		return "", domain.ErrInvalidToken
	}

	var claims jwtClaims
	if err := decodeJSONSegment(parts[1], &claims); err != nil {
		return "", domain.ErrInvalidToken
	}

	if claims.Subject == "" || m.now().UTC().Unix() > claims.ExpiresAt {
		return "", domain.ErrInvalidToken
	}

	return claims.Subject, nil
}

func (m *HMACTokenManager) sign(value string) string {
	mac := hmac.New(sha256.New, m.secret)
	mac.Write([]byte(value))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

func encodeJSONSegment(value any) (string, error) {
	bytes, err := json.Marshal(value)
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

func decodeJSONSegment(segment string, target any) error {
	bytes, err := base64.RawURLEncoding.DecodeString(segment)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, target)
}
