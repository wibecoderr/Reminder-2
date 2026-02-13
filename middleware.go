package utils

import (
	"context"
	"net/http"
	"strings"

	dbhelper "github.com/wibecoderr/Reminder-2.git/database"
	"github.com/wibecoderr/Reminder-2.git/model"
)

type contextKey struct{}

var usercontextKey = contextKey{}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			RespondError(w, http.StatusUnauthorized, nil, "missing authorization header")
			return
		}

		// Verify Bearer prefix
		if !strings.HasPrefix(authHeader, "Bearer ") {
			RespondError(w, http.StatusUnauthorized, nil, "invalid authorization format, expected 'Bearer <token>'")
			return
		}

		// Extract token
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenStr == "" {
			RespondError(w, http.StatusUnauthorized, nil, "invalid authorization format, expected 'Bearer <token>'")
			return
		}

		// Verify JWT and extract claims
		userID, sessionID, err := VerifyJWT(tokenStr)
		if err != nil {
			RespondError(w, http.StatusUnauthorized, err, "invalid or expired token")
			return
		}

		// Verify session exists in database and is not expired
		dbUserID, err := dbhelper.GetUserIDBySession(sessionID)
		if err != nil {
			RespondError(w, http.StatusUnauthorized, err, "session not found or expired")
			return
		}

		// Verify userID from JWT matches userID from database
		if dbUserID != userID {
			RespondError(w, http.StatusUnauthorized, nil, "invalid or expired token")
			return
		}

		// Set user context
		user := &model.UserCxt{
			UserId:    userID,
			SessionId: sessionID,
		}

		ctx := context.WithValue(r.Context(), usercontextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func UserContext(r *http.Request) *model.UserCxt {
	user, _ := r.Context().Value(usercontextKey).(*model.UserCxt)
	return user
}
