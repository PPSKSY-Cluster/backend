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
	JWTAccessKeypair  *rsa.PrivateKey
	JWTRefreshKeypair *rsa.PrivateKey
}

var authInstance Auth

func InitAuth() error {
	keypair, err := generateKeyPair()
	if err != nil {
		return err
	}
	authInstance.JWTAccessKeypair = keypair

	keypair, err = generateKeyPair()
	if err != nil {
		return err
	}
	authInstance.JWTRefreshKeypair = keypair

	return nil
}

// Helper to generate a JWT token on login
func CheckCredentials() func(username, password string) (db.User, string, error) {
	return func(username, password string) (db.User, string, error) {
		user, err := db.GetUserWithCredentials(username)
		if err != nil {
			return db.User{}, "", err
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err != nil {
			return db.User{}, "", err
		}

		token := jwt.New(jwt.SigningMethodRS256)

		claims := token.Claims.(jwt.MapClaims)
		claims["username"] = username
		claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // expires in 24 hours
		t, err := token.SignedString(authInstance.JWTRefreshKeypair)
		if err != nil {
			return db.User{}, "", err
		}

		user.Password = ""
		return user, t, nil
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
			return &authInstance.JWTAccessKeypair.PublicKey, nil
		})

		if err != nil {
			c.JSON(bson.M{"Message": err.Error()})
			return c.SendStatus(401)
		}

		return c.Next()
	}
}

// takes a refresh token as parameter and returns a new access token
func RefreshAccessToken(token string) (string, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		return &authInstance.JWTRefreshKeypair.PublicKey, nil
	})

	if err != nil {
		return "", err
	}

	claims["exp"] = time.Now().Add(time.Hour * 1).Unix() // expires in 1 hour
	accessToken := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	accessStr, err := accessToken.SignedString(authInstance.JWTAccessKeypair)
	if err != nil {
		return "", err
	}

	return accessStr, nil
}
