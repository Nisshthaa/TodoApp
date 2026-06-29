package middlewares

import (
	"context"
	"errors"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/yourusername/my-project/database/dbHelper"
	"github.com/yourusername/my-project/models"
	"github.com/yourusername/my-project/utils"
)

type ContextKeys string

const (
	userContext ContextKeys = "userContext"
)

func Authenticate(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		apiKey := r.Header.Get("x-api-key")
		if apiKey == "" {
			utils.RespondError(w, http.StatusUnauthorized, nil, "token header missing")
			return
		}

		token, err := jwt.Parse(apiKey, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid signing ng method")
			}
			return []byte(os.Getenv("JWT_SECRET_KEY")), nil
		})
		
		if err != nil || !token.Valid {
			utils.RespondError(w, http.StatusUnauthorized, err, "invalid token")
			return
		}

		claimValues, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			utils.RespondError(w, http.StatusUnauthorized, nil, "invalid token claims")
			return
		}

		sessionID := claimValues["sessionId"].(string)
		archivedAt, err := dbHelper.GetArchivedAt(sessionID)
		if err != nil {
			utils.RespondError(w, http.StatusInternalServerError, err, "internal server error")
			return
		}

		if archivedAt != nil {
			utils.RespondError(w, http.StatusUnauthorized, nil, "invalid token user already logged out")
			return
		}

		user := &models.UserCtx{
			UserID:    claimValues["userId"].(string),
			SessionID: sessionID,
		}

		ctx := context.WithValue(r.Context(), userContext, user)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func UserContext(r *http.Request) *models.UserCtx {
	user := r.Context().Value(userContext).(*models.UserCtx)
	return user
}
