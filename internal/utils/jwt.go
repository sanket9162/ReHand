package utils

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/sanket9162/rehand/internal/config"
)

const (
	TokenTypeAccess  = "access"
	TokenTypeRefresh = "refresh"

	accessCookieName  = "access_token"
	refreshCookieName = "refresh_token"
)

type JWT struct {
	cfg config.Auth
}

func NewJWT(cfg config.Auth) *JWT {
	return &JWT{cfg: cfg}
}

type JwtUser struct {
	ID uuid.UUID `json:"id"`
}

type TokenPair struct {
	Token        string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Claims struct {
	ID        uuid.UUID `json:"id"`
	TokenType string    `json:"token_type"`
	jwt.RegisteredClaims
}

// GenerateTokenPair creates a signed access and refresh token pair for user.
func (j *JWT) GenerateTokenPair(user *JwtUser) (TokenPair, error) {
	accessToken, err := j.generateToken(user, TokenTypeAccess, j.cfg.TokenExpiry)
	if err != nil {
		return TokenPair{}, err
	}

	refreshToken, err := j.generateToken(user, TokenTypeRefresh, j.cfg.RefreshExpiry)
	if err != nil {
		return TokenPair{}, err
	}

	return TokenPair{
		Token:        accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (j *JWT) generateToken(user *JwtUser, tokenType string, ttl time.Duration) (string, error) {
	now := time.Now().UTC()
	claims := Claims{
		ID:        user.ID,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID.String(),
			Audience:  jwt.ClaimStrings{j.cfg.Audience},
			Issuer:    j.cfg.Issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(j.cfg.Secret))
	if err != nil {
		return "", fmt.Errorf("sign %s token: %w", tokenType, err)
	}

	return signedToken, nil
}

func (j *JWT) GetAccessTokenCookie(accessToken string) *http.Cookie {
	return j.buildCookie(accessCookieName, accessToken, j.cfg.TokenExpiry, false)
}

func (j *JWT) GetRefreshCookie(refreshToken string) *http.Cookie {
	return j.buildCookie(refreshCookieName, refreshToken, j.cfg.RefreshExpiry, false)
}

func (j *JWT) GetExpiredAccessTokenCookie() *http.Cookie {
	return j.buildCookie(accessCookieName, "", 0, true)
}

func (j *JWT) GetExpiredRefreshCookie() *http.Cookie {
	return j.buildCookie(refreshCookieName, "", 0, true)
}

func (j *JWT) buildCookie(name, value string, ttl time.Duration, expired bool) *http.Cookie {
	cookie := &http.Cookie{
		Name:     name,
		Path:     j.cfg.CookiePath,
		Value:    value,
		Domain:   j.cfg.CookieDomain,
		HttpOnly: true,
		Secure:   j.cfg.CookieSecure,
		SameSite: http.SameSiteLaxMode,
	}

	if expired {
		cookie.Expires = time.Unix(0, 0)
		cookie.MaxAge = -1
		return cookie
	}

	cookie.Expires = time.Now().UTC().Add(ttl)
	cookie.MaxAge = int(ttl.Seconds())
	return cookie
}

// GetTokenFromHeader extracts a Bearer token from the Authorization header.
func (j *JWT) GetTokenFromHeader(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("no auth header")
	}

	token, ok := strings.CutPrefix(authHeader, "Bearer ")
	if !ok {
		return "", errors.New("invalid auth header")
	}

	token = strings.TrimSpace(token)
	if token == "" {
		return "", errors.New("empty auth token")
	}

	return token, nil
}

// ValidateToken verifies the JWT signature, issuer, audience, expiry, and expected token type.
func (j *JWT) ValidateToken(tokenString, expectedType string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (any, error) {
			return []byte(j.cfg.Secret), nil
		},
		jwt.WithIssuer(j.cfg.Issuer),
		jwt.WithAudience(j.cfg.Audience),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	)
	if err != nil {
		return nil, err
	}

	if !token.Valid || claims.TokenType != expectedType {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
