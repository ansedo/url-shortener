package middlewares

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"github.com/google/uuid"
	"net/http"
)

type CookieCtxType string

const (
	CookieCtxName CookieCtxType = "UID"

	hmacKey    = "One Flew over the Cuckoo's Nest"
	cookieName = "token"
)

func getNewToken() string {
	uid := uuid.New().String()
	h := hmac.New(sha256.New, []byte(hmacKey))
	h.Write([]byte(uid))
	return hex.EncodeToString(h.Sum(nil)) + uid
}

func getTokenFromRequest(r *http.Request) (string, error) {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

func isTokenValid(token string) bool {
	h := hmac.New(sha256.New, []byte(hmacKey))
	if len(token) < h.BlockSize() {
		return false
	}
	h.Write([]byte(token[h.BlockSize():]))
	sign, err := hex.DecodeString(token[:h.BlockSize()])
	if err != nil {
		return false
	}
	return hmac.Equal(sign, h.Sum(nil))
}

func getUIDFromValidToken(token string) string {
	return token[hmac.New(sha256.New, []byte(hmacKey)).BlockSize():]
}

func Cookie(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := getTokenFromRequest(r)
		if err != nil || !isTokenValid(token) {
			token = getNewToken()
		}
		http.SetCookie(w, &http.Cookie{Name: cookieName, Value: token})
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), CookieCtxName, getUIDFromValidToken(token))))
	})
}
