package ocrmserver

import (
	"context"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/sergazyyev/crmlibrary/ocrmmodel"
	"net/http"
	"strings"
	"time"
)

const (
	ctxClaimsKey ctxKey = iota
	ctxRequestIdKey
)

type ctxKey int8

func (server *BaseServer) authenticateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !server.UseAuthMidd {
			next.ServeHTTP(w, r)
			return
		}
		ignorePaths := strings.Join(server.AuthIgnorePaths, "")
		if server.AuthIgnorePaths != nil && strings.Contains(ignorePaths, r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			server.Error(w, http.StatusUnauthorized, errors.New("not authenticated"))
			return
		}
		claims := &ocrmmodel.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (i interface{}, e error) {
			return server.JwtKey, nil
		})
		if (token != nil && !token.Valid) || (err != nil && err == jwt.ErrSignatureInvalid) {
			server.Error(w, http.StatusUnauthorized, errors.New("not authenticated"))
			return
		}
		if err != nil {
			server.Logger.Errorf("cant parse token with claim, err %v", err)
			server.Error(w, http.StatusInternalServerError, err)
			return
		}
		expTime := time.Now().Add(time.Duration(server.JwtTokenLiveMinutes) * time.Minute)
		claims.ExpiresAt = expTime.Unix()
		token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err = token.SignedString(server.JwtKey)
		if err != nil {
			server.Logger.Errorf("cant get jwtString for user %s, err: %v", claims.User.Username, err)
			server.Error(w, http.StatusInternalServerError, err)
			return
		}
		w.Header().Set("Authorization", tokenString)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxClaimsKey, claims)))
	})
}

func (server *BaseServer) loggingRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		server.Logger.Tracef("started %s %s remote_address=%s request_id=%s", r.Method, r.RequestURI, r.RemoteAddr, r.Context().Value(ctxRequestIdKey))
		now := time.Now()
		rw := &responseWriter{w, http.StatusOK}
		next.ServeHTTP(rw, r)
		server.Logger.Tracef("completed with %d %s in %v remote_address=%s request_id=%s",
			rw.statusCode,
			http.StatusText(rw.statusCode),
			time.Now().Sub(now),
			r.RemoteAddr,
			r.Context().Value(ctxRequestIdKey))
	})
}

func (server *BaseServer) setRequestId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()
		w.Header().Set("X-Request-Id", id)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxRequestIdKey, id)))
	})
}