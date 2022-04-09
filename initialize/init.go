package initialize

import (
	"github.com/hamed-lohi/user-manage/db"
	"github.com/hamed-lohi/user-manage/entity/user"
	"github.com/labstack/echo/v4"

	echoSwagger "github.com/swaggo/echo-swagger"
)

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

	registerHandlers(v1, dp)
	user.Seed(dp)
	e.Logger.Fatal(e.Start("127.0.0.1:8585"))

}

func registerHandlers(v1 *echo.Group, dp *db.DBProvider) {
	user.RegisterHandlers(v1, dp)
	// product.RegisterHandlers(v1, dp)

}
