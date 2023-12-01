package auth

import (
	"context"
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	firebase2 "ookp/src/providers/firebase"
)

type AuthMiddleware struct {
	auth *auth.Client
}

func NewAuthMiddleware(firebase firebase2.FirebaseProvider) *AuthMiddleware {
	client, err := firebase.GetApp().Auth(context.Background())
	if err != nil {
		panic(err)
	}

	return &AuthMiddleware{
		auth: client,
	}
}

func (md *AuthMiddleware) Validate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authToken, err := ctx.Cookie("token")
		if err != nil {
			ctx.Redirect(302, "/login")
			ctx.Abort()
		}

		token, err := md.auth.VerifySessionCookieAndCheckRevoked(ctx, authToken)
		if err != nil {
			ctx.Redirect(302, "/login")
			ctx.Abort()
		}

		ctx.Set("token", token)
	}

}
