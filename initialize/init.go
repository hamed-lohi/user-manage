package initialize

import (
	"log"
	"net/http"
	_ "net/http/pprof"

	"github.com/hamed-lohi/user-manage/db"
	"github.com/hamed-lohi/user-manage/entity/user"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	echo_log "github.com/labstack/gommon/log"
	"gopkg.in/go-playground/validator.v9"

	echoSwagger "github.com/swaggo/echo-swagger"
)

type Validator struct {
	validator *validator.Validate
}

func NewValidator() *Validator {
	return &Validator{
		validator: validator.New(),
	}
}

func (v *Validator) Validate(i interface{}) error {
	return v.validator.Struct(i)
}

func NewEcho() *echo.Echo {
	e := echo.New()
	e.Logger.SetLevel(echo_log.DEBUG)
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Logger())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))
	e.Validator = NewValidator()
	return e
}

func InitializeWebServer() {

	e := NewEcho()

	// // Group level middleware
	// g := e.Group("/admin")
	// g.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
	// if username == "joe" && password == "secret" {
	// 	return true, nil
	// }
	// return false, nil
	// }))

	e.GET("/swagger/*", echoSwagger.WrapHandler) //

	v1 := e.Group("/api")

	dp := db.NewDBProvider()

	//us := store.NewUserStore(dp)
	//h := handler.NewHandler(us)
	//h.Register(v1)

	//e.GET("/debug/pprof/*", echo.WrapHandler(http.DefaultServeMux))

	registerHandlers(v1, dp)
	user.Seed(dp)

	go func() {

		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	e.Logger.Fatal(e.Start("127.0.0.1:8585"))

}

func registerHandlers(v1 *echo.Group, dp *db.DBProvider) {
	user.RegisterHandlers(v1, dp)
	// product.RegisterHandlers(v1, dp)

}
