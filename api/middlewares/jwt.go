package middlewares

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mhg14/hotel-reservation/db"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func JWTAuthentication(userStore db.UserStore) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token, ok := c.GetReqHeaders()["X-Api-Token"]

		if !ok {
			fmt.Printf("Unauthorized")
		}

		claims, err := validateToken(token)
		if err != nil {
			return err
		}

		expiresFloat := claims["expires"].(float64)
		expires := int64(expiresFloat)
		if time.Now().Unix() > expires {
			return fmt.Errorf("token expired")
		}

		userID, _ := primitive.ObjectIDFromHex(claims["id"].(string))
		user, err := userStore.GetUserById(c.Context(), userID)
		fmt.Println(user)
		if err != nil {
			return fmt.Errorf("unauthorized")
		}
		// set the current authenticated user to the context
		c.Context().SetUserValue("user", user)
		return c.Next()
	}

}

func validateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("invalid signing method", token.Header["alg"])
			return fmt.Errorf("unauthorized"), nil
		}
		secret := os.Getenv("JWT_SECRET")
		fmt.Println("NEVER PRINT YOUR SECRET", secret)
		return []byte(secret), nil
	})

	if err != nil {
		fmt.Println("Failed to parse jwt token:", err)
		return nil, fmt.Errorf("unauthorized")
	}

	if !token.Valid {
		fmt.Println("invalid token")
		return nil, fmt.Errorf("unauthorized")
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return nil, fmt.Errorf("unauthorized")
	}
	return claims, nil
}
