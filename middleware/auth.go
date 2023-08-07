// middleware/auth.go
package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"project/util"

	"github.com/gin-gonic/gin"
)

var jwtKey = []byte("my_secret_key")

// Authenticate middleware checks for a valid JWT token in the Authorization header.
// func Authenticate(next http.HandlerFunc) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		tokenString := r.Header.Get("Authorization")
// 		if tokenString == "" {
// 			util.WriteError(w, http.StatusUnauthorized, "Unauthorized")
// 			return
// 		}

// 		userID, err := util.VerifyJWT(tokenString, jwtKey)
// 		if err != nil {
// 			util.WriteError(w, http.StatusUnauthorized, "Unauthorized")
// 			return
// 		}

// 		ctx := r.Context()
// 		ctx = context.WithValue(ctx, "userID", userID)
// 		next.ServeHTTP(w, r.WithContext(ctx))
// 	}
// }

// AuthMiddleware is a middleware function that checks for a valid JWT in the request header.
// If the token is present and valid, it sets the userID in the context for the next handler to use.
// Otherwise, it returns an error response.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// fmt.Println("PATH ::: ", c.Request.URL.Path)
		// if c.Request.URL.Path == "/home" || c.Request.URL.Path == "/add" {
		// 	c.Next()
		// 	return
		// }

		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			util.WriteError(c.Writer, http.StatusUnauthorized, "JWT token missing")

		}
		splitToken := strings.Split(tokenString, "Bearer ")
		tokenString = splitToken[1]
		fmt.Println(tokenString)

		if tokenString == "" {
			util.WriteError(c.Writer, http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}

		userID, err := util.VerifyJWT(tokenString, jwtKey)
		fmt.Println(userID)
		if err != nil {
			util.WriteError(c.Writer, http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}

		c.Set("userID", userID)
		c.Next()
	}
}
