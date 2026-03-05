package auth

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
)

type (
	ctxKeyClaims struct{}
	Claims       map[string]any
)

func ContextWithClaims(ctx context.Context, claims Claims) context.Context {
	return context.WithValue(ctx, ctxKeyClaims{}, claims)
}

func ClaimsFromContext(ctx context.Context) (Claims, error) {
	claims, ok := ctx.Value(ctxKeyClaims{}).(Claims)
	if !ok {
		return nil, errors.New("auth claims not found in context")
	}
	return claims, nil
}

// ParseClaimsPayload decodes a base64-url payload token into a generic claims map.
// The token itself is expected to be a JSON payload only.
func ParseClaimsPayload(token string) (Claims, error) {
	payloadBytes, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil {
		return nil, fmt.Errorf("payload decode failed: %w", err)
	}

	claims := Claims{}
	if err = json.Unmarshal(payloadBytes, &claims); err != nil {
		return nil, fmt.Errorf("payload json decode failed: %w", err)
	}
	return claims, nil
}

func BuildClaimsPayload(claims Claims) string {
	payloadBytes, _ := json.Marshal(claims)
	return base64.RawURLEncoding.EncodeToString(payloadBytes)
}
