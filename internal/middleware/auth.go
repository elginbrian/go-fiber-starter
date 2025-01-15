package middleware

import (
	"errors"
	"fiber-starter/pkg/response"
	"log"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

func TokenValidationMiddleware(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return response.Error(c, "Authorization header is required", fiber.StatusUnauthorized)
		}

		parts := strings.Split(authHeader, "Bearer ")
		if len(parts) != 2 {
			return response.Error(c, "Bearer token is required", fiber.StatusUnauthorized)
		}

		tokenString := parts[1]
		if secret == "" {
			log.Println("Error: JWT_SECRET_KEY is not set in the environment variables")
			return response.Error(c, "Internal server error: Missing JWT secret", fiber.StatusInternalServerError)
		}
		log.Println("JWT Secret Key loaded successfully")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(secret), nil
		})

		if err != nil {
			log.Println("Error parsing token:", err)
			return response.Error(c, "Invalid token: "+err.Error(), fiber.StatusUnauthorized)
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			log.Println("Token is valid")

			if exp, ok := claims["exp"].(float64); ok {
				if time.Unix(int64(exp), 0).Before(time.Now()) {
					log.Println("Token has expired")
					return response.Error(c, "Token has expired", fiber.StatusUnauthorized)
				}
			} else {
				log.Println("Invalid token claims")
				return response.Error(c, "Invalid token claims", fiber.StatusUnauthorized)
			}

			if userID, ok := claims["user_id"].(float64); ok {
				log.Printf("User ID extracted from token: %v", userID)
				c.Locals("user_id", int(userID))
			} else {
				log.Println("User ID not found in token")
				return response.Error(c, "User ID not found in token", fiber.StatusUnauthorized)
			}
		} else {
			log.Println("Invalid or expired token")
			return response.Error(c, "Invalid or expired token", fiber.StatusUnauthorized)
		}

		log.Println("Proceeding to the next handler")
		return c.Next()
	}
}