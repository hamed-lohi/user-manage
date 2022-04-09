package user

import (
	"github.com/hamed-lohi/user-manage/identity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type userResponse struct {
	User struct {
		ID       primitive.ObjectID `json:"id"`
		Username string             `json:"username"`
		Email    string             `json:"email"`
		Bio      *string            `json:"bio"`
		Image    *string            `json:"image"`
		Roles    []identity.Role    `json:"roles"`
		Token    string             `json:"token"`
	} `json:"user"`
}

func newUserResponse(u *User, hasToken bool) *userResponse {
	r := new(userResponse)
	r.User.ID = u.ID
	r.User.Username = u.Username
	r.User.Email = u.Email
	r.User.Bio = u.Bio
	r.User.Roles = u.Roles

	//r.User.Image = u.Image
	if hasToken {
		r.User.Token = identity.GenerateJWT(u.ID, u.Roles)
	}

	return r
}

type usersResponse struct {
	ID       primitive.ObjectID `json:"id"`
	Username string             `json:"username"`
	Email    string             `json:"email"`
	Bio      *string            `json:"bio"`
	//Image    *string      `json:"image"`
	Roles []identity.Role `json:"roles"`
}

type userListResponse struct {
	Users []usersResponse `json:"users"`
}

func newUserListResponse(users []User) *userListResponse {
	r := new(userListResponse)
	cr := usersResponse{}
	r.Users = make([]usersResponse, 0)
	for _, i := range users {
		cr.Username = i.Username
		cr.Email = i.Email
		cr.ID = i.ID
		cr.Bio = i.Bio
		cr.Roles = i.Roles

		r.Users = append(r.Users, cr)
	}
	return r
}
