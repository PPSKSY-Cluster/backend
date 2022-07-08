package auth

import (
	"crypto/rsa"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PPSKSY-Cluster/backend/db"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	JWTKeypair *rsa.PrivateKey
}

var authInstance Auth

func InitAuth() error {
	keypair, err := generateKeyPair()
	if err != nil {
		return err
	}
	authInstance.JWTKeypair = keypair

	return nil
}

// Helper to generate a JWT token on login
func CheckCredentials() func(username, password string) (string, error) {
	return func(username, password string) (string, error) {
		user, err := db.GetUserCredentials(username)
		if err != nil {
			return "", err
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err != nil {
			return "", err
		}

		token := jwt.New(jwt.SigningMethodRS256)

		claims := token.Claims.(jwt.MapClaims)
		claims["username"] = username
		claims["admin"] = true
		claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // expires in 24 hours
		t, err := token.SignedString(authInstance.JWTKeypair)
		if err != nil {
			return "", err
		}

		return t, nil
	}
}

// Helper to generate the hashed password on user creation
func HashPW(password string) (string, error) {
	cost, err := strconv.Atoi(os.Getenv("BCRYPT_COST"))
	if err != nil {
		return "", err
	}

	hashedPW, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}

	return string(hashedPW), nil
}

// Middleware that checks if the user is authenticated
func CheckToken() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		token := c.Get("Authorization")
		token = strings.Replace(token, "Bearer ", "", 1)
		if token == "" {
			c.JSON(bson.M{"Message": "No JWT token provided"})
			return c.SendStatus(401)
		}

		_, err := jwt.ParseWithClaims(token, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
			return &authInstance.JWTKeypair.PublicKey, nil
		})

		if err != nil {
			c.JSON(bson.M{"Message": err.Error()})
			return c.SendStatus(401)
		}

		return c.Next()
	}
}
