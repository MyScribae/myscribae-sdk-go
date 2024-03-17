package provider

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
)

type SubscriberToken struct {
	Subject       string
	Expiration    time.Time
	Issuer        string
	IssuedAt      time.Time
	ScriptsClaims []ScriptClaim
}

var (
	ErrInvalidSubscriberToken = errors.New("invalid subscriber token")
	ErrTokenMissingClaims     = errors.New("token missing claims")
	ErrMissingSubject         = errors.New("missing subject")
	ErrMissingExpiration      = errors.New("missing expiration")
	ErrInvalidExpiration      = errors.New("invalid expiration")
	ErrExpiredToken           = errors.New("expired token")
	ErrMissingIssuer          = errors.New("missing issuer")
	ErrMissingIssuedAt        = errors.New("missing issued at")
	ErrInvalidIssuedAt        = errors.New("invalid issued at")
	ErrTokenNotYetEffective   = errors.New("token not yet effective")
)

func NewSubscriberToken(token *jwt.Token) (*SubscriberToken, error) {
	if token.Valid {
		return nil, ErrInvalidSubscriberToken
	}

	if err := token.Claims.Valid(); err != nil {
		log.Printf("Token claims are invalid: %v", err)
		return nil, ErrTokenMissingClaims
	}

	var (
		claims           = token.Claims.(jwt.MapClaims)
		sub       string = claims["sub"].(string)
		exp       string = claims["exp"].(string)
		iss       string = claims["iss"].(string)
		iat       string = claims["iat"].(string)
		claimsRaw string = claims["claims"].(string)
	)

	log.Printf("Subscriber token: sub=%s, exp=%s, iss=%s, iat=%s", sub, exp, iss, iat)
	if sub == "" {
		return nil, ErrMissingSubject
	}

	if exp == "" {
		return nil, ErrMissingExpiration
	}

	expVal, err := strconv.Atoi(exp)
	if err != nil {
		log.Printf("Failed to parse expiration (%s): %v", exp, err)
		return nil, ErrInvalidExpiration
	}

	expTime := time.Unix(int64(expVal), 0)
	if expTime.Before(time.Now()) {
		return nil, ErrExpiredToken
	}

	if iss == "" {
		return nil, ErrMissingIssuer
	}

	if iat == "" {
		return nil, ErrMissingIssuedAt
	}
	iatVal, err := strconv.Atoi(iat)
	if err != nil {
		log.Printf("Failed to parse issued at (%s): %v", iat, err)
		return nil, ErrInvalidIssuedAt
	}
	iatTime := time.Unix(int64(iatVal), 0)
	if iatTime.After(time.Now()) {
		return nil, ErrTokenNotYetEffective
	}

	// get claims from claims string
	var scriptClaims []ScriptClaim
	if err := json.Unmarshal([]byte(claimsRaw), &scriptClaims); err != nil {
		log.Printf("Failed to unmarshal claims: %v", err)
		return nil, fmt.Errorf("failed to unmarshal claims: %w", err)
	}

	return &SubscriberToken{
		Subject:       sub,
		Expiration:    expTime,
		Issuer:        iss,
		IssuedAt:      iatTime,
		ScriptsClaims: scriptClaims,
	}, nil
}
