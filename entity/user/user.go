package user

import (
	"errors"
	"fmt"
	"log"
	"os/user"

	"github.com/hamed-lohi/user-manage/db"
	"github.com/hamed-lohi/user-manage/identity"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var store *UserStore

type userList struct {
	Users []User `json:"users"`
}

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username string             `bson:"username,omitempty"`
	Email    string             `bson:"email,omitempty"`
	Password string             `bson:"password"`
	Bio      *string            //`bson:"bio"`
	Roles    []identity.Role    `bson:"roles,omitempty"`
	// Image      *string
}

func (u *User) SetPassword(plain string) error {
	h, err := u.HashPassword(plain)
	u.Password = h
	return err
}

func (u *User) HashPassword(plain string) (string, error) {
	if len(plain) == 0 {
		return "", errors.New("password should not be empty")
	}
	h, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	return string(h), err
}

func (u *User) CheckPassword(plain string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plain))
	return err == nil
}

func RegisterHandlers(v1 *echo.Group, dp *db.DBProvider) {

	store = NewStore(dp)

	guestUsers := v1.Group("/users")
	guestUsers.POST("", SignUp)
	guestUsers.POST("/login", Login)

	jwtMiddleware := identity.JWT(identity.JWTSecret)
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
	user.GET("", ListUser, identity.CheckAccessByRole(identity.Admin))
	user.POST("", InsertUser, identity.CheckAccessByRole(identity.Admin))
	user.DELETE("/:id", DeleteUser, identity.CheckAccessByRole(identity.Admin))
	user.PUT("", UpdateProfile)
	user.PUT("/:id", UpdateUser, identity.CheckAccessByRole(identity.Admin))
	user.GET("/info", CurrentUser)
}

func Seed(dp *db.DBProvider) {

	//defer cl.Dispose()
	//coll := dp.GetCollection(Users)

	var result user.User
	err := store.collection.FindOne(dp.Context, bson.M{"username": "Admin"}).Decode(&result) // context.TODO()
	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in
		// the collection.
		if err == mongo.ErrNoDocuments {
			admin := &User{
				//ID:       ,
				Username: "Admin",
				Email:    "admin@gmail.com",
				//Password: "aaaa",
				Bio:   new(string),
				Roles: []identity.Role{identity.Admin},
			}
			admin.SetPassword("aaa")
			store.collection.InsertOne(dp.Context, admin)
			return
		}
		log.Fatal(err)
	}
	fmt.Printf("found document %v", result)

	// if err := cl.Disconnect(context.TODO()); err != nil {
	// 	panic(err)
	// }
}
