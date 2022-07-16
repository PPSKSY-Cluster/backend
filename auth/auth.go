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
	"go.mongodb.org/mongo-driver/bson/primitive"
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
		claims["userid"] = user.ID
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

// Middleware that checks if the user is authenticated and authorized under
// under given restriction
func CheckToken(restrictedTo db.UserType) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		token := c.Get("Authorization")
		token = strings.Replace(token, "Bearer ", "", 1)
		if token == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "No JWT token provided")
		}

		claims := jwt.MapClaims{}
		_, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
			return &authInstance.JWTAccessKeypair.PublicKey, nil
		})

		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, err.Error())
		}

		id, _ := primitive.ObjectIDFromHex(claims["userid"].(string))
		user, err := db.GetUserById(id)
		if err != nil {
			c.JSON(bson.M{"Message": "Could not find user"})
			return c.SendStatus(404)
		}

		if !userTypeIncludes(user.Type, restrictedTo) {
			c.JSON(bson.M{"Message": "Not authorized to access this route"})
			return c.SendStatus(401)
		}

		c.Locals("jwtUserId", user.ID)
		c.Locals("jwtUserType", user.Type)

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

// does the given type include the rights of the expected type
func userTypeIncludes(givenType db.UserType, expectedType db.UserType) bool {
	hasSuperAdminRights := db.SuperAdminUT == givenType
	hasAdminRights := hasSuperAdminRights || db.AdminUT == givenType
	switch {
	case expectedType == db.UserUT:
		return true
	case expectedType == db.AdminUT && hasAdminRights:
		return true
	case expectedType == db.SuperAdminUT && hasSuperAdminRights:
		return true
	default:
		return false
	}
}
