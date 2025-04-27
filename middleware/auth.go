package middleware

import (
	"fmt"
	"simple-golang-tdd/utils"

	"github.com/gin-gonic/gin"
)

// / JWTAuthMiddleware is the middleware to check JWT validity and extract the username
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token from the Authorization header
		tokenString, err := utils.ExtractTokenFromHeader(c)
		if err != nil {
			utils.ErrorResponse(c, 401, err.Error()) // Send error response if token is missing or malformed
			c.Abort()                                // Abort further processing
			return
		}

		// Validate the token and get claims
		claims, err := utils.ValidateToken(tokenString, "access")
		if err != nil {
			utils.ErrorResponse(c, 401, fmt.Sprintf("Unauthorized: %v", err))
			c.Abort() // Abort further processing
			return
		}

		// Extract the username from the token claims (sub field)
		username, ok := claims["sub"].(string)
		fmt.Println("Username from token:", username)
		if !ok {
			utils.ErrorResponse(c, 401, "Unauthorized: Invalid username in token")
			c.Abort() // Abort further processing
			return
		}

		// Optionally: Store the username in the context for further use in your handlers
		c.Set("username", username)

		// Continue to the next handler
		c.Next()
	}
}
