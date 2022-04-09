package identity

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/hamed-lohi/user-manage/customerror"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Role uint

const (
	//Unknown Role = iota
	Guest Role = iota + 1
	Member
	Moderator
	Admin
)

type (
	JWTConfig struct {
		Skipper    Skipper
		SigningKey interface{}
	}
	Skipper      func(c echo.Context) bool
	jwtExtractor func(echo.Context) (string, error)
)

var (
	ErrJWTMissing = echo.NewHTTPError(http.StatusUnauthorized, "missing or malformed jwt")
	ErrJWTInvalid = echo.NewHTTPError(http.StatusForbidden, "invalid or expired jwt")
)

var JWTSecret = []byte("!-!SECRET!-!")

func GenerateJWT(id primitive.ObjectID, roles []Role) string {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = id
	claims["roles"] = roles
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	t, _ := token.SignedString(JWTSecret)
	return t
}

func JWT(key interface{}) echo.MiddlewareFunc {
	c := JWTConfig{}
	c.SigningKey = key
	return JWTWithConfig(c)
}

func JWTWithConfig(config JWTConfig) echo.MiddlewareFunc {
	extractor := jwtFromHeader("Authorization", "Token")
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			auth, err := extractor(c)
			if err != nil {
				if config.Skipper != nil {
					if config.Skipper(c) {
						return next(c)
					}
				}
				return c.JSON(http.StatusUnauthorized, customerror.NewError(err))
			}
			token, err := jwt.Parse(auth, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return config.SigningKey, nil
			})
			if err != nil {
				return c.JSON(http.StatusForbidden, customerror.NewError(ErrJWTInvalid))
			}
			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				userID, _ := primitive.ObjectIDFromHex(claims["id"].(string))
				userRoles := claims["roles"].([]interface{})
				c.Set("user", userID)
				c.Set("roles", userRoles)
				return next(c)
			}
			return c.JSON(http.StatusForbidden, customerror.NewError(ErrJWTInvalid))
		}
	}
}

// jwtFromHeader returns a `jwtExtractor` that extracts token from the request header.
func jwtFromHeader(header string, authScheme string) jwtExtractor {
	return func(c echo.Context) (string, error) {
		auth := c.Request().Header.Get(header)

		l := len(authScheme)
		if len(auth) > l+1 && auth[:l] == authScheme {
			return auth[l+1:], nil
		}
		return "", ErrJWTMissing
	}
}

// Route level middleware
func CheckAccessByRole(r Role) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			//fmt.Println(c.Get("roles"))
			roles, ok := c.Get("roles").([]interface{})

			if !ok {
				return c.JSON(http.StatusForbidden, "not access")
			}

			for _, rol := range roles {
				if Role(rol.(float64)) >= r {
					return next(c)
				}
			}

			return c.JSON(http.StatusForbidden, "not access")
		}
	}
}

// --------------------------------------------------------------------------------- not used

func CheckAccess(r Role) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			readCookie("Auth", c)
			return nil
		}
	}

}

func CheckloginInfo(r Role) middleware.BasicAuthValidator {
	return func(username, password string, c echo.Context) (bool, error) {
		if username == "joe" && password == "secret" {
			return true, nil
		}
		return false, nil
	}

}

func Cookie(key interface{}) echo.MiddlewareFunc {
	// c := JWTConfig{}
	// c.SigningKey = key
	// return JWTWithConfig(c)
	return nil
}

func writeCookie(cName string, c echo.Context) error {
	cookie := new(http.Cookie)
	cookie.Name = cName
	cookie.Value = "jon"
	cookie.Expires = time.Now().Add(24 * time.Hour)
	c.SetCookie(cookie)
	return c.String(http.StatusOK, "write a cookie")
}

func readCookie(cName string, c echo.Context) error {
	cookie, err := c.Cookie(cName)
	if err != nil {
		return err
	}
	fmt.Println(cookie.Name)
	fmt.Println(cookie.Value)
	return c.String(http.StatusOK, "read a cookie")
}

func readAllCookies(c echo.Context) error {
	for _, cookie := range c.Cookies() {
		fmt.Println(cookie.Name)
		fmt.Println(cookie.Value)
	}
	return c.String(http.StatusOK, "read all the cookies")
}
