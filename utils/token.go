package utils

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/hamed-lohi/user-management/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var JWTSecret = []byte("!-!SECRET!-!")

func GenerateJWT(id primitive.ObjectID, roles []model.Role) string {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = id
	claims["roles"] = roles
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	t, _ := token.SignedString(JWTSecret)
	return t
}
