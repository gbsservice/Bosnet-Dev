package middleware

import (
	"net/http"
	"api_kino/app/provider"
	"api_kino/config/app"
	"api_kino/config/constant"
	"api_kino/config/database"
	"api_kino/service/jwt_auth"
	"api_kino/service/web"
	"time"

	"github.com/gin-gonic/gin"
)

func Api() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(constant.RequestTime, time.Now())
	}
}

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Request.Header.Get("Authorization")
		jwtToken, err := jwt_auth.ValidateToken(tokenString)
		if err != nil {
			web.Response(c, http.StatusUnauthorized, web.H{
				Error: constant.ErrorLogin,
			})
			c.Abort()
			return
		}
		db := database.DB
		auth, err := provider.GetUser(db, jwtToken.UserID)
		if err != nil {
			web.Response(c, http.StatusUnauthorized, web.H{
				Error: constant.ErrorLogin,
			})
			c.Abort()
			return
		}
		// if auth.Password != jwtToken.Password {
		// 	web.Response(c, http.StatusUnauthorized, web.H{
		// 		Error: constant.ErrorLogin,
		// 	})
		// 	c.Abort()
		// 	return
		// }
		c.Set(constant.JwtClaim, jwtToken)
		c.Set(constant.Auth, auth)
	}
}

func AuthKey() gin.HandlerFunc {
	return func(c *gin.Context) {
		authKey := c.Request.Header.Get("Key")
		if authKey != app.Config().AuthApiKey {
			web.Response(c, http.StatusUnauthorized, web.H{
				Error: constant.ErrorToken,
			})
			c.Abort()
			return
		}
	}
}

func AuthToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Request.Header.Get("Authorization")
		jwtToken, err := jwt_auth.ValidateToken(tokenString)
		if err != nil {
			web.Response(c, http.StatusUnauthorized, web.H{
				Error: constant.ErrorLogin,
			})
			c.Abort()
			return
		}
		c.Set(constant.JwtClaim, jwtToken)
	}
}
