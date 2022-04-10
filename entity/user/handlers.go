package user

import (
	"net/http"

	"github.com/hamed-lohi/user-manage/customerror"
	"github.com/hamed-lohi/user-manage/identity"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Login godoc
// @Summary Login for existing user
// @Description Login for existing user
// @ID login
// @Tags user
// @Accept  json
// @Produce  json
// @Param user body userLoginRequest true "Credentials to use"
// @Success 200 {object} userResponse
// @Failure 400 {object} customerror.Error
// @Failure 401 {object} customerror.Error
// @Failure 422 {object} customerror.Error
// @Failure 404 {object} customerror.Error
// @Failure 500 {object} customerror.Error
// @Router /users/login [post]
func Login(c echo.Context) error {
	req := &userLoginRequest{}
	if err := req.bind(c); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, customerror.NewError(err))
	}
	u, err := store.GetByEmail(req.User.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, customerror.NewError(err))
	}
	if u == nil {
		return c.JSON(http.StatusForbidden, customerror.AccessForbidden())
	}
	if !u.CheckPassword(req.User.Password) {
		return c.JSON(http.StatusForbidden, customerror.AccessForbidden())
	}

	result := echo.Map{
		"id":       u.ID,
		"username": u.Username,
		"email":    u.Email,
		"bio":      u.Bio,
		"roles":    u.Roles,
		"token":    identity.GenerateJWT(u.ID, u.Roles),
	}

	return c.JSON(http.StatusOK, result) // newUserResponse(u, true)
}

// SignUp godoc
// @Summary Register a new user
// @Description Register a new user
// @ID sign-up
// @Tags user
// @Accept  json
// @Produce  json
// @Param user body userRegisterRequest true "User info for registration"
// @Success 201 {object} userResponse
// @Failure 400 {object} customerror.Error
// @Failure 404 {object} customerror.Error
// @Failure 500 {object} customerror.Error
// @Router /users [post]
func SignUp(c echo.Context) error {
	var u User
	req := &userRegisterRequest{}
	if err := req.bind(c, &u); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, customerror.NewError(err))
	}
	u.Roles = []identity.Role{identity.Guest}
	if err := store.Create(&u); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, customerror.NewError(err))
	}

	result := echo.Map{
		"id":       u.ID,
		"username": u.Username,
		"email":    u.Email,
		"bio":      u.Bio,
		"roles":    u.Roles,
		"token":    identity.GenerateJWT(u.ID, u.Roles),
	}

	return c.JSON(http.StatusCreated, result) // newUserResponse(&u, true)
}

// SignUp godoc
// @Summary Add new user
// @Description Add new user
// @ID add-user
// @Tags user
// @Accept  json
// @Produce  json
// @Param user body userRegisterRequest true "User info for registration"
// @Success 201 {object} userResponse
// @Failure 400 {object} customerror.Error
// @Failure 404 {object} customerror.Error
// @Failure 500 {object} customerror.Error
// @Security ApiKeyAuth
// @Router /user [post]
func InsertUser(c echo.Context) error {
	var u User
	req := &userInsertRequest{}
	if err := req.bind(c, &u); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, customerror.NewError(err))
	}
	if err := store.Create(&u); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, customerror.NewError(err))
	}

	// result := echo.Map{
	// 	"id":       u.ID,
	// 	"username": u.Username,
	// 	"email":    u.Email,
	// 	"bio":      u.Bio,
	// 	"roles":    u.Roles,
	// }

	return c.JSON(http.StatusCreated, newUserResponse(&u, false))
}

// UpdateProfile godoc
// @Summary Update profile
// @Description Update user information for current user
// @ID update-profile
// @Tags user
// @Accept  json
// @Produce  json
// @Param user body userUpdateRequest true "User details to update. At least **one** field is required."
// @Success 200 {object} userResponse
// @Failure 400 {object} customerror.Error
// @Failure 401 {object} customerror.Error
// @Failure 422 {object} customerror.Error
// @Failure 404 {object} customerror.Error
// @Failure 500 {object} customerror.Error
// @Security ApiKeyAuth
// @Router /user [put]
func UpdateProfile(c echo.Context) error {
	u, err := store.GetByID(userIDFromToken(c))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, customerror.NewError(err))
	}

	return updateUser(c, u)
}

// UpdateUser godoc
// @Summary Update current user
// @Description Update user information for current user
// @ID update-user
// @Tags user
// @Accept  json
// @Produce  json
// @Param        id   path      string  true  "User ID"
// @Param user body userUpdateRequest true "User details to update. At least **one** field is required."
// @Success 200 {object} userResponse
// @Failure 400 {object} customerror.Error
// @Failure 401 {object} customerror.Error
// @Failure 422 {object} customerror.Error
// @Failure 404 {object} customerror.Error
// @Failure 500 {object} customerror.Error
// @Security ApiKeyAuth
// @Router /user/{id} [put]
func UpdateUser(c echo.Context) error {

	objId, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, customerror.NewError(err))
	}
	u, err := store.GetByID(objId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, customerror.NewError(err))
	}

	return updateUser(c, u)
}

func updateUser(c echo.Context, u *User) error {

	if u == nil {
		return c.JSON(http.StatusNotFound, customerror.NotFound())
	}
	req := newUserUpdateRequest()
	req.populate(u)
	if err := req.bind(c, u); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, customerror.NewError(err))
	}
	if err := store.Update(u); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, customerror.NewError(err))
	}

	// result := echo.Map{
	// 	"id":       u.ID,
	// 	"username": u.Username,
	// 	"email":    u.Email,
	// 	"bio":      u.Bio,
	// 	"roles":    u.Roles,
	// }

	return c.JSON(http.StatusOK, newUserResponse(u, false))
}

// DeleteUser godoc
// @Summary Delete a user
// @Description Delete a user
// @ID delete-user
// @Tags user
// @Accept  json
// @Produce  json
// @Param        id   path      string  true  "User ID"
// @Success 201 {object} userResponse
// @Failure 400 {object} customerror.Error
// @Failure 404 {object} customerror.Error
// @Failure 500 {object} customerror.Error
// @Security ApiKeyAuth
// @Router /user/{id} [delete]
func DeleteUser(c echo.Context) error {
	//id := []byte(c.Param("id"))
	objId, _ := primitive.ObjectIDFromHex(c.Param("id"))
	var u User
	if err := store.Delete(objId); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, customerror.NewError(err))
	}

	return c.JSON(http.StatusOK, newUserResponse(&u, false))
}

// CurrentUser godoc
// @Summary Get the current user
// @Description Gets the currently logged-in user
// @ID current-user
// @Tags user
// @Accept  json
// @Produce  json
// @Success 200 {object} userResponse
// @Failure 400 {object} customerror.Error
// @Failure 401 {object} customerror.Error
// @Failure 422 {object} customerror.Error
// @Failure 404 {object} customerror.Error
// @Failure 500 {object} customerror.Error
// @Security ApiKeyAuth
// @Router /user/info [get]
func CurrentUser(c echo.Context) error {

	u, err := store.GetByID(userIDFromToken(c))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, customerror.NewError(err))
	}
	if u == nil {
		return c.JSON(http.StatusNotFound, customerror.NotFound())
	}

	return c.JSON(http.StatusOK, newUserResponse(u, false))
}

// UserList godoc
// @Summary Users list
// @Description Users list for Admin
// @ID list-user
// @Tags user
// @Accept  json
// @Produce  json
// @Success 200 {object} []User
// @Failure 400 {object} customerror.Error
// @Failure 401 {object} customerror.Error
// @Failure 422 {object} customerror.Error
// @Failure 404 {object} customerror.Error
// @Failure 500 {object} customerror.Error
// @Security ApiKeyAuth
// @Router /user [get]
func ListUser(c echo.Context) error {
	users, err := store.GetUserList()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, customerror.NewError(err))
	}
	if users == nil {
		return c.JSON(http.StatusNotFound, customerror.NotFound())
	}

	r := new(userList)
	cr := User{}
	r.Users = make([]User, 0)
	for _, i := range *users {
		cr.Username = i.Username
		cr.Email = i.Email
		cr.ID = i.ID
		cr.Bio = i.Bio
		cr.Roles = i.Roles

		r.Users = append(r.Users, cr)
	}

	return c.JSON(http.StatusOK, r)
}

func userIDFromToken(c echo.Context) primitive.ObjectID {

	objId, ok := c.Get("user").(primitive.ObjectID)
	if !ok {
		// log.Fatal(err)
		return primitive.ObjectID{}
	}
	return objId
}
