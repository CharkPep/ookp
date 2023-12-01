package auth

import (
	"context"
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"gorm.io/gorm"
	"net/http"
	"ookp/src/model"
	"ookp/src/providers/datasource"
	firebase2 "ookp/src/providers/firebase"
	"time"
)

//DOTO refactor to service layer/ usecase
//Should not be like this

type AuthController struct {
	auth *auth.Client
	db   *gorm.DB
}

func NewAuthController(db *datasource.Datasource, firebase *firebase2.FirebaseProvider) *AuthController {
	client, err := firebase.GetApp().Auth(context.Background())
	if err != nil {
		panic(err)
	}
	return &AuthController{
		auth: client,
		db:   db.DB,
	}
}

func (a *AuthController) login(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "login.html", gin.H{
		"title": "Login",
	})
}

func (a *AuthController) register(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "register.html", gin.H{
		"title": "Register",
	})
}

func (a *AuthController) createSession(ctx *gin.Context) {
	idtoken := ctx.GetHeader("Authorization")
	if idtoken == "" {
		http.Error(ctx.Writer, "Unauthorized", http.StatusUnauthorized)
		ctx.Abort()
	}
	expiresIn := time.Hour * 24 * 5

	decoded, err := a.auth.VerifyIDToken(ctx, idtoken)
	if err != nil {
		http.Error(ctx.Writer, "Unauthorized", http.StatusUnauthorized)
		ctx.Abort()
	}

	user := &model.User{
		UID: decoded.UID,
	}

	a.db.Transaction(func(tx *gorm.DB) error {
		result := tx.First(user)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				tx.Create(user)
			} else {
				return result.Error
			}
		}

		return nil
	})

	cookie, err := a.auth.SessionCookie(ctx, idtoken, expiresIn)
	if err != nil {
		http.Error(ctx.Writer, "Unauthorized", http.StatusUnauthorized)
		return
	}

	ctx.SetCookie("token", cookie, int(expiresIn.Seconds()), "/", "/", false, true)
}

func (a *AuthController) logout(ctx *gin.Context) {
	cookie, err := ctx.Cookie("token")
	if err != nil {
		ctx.Redirect(302, "/login")
		ctx.Abort()
	}

	a.auth.RevokeRefreshTokens(context.Background(), cookie)
	ctx.SetCookie("token", "", -1, "/", "/", false, true)
	ctx.Redirect(302, "/login")
}

var Module = fx.Module("auth-controller",
	fx.Provide(NewAuthController),
	fx.Invoke(func(authController *AuthController, app *model.App) {
		app.Router.GET("/login", authController.login)
		app.Router.GET("/register", authController.register)
		app.Router.POST("/login", authController.createSession)
		app.Router.GET("/logout", authController.logout)
	}))
