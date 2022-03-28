package handler

import (
	mdl "github.com/hamed-lohi/user-management/middleware"
	"github.com/hamed-lohi/user-management/model"
	"github.com/hamed-lohi/user-management/utils"
	"github.com/labstack/echo/v4"
)

func (h *Handler) Register(v1 *echo.Group) {

	guestUsers := v1.Group("/users")
	guestUsers.POST("", h.SignUp)
	guestUsers.POST("/login", h.Login)

	jwtMiddleware := mdl.JWT(utils.JWTSecret)
	user := v1.Group("/user", jwtMiddleware)

	// user.Use(middleware.JWTWithConfig(
	// 	middleware.JWTConfig{
	// 		Skipper: func(c echo.Context) bool {
	// 			if c.Request().Method == "GET" && c.Path() != "/api/user" {
	// 				return true
	// 			}
	// 			return false
	// 		},
	// 		SigningKey: utils.JWTSecret,
	// 	},
	// ))

	user.GET("", h.ListUser, mdl.CheckAccessByRole(model.Admin))
	user.POST("", h.InsertUser, mdl.CheckAccessByRole(model.Admin))
	user.DELETE("/:id", h.DeleteUser, mdl.CheckAccessByRole(model.Admin))
	user.PUT("", h.UpdateProfile)
	user.PUT("/:id", h.UpdateUser, mdl.CheckAccessByRole(model.Admin))
	user.GET("/info", h.CurrentUser)
}
