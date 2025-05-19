package middleware

import (
	"errors" // Import errors
	"gatherly-app/service"

	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// --- Keep your AuthMiddleware interface definition ---
type AuthMiddleware interface {
	RequireToken() gin.HandlerFunc
}

// --- Keep your authMiddleware struct definition ---
type authMiddleware struct {
	jwtService service.JwtService
}

// --- Keep your NewAuthMiddleware constructor ---
func NewAuthMiddleware(jwtService service.JwtService) AuthMiddleware {
	return &authMiddleware{jwtService: jwtService}
}

// --- Replace your RequireToken function with this one ---
func (m *authMiddleware) RequireToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("AuthMiddleware: Running...")

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Println("AuthMiddleware: Authorization header missing.")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			log.Println("AuthMiddleware: Bearer prefix missing.")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Bearer token is required"})
			return
		}
		// log.Printf("AuthMiddleware: Processing token: %s...\n", tokenString[:min(10, len(tokenString))])

        // --- MODIFIED SECTION ---
		// Call the updated ValidateToken which returns *Claims, error
		claims, err := m.jwtService.ValidateToken(tokenString)
		if err != nil {
			// Log the specific validation/parsing error received from ValidateToken
			log.Printf("AuthMiddleware: Token validation/parsing failed: %v\n", err)
			// Return a generic error message to the client
			// Check if the error indicates expiration for a potentially different message
			if errors.Is(err, errors.New("token has expired")) { // Check against the specific error from ValidateToken
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token has expired"})
			} else {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			}
			return
		}
        // --- END MODIFIED SECTION ---

		// If we get here, claims are non-nil and valid
		log.Printf("AuthMiddleware: Claims retrieved successfully. UserID: %d, Email: %s, Role: %s, Latitude: %v, Longitude: %v\n", claims.UserID, claims.Email, claims.Role, claims.Latitude, claims.Longitude)

        // Check if claims.UserID is zero, which might indicate an issue if IDs start from 1
        if claims.UserID == 0 {
             log.Printf("AuthMiddleware: Warning - Parsed UserID is 0. Claims: %+v\n", claims)
        }

		c.Set("userID", claims.UserID) // Set context value
		c.Set("userEmail", claims.Email) // Optional: Set other values
        c.Set("userRole", claims.Role) // Optional: Set other values
		c.Set("userLat", claims.Latitude)
		c.Set("userLon", claims.Longitude)
		log.Println("AuthMiddleware: Context values set (userID, userEmail, userRole).")

		log.Println("AuthMiddleware: Calling c.Next()")
		c.Next() // Continue to the next handler (your controller)
	}
}

