package middleware

import (
	"Oauth/database/helper"
	"Oauth/handler"
	"Oauth/models"
	"Oauth/utilities"
	"context"
	"encoding/json"
	"fmt"
	"github.com/form3tech-oss/jwt-go"
	"github.com/go-chi/chi/v5"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"net/http"
	"runtime/debug"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")

		claims := models.Claims{}

		tkn, err1 := jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (interface{}, error) {
			return handler.JwtKey, nil
		})
		if err1 != nil {
			if err1 == jwt.ErrSignatureInvalid {
				logrus.Printf("Signature invalid:%v", err1)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			logrus.Printf("ParseErr:%v", err1)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if !tkn.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			logrus.Printf("token is invalid")
			return
		}

		_, err := helper.CheckSession(claims.ID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logrus.Printf("session expired:%v", err)
			return
		}

		userID := claims.ID
		role := claims.Role

		value := models.ContextValues{ID: userID, Role: role}
		ctx := context.WithValue(r.Context(), utilities.UserContextKey, value)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contextValues, ok := r.Context().Value(utilities.UserContextKey).(models.ContextValues)

		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			logrus.Printf("AdminMiddleware:Context for ID:%v", ok)
			return
		}

		if contextValues.Role != "seller" {
			w.WriteHeader(http.StatusUnauthorized)
			logrus.Printf("Role invalid")
			_, err := w.Write([]byte("ERROR: Role mismatch"))

			if err != nil {
				return
			}

			return
		}

		next.ServeHTTP(w, r)
	})
}

// corsOptions setting up routes for cors
func corsOptions() *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Access-Token", "importDate", "X-Client-Version", "Cache-Control", "Pragma", "x-started-at", "x-api-key", "token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
	})
}

// CommonMiddlewares middleware common for all routes
func CommonMiddlewares() chi.Middlewares {
	return chi.Chain(
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Add("Content-Type", "application/json")
				next.ServeHTTP(w, r)
			})
		},
		corsOptions().Handler,
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				defer func() {
					err := recover()
					if err != nil {
						logrus.Errorf("Request Panic err: %v", err)
						jsonBody, _ := json.Marshal(map[string]string{
							"error": "There was an internal server error",
							"trace": fmt.Sprintf("%+v", err),
							"stack": string(debug.Stack()),
						})
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(http.StatusInternalServerError)
						_, err := w.Write(jsonBody)
						if err != nil {
							logrus.Errorf("Failed to send response from middleware with error: %+v", err)
						}
					}
				}()
				next.ServeHTTP(w, r)
			})
		},
	)
}
